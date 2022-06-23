// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package renderer

import (

	// "reflect"
	"os"
	"reflect"
	"testing"

	v1 "github.com/stolostron/multiclusterhub-operator/api/v1"
	"github.com/stolostron/multiclusterhub-operator/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	chartsDir  = "../templates/charts/toggle"
	chartsPath = "../templates/charts/toggle/insights"
	crdsDir    = "../templates/crds"
)

func TestRender(t *testing.T) {

	proxyList := []string{"insights-client"}
	mchNodeSelector := map[string]string{"select": "test"}
	mchImagePullSecret := "test"
	mchNamespace := "default"
	mchTolerations := []corev1.Toleration{
		{
			Key:      "dedicated",
			Operator: "Exists",
			Effect:   "NoSchedule",
		},
	}
	testMCH := &v1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testmch",
			Namespace: mchNamespace,
		},
		Spec: v1.MultiClusterHubSpec{
			NodeSelector:    mchNodeSelector,
			ImagePullSecret: mchImagePullSecret,
			Tolerations:     mchTolerations,
		},
	}
	containsHTTP := false
	containsHTTPS := false
	containsNO := false
	os.Setenv("POD_NAMESPACE", "default")
	os.Setenv("HTTP_PROXY", "test1")
	os.Setenv("HTTPS_PROXY", "test2")
	os.Setenv("NO_PROXY", "test3")

	testImages := map[string]string{}
	for _, v := range utils.GetTestImages() {
		testImages[v] = "quay.io/test/test:Test"
	}
	// multiple charts
	chartsDir := chartsDir
	templates, errs := RenderCharts(chartsDir, testMCH, testImages)
	if len(errs) > 0 {
		for _, err := range errs {
			t.Logf(err.Error())
		}
		t.Fatalf("failed to retrieve templates")
		if len(templates) == 0 {
			t.Fatalf("Unable to render templates")
		}
	}

	for _, template := range templates {
		if template.GetKind() == "Deployment" {
			deployment := &appsv1.Deployment{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(template.Object, deployment)
			if err != nil {
				t.Fatalf(err.Error())
			}

			selectorEquality := reflect.DeepEqual(deployment.Spec.Template.Spec.NodeSelector, mchNodeSelector)
			if !selectorEquality {
				t.Fatalf("Node Selector did not propagate to the deployments use")
			}
			secretEquality := reflect.DeepEqual(deployment.Spec.Template.Spec.ImagePullSecrets[0].Name, mchImagePullSecret)
			if !secretEquality {
				t.Fatalf("Image Pull Secret did not propagate to the deployments use")
			}
			tolerationEquality := reflect.DeepEqual(deployment.Spec.Template.Spec.Tolerations, mchTolerations)
			if !tolerationEquality {
				t.Fatalf("Toleration did not propagate to the deployments use")
			}
			if deployment.ObjectMeta.Namespace != mchNamespace {
				t.Fatalf("Namespace did not propagate to the deployments use")
			}
			if utils.Contains(proxyList, deployment.ObjectMeta.Name) {
				for _, proxyVar := range deployment.Spec.Template.Spec.Containers[0].Env {
					switch proxyVar.Name {
					case "HTTP_PROXY":
						containsHTTP = true
						if proxyVar.Value != "test1" {
							t.Fatalf("HTTP_PROXY not propagated")
						}
					case "HTTPS_PROXY":
						containsHTTPS = true
						if proxyVar.Value != "test2" {
							t.Fatalf("HTTPS_PROXY not propagated")
						}
					case "NO_PROXY":
						containsNO = true
						if proxyVar.Value != "test3" {
							t.Fatalf("NO_PROXY not propagated")
						}
					}

				}

				if !containsHTTP || !containsHTTPS || !containsNO {
					t.Fatalf("proxy variables not set in %s", deployment.ObjectMeta.Name)
				}
			}
			containsHTTP = false
			containsHTTPS = false
			containsNO = false
		}

	}

	// single chart
	singleChartTestImages := map[string]string{}
	for _, v := range utils.GetTestImages() {
		singleChartTestImages[v] = "quay.io/test/test:Test"
	}
	chartsPath := chartsPath
	singleChartTemplates, errs := RenderChart(chartsPath, testMCH, singleChartTestImages)
	if len(errs) > 0 {
		for _, err := range errs {
			t.Logf(err.Error())
		}
		t.Fatalf("failed to retrieve templates")
		if len(singleChartTemplates) == 0 {
			t.Fatalf("Unable to render templates")
		}
	}
	for _, template := range singleChartTemplates {
		if template.GetKind() == "Deployment" {
			deployment := &appsv1.Deployment{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(template.Object, deployment)
			if err != nil {
				t.Fatalf(err.Error())
			}

			selectorEquality := reflect.DeepEqual(deployment.Spec.Template.Spec.NodeSelector, mchNodeSelector)
			if !selectorEquality {
				t.Fatalf("Node Selector did not propagate to the deployments use")
			}
			secretEquality := reflect.DeepEqual(deployment.Spec.Template.Spec.ImagePullSecrets[0].Name, mchImagePullSecret)
			if !secretEquality {
				t.Fatalf("Image Pull Secret did not propagate to the deployments use")
			}
			tolerationEquality := reflect.DeepEqual(deployment.Spec.Template.Spec.Tolerations, mchTolerations)
			if !tolerationEquality {
				t.Fatalf("Toleration did not propagate to the deployments use")
			}
			if deployment.ObjectMeta.Namespace != mchNamespace {
				t.Fatalf("Namespace did not propagate to the deployments use")
			}

			if utils.Contains(proxyList, deployment.ObjectMeta.Name) {
				for _, proxyVar := range deployment.Spec.Template.Spec.Containers[0].Env {
					switch proxyVar.Name {
					case "HTTP_PROXY":
						containsHTTP = true
						if proxyVar.Value != "test1" {
							t.Fatalf("HTTP_PROXY not propagated")
						}
					case "HTTPS_PROXY":
						containsHTTPS = true
						if proxyVar.Value != "test2" {
							t.Fatalf("HTTPS_PROXY not propagated")
						}
					case "NO_PROXY":
						containsNO = true
						if proxyVar.Value != "test3" {
							t.Fatalf("NO_PROXY not propagated")
						}
					}
				}

				if !containsHTTP || !containsHTTPS || !containsNO {
					t.Fatalf("proxy variables not set")
				}
			}
			containsHTTP = false
			containsHTTPS = false
			containsNO = false
		}

	}

	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")
	os.Unsetenv("NO_PROXY")
	os.Unsetenv("POD_NAMESPACE")

}

func TestRenderCRDs(t *testing.T) {
	tests := []struct {
		name   string
		crdDir string
		want   []error
	}{
		{
			name:   "Render CRDs directory",
			crdDir: crdsDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, errs := RenderCRDs(tt.crdDir)
			if errs != nil && len(errs) > 1 {
				t.Errorf("RenderCRDs() got = %v, want %v", errs, nil)
			}

			for _, u := range got {
				kind := "CustomResourceDefinition"
				apiVersion := "apiextensions.k8s.io/v1"
				if u.GetKind() != kind {
					t.Errorf("RenderCRDs() got Kind = %v, want Kind %v", errs, kind)
				}

				if u.GetAPIVersion() != apiVersion {
					t.Errorf("RenderCRDs() got apiversion = %v, want apiversion %v", errs, apiVersion)
				}
			}
		})
	}
}

func testFailures(t *testing.T) {
	os.Setenv("CRD_OVERRIDE", "pkg/doesnotexist")
	_, errs := RenderCRDs(crdsDir)
	if errs == nil {
		t.Fatalf("Should have received an error")
	}
	os.Unsetenv("CRD_OVERRIDE")
}
