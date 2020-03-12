package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MultiCloudHubSpec defines the desired state of MultiCloudHub
// +k8s:openapi-gen=true
type MultiCloudHubSpec struct {
	// Version of the MultiCloud hub
	Version string `json:"version"`

	// Repository of the MultiCloud hub images
	ImageRepository string `json:"imageRepository"`

	// ImageTagSuffix of the MultiCloud hub images
	ImageTagSuffix string `json:"imageTagSuffix"`

	// Pull policy of the MultiCloud hub images
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy"`

	// Pull secret of the MultiCloud hub images
	// +optional
	ImagePullSecret string `json:"imagePullSecret,omitempty"`

	// Hostname portion of OCP Domain
	// +optional
	OCPHOST string `json:"ocpHost,omitempty"`

	// Spec of NodeSelector
	// +optional
	NodeSelector *NodeSelector `json:"nodeSelector,omitempty"`

	// Spec of foundation
	Foundation `json:"foundation"`

	// Spec of etcd
	Etcd `json:"etcd"`

	// Spec of hive
	Hive HiveConfigSpec `json:"hive"`

	// Spec of mongo
	Mongo `json:"mongo"`
}

// NodeSelector defines the desired state of NodeSelector
type NodeSelector struct {
	// Spec of OS
	// +optional
	OS string `json:"os,omitempty"`

	// Spec of CustomLabelSelector
	// +optional
	CustomLabelSelector string `json:"customLabelSelector,omitempty"`

	// Spec of CustomLabelValue
	// +optional
	CustomLabelValue string `json:"customLabelValue,omitempty"`
}

// Foundation defines the desired state of MultiCloudHub foundation components
type Foundation struct {
	// Spec of apiserver
	// +optional
	Apiserver `json:"apiserver,omitempty"`

	// Spec of controller
	// +optional
	Controller `json:"controller,omitempty"`
}

type Apiserver struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Certificates of API server
	// +optional
	ApiserverSecret string `json:"apiserverSecret,omitempty"`

	// Certificates of Klusterlet
	// +optional
	KlusterletSecret string `json:"klusterletSecret,omitempty"`

	// Configuration of the pod
	// +optional
	Configuration map[string]string `json:"configuration,omitempty"`
}

type Controller struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Configuration of the pod
	// +optional
	Configuration map[string]string `json:"configuration,omitempty"`
}

// Etcd defines the desired state of etcd
type Etcd struct {
	// Endpoints of etcd
	Endpoints string `json:"endpoints"`

	// Secret of etcd
	// +optional
	Secret string `json:"secret,omitempty"`
}

type HiveConfigSpec struct {
	// ManagedDomains is the list of DNS domains that are allowed to be used by the 'managedDNS' feature.
	// When specifying 'managedDNS: true' in a ClusterDeployment, the ClusterDeployment's
	// baseDomain must be a direct child of one of these domains, otherwise the
	// ClusterDeployment creation will result in a validation error.
	// +optional
	ManagedDomains []string `json:"managedDomains,omitempty"`

	// ExternalDNS specifies configuration for external-dns if it is to be deployed by
	// Hive. If absent, external-dns will not be deployed.
	// +optional
	ExternalDNS *ExternalDNSConfig `json:"externalDNS,omitempty"`

	// AdditionalCertificateAuthorities is a list of references to secrets in the
	// 'hive' namespace that contain an additional Certificate Authority to use when communicating
	// with target clusters. These certificate authorities will be used in addition to any self-signed
	// CA generated by each cluster on installation.
	// +optional
	AdditionalCertificateAuthorities []corev1.LocalObjectReference `json:"additionalCertificateAuthorities,omitempty"`

	// GlobalPullSecret is used to specify a pull secret that will be used globally by all of the cluster deployments.
	// For each cluster deployment, the contents of GlobalPullSecret will be merged with the specific pull secret for
	// a cluster deployment(if specified), with precedence given to the contents of the pull secret for the cluster deployment.
	// +optional
	GlobalPullSecret *corev1.LocalObjectReference `json:"globalPullSecret,omitempty"`

	// Backup specifies configuration for backup integration.
	// If absent, backup integration will be disabled.
	// +optional
	Backup BackupConfig `json:"backup,omitempty"`

	// FailedProvisionConfig is used to configure settings related to handling provision failures.
	FailedProvisionConfig FailedProvisionConfig `json:"failedProvisionConfig"`

	// MaintenanceMode can be set to true to disable the hive controllers in situations where we need to ensure
	// nothing is running that will add or act upon finalizers on Hive types. This should rarely be needed.
	// Sets replicas to 0 for the hive-controllers deployment to accomplish this.
	MaintenanceMode *bool `json:"maintenanceMode,omitempty"`
}

