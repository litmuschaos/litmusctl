package k8s

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/litmuschaos/litmusctl/tmp-pkg/constants"
	v1 "k8s.io/api/core/v1"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

// ValidNs takes a valid namespace as input from user
func ValidNs(mode string, label string, kubeconfig *string) (string, bool) {
start:
	var (
		namespace string
		nsExists  bool
	)

	if mode == "namespace" {
		fmt.Print("📁 Enter the namespace (existing) [", constants.DefaultNs, "]: ")
		fmt.Scanln(&namespace)

	} else if mode == "cluster" {
		fmt.Print("📁 Enter the namespace (new or existing) [", constants.DefaultNs, "]: ")
		fmt.Scanln(&namespace)
	} else {
		fmt.Printf("\n 🚫 No mode selected \n")
		os.Exit(1)
	}

	if namespace == "" {
		namespace = constants.DefaultNs
	}
	ok, err := NsExists(namespace, kubeconfig)
	if err != nil {
		fmt.Printf("\n 🚫 Namespace existence check failed: {%s}\n", err.Error())
		os.Exit(1)
	}
	if ok {
		if PodExists(namespace, label, kubeconfig) {
			fmt.Println("🚫 Subscriber already present. Please enter a different namespace")
			goto start
		} else {
			nsExists = true
			fmt.Println("👍 Continuing with", namespace, "namespace")
		}
	} else {
		if val, _ := CheckSAPermissions("create", "namespace", false, kubeconfig); !val {
			fmt.Println("🚫 You don't have permissions to create a namespace.\n🙄 Please enter an existing namespace.")
			goto start
		}
		nsExists = false
	}

	return namespace, nsExists
}

// CreateNs creates the given namespace
func CreateNs(namespace string, kubeconfig *string) {
	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s", namespace)}}
	_, newErr := clientset.CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
	if newErr != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(namespace, "namespace created successfully")
}
