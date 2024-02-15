// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package utils

import (
	"fmt"
	"strings"

	operatorsv1 "github.com/stolostron/multiclusterhub-operator/api/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	/*
		AnnotationMCHPause is an annotation used in multiclusterhub to identify if the multiclusterhub is paused or not.
	*/
	AnnotationMCHPause = "mch-pause"

	/*
		AnnotationImageRepo is an annotation used in multiclusterhub to specify a custom image repository to use.
	*/
	AnnotationImageRepo = "mch-imageRepository"

	/*
		AnnotationImageOverridesCM is an annotation used in multiclusterhub to specify a custom ConfigMap containing
		image overrides.
	*/
	AnnotationImageOverridesCM = "mch-imageOverridesCM"

	/*
		AnnotationLimitOverridesCM is an annotation used in multiclusterhub to specify a custom ConfigMap
		containing resource template overrides.
	*/
	AnnotationTemplateOverridesCM = "operator.multicluster.openshift.io/template-override-cm"

	/*
		AnnotationConfiguration is an annotation used in a resource's metadata to identify the configuration
		last used to create it.
	*/
	AnnotationConfiguration = "installer.open-cluster-management.io/last-applied-configuration"

	/*
		AnnotationMCESubscriptionSpec is an annotation used in multiclusterhub to identify the subscription spec
		last used to create the multiclustengine.
	*/
	AnnotationMCESubscriptionSpec = "installer.open-cluster-management.io/mce-subscription-spec"

	/*
		AnnotationOADPSubscriptionSpec is an annotation used to override the OADP subscription used in cluster-backup.
	*/
	AnnotationOADPSubscriptionSpec = "installer.open-cluster-management.io/oadp-subscription-spec"

	/*
		AnnotationIgnoreOCPVersion is an annotation used to indicate the operator should not check the OpenShift
		Container Platform (OCP) version before proceeding when set.
	*/
	AnnotationIgnoreOCPVersion = "ignoreOCPVersion"

	/*
		AnnotationReleaseVersion is an annotation used to indicate the release version that should be applied to all
		resources managed by the MCH operator.
	*/
	AnnotationReleaseVersion = "installer.open-cluster-management.io/release-version"

	/*
		AnnotationKubeconfig is an annotation used to specify the secret name residing in targetcontaining the
		kubeconfig to access the remote cluster.
	*/
	AnnotationKubeconfig = "mch-kubeconfig"
)

// IsPaused returns true if the multiclusterhub instance is labeled as paused, and false otherwise
func IsPaused(instance *operatorsv1.MultiClusterHub) bool {
	a := instance.GetAnnotations()
	if a == nil {
		return false
	}

	if a[AnnotationMCHPause] != "" && strings.EqualFold(a[AnnotationMCHPause], "true") {
		return true
	}

	return false
}

// AnnotationsMatch returns true if all annotation values used by the operator match
func AnnotationsMatch(old, new map[string]string) bool {
	return old[AnnotationMCHPause] == new[AnnotationMCHPause] &&
		old[AnnotationImageRepo] == new[AnnotationImageRepo] &&
		old[AnnotationImageOverridesCM] == new[AnnotationImageOverridesCM] &&
		old[AnnotationMCESubscriptionSpec] == new[AnnotationMCESubscriptionSpec] &&
		old[AnnotationOADPSubscriptionSpec] == new[AnnotationOADPSubscriptionSpec]
}

// getAnnotation returns the annotation value for a given key, or an empty string if not set
func getAnnotation(instance *operatorsv1.MultiClusterHub, key string) string {
	a := instance.GetAnnotations()
	if a == nil {
		return ""
	}
	return a[key]
}

// GetImageRepository returns the image repo annotation, or an empty string if not set
func GetImageRepository(instance *operatorsv1.MultiClusterHub) string {
	return getAnnotation(instance, AnnotationImageRepo)
}

// GetImageOverridesConfigmapName returns the images override configmap annotation value, or an empty string if not set
func GetImageOverridesConfigmapName(instance *operatorsv1.MultiClusterHub) string {
	return getAnnotation(instance, AnnotationImageOverridesCM)
}

// GetTemplateOverridesConfigmapName returns the templates override configmap annotation value, or an empty string if not set
func GetTemplateOverridesConfigmapName(instance *operatorsv1.MultiClusterHub) string {
	return getAnnotation(instance, AnnotationTemplateOverridesCM)
}

func OverrideImageRepository(imageOverrides map[string]string, imageRepo string) map[string]string {
	for imageKey, imageRef := range imageOverrides {
		image := strings.LastIndex(imageRef, "/")
		imageOverrides[imageKey] = fmt.Sprintf("%s%s", imageRepo, imageRef[image:])
	}
	return imageOverrides
}

func GetMCEAnnotationOverrides(instance *operatorsv1.MultiClusterHub) string {
	return getAnnotation(instance, AnnotationMCESubscriptionSpec)
}

func GetOADPAnnotationOverrides(instance *operatorsv1.MultiClusterHub) string {
	return getAnnotation(instance, AnnotationOADPSubscriptionSpec)
}

// ShouldIgnoreOCPVersion returns true if the instance is annotated to skip
// the minimum OCP version requirement
func ShouldIgnoreOCPVersion(instance *operatorsv1.MultiClusterHub) bool {
	a := instance.GetAnnotations()
	if a == nil {
		return false
	}

	if _, ok := a[AnnotationIgnoreOCPVersion]; ok {
		return true
	}
	return false
}

// GetHostedCredentialsSecret returns the secret namespacedName containing the kubeconfig
// to access the hosted cluster
func GetHostedCredentialsSecret(mch *operatorsv1.MultiClusterHub) (types.NamespacedName, error) {
	nn := types.NamespacedName{}
	if mch.Annotations == nil || mch.Annotations[AnnotationKubeconfig] == "" {
		return nn, fmt.Errorf("no kubeconfig secret annotation defined in %s", mch.Name)
	}
	nn.Name = mch.Annotations[AnnotationKubeconfig]
	nn.Namespace = mch.Namespace
	return nn, nil
}
