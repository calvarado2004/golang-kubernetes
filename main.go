package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	//
	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// main function, we will create a Deployment and a PVC with a Portworx Sharedv4 volume
func main() {

	log.Printf("---Starting Kubernetes external client!---")

	createSharedv4PVC("default", "nginx", "px-csi-db", "2Gi")

	createDeploymentWithPVC("default", "nginx", "nginx", "1.23.2", 80, 3)

	var labels = make(map[string]string)

	labels["app"] = "nginx"

	err := waitForPods(labels)

	if err != nil {
		log.Printf("Error waiting for pods: %v", err)
		return
	}

	listPods("default")

	var test []string

	updateOptions := metav1.UpdateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		FieldManager:    "k8s-external-client",
		FieldValidation: metav1.FieldValidationStrict,
		DryRun:          test,
	}

	updateOptions.DeepCopy()

	err = updateDeployment("default", "nginx-deployment", updateOptions)
	if err != nil {
		log.Printf("Error updating deployment: %v", err)
		return
	}

}
