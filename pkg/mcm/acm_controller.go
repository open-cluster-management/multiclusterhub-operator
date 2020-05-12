// Copyright (c) 2020 Red Hat, Inc.

package mcm

import (
	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ACMControllerName is the name of the acm controller deployment
const ACMControllerName string = "acm-controller"

// ACMControllerDeployment creates the deployment for the acm controller
func ACMControllerDeployment(m *operatorsv1beta1.MultiClusterHub, cache utils.CacheSpec) *appsv1.Deployment {
	replicas := getReplicaCount(m)
	mode := int32(420)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ACMControllerName,
			Namespace: m.Namespace,
			Labels:    defaultLabels(ACMControllerName),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: defaultLabels(ACMControllerName),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: defaultLabels(ACMControllerName),
				},
				Spec: corev1.PodSpec{
					ImagePullSecrets:   []corev1.LocalObjectReference{{Name: m.Spec.ImagePullSecret}},
					ServiceAccountName: ServiceAccount,
					NodeSelector:       m.Spec.NodeSelector,
					Affinity:           utils.DistributePods("app", ACMControllerName),
					Volumes: []corev1.Volume{
						{
							Name: "klusterlet-certs",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{DefaultMode: &mode, SecretName: utils.KlusterletSecretName},
							},
						},
					},
					Containers: []corev1.Container{{
						Image:           Image(cache),
						ImagePullPolicy: utils.GetImagePullPolicy(m),
						Name:            ACMControllerName,
						Args: []string{
							"/acm-controller",
							"--klusterlet-cafile=/var/run/klusterlet/ca.crt",
							"--max-qps=100.0",
							"--max-burst=200",
						},
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse("200m"),
								v1.ResourceMemory: resource.MustParse("256Mi"),
							},
							Limits: v1.ResourceList{
								v1.ResourceMemory: resource.MustParse("2048Mi"),
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "klusterlet-certs", MountPath: "/var/run/klusterlet"},
						},
					}},
				},
			},
		},
	}

	dep.SetOwnerReferences([]metav1.OwnerReference{
		*metav1.NewControllerRef(m, m.GetObjectKind().GroupVersionKind()),
	})
	return dep
}
