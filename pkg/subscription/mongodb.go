package subscription

import (
	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// MongoDB overrides the multicluster-mongodb chart
func MongoDB(m *operatorsv1alpha1.MultiClusterHub, cache utils.CacheSpec) *unstructured.Unstructured {
	sub := &Subscription{
		Name:      "multicluster-mongodb",
		Namespace: m.Namespace,
		Overrides: map[string]interface{}{
			"imageTagPostfix": imageSuffix(m),
			"imagePullSecrets": []string{
				m.Spec.ImagePullSecret,
			},
			"network_ip_version": networkVersion(m),
			"auth": map[string]interface{}{
				"enabled":             true,
				"existingAdminSecret": "mongodb-admin",
			},
			"image": map[string]interface{}{
				"repository": m.Spec.ImageRepository,
				"pullPolicy": m.Spec.ImagePullPolicy,
			},
			"installImage": map[string]interface{}{
				"repository": m.Spec.ImageRepository,
				"pullPolicy": m.Spec.ImagePullPolicy,
			},
			"persistentVolume": map[string]interface{}{
				"accessModes": []string{
					"ReadWriteOnce",
				},
				"enabled":      true,
				"size":         m.Spec.Mongo.Storage,
				"storageClass": m.Spec.Mongo.StorageClass,
			},
			"curl": map[string]interface{}{
				"image": map[string]interface{}{
					"repository": m.Spec.ImageRepository,
					"pullPolicy": m.Spec.ImagePullPolicy,
				},
			},
			"replicas": m.Spec.Mongo.ReplicaCount,
			"tls": map[string]interface{}{
				"casecret": "multicloud-ca-cert",
				"issuer":   "multicloud-ca-issuer",
				"enabled":  true,
			},
			"hubconfig": map[string]interface{}{
				"replicaCount": m.Spec.Mongo.ReplicaCount,
				"nodeSelector": m.Spec.NodeSelector,
			},
		},
	}

	if cache.ImageShaDigests != nil {
		sub.Overrides["imageShaDigests"] = cache.ImageShaDigests
	}

	return newSubscription(m, sub)
}