// HiveConfigStatus defines the observed state of Hive
type HiveConfigStatus struct {
	// AggregatorClientCAHash keeps an md5 hash of the aggregator client CA
	// configmap data from the openshift-config-managed namespace. When the configmap changes,
	// admission is redeployed.
	AggregatorClientCAHash string `json:"aggregatorClientCAHash,omitempty"`
}

// BackupConfig contains settings for the Velero backup integration.
type BackupConfig struct {
	// Velero specifies configuration for the Velero backup integration.
	// +optional
	Velero VeleroBackupConfig `json:"velero,omitempty"`

	// MinBackupPeriodSeconds specifies that a minimum of MinBackupPeriodSeconds will occur in between each backup.
	// This is used to rate limit backups. This potentially batches together multiple changes into 1 backup.
	// No backups will be lost as changes that happen during this interval are queued up and will result in a
	// backup happening once the interval has been completed.
	// +optional
	MinBackupPeriodSeconds *int `json:"minBackupPeriodSeconds,omitempty"`
}

// VeleroBackupConfig contains settings for the Velero backup integration.
type VeleroBackupConfig struct {
	// Enabled dictates if Velero backup integration is enabled.
	// If not specified, the default is disabled.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}

// FailedProvisionConfig contains settings to control behavior undertaken by Hive when an installation attempt fails.
type FailedProvisionConfig struct {

	// SkipGatherLogs disables functionality that attempts to gather full logs from the cluster if an installation
	// fails for any reason. The logs will be stored in a persistent volume for up to 7 days.
	SkipGatherLogs bool `json:"skipGatherLogs,omitempty"`
}

// ExternalDNSConfig contains settings for running external-dns in a Hive
// environment.
type ExternalDNSConfig struct {

	// AWS contains AWS-specific settings for external DNS
	// +optional
	AWS *ExternalDNSAWSConfig `json:"aws,omitempty"`

	// GCP contains GCP-specific settings for external DNS
	// +optional
	GCP *ExternalDNSGCPConfig `json:"gcp,omitempty"`

	// As other cloud providers are supported, additional fields will be
	// added for each of those cloud providers. Only a single cloud provider
	// may be configured at a time.
}

// ExternalDNSAWSConfig contains AWS-specific settings for external DNS
type ExternalDNSAWSConfig struct {
	// Credentials references a secret that will be used to authenticate with
	// AWS Route53. It will need permission to manage entries in each of the
	// managed domains for this cluster.
	// Secret should have AWS keys named 'aws_access_key_id' and 'aws_secret_access_key'.
	// +optional
	Credentials corev1.LocalObjectReference `json:"credentials,omitempty"`
}

// ExternalDNSGCPConfig contains GCP-specific settings for external DNS
type ExternalDNSGCPConfig struct {
	// Credentials references a secret that will be used to authenticate with
	// GCP DNS. It will need permission to manage entries in each of the
	// managed domains for this cluster.
	// Secret should have a key names 'osServiceAccount.json'.
	// The credentials must specify the project to use.
	// +optional
	Credentials corev1.LocalObjectReference `json:"credentials,omitempty"`
}

// Mongo defines the desired state of mongo
type Mongo struct {
	// Endpoints of mongo
	Endpoints string `json:"endpoints"`

	// Replica set of mongo
	ReplicaSet string `json:"replicaSet"`

	// User secret of mongo
	// +optional
	UserSecret string `json:"userSecret,omitempty"`

	// TLS secret of mongo
	// +optional
	TLSSecret string `json:"tlsSecret,omitempty"`

	// CA secret of mongo
	// +optional
	CASecret string `json:"caSecret,omitempty"`
}

// MultiCloudHubStatus defines the observed state of MultiCloudHub
// +k8s:openapi-gen=true
type MultiCloudHubStatus struct {
	// Represents the running phase of the MultiCloudHub
	Phase string `json:"phase"`

	// Represents the status of each deployment
	// +optional
	Deployments []DeploymentResult `json:"deployments,omitempty"`
}

// DeploymentResult defines the observed state of Deployment
type DeploymentResult struct {
	// Name of the deployment
	Name string `json:"name"`

	// The most recently observed status of the Deployment
	Status appsv1.DeploymentStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MultiCloudHub is the Schema for the multicloudhubs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=multicloudhubs,scope=Namespaced
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="Multicloudhub Operator"
type MultiCloudHub struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiCloudHubSpec   `json:"spec,omitempty"`
	Status MultiCloudHubStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MultiCloudHubList contains a list of MultiCloudHub
type MultiCloudHubList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiCloudHub `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MultiCloudHub{}, &MultiCloudHubList{})
}
