package main

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
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

// createSharedv4PVC creates a PVC using the given storage class, ReadWriteMany
func createSharedv4PVC(namespace string, appName string, storageClassName string, size string) {

	log.Printf("---Creating PVC on %s namespace of app %s ---", namespace, appName)

	clientset, err := getClient()
	if err != nil {
		log.Printf("Error getting client: %v", err)
		return
	}

	pvc := clientset.CoreV1().PersistentVolumeClaims(namespace)

	_, err = pvc.Create(context.Background(), &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName + "-pvc",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(size),
				},
			},
			StorageClassName: &storageClassName,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating PVC %s", err)
		return
	}

	log.Printf("PVC %s created successfully!", appName+"-pvc")

}

// createDeployment creates a deployment with a PVC and Pod Anti-Affinity rules
func createDeploymentWithPVC(namespace string, appName string, appImage string, tagImage string, containerPort int32, replicas int32) {

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
					Affinity: &corev1.Affinity{
						PodAntiAffinity: &corev1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
								{
									LabelSelector: &metav1.LabelSelector{
										MatchExpressions: []metav1.LabelSelectorRequirement{
											{
												Key:      "app",
												Operator: metav1.LabelSelectorOpIn,
												Values:   []string{appName},
											},
										},
									},
									TopologyKey: "kubernetes.io/hostname",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{

							Name:  appName,
							Image: appImage + ":" + tagImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: containerPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      appName + "-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: appName + "-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: appName + "-pvc",
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

// main function, we will create a Deployment and a PVC with a Portworx Sharedv4 volume
func main() {

	log.Printf("---Starting Kubernetes external client!---")

	createSharedv4PVC("default", "nginx", "portworx-sharedv4-csi", "2Gi")

	createDeploymentWithPVC("default", "nginx", "nginx", "1.19.0", 80, 3)

	log.Printf("Waiting for 10 seconds to create the deployment")

	time.Sleep(10 * time.Second)

	listPods("default")

}
