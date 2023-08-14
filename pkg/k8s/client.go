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
	"os"
	"path/filepath"

	"github.com/litmuschaos/litmusctl/pkg/utils"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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
