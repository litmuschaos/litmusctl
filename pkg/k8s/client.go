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
	"flag"
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

//getKubeConfig setup the config for access cluster resource
func GetKubeConfig() (*rest.Config, error) {
	// Use in-cluster config if kubeconfig path is not specified
	KubeConfig := os.Getenv("KUBECONFIG")
	// Use in-cluster config if kubeconfig path is not specified
	if KubeConfig == "" {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags("", KubeConfig)
}

func GetKubecConfig1()(*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"),"")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println("Could not get config")
		return nil, err
	}

	return cfg, nil
}

// Returns a new kubernetes client set
func ClientSet(kubeconfig *string) (*kubernetes.Clientset, error) {
	if *kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kcfg := filepath.Join(home, ".kube", "config")
			kubeconfig = &kcfg
		} else {
			utils.Red.Println("ERROR: Clientset generation failed!")
			os.Exit(1)
		}
	}

	// create the config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		utils.Red.Println("ERROR: ", err.Error())
		os.Exit(1)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.Red.Println("ERROR: ", err.Error())
		os.Exit(1)
	}
	return clientset, err
}

//This function returns a dynamic client and discovery client
func GetDynamicAndDiscoveryClient() (discovery.DiscoveryInterface, dynamic.Interface, error) {
	// returns a config object which uses the service account kubernetes gives to pods
	config, err := GetKubecConfig1()
	if err != nil {
		fmt.Println("ERROR2", err)
		return nil, nil, err
	}
	// NewDiscoveryClientForConfig creates a new DiscoveryClient for the given config
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	// NewForConfig creates a new dynamic client or returns an error.
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return discoveryClient, dynamicClient, nil
}
