package common

import (
	"context"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/common/k8s"
	"github.com/litmuschaos/litmusctl/pkg/constants"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// discoverPlatform determines the host platform and returns it
func DiscoverPlatform() string {
	if ok, _ := IsAWSPlatform(); ok {
		return "AWS"
	}
	if ok, _ := IsGKEPlatform(); ok {
		return "GKE"
	}
	if ok, _ := IsOpenshiftPlatform(); ok {
		return "Openshift"
	}
	if ok, _ := k8s.NsExists("cattle-system"); ok {
		return "Rancher"
	}
	return constants.DefaultPlatform
}

// IsAWSPlatform determines if the host platform is AWS
// by checking the ProviderID inside node spec
//
// Sample node custom resource of an AWS node
// {
//     "apiVersion": "v1",
//     "kind": "Node",
//     ....
//     "spec": {
//         "providerID": "aws:///us-east-2b/i-0bf24d83f4b993738"
//     }
//   }
// }
func IsAWSPlatform() (bool, error) {
	clientset, err := k8s.ClientSet()
	if err != nil {
		return false, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return false, err
	}
	if strings.HasPrefix(nodeList.Items[0].Spec.ProviderID, constants.AWSIdentifier) {
		return true, nil
	}
	return false, nil
}

// IsGKEPlatform determines if the host platform is GKE
// by checking the ProviderID inside node spec
//
// Sample node custom resource of an GKE node
// {
//     "apiVersion": "v1",
//     "kind": "Node",
//     ....
//     "spec": {
//         "providerID": ""
//     }
//   }
// }
func IsGKEPlatform() (bool, error) {
	clientset, err := k8s.ClientSet()
	if err != nil {
		return false, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return false, err
	}
	if strings.HasPrefix(nodeList.Items[0].Spec.ProviderID, constants.GKEIdentifier) {
		return true, nil
	}
	return false, nil
}

// IsOpenshiftPlatform determines if the host platform
// is Openshift by checking "node.openshift.io/os_id"
// label on the nodes
//
// Sample node custom resource of an Openshift node
// {
//     "apiVersion": "v1",
//     "kind": "Node",
//     "metadata": {
//         "labels": {
//             "node.openshift.io/os_id": "rhcos"
//         }
//    }
//    ....
// }
func IsOpenshiftPlatform() (bool, error) {
	clientset, err := k8s.ClientSet()
	if err != nil {
		return false, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{
		LabelSelector: constants.OpenshiftIdentifier,
	})
	if err != nil {
		return false, err
	}
	if len(nodeList.Items) > 0 {
		return true, nil
	}
	return false, nil
}
