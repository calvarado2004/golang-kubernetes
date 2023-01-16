package main

import (
	"log"
	"time"

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

	createSharedv4PVC("default", "nginx", "portworx-sharedv4-csi", "2Gi")

	createDeploymentWithPVC("default", "nginx", "nginx", "1.19.0", 80, 3)

	log.Printf("Waiting for 10 seconds to create the deployment")

	time.Sleep(10 * time.Second)

	listPods("default")

}
