// Copyright (c) 2020 Red Hat, Inc.

package multiclusterhub

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	subrelv1 "github.com/open-cluster-management/multicloud-operators-subscription-release/pkg/apis/apps/v1"
	operatorsv1 "github.com/open-cluster-management/multicloudhub-operator/pkg/apis/operator/v1"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/foundation"
	"github.com/open-cluster-management/multicloudhub-operator/pkg/helmrepo"
	"github.com/open-cluster-management/multicloudhub-operator/version"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var deployments = []types.NamespacedName{
	{Name: helmrepo.HelmRepoName, Namespace: "open-cluster-management"},
	{Name: foundation.OCMControllerName, Namespace: "open-cluster-management"},
	{Name: foundation.OCMProxyServerName, Namespace: "open-cluster-management"},
	{Name: foundation.WebhookName, Namespace: "open-cluster-management"},
}

var appsubs = []types.NamespacedName{
	{Name: "application-chart-sub", Namespace: "open-cluster-management"},
	{Name: "cert-manager-sub", Namespace: "open-cluster-management"},
	{Name: "cert-manager-webhook-sub", Namespace: "open-cluster-management"},
	{Name: "configmap-watcher-sub", Namespace: "open-cluster-management"},
	{Name: "console-chart-sub", Namespace: "open-cluster-management"},
	{Name: "grc-sub", Namespace: "open-cluster-management"},
	{Name: "kui-web-terminal-sub", Namespace: "open-cluster-management"},
	{Name: "management-ingress-sub", Namespace: "open-cluster-management"},
	{Name: "rcm-sub", Namespace: "open-cluster-management"},
	{Name: "search-prod-sub", Namespace: "open-cluster-management"},
	{Name: "topology-sub", Namespace: "open-cluster-management"},
}

var unknownStatus = operatorsv1.StatusCondition{
	Type:               "Unknown",
	Status:             metav1.ConditionUnknown,
	LastUpdateTime:     metav1.Now(),
	LastTransitionTime: metav1.Now(),
	Reason:             "No conditions available",
	Message:            "No conditions available",
}

// UpdateStatus updates status
func (r *ReconcileMultiClusterHub) UpdateStatus(m *operatorsv1.MultiClusterHub) (reconcile.Result, error) {
	oldStatus := m.Status
	newStatus := m.Status.DeepCopy()
	newStatus.DesiredVersion = version.Version

	components := make(map[string]operatorsv1.StatusCondition)

	deployment := &appsv1.Deployment{}
	for i, d := range deployments {
		r.client.Get(context.TODO(), deployments[i], deployment)
		components[d.Name] = mapDeployment(deployment)
	}

	for _, d := range appsubs {
		components[d.Name] = unknownStatus
	}

	hrList := &subrelv1.HelmReleaseList{}
	r.client.List(context.TODO(), hrList)
	for _, hr := range hrList.Items {
		if _, ok := components[hr.OwnerReferences[0].Name]; ok {
			components[hr.OwnerReferences[0].Name] = mapHelmRelease(&hr)
		}
	}

	newStatus.Phase = aggregateStatus(components)

	newStatus.CurrentVersion = oldStatus.CurrentVersion
	if newStatus.Phase == operatorsv1.HubRunning {
		newStatus.CurrentVersion = version.Version
	}

	m.Status = *newStatus

	return r.applyStatus(m)
}

func (r *ReconcileMultiClusterHub) applyStatus(m *operatorsv1.MultiClusterHub) (reconcile.Result, error) {
	err := r.client.Status().Update(context.TODO(), m)
	if err != nil {
		if errors.IsConflict(err) {
			// Error from object being modified is normal behavior and should not be treated like an error
			log.Info("Failed to update status", "Reason", "Object has been modified")
			return reconcile.Result{RequeueAfter: resyncPeriod}, nil
		}

		log.Error(err, fmt.Sprintf("Failed to update %s/%s status ", m.Namespace, m.Name))
		return reconcile.Result{}, err
	}

	if m.Status.Phase != operatorsv1.HubRunning {
		return reconcile.Result{RequeueAfter: resyncPeriod}, nil
	} else {
		return reconcile.Result{}, nil
	}
}

func successfulDeploy(d *appsv1.Deployment) bool {
	latest := latestDeployCondition(d.Status.Conditions)
	return latest.Type == appsv1.DeploymentAvailable && latest.Status == corev1.ConditionTrue
}

func latestDeployCondition(conditions []appsv1.DeploymentCondition) appsv1.DeploymentCondition {
	if len(conditions) < 1 {
		return appsv1.DeploymentCondition{}
	}
	latest := conditions[0]
	for i := range conditions {
		if conditions[i].LastTransitionTime.Time.After(latest.LastTransitionTime.Time) {
			latest = conditions[i]
		}
	}
	return latest
}

func mapDeployment(ds *appsv1.Deployment) operatorsv1.StatusCondition {
	if len(ds.Status.Conditions) < 1 {
		return unknownStatus
	}

	dcs := latestDeployCondition(ds.Status.Conditions)
	ret := operatorsv1.StatusCondition{
		Type:               string(dcs.Type),
		Status:             metav1.ConditionStatus(string(dcs.Status)),
		LastUpdateTime:     dcs.LastUpdateTime,
		LastTransitionTime: dcs.LastTransitionTime,
		Reason:             dcs.Reason,
		Message:            dcs.Message,
	}
	if successfulDeploy(ds) {
		ret.Message = ""
	}

	return ret
}

func successfulHelmRelease(hr *subrelv1.HelmRelease) bool {
	latest := latestHelmReleaseCondition(hr.Status.Conditions)
	return latest.Type == subrelv1.ConditionDeployed && latest.Status == subrelv1.StatusTrue
}

func latestHelmReleaseCondition(conditions []subrelv1.HelmAppCondition) subrelv1.HelmAppCondition {
	if len(conditions) < 1 {
		return subrelv1.HelmAppCondition{}
	}
	latest := conditions[0]
	for i := range conditions {
		if conditions[i].LastTransitionTime.Time.After(latest.LastTransitionTime.Time) {
			latest = conditions[i]
		}
	}
	return latest
}

func mapHelmRelease(hr *subrelv1.HelmRelease) operatorsv1.StatusCondition {
	if len(hr.Status.Conditions) < 1 {
		return unknownStatus
	}

	condition := latestHelmReleaseCondition(hr.Status.Conditions)
	ret := operatorsv1.StatusCondition{
		Type:               string(condition.Type),
		Status:             metav1.ConditionStatus(condition.Status),
		LastUpdateTime:     metav1.Now(),
		LastTransitionTime: condition.LastTransitionTime,
		Reason:             string(condition.Reason),
		Message:            condition.Message,
	}
	// Ignore success messages
	if !isErrorType(ret.Type) {
		ret.Message = ""
	}
	return ret
}

func successfulComponent(sc operatorsv1.StatusCondition) bool {
	return (sc.Status == metav1.ConditionTrue) && (sc.Type == "Available" || sc.Type == "Deployed")
}

func aggregateStatus(components map[string]operatorsv1.StatusCondition) operatorsv1.HubPhaseType {
	for k, val := range components {
		if !successfulComponent(val) {
			log.Info("Waiting on", "name", k)
			return operatorsv1.HubPending
		}
	}
	return operatorsv1.HubRunning
}

func isErrorType(cr string) bool {
	return cr == string(subrelv1.ReasonInstallError) ||
		cr == string(subrelv1.ReasonUpdateError) ||
		cr == string(subrelv1.ReasonReconcileError) ||
		cr == string(subrelv1.ReasonUninstallError)
}
