package mcm

import (
	"testing"

	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWebhookDeployment(t *testing.T) {
	empty := &operatorsv1beta1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1beta1.MultiClusterHubSpec{
			ImageRepository: "",
			ImagePullPolicy: "",
			ImagePullSecret: "",
			ImageTagSuffix:  "",
			Mongo:           operatorsv1beta1.Mongo{},
		},
	}

	cs := utils.CacheSpec{
		IngressDomain:   "testIngress",
		ImageShaDigests: map[string]string{},
	}

	t.Run("MCH with empty fields", func(t *testing.T) {
		_ = WebhookDeployment(empty, cs)
	})

	essentialsOnly := &operatorsv1beta1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1beta1.MultiClusterHubSpec{
			ImageRepository: "test",
			ImagePullPolicy: "test",
			ImageTagSuffix:  "test",
		},
	}
	t.Run("MCH with only required values", func(t *testing.T) {
		_ = WebhookDeployment(essentialsOnly, cs)
	})
}

func TestWebhookService(t *testing.T) {
	mch := &operatorsv1beta1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testName",
			Namespace: "testNS",
		},
	}

	t.Run("Create service", func(t *testing.T) {
		s := WebhookService(mch)
		if ns := s.Namespace; ns != "testNS" {
			t.Errorf("expected namespace %s, got %s", "testNS", ns)
		}
		if ref := s.GetOwnerReferences(); ref[0].Name != "testName" {
			t.Errorf("expected ownerReference %s, got %s", "testName", ref[0].Name)
		}
	})
}
