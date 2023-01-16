package main

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

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
