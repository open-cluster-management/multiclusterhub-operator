// Copyright (c) 2020 Red Hat, Inc.

package mcm

import (
	"testing"

	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestControllerDeployment(t *testing.T) {
	empty := &operatorsv1beta1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec: operatorsv1beta1.MultiClusterHubSpec{
			ImagePullSecret: "",
			Mongo:           operatorsv1beta1.Mongo{},
		},
	}

	cs := utils.CacheSpec{
		IngressDomain:  "testIngress",
		ImageOverrides: map[string]string{},
	}

	t.Run("MCH with empty fields", func(t *testing.T) {
		_ = ControllerDeployment(empty, cs)
	})

	essentialsOnly := &operatorsv1beta1.MultiClusterHub{
		ObjectMeta: metav1.ObjectMeta{Namespace: "test"},
		Spec:       operatorsv1beta1.MultiClusterHubSpec{},
	}
	t.Run("MCH with only required values", func(t *testing.T) {
		_ = ControllerDeployment(essentialsOnly, cs)
	})
}
