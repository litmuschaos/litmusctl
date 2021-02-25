package k8s

import (
	"context"
	"fmt"
	"log"
	"os"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	discovery "k8s.io/client-go/discovery"
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

func CheckSAPermissions(verb, resource string, print bool) (bool, error) {

	var o CanIOptions
	o.Verb = verb
	o.Resource.Resource = resource
	client, err := ClientSet()
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
		if print {
			fmt.Println("🔑 ", resource, "- ✅")
		}
	} else {
		if print {
			fmt.Println("🔑 ", resource, "- ❌")
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

func ValidateSAPermissions(mode string) {
	var pems [2]bool
	var err error
	if mode == "cluster" {
		resources := [2]string{"clusterrole", "clusterrolebinding"}
		i := 0
		for _, resource := range resources {
			pems[i], err = CheckSAPermissions("create", resource, true)
			if err != nil {
				fmt.Println(err)
			}
			i++
		}
	} else {
		resources := [2]string{"role", "rolebinding"}
		i := 0
		for _, resource := range resources {
			pems[i], err = CheckSAPermissions("create", resource, true)
			if err != nil {
				fmt.Println(err)
			}
			i++
		}
	}

	for _, pem := range pems {
		if !pem {
			fmt.Println("\n🚫 You don't have sufficient permissions.\n🙄 Please use a service account with sufficient permissions.")
			os.Exit(1)
		}
	}
	fmt.Println("\n🌟 Sufficient permissions. Registering Agent")
}
