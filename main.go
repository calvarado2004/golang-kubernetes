package main

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"

	//
	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// getClient returns a Kubernetes client
func getClient() (*kubernetes.Clientset, error) {
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

// createDeployment creates a deployment
func createDeployment(namespace string, appName string, appImage string, tagImage string, containerPort int32, replicas int32) {

	var int32Ptr = func(i int32) *int32 { return &i }

	log.Printf("---Creating deployment on %s namespace of app %s ---", namespace, appName)

	clientset, err := getClient()
	if err != nil {
		log.Printf("Error getting client: %v", err)
		return
	}

	deployment := clientset.AppsV1().Deployments(namespace)

	deploy, err := deployment.Create(context.Background(), &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{

							Name:  appName,
							Image: appImage + ":" + tagImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerPort,
								},
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating deployment %s", err)
		return
	}

	log.Printf("Deployment %s created successfully!", deploy.Name)

}

func main() {

	log.Printf("---Starting Kubernetes external client!---")

	listPods("portworx-client")

	createDeployment("default", "nginx", "nginx", "1.19.0", 80, 3)

}
