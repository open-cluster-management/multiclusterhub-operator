// Copyright (c) 2020 Red Hat, Inc.

package subscription

import (
	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// Console overrides the console-chart chart
func Console(m *operatorsv1beta1.MultiClusterHub, cache utils.CacheSpec) *unstructured.Unstructured {
	sub := &Subscription{
		Name:      "console-chart",
		Namespace: m.Namespace,
		Overrides: map[string]interface{}{
			"pullSecret":   m.Spec.ImagePullSecret,
			"ocpingress":   cache.IngressDomain,
			"cfcRouterUrl": "https://management-ingress:443",
			"hubconfig": map[string]interface{}{
				"replicaCount": utils.DefaultReplicaCount(m),
				"nodeSelector": m.Spec.NodeSelector,
			},
			"global": map[string]interface{}{
				"imageOverrides": cache.ImageOverrides,
				"pullPolicy":     utils.GetImagePullPolicy(m),
			},
		},
	}

	return newSubscription(m, sub)
}
