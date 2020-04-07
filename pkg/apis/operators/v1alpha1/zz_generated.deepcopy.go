// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupConfig) DeepCopyInto(out *BackupConfig) {
	*out = *in
	out.Velero = in.Velero
	if in.MinBackupPeriodSeconds != nil {
		in, out := &in.MinBackupPeriodSeconds, &out.MinBackupPeriodSeconds
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupConfig.
func (in *BackupConfig) DeepCopy() *BackupConfig {
	if in == nil {
		return nil
	}
	out := new(BackupConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentResult) DeepCopyInto(out *DeploymentResult) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentResult.
func (in *DeploymentResult) DeepCopy() *DeploymentResult {
	if in == nil {
		return nil
	}
	out := new(DeploymentResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Etcd) DeepCopyInto(out *Etcd) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Etcd.
func (in *Etcd) DeepCopy() *Etcd {
	if in == nil {
		return nil
	}
	out := new(Etcd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSAWSConfig) DeepCopyInto(out *ExternalDNSAWSConfig) {
	*out = *in
	out.Credentials = in.Credentials
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSAWSConfig.
func (in *ExternalDNSAWSConfig) DeepCopy() *ExternalDNSAWSConfig {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSAWSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSConfig) DeepCopyInto(out *ExternalDNSConfig) {
	*out = *in
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(ExternalDNSAWSConfig)
		**out = **in
	}
	if in.GCP != nil {
		in, out := &in.GCP, &out.GCP
		*out = new(ExternalDNSGCPConfig)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSConfig.
func (in *ExternalDNSConfig) DeepCopy() *ExternalDNSConfig {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDNSGCPConfig) DeepCopyInto(out *ExternalDNSGCPConfig) {
	*out = *in
	out.Credentials = in.Credentials
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDNSGCPConfig.
func (in *ExternalDNSGCPConfig) DeepCopy() *ExternalDNSGCPConfig {
	if in == nil {
		return nil
	}
	out := new(ExternalDNSGCPConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FailedProvisionConfig) DeepCopyInto(out *FailedProvisionConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FailedProvisionConfig.
func (in *FailedProvisionConfig) DeepCopy() *FailedProvisionConfig {
	if in == nil {
		return nil
	}
	out := new(FailedProvisionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HiveConfigSpec) DeepCopyInto(out *HiveConfigSpec) {
	*out = *in
	if in.ExternalDNS != nil {
		in, out := &in.ExternalDNS, &out.ExternalDNS
		*out = new(ExternalDNSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.AdditionalCertificateAuthorities != nil {
		in, out := &in.AdditionalCertificateAuthorities, &out.AdditionalCertificateAuthorities
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.GlobalPullSecret != nil {
		in, out := &in.GlobalPullSecret, &out.GlobalPullSecret
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	in.Backup.DeepCopyInto(&out.Backup)
	out.FailedProvisionConfig = in.FailedProvisionConfig
	if in.MaintenanceMode != nil {
		in, out := &in.MaintenanceMode, &out.MaintenanceMode
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HiveConfigSpec.
func (in *HiveConfigSpec) DeepCopy() *HiveConfigSpec {
	if in == nil {
		return nil
	}
	out := new(HiveConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HiveConfigStatus) DeepCopyInto(out *HiveConfigStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HiveConfigStatus.
func (in *HiveConfigStatus) DeepCopy() *HiveConfigStatus {
	if in == nil {
		return nil
	}
	out := new(HiveConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mongo) DeepCopyInto(out *Mongo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mongo.
func (in *Mongo) DeepCopy() *Mongo {
	if in == nil {
		return nil
	}
	out := new(Mongo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiClusterHub) DeepCopyInto(out *MultiClusterHub) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiClusterHub.
func (in *MultiClusterHub) DeepCopy() *MultiClusterHub {
	if in == nil {
		return nil
	}
	out := new(MultiClusterHub)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MultiClusterHub) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiClusterHubList) DeepCopyInto(out *MultiClusterHubList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MultiClusterHub, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiClusterHubList.
func (in *MultiClusterHubList) DeepCopy() *MultiClusterHubList {
	if in == nil {
		return nil
	}
	out := new(MultiClusterHubList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MultiClusterHubList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiClusterHubSpec) DeepCopyInto(out *MultiClusterHubSpec) {
	*out = *in
	if in.ReplicaCount != nil {
		in, out := &in.ReplicaCount, &out.ReplicaCount
		*out = new(int32)
		**out = **in
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = new(NodeSelector)
		**out = **in
	}
	in.Hive.DeepCopyInto(&out.Hive)
	out.Mongo = in.Mongo
	out.Etcd = in.Etcd
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiClusterHubSpec.
func (in *MultiClusterHubSpec) DeepCopy() *MultiClusterHubSpec {
	if in == nil {
		return nil
	}
	out := new(MultiClusterHubSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiClusterHubStatus) DeepCopyInto(out *MultiClusterHubStatus) {
	*out = *in
	if in.Deployments != nil {
		in, out := &in.Deployments, &out.Deployments
		*out = make([]DeploymentResult, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MultiClusterHubStatus.
func (in *MultiClusterHubStatus) DeepCopy() *MultiClusterHubStatus {
	if in == nil {
		return nil
	}
	out := new(MultiClusterHubStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeSelector) DeepCopyInto(out *NodeSelector) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeSelector.
func (in *NodeSelector) DeepCopy() *NodeSelector {
	if in == nil {
		return nil
	}
	out := new(NodeSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroBackupConfig) DeepCopyInto(out *VeleroBackupConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroBackupConfig.
func (in *VeleroBackupConfig) DeepCopy() *VeleroBackupConfig {
	if in == nil {
		return nil
	}
	out := new(VeleroBackupConfig)
	in.DeepCopyInto(out)
	return out
}
