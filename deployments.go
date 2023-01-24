package main

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

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

	var labels = make(map[string]string)

	labels["app"] = appName

	ObjectMetaVar := metav1.ObjectMeta{
		Name:   appName + "-deployment",
		Labels: labels,
	}

	LabelSelectorRequirementVar := metav1.LabelSelectorRequirement{
		Key:      "app",
		Operator: "In",
		Values:   []string{appName},
	}

	PodAffinityTermVar := corev1.PodAffinityTerm{
		LabelSelector: &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				LabelSelectorRequirementVar,
			},
		},
		TopologyKey: "kubernetes.io/hostname",
	}

	AffinityVar := corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				PodAffinityTermVar,
			},
		},
	}

	VolumeMountVar := corev1.VolumeMount{
		Name:      appName + "-volume",
		MountPath: "/data",
	}

	ContainerPortVar := corev1.ContainerPort{
		ContainerPort: containerPort,
	}

	ResourcesVar := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("10m"),
			corev1.ResourceMemory: resource.MustParse("10Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("100Mi"),
		},
	}

	ContainerVar := corev1.Container{
		Name:            appName,
		Image:           appImage + ":" + tagImage,
		ImagePullPolicy: corev1.PullAlways,
		Ports: []corev1.ContainerPort{
			ContainerPortVar,
		},
		Resources: ResourcesVar,
		VolumeMounts: []corev1.VolumeMount{
			VolumeMountVar,
		},
	}

	VolumeVar := corev1.Volume{
		Name: appName + "-volume",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: appName + "-pvc",
			},
		},
	}

	DeploymentSpecVar := appsv1.DeploymentSpec{
		Replicas: int32Ptr(replicas),
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Affinity:      &AffinityVar,
				SchedulerName: "stork",
				Containers: []corev1.Container{
					ContainerVar,
				},
				Volumes: []corev1.Volume{
					VolumeVar,
				},
			},
		},
	}

	DeploymentStruct := &appsv1.Deployment{
		ObjectMeta: ObjectMetaVar,
		Spec:       DeploymentSpecVar,
	}

	deploy, err := deployment.Create(context.Background(), DeploymentStruct, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating deployment %s", err)
		return
	}

	log.Printf("Deployment %s created successfully!", deploy.Name)

}
