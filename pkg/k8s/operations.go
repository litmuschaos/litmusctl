/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package k8s

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/sirupsen/logrus"
	authorizationv1 "k8s.io/api/authorization/v1"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
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
	Verb      string
	Resource  string
	Print     bool
	Namespace string
}

func CheckSAPermissions(params CheckSAPermissionsParams, kubeconfig *string) (bool, error) {

	var o CanIOptions
	o.Verb = params.Verb
	o.Resource.Resource = params.Resource
	o.Namespace = params.Namespace
	client, err := ClientSet(kubeconfig)
	if err != nil {
		return false, err
	}

	AuthClient := client.AuthorizationV1()

	sar := &authorizationv1.SelfSubjectAccessReview{
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
			utils.White_B.Print("\n🔑 ", params.Resource, " ✅")
		}
	} else {
		if params.Print {
			utils.White_B.Print("\n🔑 ", params.Resource, " ❌")
		}
		if len(response.Status.Reason) > 0 {
			utils.White_B.Println(response.Status.Reason)
		}
		if len(response.Status.EvaluationError) > 0 {
			utils.Red.Println(response.Status.EvaluationError)
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
		utils.White_B.Print("\nEnter the namespace (existing namespace) [Default: ", utils.DefaultNs, "]: ")
		fmt.Scanln(&namespace)

	} else if mode == "cluster" {
		utils.White_B.Print("\nEnter the namespace (new or existing namespace) [Default: ", utils.DefaultNs, "]: ")
		fmt.Scanln(&namespace)
	} else {
		utils.Red.Printf("\n 🚫 No mode selected \n")
		os.Exit(1)
	}

	if namespace == "" {
		namespace = utils.DefaultNs
	}
	ok, err := NsExists(namespace, kubeconfig)
	if err != nil {
		utils.Red.Printf("\n 🚫 Namespace existence check failed: {%s}\n", err.Error())
		os.Exit(1)
	}
	if ok {
		if podExists(podExistsParams{namespace, label}, kubeconfig) {
			utils.Red.Println("\n🚫 There is a Chaos Infra already present in this namespace. Please enter a different namespace")
			goto start
		} else {
			nsExists = true
			utils.White_B.Println("👍 Continuing with", namespace, "namespace")
		}
	} else {
		if val, _ := CheckSAPermissions(CheckSAPermissionsParams{"create", "namespace", false, namespace}, kubeconfig); !val {
			utils.Red.Println("🚫 You don't have permissions to create a namespace.\n Please enter an existing namespace.")
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
		utils.White_B.Println("💡 Connecting Chaos Infra to ChaosCenter.")
		if p.Status.Phase == "Running" {
			utils.White_B.Println("🏃 Chaos Infra is running!!")
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
	utils.White_B.Print("\nEnter service account [Default: ", utils.DefaultSA, "]: ")
	fmt.Scanln(&sa)
	if sa == "" {
		sa = utils.DefaultSA
	}
	if SAExists(SAExistsParams{namespace, sa}, kubeconfig) {
		utils.White_B.Print("\n👍 Using the existing service account")
		return sa, true
	}
	return sa, false
}

// Token: Authorization token
// EndPoint: Endpoint in .litmusconfig
// YamlPath: Path of yaml file
type ApplyYamlParams struct {
	Token    string
	Endpoint string
	YamlPath string
}

func ApplyYaml(params ApplyYamlParams, kubeconfig string, isLocal bool) (output string, err error) {
	path := params.YamlPath
	if !isLocal {
		path = fmt.Sprintf("%s%s/%s.yaml", params.Endpoint, params.YamlPath, params.Token)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			return "", err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		resp_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		err = ioutil.WriteFile("chaos-infra-manifest.yaml", resp_body, 0644)
		if err != nil {
			return "", err
		}
		path = "chaos-infra-manifest.yaml"
	}

	args := []string{"kubectl", "apply", "-f", path}
	if kubeconfig != "" {
		args = append(args, []string{"--kubeconfig", kubeconfig}...)
	} else {
		args = []string{"kubectl", "apply", "-f", path}
	}

	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()

	// err, can have exit status 1
	if err != nil {
		// if we get standard error then, return the same
		if errStr != "" {
			return "", fmt.Errorf(errStr)
		}

		// if not standard error found, return error
		return "", err
	}

	// If no error found, return standard output
	return outStr, nil
}

// UpgradeInfra upgrades the Chaos Infrastructure using the provided manifest and kubeconfig with the help of client-go library.
func UpgradeInfra(manifest []byte, kubeconfig string) (string, error) {

	// Get Kubernetes and dynamic clients along with the configuration.
	_, kubeClient, dynamicClient, err := getClientAndConfig(kubeconfig)
	if err != nil {
		return "", err
	}

	// Decode the manifest into a list of Unstructured resources.
	resources, err := decodeManifest(manifest)
	if err != nil {
		return "", err
	}

	// Apply the decoded resources using the dynamic client and Kubernetes client.
	err = applyResources(resources, dynamicClient, kubeClient)
	if err != nil {
		return "", err
	}

	logrus.Println("🚀 Successfully Upgraded Chaos Infra")

	return "Success", nil

}

// retrieves the Kubernetes and dynamic clients along with the configuration.
func getClientAndConfig(kubeconfig string) (*rest.Config, *kubernetes.Clientset, dynamic.Interface, error) {
	var config *rest.Config
	var dynamicClient dynamic.Interface

	// If kubeconfig is provided, use it to create the configuration and dynamic client.
	if kubeconfig != "" {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to build config from flags: %w", err)
		}
	} else {
		// Use the default kubeconfig file at $HOME/.kube/config.
		home := homedir.HomeDir()
		defaultKubeconfig := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(defaultKubeconfig); !os.IsNotExist(err) {
			var err error
			config, err = clientcmd.BuildConfigFromFlags("", defaultKubeconfig)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to build config from flags: %w", err)
			}
		} else {
			// Return an error if kubeconfig is not provided and the default file doesn't exist
			return nil, nil, nil, errors.New("kubeconfig not provided, and default kubeconfig not found")
		}
	}

	// Create the dynamic client using the configuration.
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Create a Kubernetes client using the configuration.
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create kube client: %w", err)
	}

	return config, kubeClient, dynamicClient, nil
}

