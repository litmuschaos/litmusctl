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
package infra_ops

import (
	"context"
	"fmt"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// - Entering any character other than numbers returns 0. Input validation need to be done.
// - If input is given as "123abc", "abc" will be used for next user input. Buffer need to be read completely.
// - String literals like "AWS" are used at multiple places. Need to be changed to constants.
func GetPlatformName(kubeconfig *string) string {
	var platform int
	discoveredPlatform := DiscoverPlatform(kubeconfig)
	utils.White_B.Println("\nPlatform List: ")
	utils.White_B.Println(utils.PlatformList)
	utils.White_B.Print("\nSelect a platform [Default: ", discoveredPlatform, "] [Range: 1-5]: ")
	fmt.Scanln(&platform)
	switch platform {
	case 0:
		return discoveredPlatform
	case 1:
		return "AWS"
	case 2:
		return "GKE"
	case 3:
		return "Openshift"
	case 4:
		return "Rancher"
	default:
		return utils.DefaultPlatform
	}
}

// discoverPlatform determines the host platform and returns it
func DiscoverPlatform(kubeconfig *string) string {
	if ok, _ := IsAWSPlatform(kubeconfig); ok {
		return "AWS"
	}
	if ok, _ := IsGKEPlatform(kubeconfig); ok {
		return "GKE"
	}
	if ok, _ := IsOpenshiftPlatform(kubeconfig); ok {
		return "Openshift"
	}
	if ok, _ := k8s.NsExists("cattle-system", kubeconfig); ok {
		return "Rancher"
	}
	return utils.DefaultPlatform
}

// IsAWSPlatform determines if the host platform is AWS
// by checking the ProviderID inside node spec
//
// Sample node custom resource of an AWS node
//
//	{
//	    "apiVersion": "v1",
//	    "kind": "Node",
//	    ....
//	    "spec": {
//	        "providerID": "aws:///us-east-2b/i-0bf24d83f4b993738"
//	    }
//	  }
//	}
func IsAWSPlatform(kubeconfig *string) (bool, error) {
	clientset, err := k8s.ClientSet(kubeconfig)
	if err != nil {
		return false, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil || len(nodeList.Items) == 0 {
		return false, err
	}
	if strings.HasPrefix(nodeList.Items[0].Spec.ProviderID, utils.AWSIdentifier) {
		return true, nil
	}
	return false, nil
}

// IsGKEPlatform determines if the host platform is GKE
// by checking the ProviderID inside node spec
//
// Sample node custom resource of an GKE node
//
//	{
//	    "apiVersion": "v1",
//	    "kind": "Node",
//	    ....
//	    "spec": {
//	        "providerID": ""
//	    }
//	  }
//	}
func IsGKEPlatform(kubeconfig *string) (bool, error) {
	clientset, err := k8s.ClientSet(kubeconfig)
	if err != nil {
		return false, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil || len(nodeList.Items) == 0 {
		return false, err
	}
	if strings.HasPrefix(nodeList.Items[0].Spec.ProviderID, utils.GKEIdentifier) {
		return true, nil
	}
	return false, nil
}

// IsOpenshiftPlatform determines if the host platform
// is Openshift by checking "node.openshift.io/os_id"
// label on the nodes
//
// Sample node custom resource of an Openshift node
//
//	{
//	    "apiVersion": "v1",
//	    "kind": "Node",
//	    "metadata": {
//	        "labels": {
//	            "node.openshift.io/os_id": "rhcos"
//	        }
//	   }
//	   ....
//	}
func IsOpenshiftPlatform(kubeconfig *string) (bool, error) {
	clientset, err := k8s.ClientSet(kubeconfig)
	if err != nil {
		return false, err
	}
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{
		LabelSelector: utils.OpenshiftIdentifier,
	})
	if err != nil {
		return false, err
	}
	if len(nodeList.Items) > 0 {
		return true, nil
	}
	return false, nil
}
