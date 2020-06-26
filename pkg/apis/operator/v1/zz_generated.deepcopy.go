// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
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
		*out = make([]corev1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.GlobalPullSecret != nil {
		in, out := &in.GlobalPullSecret, &out.GlobalPullSecret
		*out = new(corev1.LocalObjectReference)
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
func (in *IngressSpec) DeepCopyInto(out *IngressSpec) {
	*out = *in
	if in.SSLCiphers != nil {
		in, out := &in.SSLCiphers, &out.SSLCiphers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpec.
func (in *IngressSpec) DeepCopy() *IngressSpec {
	if in == nil {
		return nil
	}
	out := new(IngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MultiClusterHub) DeepCopyInto(out *MultiClusterHub) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
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
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Hive.DeepCopyInto(&out.Hive)
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.Overrides = in.Overrides
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
func (in *Overrides) DeepCopyInto(out *Overrides) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Overrides.
func (in *Overrides) DeepCopy() *Overrides {
	if in == nil {
		return nil
	}
	out := new(Overrides)
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