// decodes the manifest byte slice into a list of Unstructured resources.
func decodeManifest(manifest []byte) ([]*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), len(manifest))

	// Split the resource file into separate YAML documents.
	resources := []*unstructured.Unstructured{}
	for {
		resourcestr := &unstructured.Unstructured{}
		if err := decoder.Decode(resourcestr); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("error in decoding resource: %w", err)
		}
		resources = append(resources, resourcestr)
	}

	return resources, nil
}

// applies the decoded resources using the dynamic client and Kubernetes client.
func applyResources(resources []*unstructured.Unstructured, dynamicClient dynamic.Interface, kubeClient *kubernetes.Clientset) error {
	for _, resource := range resources {
		logrus.Infof("Applying resource: %s , kind: %s", resource.GetName(), resource.GetKind())

		gvk := resource.GroupVersionKind()
		// a mapper for REST mapping.
		mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(kubeClient.Discovery()))
		// Map GVK to GVR using the REST mapper.
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("error in resource gvk to gvr mapping: %w", err)
		}

		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			// Namespaced resources should specify the namespace
			dr = dynamicClient.Resource(mapping.Resource).Namespace(resource.GetNamespace())
		} else {
			// For cluster-wide resources
			dr = dynamicClient.Resource(mapping.Resource)
		}

		// Apply the resource using the dynamic client.
		_, err = dr.Apply(context.TODO(), resource.GetName(), resource, metav1.ApplyOptions{
			Force:        true,
			FieldManager: "application/apply-patch",
		})
		if err != nil {
			return fmt.Errorf("error in applying resource: %w", err)
		}

		logrus.Info("Resource applied successfully")
	}

	return nil
}

// GetConfigMap returns config map for a given name and namespace
func GetConfigMap(c context.Context, name string, namespace string) (map[string]string, error) {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("configmap", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeconfig = flag.String("configmap", "", "")
	}
	flag.Parse()

	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		return nil, err
	}
	x, err := clientset.CoreV1().ConfigMaps(namespace).Get(c, name, metav1.GetOptions{})
	if err != nil {
		return nil, err

	}
	return x.Data, nil
}
