package k8s

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/utils"
	authorizationv1 "k8s.io/api/authorization/v1"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

type CanIOptions struct {
	NoHeaders       bool
	Namespace       string
	AuthClient      authorizationv1client.AuthorizationV1Interface
	DiscoveryClient discovery.DiscoveryInterface

	Verb         string
	Resource     schema.GroupVersionResource
	Subresource  string
	ResourceName string
}

// NsExists checks if the given namespace already exists
func NsExists(namespace string, kubeconfig *string) (bool, error) {
	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		return false, err
	}
	ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if k8serror.IsNotFound(err) {
		return false, nil
	}

	if err == nil && ns != nil {
		return true, nil
	}

	return false, err
}

type CheckSAPermissionsParams struct {
	Verb     string
	Resource string
	Print    bool
}

func CheckSAPermissions(params CheckSAPermissionsParams, kubeconfig *string) (bool, error) {

	var o CanIOptions
	o.Verb = params.Verb
	o.Resource.Resource = params.Resource
	client, err := ClientSet(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	AuthClient := client.AuthorizationV1()

	var sar *authorizationv1.SelfSubjectAccessReview
	sar = &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace:   o.Namespace,
				Verb:        o.Verb,
				Group:       o.Resource.Group,
				Resource:    o.Resource.Resource,
				Subresource: o.Subresource,
				Name:        o.ResourceName,
			},
		},
	}

	response, err := AuthClient.SelfSubjectAccessReviews().Create(context.TODO(), sar, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}

	if response.Status.Allowed {
		if params.Print {
			fmt.Println("🔑 ", params.Resource, "- ✅")
		}
	} else {
		if params.Print {
			fmt.Println("🔑 ", params.Resource, "- ❌")
		}
		if len(response.Status.Reason) > 0 {
			fmt.Println(response.Status.Reason)
		}
		if len(response.Status.EvaluationError) > 0 {
			fmt.Println(response.Status.EvaluationError)
		}
	}

	return response.Status.Allowed, nil
}

// ValidNs takes a valid namespace as input from user
func ValidNs(mode string, label string, kubeconfig *string) (string, bool) {
start:
	var (
		namespace string
		nsExists  bool
	)

	if mode == "namespace" {
		fmt.Print("📁 Enter the namespace (existing) [", utils.DefaultNs, "]: ")
		fmt.Scanln(&namespace)

	} else if mode == "cluster" {
		fmt.Print("📁 Enter the namespace (new or existing) [", utils.DefaultNs, "]: ")
		fmt.Scanln(&namespace)
	} else {
		fmt.Printf("\n 🚫 No mode selected \n")
		os.Exit(1)
	}

	if namespace == "" {
		namespace = utils.DefaultNs
	}
	ok, err := NsExists(namespace, kubeconfig)
	if err != nil {
		fmt.Printf("\n 🚫 Namespace existence check failed: {%s}\n", err.Error())
		os.Exit(1)
	}
	if ok {
		if podExists(podExistsParams{namespace, label}, kubeconfig) {
			fmt.Println("🚫 Subscriber already present. Please enter a different namespace")
			goto start
		} else {
			nsExists = true
			fmt.Println("👍 Continuing with", namespace, "namespace")
		}
	} else {
		if val, _ := CheckSAPermissions(CheckSAPermissionsParams{"create", "namespace", false}, kubeconfig); !val {
			fmt.Println("🚫 You don't have permissions to create a namespace.\n🙄 Please enter an existing namespace.")
			goto start
		}
		nsExists = false
	}

	return namespace, nsExists
}

type WatchPodParams struct {
	Namespace string
	Label     string
}

// WatchPod watches for the pod status
func WatchPod(params WatchPodParams, kubeconfig *string) {
	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	watch, err := clientset.CoreV1().Pods(params.Namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: params.Label,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	for event := range watch.ResultChan() {
		p, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Fatal("unexpected type")
		}
		fmt.Println("💡 Connecting agent to Litmus Portal.")
		if p.Status.Phase == "Running" {
			fmt.Println("🏃 Agents running!!")
			watch.Stop()
			break
		}
	}
}

type PodList struct {
	Items []string
}

type podExistsParams struct {
	Namespace string
	Label     string
}

// PodExists checks if the pod with the given label already exists in the given namespace
func podExists(params podExistsParams, kubeconfig *string) bool {
	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		log.Fatal(err)
		return false
	}
	watch, err := clientset.CoreV1().Pods(params.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: params.Label,
	})
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	if len(watch.Items) >= 1 {
		return true
	}

	return false
}

type SAExistsParams struct {
	Namespace      string
	Serviceaccount string
}

// SAExists checks if the given service account exists in the given namespace
func SAExists(params SAExistsParams, kubeconfig *string) bool {
	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf("serviceaccounts \"%s\" not found", params.Serviceaccount)
	_, newErr := clientset.CoreV1().ServiceAccounts(params.Namespace).Get(context.TODO(), params.Serviceaccount, metav1.GetOptions{})
	if newErr != nil {
		if newErr.Error() == msg {
			return false
		}
		log.Fatal(newErr)
	}
	return true
}

// ValidSA gets a valid service account as input
func ValidSA(namespace string, kubeconfig *string) (string, bool) {
	var sa string
	fmt.Print("🔑 Enter service account [", utils.DefaultSA, "]: ")
	fmt.Scanln(&sa)
	if sa == "" {
		sa = utils.DefaultSA
	}
	if SAExists(SAExistsParams{namespace, sa}, kubeconfig) {
		fmt.Println("👍 Using the existing service account")
		return sa, true
	}
	return sa, false
}
