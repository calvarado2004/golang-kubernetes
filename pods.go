package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// listPods lists all pods in the given namespace
func listPods(namespace string) {

	log.Printf("---Listing pods on %s namespace---", namespace)

	clientset, err := getClient()
	if err != nil {
		panic(err.Error())
	}

	list, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return
	}
	for _, pod := range list.Items {
		println(pod.Name)
	}

}
