package rendering

import (
	"os"
	"path"
	"reflect"
	"testing"

	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/rendering/templates"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestRender(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working dir %v", err)
	}
	templatesPath := path.Join(path.Dir(path.Dir(wd)), "templates")
	os.Setenv(templates.TemplatesPathEnvVar, templatesPath)
	defer os.Unsetenv(templates.TemplatesPathEnvVar)

	var replicas int32 = 1
	mchcr := &operatorsv1alpha1.MultiClusterHub{
		TypeMeta:   metav1.TypeMeta{Kind: "MultiClusterHub"},
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1alpha1.MultiClusterHubSpec{
			Version:         "latest",
			ImageRepository: "quay.io/open-cluster-management",
			ImagePullPolicy: "Always",
			ImagePullSecret: "test",
			NodeSelector: &operatorsv1alpha1.NodeSelector{
				OS:                  "test",
				CustomLabelSelector: "test",
				CustomLabelValue:    "test",
			},
			Foundation: operatorsv1alpha1.Foundation{
				Apiserver: operatorsv1alpha1.Apiserver{
					Replicas: &replicas,
					Configuration: map[string]string{
						"test": "test",
					},
				},
				Controller: operatorsv1alpha1.Controller{
					Replicas: &replicas,
					Configuration: map[string]string{
						"test": "test",
					},
				},
			},
			Mongo: operatorsv1alpha1.Mongo{
				Endpoints:  "test",
				ReplicaSet: "test",
			},
		},
	}

	renderer := NewRenderer(mchcr)
	objs, err := renderer.Render(nil)
	if err != nil {
		t.Fatalf("failed to render multiclusterhub %v", err)
	}

	printObjs(t, objs)
}

func printObjs(t *testing.T, objs []*unstructured.Unstructured) {
	for _, obj := range objs {
		t.Log(obj)
	}
}

func Test_addInstallerLabel(t *testing.T) {
	installer := "example-installer"

	t.Run("Should add label when none exist", func(t *testing.T) {
		u := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "apps.open-cluster-management.io/v1",
				"kind":       "Channel",
			},
		}
		want := map[string]string{
			"installer": installer,
		}

		addInstallerLabel(u, installer)
		if got := u.GetLabels(); !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
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
		want := map[string]string{
			"hello":     "world",
			"installer": installer,
		}

		addInstallerLabel(u, installer)
		if got := u.GetLabels(); !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
