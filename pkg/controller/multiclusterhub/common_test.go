// Copyright (c) 2020 Red Hat, Inc.

package multiclusterhub

import (
	"context"
	"fmt"
	"os"
	"testing"

	operatorsv1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operator/v1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/channel"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/foundation"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/helmrepo"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/subscription"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func kind(kind string) schema.GroupKind {
	return schema.GroupKind{Group: "", Kind: kind}
}

func Test_ensureDeployment(t *testing.T) {
	r, err := getTestReconciler(full_mch)
	if err != nil {
		t.Fatalf("Failed to create test reconciler")
	}

	cacheSpec := CacheSpec{
		IngressDomain: "apps.smart-buck.dev01.red-chesterfield.com",
		ImageOverrides: map[string]string{
			"application_ui": "quay.io/open-cluster-management/application-ui@sha256:c740fc7bac067f003145ab909504287360564016b7a4a51b7ad4987aca123ac1",
			"console_api":    "quay.io/open-cluster-management/console-api@sha256:3ef1043b4e61a09b07ff37f9ad8fc6e707af9813936cf2c0d52f2fa0e489c75f",
			"rcm_controller": " quay.io/open-cluster-management/rcm-controller@sha256:8fab4d788241bf364dbc1b8c1ea5ccf18d3145a640dbd456b0dc7ba204e36819",
		},
	}

	tests := []struct {
		Name       string
		MCH        *operatorsv1.MultiClusterHub
		Deployment *appsv1.Deployment
		Result     error
	}{
		{
			Name:       "Test: EnsureDeployment - Multiclusterhub-repo",
			MCH:        full_mch,
			Deployment: helmrepo.Deployment(full_mch, cacheSpec.ImageOverrides),
			Result:     nil,
		},
		{
			Name:       "Test: EnsureDeployment - Webhook",
			MCH:        full_mch,
			Deployment: foundation.WebhookDeployment(full_mch, cacheSpec.ImageOverrides),
			Result:     nil,
		},
		{
			Name:       "Test: EnsureDeployment - Empty Deployment",
			MCH:        full_mch,
			Deployment: &appsv1.Deployment{},
			Result:     errors.NewInvalid(kind("Test"), "", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err = r.ensureDeployment(tt.MCH, tt.Deployment)

			if tt.Result != nil {
				// Check if error matches desired error

				if errors.ReasonForError(err) != errors.ReasonForError(tt.Result) {
					t.Fatalf("ensureDeployment() error = %v, wantErr %v", err, tt.Result)
				}
			} else {
				if err != nil {
					t.Fatalf("ensureDeployment() error = %v, wantErr %v", err, tt.Result)
				}

				deploy := &appsv1.Deployment{}
				err = r.client.Get(context.TODO(), types.NamespacedName{
					Name:      tt.Deployment.Name,
					Namespace: tt.Deployment.Namespace,
				}, deploy)

				if err != tt.Result {
					t.Fatalf("Could not find created '%s' deployment: %s", tt.Deployment.Name, err.Error())
				}
			}
		})
	}

}

func Test_ensureService(t *testing.T) {
	r, err := getTestReconciler(full_mch)
	if err != nil {
		t.Fatalf("Failed to create test reconciler")
	}

	tests := []struct {
		Name    string
		MCH     *operatorsv1.MultiClusterHub
		Service *corev1.Service
		Result  error
	}{
		{
			Name:    "Test: ensureService - Multiclusterhub-repo",
			MCH:     full_mch,
			Service: helmrepo.Service(full_mch),
			Result:  nil,
		},
		{
			Name:    "Test: ensureService - Webhook",
			MCH:     full_mch,
			Service: foundation.WebhookService(full_mch),
			Result:  nil,
		},
		{
			Name:    "Test: ensureService - Empty service",
			MCH:     full_mch,
			Service: &corev1.Service{},
			Result:  errors.NewInvalid(kind("Test"), "", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err = r.ensureService(tt.MCH, tt.Service)

			if tt.Result != nil {
				// Check if error matches desired error
				if errors.ReasonForError(err) != errors.ReasonForError(tt.Result) {
					t.Fatalf("ensureService() error = %v, wantErr %v", err, tt.Result)
				}
			} else {
				if err != nil {
					t.Fatalf("ensureService() error = %v, wantErr %v", err, tt.Result)
				}

				service := &corev1.Service{}
				err = r.client.Get(context.TODO(), types.NamespacedName{
					Name:      tt.Service.Name,
					Namespace: tt.Service.Namespace,
				}, service)
				if err != tt.Result {
					t.Fatalf("Could not find created '%s' service: %s", tt.Service.Name, err.Error())
				}
			}
		})
	}
}

func Test_ensureChannel(t *testing.T) {
	r, err := getTestReconciler(full_mch)
	if err != nil {
		t.Fatalf("Failed to create test reconciler")
	}

	tests := []struct {
		Name    string
		MCH     *operatorsv1.MultiClusterHub
		Channel *unstructured.Unstructured
		Result  error
	}{
		{
			Name:    "Test: ensureChannel - Multiclusterhub-repo",
			MCH:     full_mch,
			Channel: channel.Channel(full_mch),
			Result:  nil,
		},
		{
			Name:    "Test: ensureChannel - Empty channel",
			MCH:     full_mch,
			Channel: &unstructured.Unstructured{},
			Result:  fmt.Errorf("Object 'Kind' is missing in 'unstructured object has no kind'"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err = r.ensureChannel(tt.MCH, tt.Channel)
			if !errorEquals(err, tt.Result) {
				t.Errorf("ensureChannel() error = %v, wantErr %v", err, tt.Result)
			}
		})

		// TODO: Check Channel is created in the fake client
	}
}

func Test_ensureSubscription(t *testing.T) {
	os.Setenv("UNIT_TEST", "true")
	defer os.Unsetenv("UNIT_TEST")

	r, err := getTestReconciler(full_mch)
	if err != nil {
		t.Fatalf("Failed to create test reconciler")
	}

	cacheSpec := CacheSpec{
		IngressDomain: "apps.smart-buck.dev01.red-chesterfield.com",
		ImageOverrides: map[string]string{
			"application_ui": "quay.io/open-cluster-management/application-ui@sha256:c740fc7bac067f003145ab909504287360564016b7a4a51b7ad4987aca123ac1",
			"console_api":    "quay.io/open-cluster-management/console-api@sha256:3ef1043b4e61a09b07ff37f9ad8fc6e707af9813936cf2c0d52f2fa0e489c75f",
			"rcm_controller": " quay.io/open-cluster-management/rcm-controller@sha256:8fab4d788241bf364dbc1b8c1ea5ccf18d3145a640dbd456b0dc7ba204e36819",
		},
	}

	tests := []struct {
		Name         string
		MCH          *operatorsv1.MultiClusterHub
		Subscription *unstructured.Unstructured
		Result       error
	}{
		{
			Name:         "Test: ensureSubscription - Cert-manager",
			MCH:          full_mch,
			Subscription: subscription.CertManager(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Cert-webhook",
			MCH:          full_mch,
			Subscription: subscription.CertWebhook(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Config-watcher",
			MCH:          full_mch,
			Subscription: subscription.ConfigWatcher(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Management-ingress",
			MCH:          full_mch,
			Subscription: subscription.ManagementIngress(full_mch, cacheSpec.ImageOverrides, cacheSpec.IngressDomain),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Application-UI",
			MCH:          full_mch,
			Subscription: subscription.ApplicationUI(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Console",
			MCH:          full_mch,
			Subscription: subscription.Console(full_mch, cacheSpec.ImageOverrides, cacheSpec.IngressDomain),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - GRC",
			MCH:          full_mch,
			Subscription: subscription.GRC(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - KUI",
			MCH:          full_mch,
			Subscription: subscription.KUIWebTerminal(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - RCM",
			MCH:          full_mch,
			Subscription: subscription.RCM(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Search",
			MCH:          full_mch,
			Subscription: subscription.Search(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Topology",
			MCH:          full_mch,
			Subscription: subscription.Topology(full_mch, cacheSpec.ImageOverrides),
			Result:       nil,
		},
		{
			Name:         "Test: ensureSubscription - Empty Sub",
			MCH:          full_mch,
			Subscription: &unstructured.Unstructured{},
			Result:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err = r.ensureSubscription(tt.MCH, tt.Subscription)
			if !errorEquals(err, tt.Result) {
				t.Errorf("ensureSubscription() error = %v, wantErr %v", err, tt.Result)
			}

			// TODO: Check Subscription is created in the fake client
		})
	}
}

func Test_ensureClusterManager(t *testing.T) {
	r, err := getTestReconciler(full_mch)
	if err != nil {
		t.Fatalf("Failed to create test reconciler")
	}

	imageOverrides := map[string]string{
		"registration": "quay.io/open-cluster-management/registration@sha256:fe95bca419976ca8ffe608bc66afcead6ef333b863f22be55df57c89ded75dda",
	}

	tests := []struct {
		Name           string
		MCH            *operatorsv1.MultiClusterHub
		ClusterManager *unstructured.Unstructured
		Result         error
	}{
		{
			Name:           "Test: ensureClusterManager - ClusterManager",
			MCH:            full_mch,
			ClusterManager: foundation.ClusterManager(full_mch, imageOverrides),
			Result:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err = r.ensureClusterManager(tt.MCH, tt.ClusterManager)
			if !errorEquals(err, tt.Result) {
				t.Fatalf("Failed to ensure ClusterManager: %s", err)
			}
		})
	}
}
