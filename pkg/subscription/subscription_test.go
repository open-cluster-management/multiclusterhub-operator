package subscription

import (
	"reflect"
	"testing"

	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func TestValidate(t *testing.T) {
	mch := &operatorsv1alpha1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1alpha1.MultiClusterHubSpec{
			Version:         "latest",
			ImageRepository: "quay.io/open-cluster-management",
			ImagePullPolicy: corev1.PullAlways,
			ImagePullSecret: "test",
			ReplicaCount:    1,
			NodeSelector: &operatorsv1alpha1.NodeSelector{
				OS:                  "test",
				CustomLabelSelector: "test",
				CustomLabelValue:    "test",
			},
			Mongo: operatorsv1alpha1.Mongo{},
		},
	}
	// 1. Valid mch
	sub := KUIWebTerminal(mch)

	// 2. Modified ImagePullSecret
	mch1 := mch.DeepCopy()
	mch1.Spec.ImagePullSecret = "notTest"
	sub1 := KUIWebTerminal(mch1)

	// 3. Modified ImagePullPolicy
	mch2 := mch.DeepCopy()
	mch2.Spec.ImagePullPolicy = corev1.PullNever
	sub2 := KUIWebTerminal(mch2)

	// 4. Modified ImageRepository
	mch3 := mch.DeepCopy()
	mch3.Spec.ImageRepository = "notquay.io/closed-cluster-management"
	sub3 := KUIWebTerminal(mch3)

	// 5. Modified ReplicaCount
	mch4 := mch.DeepCopy()
	mch4.Spec.ReplicaCount = 2
	sub4 := KUIWebTerminal(mch4)

	type args struct {
		found *unstructured.Unstructured
		want  *unstructured.Unstructured
	}
	tests := []struct {
		name  string
		args  args
		want  *unstructured.Unstructured
		want1 bool
	}{
		{
			name:  "Valid subscription",
			args:  args{sub, sub},
			want:  nil,
			want1: false,
		},
		{
			name:  "Modified ImagePullSecret",
			args:  args{sub, sub1},
			want:  sub1,
			want1: true,
		},
		{
			name:  "Modified ImagePullPolicy",
			args:  args{sub, sub2},
			want:  sub2,
			want1: true,
		},
		{
			name:  "Modified ImageRepository",
			args:  args{sub, sub3},
			want:  sub3,
			want1: true,
		},
		{
			name:  "Modified ReplicaCount",
			args:  args{sub, sub4},
			want:  sub4,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Validate(tt.args.found, tt.args.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Validate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSubscriptions(t *testing.T) {
	mch := &operatorsv1alpha1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1alpha1.MultiClusterHubSpec{
			Version:         "latest",
			ImageRepository: "quay.io/open-cluster-management",
			ImagePullPolicy: corev1.PullAlways,
			ImagePullSecret: "test",
			ReplicaCount:    1,
			NodeSelector: &operatorsv1alpha1.NodeSelector{
				OS:                  "test",
				CustomLabelSelector: "test",
				CustomLabelValue:    "test",
			},
			Mongo: operatorsv1alpha1.Mongo{},
		},
	}

	cache := utils.CacheSpec{
		IngressDomain: "testIngress",
	}

	tests := []struct {
		name string
		got  *unstructured.Unstructured
	}{
		{"ApplicationUI subscription", ApplicationUI(mch)},
		{"CertManager subscription", CertManager(mch)},
		{"CertWebhook subscription", CertWebhook(mch)},
		{"ConfigWatcher subscription", ConfigWatcher(mch)},
		{"Console subscription", Console(mch, cache)},
		{"GRC subscription", GRC(mch)},
		{"KUIWebTerminal subscription", KUIWebTerminal(mch)},
		{"ManagementIngress subscription", ManagementIngress(mch, cache)},
		{"MongoDB subscription", MongoDB(mch)},
		{"RCM subscription", RCM(mch)},
		{"Search subscription", Search(mch)},
		{"Topology subscription", Topology(mch)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.got.MarshalJSON(); err != nil {
				t.Error("Subscription does not marshal properly")
			}
			_, err := yaml.Marshal(tt.got.Object["spec"])
			if err != nil {
				t.Error("Issue parsing subscription values")
			}
		})
	}
}
