// Copyright (c) 2020 Red Hat, Inc.

package subscription

import (
	subalpha1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	operatorsv1beta1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1beta1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
)

// MongoDB overrides the multicluster-mongodb chart
func MongoDB(m *operatorsv1beta1.MultiClusterHub, cache utils.CacheSpec) *subalpha1.Subscription {
	sub := &Subscription{
		Name:      "multicluster-mongodb",
		Namespace: m.Namespace,
		Overrides: map[string]interface{}{
			"imagePullSecrets": []string{
				m.Spec.ImagePullSecret,
			},
			"network_ip_version": networkVersion(m),
			"auth": map[string]interface{}{
				"enabled":             true,
				"existingAdminSecret": "mongodb-admin",
			},
			"persistentVolume": map[string]interface{}{
				"accessModes": []string{
					"ReadWriteOnce",
				},
				"enabled":      true,
				"size":         m.Spec.Mongo.Storage,
				"storageClass": m.Spec.Mongo.StorageClass,
			},
			"replicas": utils.DefaultReplicaCount(m),
			"tls": map[string]interface{}{
				"casecret": "multicloud-ca-cert",
				"issuer":   "multicloud-ca-issuer",
				"enabled":  true,
			},
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
