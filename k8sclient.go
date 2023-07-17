package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// getClient returns a Kubernetes clientset based on the environment
func getClient() (*kubernetes.Clientset, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		// We are running inside a cluster
		return getClientInternal()
	} else {
		// We are running outside a cluster
		return getClientExternal()
	}
}

// getClientExternal returns a Kubernetes client when running outside a Kubernetes cluster
func getClientExternal() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
		return nil, err
	}

	return clientset, nil

}

// getClientInternal returns a Kubernetes client when running inside a Kubernetes cluster
func getClientInternal() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
		return nil, err
	}

	return clientset, nil
}
