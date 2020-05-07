// Copyright (c) 2020 Red Hat, Inc.

package subscription

import (
	subalpha1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
)

// CertManager overrides the cert-manager chart
func CertManager(m *operatorsv1beta1.MultiClusterHub, cache utils.CacheSpec) *subalpha1.Subscription {
	sub := &Subscription{
		Name:      "cert-manager",
		Namespace: utils.CertManagerNS(m),
		Overrides: map[string]interface{}{
			"imagePullSecret": m.Spec.ImagePullSecret,
			"global": map[string]interface{}{
				"isOpenshift":    true,
				"imageOverrides": cache.ImageOverrides,
				"pullPolicy":     utils.GetImagePullPolicy(m),
			},
			"serviceAccount": map[string]interface{}{
				"create": true,
				"name":   "cert-manager",
			},
			"extraEnv": []map[string]interface{}{
				{
					"name":  "OWNED_NAMESPACE",
					"value": utils.CertManagerNS(m),
				},
			},
			"hubconfig": map[string]interface{}{
				"replicaCount": utils.DefaultReplicaCount(m),
				"nodeSelector": m.Spec.NodeSelector,
			},
		},
	}

	return newSubscription(m, sub)
}

// CertWebhook overrides the cert-manager-webhook chart
func CertWebhook(m *operatorsv1beta1.MultiClusterHub, cache utils.CacheSpec) *subalpha1.Subscription {
	sub := &Subscription{
		Name:      "cert-manager-webhook",
		Namespace: utils.CertManagerNS(m),
		Overrides: map[string]interface{}{
			"pkiNamespace": m.Namespace,
			"global": map[string]interface{}{
				"pullSecret":     m.Spec.ImagePullSecret,
				"imageOverrides": cache.ImageOverrides,
			},
			"serviceAccount": map[string]interface{}{
				"create": true,
				"name":   "cert-manager-webhook",
			},
			"hubconfig": map[string]interface{}{
				"replicaCount": utils.DefaultReplicaCount(m),
				"nodeSelector": m.Spec.NodeSelector,
			},
		},
	}

	cainjector := map[string]interface{}{
		"serviceAccount": map[string]interface{}{
			"create": false,
			"name":   "default",
		},
		"hubconfig": map[string]interface{}{
			"replicaCount": utils.DefaultReplicaCount(m),
			"nodeSelector": m.Spec.NodeSelector,
		},
	}

	sub.Overrides["cainjector"] = cainjector

	return newSubscription(m, sub)
}

// ConfigWatcher overrides the configmap-watcher chart
func ConfigWatcher(m *operatorsv1beta1.MultiClusterHub, cache utils.CacheSpec) *subalpha1.Subscription {
	sub := &Subscription{
		Name:      "configmap-watcher",
		Namespace: utils.CertManagerNS(m),
		Overrides: map[string]interface{}{
			"global": map[string]interface{}{
				"pullSecret":     m.Spec.ImagePullSecret,
				"imageOverrides": cache.ImageOverrides,
				"pullPolicy":     utils.GetImagePullPolicy(m),
			},
			"serviceAccount": map[string]interface{}{
				"create": true,
				"name":   "cert-manager-config",
			},
			"hubconfig": map[string]interface{}{
				"replicaCount": utils.DefaultReplicaCount(m),
				"nodeSelector": m.Spec.NodeSelector,
			},
		},
	}

	return newSubscription(m, sub)
}
