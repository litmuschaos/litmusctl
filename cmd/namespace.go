package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"

	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NsExists checks if the given namespace already exists
func NsExists(namespace string) (bool, error) {
	clientset, err := ClientSet()
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

// ValidNs takes a valid namespace as input
func ValidNs() (string, bool) {
	var namespace string
	var nsExists bool
	fmt.Print("📁 Enter the namespace [", defaultNs, "]: ")
	fmt.Scanln(&namespace)
	if namespace == "" {
		namespace = defaultNs
	}
	ok, err := NsExists(namespace)
	if err != nil {
		fmt.Printf("\n Namespace existence check failed: {%s}\n", err.Error())
		os.Exit(1)
	}
	if ok {
		if PodExists(namespace, agentLabel) {
			fmt.Println("🚫 Subscriber already present. Please enter a different namespace")
			namespace, nsExists = ValidNs()
		} else {
			nsExists = true
			fmt.Println("👍 Continuing with", namespace, "namespace")
		}
	} else {
		if val, _ := CheckSAPermissions("create", "namespace", false); !val {
			fmt.Println("🚫 You don't have permissions to create a namespace.\n🙄 Please enter an existing namespace.")
			namespace, nsExists = ValidNs()
		}
		nsExists = false
	}
	return namespace, nsExists
}

// CreateNs creates the given namespace
func CreateNs(namespace string) {
	clientset, err := ClientSet()
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
