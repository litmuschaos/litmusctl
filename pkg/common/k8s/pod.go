package k8s

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WatchPod watches for the pod status
func WatchPod(namespace, label string) {
	clientset, err := ClientSet()
	if err != nil {
		log.Fatal(err)
	}
	watch, err := clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	for event := range watch.ResultChan() {
		p, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Fatal("unexpected type")
		}
		fmt.Println("💡 Connecting agent to Kubera Enterprise.")
		if p.Status.Phase == "Running" {
			fmt.Println("🏃 Agents running!!")
			watch.Stop()
			break
		}
	}
}

type PodList struct {
	Items []string
}

// PodExists checks if the pod with the given label already exists in the given namespace
func PodExists(namespace, label string) bool {
	clientset, err := ClientSet()
	if err != nil {
		log.Fatal(err)
	}
	watch, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(watch.Items) >= 1 {
		return true
	}
	return false
}
