package k8s

import (
	"context"
	"fmt"
	"log"

	"github.com/litmuschaos/litmusctl/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SAExists checks if the given service account exists in the given namespace
func SAExists(namespace, serviceaccount string, kubeconfig *string) bool {
	clientset, err := ClientSet(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf("serviceaccounts \"%s\" not found", serviceaccount)
	_, newErr := clientset.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), serviceaccount, metav1.GetOptions{})
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
	fmt.Print("🔑 Enter service account [", constants.DefaultSA, "]: ")
	fmt.Scanln(&sa)
	if sa == "" {
		sa = constants.DefaultSA
	}
	if SAExists(namespace, sa, kubeconfig) {
		fmt.Println("👍 Using the existing service account")
		return sa, true
	}
	return sa, false
}
