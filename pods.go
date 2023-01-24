package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"time"
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

// waitForPods waits for all pods to be ready
func waitForPods(deployLabels map[string]string) error {

	//time.Sleep(5 * time.Second)

	client, err := getClient()
	if err != nil {
		log.Printf("Error getting client: %v", err)
		return err
	}

	validatedLabels, err := labels.ValidatedSelectorFromSet(deployLabels)
	if err != nil {
		return fmt.Errorf("failed to validate labels: %w", err)
	}

	labelsToUse := validatedLabels.String()

	for {
		podList, err := client.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{LabelSelector: labelsToUse})
		if err != nil {
			return fmt.Errorf("failed to list pods: %w", err)
		}

		log.Printf("Waiting for %d pods to be ready", len(podList.Items))

		podsRunning := 0

		for _, pod := range podList.Items {
			if pod.Status.Phase == "Running" {
				podsRunning++
			}
		}

		if podsRunning > 0 && podsRunning == len(podList.Items) {
			break
		}

		fmt.Printf("Waiting for pods to be ready, running (%d/%d)...\n", podsRunning, len(podList.Items))

		time.Sleep(2 * time.Second)

	}

	return nil
}
