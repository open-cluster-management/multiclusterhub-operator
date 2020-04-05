package utils

import (
	"time"

	operatorsv1alpha1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operators/v1alpha1"
)

const (
	// WebhookServiceName ...
	WebhookServiceName = "multiclusterhub-operator-webhook"
	// APIServerSecretName ...
	APIServerSecretName = "mcm-apiserver-self-signed-secrets" // #nosec G101 (no confidential credentials)
	// KlusterletSecretName ...
	KlusterletSecretName = "mcm-klusterlet-self-signed-secrets" // #nosec G101 (no confidential credentials)

	// CertManagerNamespace ...
	CertManagerNamespace = "cert-manager"

	// MongoEndpoints ...
	MongoEndpoints = "multicluster-mongodb"
	// MongoReplicaSet ...
	MongoReplicaSet = "rs0"
	// MongoTLSSecret ...
	MongoTLSSecret = "multicluster-mongodb-client-cert"
	// MongoCaSecret ...
	MongoCaSecret = "multicloud-ca-cert" // #nosec G101 (no confidential credentials)

	podNamespaceEnvVar = "POD_NAMESPACE"
	apiserviceName     = "mcm-apiserver"
	rsaKeySize         = 2048
	duration365d       = time.Hour * 24 * 365

	// DefaultRepository ...
	DefaultRepository = "quay.io/open-cluster-management"
	// LatestVerison ...
	LatestVerison = "latest"
)

// CacheSpec ...
type CacheSpec struct {
	IngressDomain string
}

// CertManagerNS returns the namespace to deploy cert manager objects
func CertManagerNS(m *operatorsv1alpha1.MultiClusterHub) string {
	if m.Spec.CloudPakCompatibility {
		return "cert-manager"
	}
	return m.Namespace

}

// MchIsValid Checks if the optional default parameters need to be set
func MchIsValid(m *operatorsv1alpha1.MultiClusterHub) bool {
	if m.Spec.Version == "" {
		return false
	}

	if m.Spec.ImageRepository == "" {
		return false
	}

	if m.Spec.ImagePullPolicy == "" {
		return false
	}

	if m.Spec.Mongo.Storage == "" {
		return false
	}

	if m.Spec.Mongo.StorageClass == "" {
		return false
	}

	if m.Spec.Etcd.Storage == "" {
		return false
	}

	if m.Spec.Etcd.StorageClass == "" {
		return false
	}

	return true
}
