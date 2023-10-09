// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package resources

import (
	operatorsv1 "github.com/stolostron/multiclusterhub-operator/api/v1"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	mcev1 "github.com/stolostron/backplane-operator/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ocmapi "open-cluster-management.io/api/addon/v1alpha1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	MulticlusterhubName      = "test-mch"
	MulticlusterhubNamespace = "open-cluster-management"
	JobName                  = "test-job"
	MultiClusterEngineName   = "multiclusterengine-sample"
)

var (
	MCHLookupKey = types.NamespacedName{Name: MulticlusterhubName, Namespace: MulticlusterhubNamespace}
	MCELookupKey = types.NamespacedName{Name: MultiClusterEngineName}
)

func EmptyMCE() mcev1.MultiClusterEngine {
	return mcev1.MultiClusterEngine{
		ObjectMeta: metav1.ObjectMeta{
			Name: MultiClusterEngineName,
		},
		Spec: mcev1.MultiClusterEngineSpec{
			TargetNamespace: "multicluster-engine",
		},
	}
}

func SpecMCH() *operatorsv1.MultiClusterHub {
	mchNodeSelector := map[string]string{"select": "test"}
	mchImagePullSecret := "test"
	mchTolerations := []corev1.Toleration{
		{
			Key:      "dedicated",
			Operator: "Exists",
			Effect:   "NoSchedule",
		},
	}

	return &operatorsv1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MulticlusterhubName,
			Namespace: MulticlusterhubNamespace,
		},
		Spec: operatorsv1.MultiClusterHubSpec{
			NodeSelector:    mchNodeSelector,
			ImagePullSecret: mchImagePullSecret,
			Tolerations:     mchTolerations,
			Overrides: &operatorsv1.Overrides{
				Components: []operatorsv1.ComponentConfig{
					{
						Name:    operatorsv1.ClusterBackup,
						Enabled: false,
					},
				},
			},
		},
	}
}

func EmptyMCH() operatorsv1.MultiClusterHub {
	return operatorsv1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MulticlusterhubName,
			Namespace: MulticlusterhubNamespace,
		},
		Spec: operatorsv1.MultiClusterHubSpec{
			Overrides: &operatorsv1.Overrides{
				Components: []operatorsv1.ComponentConfig{
					{
						Name:    operatorsv1.ClusterBackup,
						Enabled: false,
					},
				},
			},
		},
	}
}

func NoComponentMCH() operatorsv1.MultiClusterHub {
	return operatorsv1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MulticlusterhubName,
			Namespace: MulticlusterhubNamespace,
		},
		Spec: operatorsv1.MultiClusterHubSpec{
			Overrides: &operatorsv1.Overrides{
				Components: []operatorsv1.ComponentConfig{
					{
						Name:    operatorsv1.Console,
						Enabled: false,
					},
					{
						Name:    operatorsv1.Insights,
						Enabled: false,
					},
					{
						Name:    operatorsv1.Search,
						Enabled: false,
					},
					{
						Name:    operatorsv1.ClusterBackup,
						Enabled: false,
					},
					{
						Name:    operatorsv1.GRC,
						Enabled: false,
					},
					{
						Name:    operatorsv1.ClusterLifecycle,
						Enabled: false,
					},
					{
						Name:    operatorsv1.MultiClusterObservability,
						Enabled: false,
					},
					{
						Name:    operatorsv1.Volsync,
						Enabled: false,
					},
				},
			},
		},
	}
}

func InsightsMCH() operatorsv1.MultiClusterHub {
	return operatorsv1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      MulticlusterhubName,
			Namespace: MulticlusterhubNamespace,
		},
		Spec: operatorsv1.MultiClusterHubSpec{
			Overrides: &operatorsv1.Overrides{
				Components: []operatorsv1.ComponentConfig{
					{
						Name:    operatorsv1.Insights,
						Enabled: true,
					},
					{
						Name:    operatorsv1.ClusterBackup,
						Enabled: false,
					},
				},
			},
		},
	}
}

func OCMNamespace() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: MulticlusterhubNamespace,
		},
	}
}

func MCENamespace() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "multicluster-engine",
			Labels: map[string]string{},
		},
	}
}

func MonitoringNamespace() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openshift-monitoring",
		},
	}
}

func SampleService(m *operatorsv1.MultiClusterHub) *corev1.Service {
	const Port = 3030

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample-service",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       int32(Port),
				TargetPort: intstr.FromInt(Port),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	return s
}

func SampleClusterManagementAddOn(component string) *ocmapi.ClusterManagementAddOn {
	addonName, err := operatorsv1.GetClusterManagementAddonName(component)
	if err != nil {
		addonName = "unknown"
	}

	addon := &ocmapi.ClusterManagementAddOn{
		ObjectMeta: metav1.ObjectMeta{
			Name: addonName,
		},
		Spec: ocmapi.ClusterManagementAddOnSpec{
			AddOnMeta: ocmapi.AddOnMeta{
				Description: "Sample addon description",
				DisplayName: component,
			},
		},
	}

	return addon
}

func SampleServiceMonitor(component string, namespace string) *promv1.ServiceMonitor {
	smName, err := operatorsv1.GetServiceMonitorName(component)
	if err != nil {
		smName = "unknown"
	}

	sm := &promv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      smName,
			Namespace: namespace,
		},
		Spec: promv1.ServiceMonitorSpec{},
	}

	return sm
}
