package cmd

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Returns a new kubernetes client set
func ClientSet() (*kubernetes.Clientset, error) {
	var kubeconfig string
	kcfg, err := rootCmd.PersistentFlags().GetString("kubeconfig")
	if err != nil {
		panic(err)
	}
	if kcfg == "" {
		home, err := homeDir()
		if err != nil {
			panic(err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = kcfg
	}
	// create the config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, err
}
