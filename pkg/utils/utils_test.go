package utils

import (
	"reflect"
	"testing"

	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNodeSelectors(t *testing.T) {
	mch := &operatorsv1alpha1.MultiClusterHub{
		Spec: operatorsv1alpha1.MultiClusterHubSpec{
			NodeSelector: &operatorsv1alpha1.NodeSelector{
				OS:                  "linux",
				CustomLabelSelector: "kubernetes.io/arch",
				CustomLabelValue:    "amd64",
			},
		},
	}
	mchNoSelector := &operatorsv1alpha1.MultiClusterHub{}
	mchEmptySelector := &operatorsv1alpha1.MultiClusterHub{
		Spec: operatorsv1alpha1.MultiClusterHubSpec{
			NodeSelector: &operatorsv1alpha1.NodeSelector{
				CustomLabelSelector: "kubernetes.io/arch",
				CustomLabelValue:    "",
			},
		},
	}

	type args struct {
		mch *operatorsv1alpha1.MultiClusterHub
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "With node selectors",
			args: args{mch},
			want: map[string]string{
				"kubernetes.io/os":   "linux",
				"kubernetes.io/arch": "amd64",
			},
		},
		{
			name: "No node selector",
			args: args{mchNoSelector},
			want: nil,
		},
		{
			name: "Empty selector value",
			args: args{mchEmptySelector},
			want: map[string]string{
				"kubernetes.io/arch": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NodeSelectors(tt.args.mch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nodeSelectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddInstallerLabel(t *testing.T) {
	name := "example-installer"
	ns := "default"

	t.Run("Should add labels when none exist", func(t *testing.T) {
		u := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "apps.open-cluster-management.io/v1",
				"kind":       "Channel",
			},
		}
		want := 2

		AddInstallerLabel(u, name, ns)
		if got := len(u.GetLabels()); got != want {
			t.Errorf("got %v labels, want %v", got, want)
		}
	})

	t.Run("Should not replace existing labels", func(t *testing.T) {
		u := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "apps.open-cluster-management.io/v1",
				"kind":       "Channel",
				"metadata": map[string]interface{}{
					"name": "channelName",
					"labels": map[string]interface{}{
						"hello": "world",
					},
				},
			},
		}
		want := 3

		AddInstallerLabel(u, name, ns)
		if got := len(u.GetLabels()); got != want {
			t.Errorf("got %v labels, want %v", got, want)
		}
	})
}

func TestContainsPullSecret(t *testing.T) {
	superset := []corev1.LocalObjectReference{{Name: "foo"}, {Name: "bar"}}
	type args struct {
		pullSecrets []corev1.LocalObjectReference
		ps          corev1.LocalObjectReference
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Contains pull secret",
			args{
				pullSecrets: superset,
				ps:          corev1.LocalObjectReference{Name: "foo"},
			},
			true,
		},
		{
			"Does not contain pull secret",
			args{
				pullSecrets: superset,
				ps:          corev1.LocalObjectReference{Name: "baz"},
			},
			false,
		},
		{
			"Empty list",
			args{
				pullSecrets: []corev1.LocalObjectReference{},
				ps:          corev1.LocalObjectReference{Name: "baz"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsPullSecret(tt.args.pullSecrets, tt.args.ps); got != tt.want {
				t.Errorf("ContainsPullSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsMap(t *testing.T) {
	superset := map[string]string{
		"hello":     "world",
		"goodnight": "moon",
		"yip":       "yip",
	}
	type args struct {
		all      map[string]string
		expected map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Superset",
			args{
				all:      superset,
				expected: map[string]string{"hello": "world", "yip": "yip"},
			},
			true,
		},
		{
			"Partial overlap",
			args{
				all:      superset,
				expected: map[string]string{"hello": "world", "greetings": "traveler"},
			},
			false,
		},
		{
			"Empty superset",
			args{
				all:      map[string]string{},
				expected: map[string]string{"yip": "yip"},
			},
			false,
		},
		{
			"Empty subset",
			args{
				all:      superset,
				expected: map[string]string{},
			},
			true,
		},
		{
			"Same keys, different values",
			args{
				all:      superset,
				expected: map[string]string{"hello": "moon", "yip": "yip"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsMap(tt.args.all, tt.args.expected); got != tt.want {
				t.Errorf("ContainsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMchIsValid(t *testing.T) {
	validMCH := &operatorsv1alpha1.MultiClusterHub{
		TypeMeta:   metav1.TypeMeta{Kind: "MultiClusterHub"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1alpha1.MultiClusterHubSpec{
			Version:         "latest",
			ImageRepository: "quay.io/open-cluster-management",
			ImagePullPolicy: "Always",
			ImagePullSecret: "test",
			ReplicaCount:    1,
			NodeSelector: &operatorsv1alpha1.NodeSelector{
				OS:                  "test",
				CustomLabelSelector: "test",
				CustomLabelValue:    "test",
			},
			Mongo: operatorsv1alpha1.Mongo{
				Storage:      "mongoStorage",
				StorageClass: "mongoStorageClass",
				ReplicaCount: 1,
			},
			Etcd: operatorsv1alpha1.Etcd{
				Storage:      "etcdStorage",
				StorageClass: "etcdStorageClass",
			},
		},
	}
	noRepo := validMCH.DeepCopy()
	noRepo.Spec.ImageRepository = ""
	noMongoReplicas := validMCH.DeepCopy()
	noMongoReplicas.Spec.Mongo.ReplicaCount = 0

	type args struct {
		m *operatorsv1alpha1.MultiClusterHub
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Valid MCH",
			args{validMCH},
			true,
		},
		{
			"Missing Image Repository",
			args{noRepo},
			false,
		},
		{
			"Zero Mongo Replicas",
			args{noMongoReplicas},
			false,
		},
		{
			"Empty object",
			args{&operatorsv1alpha1.MultiClusterHub{}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MchIsValid(tt.args.m); got != tt.want {
				t.Errorf("MchIsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
