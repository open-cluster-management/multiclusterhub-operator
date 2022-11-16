// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controllers

import (
	"context"
	"encoding/json"
	e "errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	consolev1 "github.com/openshift/api/operator/v1"

	"github.com/Masterminds/semver"
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	subv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	searchv2v1alpha1 "github.com/stolostron/search-v2-operator/api/v1alpha1"
	apixv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	subhelmv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/helmrelease/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	utils "github.com/stolostron/multiclusterhub-operator/pkg/utils"

	operatorv1 "github.com/stolostron/multiclusterhub-operator/api/v1"
	"github.com/stolostron/multiclusterhub-operator/pkg/channel"
	"github.com/stolostron/multiclusterhub-operator/pkg/helmrepo"
	"github.com/stolostron/multiclusterhub-operator/pkg/manifest"
	"github.com/stolostron/multiclusterhub-operator/pkg/multiclusterengine"
	"github.com/stolostron/multiclusterhub-operator/pkg/subscription"
	"github.com/stolostron/multiclusterhub-operator/pkg/version"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CacheSpec ...
type CacheSpec struct {
	IngressDomain    string
	ImageOverrides   map[string]string
	ImageRepository  string
	ManifestVersion  string
	ImageOverridesCM string
}

func (r *MultiClusterHubReconciler) ensureDeployment(m *operatorv1.MultiClusterHub, dep *appsv1.Deployment) (ctrl.Result, error) {
	r.Log.Info("Reconciling MultiClusterHub")

	if utils.ProxyEnvVarsAreSet() {
		dep = addProxyEnvVarsToDeployment(dep)
	}

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: m.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Client.Create(context.TODO(), dep)
		if err != nil {
			// Deployment failed
			r.Log.Error(err, "Failed to create new Deployment")
			return ctrl.Result{}, err
		}

		// Deployment was successful
		r.Log.Info("Created a new Deployment")
		condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
		SetHubCondition(&m.Status, *condition)
		return ctrl.Result{}, nil

	} else if err != nil {
		// Error that isn't due to the deployment not existing
		r.Log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Validate object based on name
	var desired *appsv1.Deployment
	var needsUpdate bool

	switch found.Name {
	case helmrepo.HelmRepoName:
		desired, needsUpdate = helmrepo.ValidateDeployment(m, r.CacheSpec.ImageOverrides, dep, found)
	default:
		r.Log.Info("Could not validate deployment; unknown name")
		return ctrl.Result{}, nil
	}

	if needsUpdate {
		err = r.Client.Update(context.TODO(), desired)
		if err != nil {
			r.Log.Error(err, "Failed to update Deployment.")
			return ctrl.Result{}, err
		}
		// Spec updated - return
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureService(m *operatorv1.MultiClusterHub, s *corev1.Service) (ctrl.Result, error) {
	svlog := r.Log.WithValues("Service.Namespace", s.Namespace, "Service.Name", s.Name)

	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: m.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			svlog.Error(err, "Failed to create new Service")
			return ctrl.Result{}, err
		}

		// Creation was successful
		svlog.Info("Created a new Service")
		condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
		SetHubCondition(&m.Status, *condition)
		return ctrl.Result{}, nil

	} else if err != nil {
		// Error that isn't due to the service not existing
		svlog.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	modified := resourcemerge.BoolPtr(false)
	existingCopy := found.DeepCopy()
	resourcemerge.EnsureObjectMeta(modified, &existingCopy.ObjectMeta, s.ObjectMeta)
	selectorSame := equality.Semantic.DeepEqual(existingCopy.Spec.Selector, s.Spec.Selector)

	typeSame := false
	requiredIsEmpty := len(s.Spec.Type) == 0
	existingCopyIsCluster := existingCopy.Spec.Type == corev1.ServiceTypeClusterIP
	if (requiredIsEmpty && existingCopyIsCluster) || equality.Semantic.DeepEqual(existingCopy.Spec.Type, s.Spec.Type) {
		typeSame = true
	}

	if selectorSame && typeSame && !*modified {
		return ctrl.Result{}, nil
	}

	existingCopy.Spec.Selector = s.Spec.Selector
	existingCopy.Spec.Type = s.Spec.Type
	err = r.Client.Update(context.TODO(), existingCopy)
	if err != nil {
		svlog.Error(err, "Failed to update Service")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureChannel(m *operatorv1.MultiClusterHub, u *unstructured.Unstructured) (ctrl.Result, error) {
	selog := r.Log.WithValues("Channel.Namespace", u.GetNamespace(), "Channel.Name", u.GetName())

	found := &unstructured.Unstructured{}
	found.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps.open-cluster-management.io",
		Kind:    "Channel",
		Version: "v1",
	})
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      u.GetName(),
		Namespace: m.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the Channel
		err = r.Client.Create(context.TODO(), u)
		if err != nil {
			// Creation failed
			selog.Error(err, "Failed to create new Channel")
			return ctrl.Result{}, err
		}

		// Creation was successful
		selog.Info("Created a new Channel")
		condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
		SetHubCondition(&m.Status, *condition)
		return ctrl.Result{}, nil

	} else if err != nil {
		// Error that isn't due to the Channel not existing
		selog.Error(err, "Failed to get Channel")
		return ctrl.Result{}, err
	}

	updated, needsUpdate := channel.Validate(m, found)
	if needsUpdate {
		selog.Info("Updating channel")
		err = r.Client.Update(context.TODO(), updated)
		if err != nil {
			// Update failed
			selog.Error(err, "Failed to update channel")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureSubscription(m *operatorv1.MultiClusterHub, u *unstructured.Unstructured) (ctrl.Result, error) {
	obLog := r.Log.WithValues("Namespace", u.GetNamespace(), "Name", u.GetName(), "Kind", u.GetKind())

	if utils.ProxyEnvVarsAreSet() {
		u = addProxyEnvVarsToSub(u)
	}

	found := &unstructured.Unstructured{}
	found.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apps.open-cluster-management.io",
		Kind:    "Subscription",
		Version: "v1",
	})
	// Try to get API group instance
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
	}, found)
	if err != nil && errors.IsNotFound(err) {

		err := r.Client.Create(context.TODO(), u)
		if err != nil {
			// Creation failed
			obLog.Error(err, "Failed to create new instance")
			return ctrl.Result{}, err
		}

		// Creation was successful
		obLog.Info("Created new object")
		condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
		SetHubCondition(&m.Status, *condition)
		return ctrl.Result{}, nil

	} else if err != nil {
		// Error that isn't due to the resource not existing
		obLog.Error(err, "Failed to get subscription")
		return ctrl.Result{}, err
	}

	// Validate object based on type
	updated, needsUpdate := subscription.Validate(found, u)
	if needsUpdate {
		obLog.Info("Updating subscription")
		// Update the resource. Skip on unit test
		err = r.Client.Update(context.TODO(), updated)
		if err != nil {
			// Update failed
			obLog.Error(err, "Failed to update object")
			return ctrl.Result{}, err
		}

		// Spec updated - return
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureNoSubscription(m *operatorv1.MultiClusterHub, u *unstructured.Unstructured) (ctrl.Result, error) {
	subLog := r.Log.WithValues("Namespace", u.GetNamespace(), "Name", u.GetName(), "Kind", u.GetKind())
	_, err := r.uninstall(m, u)
	if err != nil {
		subLog.Error(err, "Failed to uninstall subscription")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureNoDeployment(m *operatorv1.MultiClusterHub, dep *appsv1.Deployment) (ctrl.Result, error) {
	dplog := r.Log.WithValues("Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(dep)
	if err != nil {
		r.Log.Error(err, "Failed to unmarshal deployment")
		return ctrl.Result{}, err
	}
	u := &unstructured.Unstructured{Object: unstructuredMap}
	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Kind: "Deployment", Version: "v1"})

	_, err = r.uninstall(m, u)
	if err != nil {
		dplog.Error(err, "Failed to uninstall subscription")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureNoService(m *operatorv1.MultiClusterHub, s *corev1.Service) (ctrl.Result, error) {
	svlog := r.Log.WithValues("Service.Namespace", s.Namespace, "Service.Name", s.Name)

	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(s)
	if err != nil {
		r.Log.Error(err, "Failed to unmarshal deployment")
		return ctrl.Result{}, err
	}
	u := &unstructured.Unstructured{Object: unstructuredMap}
	u.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Kind: "Service", Version: "v1"})

	_, err = r.uninstall(m, u)
	if err != nil {
		svlog.Error(err, "Failed to uninstall subscription")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureNoUnstructured(m *operatorv1.MultiClusterHub, u *unstructured.Unstructured) (ctrl.Result, error) {
	subLog := r.Log.WithValues("Namespace", u.GetNamespace(), "Name", u.GetName(), "Kind", u.GetKind())
	_, err := r.uninstall(m, u)
	if err != nil {
		subLog.Error(err, "Failed to uninstall subscription")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureNoNamespace(m *operatorv1.MultiClusterHub, u *unstructured.Unstructured) (ctrl.Result, error) {
	subLog := r.Log.WithValues("Name", u.GetName(), "Kind", u.GetKind())
	gone, err := r.uninstall(m, u)
	if err != nil {
		subLog.Error(err, "Failed to uninstall namespace")
		return ctrl.Result{}, err
	}
	if gone == true {
		return ctrl.Result{}, nil
	} else {
		return ctrl.Result{RequeueAfter: resyncPeriod}, nil
	}
}

func (r *MultiClusterHubReconciler) ensureUnstructuredResource(m *operatorv1.MultiClusterHub, u *unstructured.Unstructured) (ctrl.Result, error) {
	obLog := r.Log.WithValues("Namespace", u.GetNamespace(), "Name", u.GetName(), "Kind", u.GetKind())

	found := &unstructured.Unstructured{}
	found.SetGroupVersionKind(u.GroupVersionKind())

	// Try to get API group instance
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Resource doesn't exist so create it
		err := r.Client.Create(context.TODO(), u)
		if err != nil {
			// Creation failed
			obLog.Error(err, "Failed to create new instance")
			return ctrl.Result{}, err
		}
		// Creation was successful
		obLog.Info("Created new resource")
		condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
		SetHubCondition(&m.Status, *condition)
		return ctrl.Result{}, nil

	} else if err != nil {
		// Error that isn't due to the resource not existing
		obLog.Error(err, "Failed to get resource")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureNamespace(m *operatorv1.MultiClusterHub, ns *corev1.Namespace) (ctrl.Result, error) {
	ctx := context.Background()

	r.Log.Info(fmt.Sprintf("Ensuring namespace: %s", ns.GetName()))

	existingNS := &corev1.Namespace{}
	err := r.Client.Get(ctx, types.NamespacedName{
		Name: ns.GetName(),
	}, existingNS)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(ctx, ns)
		if err != nil {
			r.Log.Info(fmt.Sprintf("Error creating namespace: %s", err.Error()))
			return ctrl.Result{Requeue: true}, nil
		}
	} else if err != nil {
		r.Log.Info(fmt.Sprintf("error locating namespace: %s. Error: %s", ns.GetName(), err.Error()))
		return ctrl.Result{Requeue: true}, nil
	}

	condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
	SetHubCondition(&m.Status, *condition)

	if existingNS.Status.Phase == corev1.NamespaceActive {
		return ctrl.Result{}, nil
	}
	r.Log.Info(fmt.Sprintf("namespace '%s' is not in an active state", ns.GetName()))
	return ctrl.Result{RequeueAfter: resyncPeriod}, nil
}

func (r *MultiClusterHubReconciler) ensureOperatorGroup(m *operatorv1.MultiClusterHub, og *olmv1.OperatorGroup) (ctrl.Result, error) {
	ctx := context.Background()

	operatorGroupList := &olmv1.OperatorGroupList{}
	err := r.Client.List(ctx, operatorGroupList, client.InNamespace(og.GetNamespace()))
	if err != nil {
		r.Log.Info(fmt.Sprintf("error listing operatorgroups in ns: %s. Error: %s", og.GetNamespace(), err.Error()))
		return ctrl.Result{Requeue: true}, nil
	}

	if len(operatorGroupList.Items) > 1 {
		r.Log.Error(fmt.Errorf("found more than one operator group in namespace %s", og.GetNamespace()), "fatal error")
		return ctrl.Result{RequeueAfter: resyncPeriod}, nil
	} else if len(operatorGroupList.Items) == 1 {
		return ctrl.Result{}, nil
	}

	r.Log.Info(fmt.Sprintf("Ensuring operator group exists in ns: %s", og.GetNamespace()))

	force := true
	err = r.Client.Patch(ctx, og, client.Apply, &client.PatchOptions{Force: &force, FieldManager: "multiclusterhub-operator"})
	if err != nil {
		r.Log.Info(fmt.Sprintf("Error: %s", err.Error()))
		return ctrl.Result{Requeue: true}, nil
	}
	condition := NewHubCondition(operatorv1.Progressing, metav1.ConditionTrue, NewComponentReason, "Created new resource")
	SetHubCondition(&m.Status, *condition)

	existingOperatorGroup := &olmv1.OperatorGroup{}
	err = r.Client.Get(ctx, types.NamespacedName{
		Name: og.GetName(),
	}, existingOperatorGroup)
	if err != nil {
		r.Log.Info(fmt.Sprintf("error locating operatorgroup: %s/%s. Error: %s", og.GetNamespace(), og.GetName(), err.Error()))
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil

}

func (r *MultiClusterHubReconciler) ensureMultiClusterEngineCR(ctx context.Context, m *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	mce, err := multiclusterengine.FindAndManageMCE(ctx, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	if mce == nil {
		// figure out if assisted service is already configured
		infraNS := ""
		configured, err := AssistedServiceConfigured(ctx, r.Client)
		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		if configured {
			ns, err := utils.FindNamespace()
			if err != nil {
				return ctrl.Result{Requeue: true}, err
			}
			infraNS = ns
		}

		mce = multiclusterengine.NewMultiClusterEngine(m, infraNS)
		err = r.Client.Create(ctx, mce)
		if err != nil {
			return ctrl.Result{Requeue: true}, fmt.Errorf("Error creating new MCE: %w", err)
		}
		return ctrl.Result{}, nil
	}

	calcMCE := multiclusterengine.RenderMultiClusterEngine(mce, m)
	err = r.Client.Update(ctx, calcMCE)
	if err != nil {
		return ctrl.Result{Requeue: true}, fmt.Errorf("Error updating MCE %s: %w", mce.Name, err)
	}
	return ctrl.Result{}, nil
}

// copies the imagepullsecret from mch to the newNS namespace
func (r *MultiClusterHubReconciler) ensurePullSecret(m *operatorv1.MultiClusterHub, newNS string) (ctrl.Result, error) {
	if m.Spec.ImagePullSecret == "" {
		// Delete imagepullsecret in MCE namespace if present
		secretList := &corev1.SecretList{}
		err := r.Client.List(
			context.TODO(),
			secretList,
			client.MatchingLabels{
				"installer.name":      m.GetName(),
				"installer.namespace": m.GetNamespace(),
			},
			client.InNamespace(newNS),
		)
		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		for i, secret := range secretList.Items {
			r.Log.Info("Deleting imagePullSecret", "Name", secret.Name, "Namespace", secret.Namespace)
			err = r.Client.Delete(context.TODO(), &secretList.Items[i])
			if err != nil {
				r.Log.Error(err, fmt.Sprintf("Error deleting imagepullsecret: %s", secret.GetName()))
				return ctrl.Result{Requeue: true}, err
			}
		}

		return ctrl.Result{}, nil
	}

	pullSecret := &corev1.Secret{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      m.Spec.ImagePullSecret,
		Namespace: m.Namespace,
	}, pullSecret)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	mceSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: pullSecret.APIVersion,
			Kind:       pullSecret.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pullSecret.Name,
			Namespace: newNS,
			Labels:    pullSecret.Labels,
		},
		Data: pullSecret.Data,
		Type: corev1.SecretTypeDockerConfigJson,
	}
	mceSecret.SetName(m.Spec.ImagePullSecret)
	mceSecret.SetNamespace(newNS)
	mceSecret.SetLabels(pullSecret.Labels)
	addInstallerLabelSecret(mceSecret, m.Name, m.Namespace)

	force := true
	err = r.Client.Patch(context.TODO(), mceSecret, client.Apply, &client.PatchOptions{Force: &force, FieldManager: "multiclusterhub-operator"})
	if err != nil {
		r.Log.Info(fmt.Sprintf("Error applying pullSecret to mce namespace: %s", err.Error()))
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// checks if imagepullsecret was created in mch namespace
func (r *MultiClusterHubReconciler) ensurePullSecretCreated(m *operatorv1.MultiClusterHub, namespace string) (ctrl.Result, error) {
	if m.Spec.ImagePullSecret == "" {
		//No imagepullsecret set, continuing
		return ctrl.Result{}, nil
	}

	pullSecret := &corev1.Secret{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      m.Spec.ImagePullSecret,
		Namespace: m.Namespace,
	}, pullSecret)

	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if pullSecret.Namespace == "" || pullSecret.Namespace != namespace {
		return ctrl.Result{Requeue: true}, fmt.Errorf("pullsecret doest not exist in namespace: %s", namespace)
	}

	return ctrl.Result{}, nil
}

// OverrideImagesFromConfigmap ...
func (r *MultiClusterHubReconciler) OverrideImagesFromConfigmap(imageOverrides map[string]string, namespace, configmapName string) (map[string]string, error) {
	r.Log.Info(fmt.Sprintf("Overriding images from configmap: %s/%s", namespace, configmapName))

	configmap := &corev1.ConfigMap{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      configmapName,
		Namespace: namespace,
	}, configmap)
	if err != nil && errors.IsNotFound(err) {
		return imageOverrides, err
	}

	if len(configmap.Data) != 1 {
		return imageOverrides, fmt.Errorf(fmt.Sprintf("Unexpected number of keys in configmap: %s", configmapName))
	}

	for _, v := range configmap.Data {

		var manifestImages []manifest.ManifestImage
		err = json.Unmarshal([]byte(v), &manifestImages)
		if err != nil {
			return nil, err
		}

		for _, manifestImage := range manifestImages {
			if manifestImage.ImageDigest != "" {
				imageOverrides[manifestImage.ImageKey] = fmt.Sprintf("%s/%s@%s", manifestImage.ImageRemote, manifestImage.ImageName, manifestImage.ImageDigest)
			} else if manifestImage.ImageTag != "" {
				imageOverrides[manifestImage.ImageKey] = fmt.Sprintf("%s/%s:%s", manifestImage.ImageRemote, manifestImage.ImageName, manifestImage.ImageTag)
			}

		}
	}

	return imageOverrides, nil
}

// Select oauth proxy image to use. If OCP 4.8 use old version. If OCP 4.9+ use new version. Set with key oauth_proxy
// before applying overrides.
func (r *MultiClusterHubReconciler) overrideOauthImage(ctx context.Context, imageOverrides map[string]string) (map[string]string, error) {
	ocpVersion, err := r.getClusterVersion(ctx)
	if err != nil {
		return nil, err
	}

	semverVersion, err := semver.NewVersion(ocpVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ocp version to semver compatible value: %w", err)
	}

	// -0 allows for prerelease builds to pass the validation.
	// If -0 is removed, developer/rc builds will not pass this check
	constraint, err := semver.NewConstraint(">= 4.9.0-0")
	if err != nil {
		return nil, fmt.Errorf("failed to set ocp version constraint: %w", err)
	}

	oauthKey := "oauth_proxy"
	oauthKeyOld := "oauth_proxy_48"
	oauthKeyNew := "oauth_proxy_49_and_up"

	if constraint.Check(semverVersion) {
		// use newer oauth image
		imageOverrides[oauthKey] = imageOverrides[oauthKeyNew]
	} else {
		// use old oauth image
		imageOverrides[oauthKey] = imageOverrides[oauthKeyOld]
	}

	return imageOverrides, nil
}

func (r *MultiClusterHubReconciler) maintainImageManifestConfigmap(mch *operatorv1.MultiClusterHub) error {
	// Define configmap
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("mch-image-manifest-%s", r.CacheSpec.ManifestVersion),
			Namespace: mch.Namespace,
		},
	}
	configmap.SetOwnerReferences([]metav1.OwnerReference{
		*metav1.NewControllerRef(mch, mch.GetObjectKind().GroupVersionKind()),
	})

	labels := make(map[string]string)
	labels["ocm-configmap-type"] = "image-manifest"
	labels["ocm-release-version"] = r.CacheSpec.ManifestVersion

	configmap.SetLabels(labels)

	// Get Configmap if it exists
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      configmap.Name,
		Namespace: configmap.Namespace,
	}, configmap)
	if err != nil && errors.IsNotFound(err) {
		// If configmap does not exist, create and return
		configmap.Data = r.CacheSpec.ImageOverrides
		err = r.Client.Create(context.TODO(), configmap)
		if err != nil {
			return err
		}
		return nil
	}

	// If cached image overrides are not equal to the configmap data, update configmap and return
	if !reflect.DeepEqual(configmap.Data, r.CacheSpec.ImageOverrides) {
		configmap.Data = r.CacheSpec.ImageOverrides
		err = r.Client.Update(context.TODO(), configmap)
		if err != nil {
			return err
		}
	}

	return nil
}

// listDeployments gets all deployments in the given namespaces
func (r *MultiClusterHubReconciler) listDeployments(namespaces []string) ([]*appsv1.Deployment, error) {
	var ret []*appsv1.Deployment

	for _, n := range namespaces {
		deployList := &appsv1.DeploymentList{}
		err := r.Client.List(context.TODO(), deployList, client.InNamespace(n))
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}

		for i := 0; i < len(deployList.Items); i++ {
			ret = append(ret, &deployList.Items[i])
		}
	}
	return ret, nil
}

// listHelmReleases gets all helmreleases in the given namespaces
func (r *MultiClusterHubReconciler) listHelmReleases(namespaces []string) ([]*subhelmv1.HelmRelease, error) {
	var ret []*subhelmv1.HelmRelease

	for _, n := range namespaces {
		hrList := &subhelmv1.HelmReleaseList{}
		err := r.Client.List(context.TODO(), hrList, client.InNamespace(n))
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}

		for i := 0; i < len(hrList.Items); i++ {
			ret = append(ret, &hrList.Items[i])
		}
	}

	return ret, nil
}

// listCustomResources gets custom resources the installer observes
func (r *MultiClusterHubReconciler) listCustomResources(m *operatorv1.MultiClusterHub) (map[string]*unstructured.Unstructured, error) {
	ret := make(map[string]*unstructured.Unstructured)

	var mceSub *unstructured.Unstructured
	gotSub, err := multiclusterengine.GetManagedMCESubscription(context.Background(), r.Client)
	if err != nil || gotSub == nil {
		mceSub = nil
	} else {
		unstructuredSub, err := runtime.DefaultUnstructuredConverter.ToUnstructured(gotSub)
		if err != nil {
			r.Log.Error(err, "Failed to unmarshal subscription")
		}
		mceSub = &unstructured.Unstructured{Object: unstructuredSub}
	}

	var mceCSV *unstructured.Unstructured
	if gotSub == nil {
		mceCSV = nil
	} else {
		mceCSV, err = r.GetCSVFromSubscription(gotSub)
		if err != nil {
			mceCSV = nil
		}
	}

	var mce *unstructured.Unstructured
	gotMCE, err := multiclusterengine.GetManagedMCE(context.Background(), r.Client)
	if err != nil || gotMCE == nil {
		mce = nil
	} else {
		unstructuredMCE, err := runtime.DefaultUnstructuredConverter.ToUnstructured(gotMCE)
		if err != nil {
			r.Log.Error(err, "Failed to unmarshal subscription")
		}
		mce = &unstructured.Unstructured{Object: unstructuredMCE}
	}

	ret["mce-sub"] = mceSub
	ret["mce-csv"] = mceCSV
	ret["mce"] = mce
	return ret, nil
}

// filterDeploymentsByRelease returns a subset of deployments with the release label value
func filterDeploymentsByRelease(deploys []*appsv1.Deployment, releaseLabel string) []*appsv1.Deployment {
	var labeledDeps []*appsv1.Deployment
	for i := range deploys {
		anno := deploys[i].GetAnnotations()
		if anno["meta.helm.sh/release-name"] == releaseLabel {
			labeledDeps = append(labeledDeps, deploys[i])
		}
	}
	return labeledDeps
}

// addInstallerLabel adds the installer name and namespace to a deployment's labels
// so it can be watched. Returns false if the labels are already present.
func addInstallerLabel(d *appsv1.Deployment, name string, ns string) bool {
	updated := false
	if d.Labels == nil {
		d.Labels = map[string]string{}
	}
	if d.Labels["installer.name"] != name {
		d.Labels["installer.name"] = name
		updated = true
	}
	if d.Labels["installer.namespace"] != ns {
		d.Labels["installer.namespace"] = ns
		updated = true
	}
	return updated
}

// addInstallerLabelSecret adds the installer name and namespace to a secret's labels
// so it can be watched. Returns false if the labels are already present.
func addInstallerLabelSecret(d *corev1.Secret, name string, ns string) bool {
	updated := false
	if d.Labels == nil {
		d.Labels = map[string]string{}
	}
	if d.Labels["installer.name"] != name {
		d.Labels["installer.name"] = name
		updated = true
	}
	if d.Labels["installer.namespace"] != ns {
		d.Labels["installer.namespace"] = ns
		updated = true
	}
	return updated
}

// getAppSubOwnedHelmReleases gets a subset of helmreleases created by the appsubs
func getAppSubOwnedHelmReleases(allHRs []*subhelmv1.HelmRelease, appsubs []types.NamespacedName) []*subhelmv1.HelmRelease {
	subMap := make(map[string]bool)
	for _, s := range appsubs {
		subMap[s.Name] = true
	}

	var ownedHRs []*subhelmv1.HelmRelease
	for _, hr := range allHRs {
		// check if this is one of our helmreleases
		owner := hr.OwnerReferences[0].Name
		if subMap[owner] {
			ownedHRs = append(ownedHRs, hr)

		}
	}
	return ownedHRs
}

// getHelmReleaseOwnedDeployments gets a subset of deployments created by the helmreleases
func getHelmReleaseOwnedDeployments(allDeps []*appsv1.Deployment, hrList []*subhelmv1.HelmRelease) []*appsv1.Deployment {
	var mchDeps []*appsv1.Deployment
	for _, hr := range hrList {
		hrDeployments := filterDeploymentsByRelease(allDeps, hr.Name)
		mchDeps = append(mchDeps, hrDeployments...)
	}
	return mchDeps
}

// labelDeployments updates deployments with installer labels if not already present
func (r *MultiClusterHubReconciler) labelDeployments(hub *operatorv1.MultiClusterHub, dList []*appsv1.Deployment) error {
	for _, d := range dList {
		// Attach installer labels so we can keep our eyes on the deployment
		if addInstallerLabel(d, hub.Name, hub.Namespace) {
			r.Log.Info("Adding installer labels to deployment", "Name", d.Name)
			err := r.Client.Update(context.TODO(), d)
			if err != nil {
				r.Log.Error(err, "Failed to update Deployment", "Name", d.Name)
				return err
			}
		}
	}
	return nil
}

// ensureSubscriptionOperatorIsRunning verifies that the subscription operator that manages helm subscriptions exists and
// is running. This validation is only intended to run during upgrade and when run as an OLM managed deployment
func (r *MultiClusterHubReconciler) ensureSubscriptionOperatorIsRunning(mch *operatorv1.MultiClusterHub, allDeps []*appsv1.Deployment) (ctrl.Result, error) {
	// skip check if not upgrading
	if mch.Status.CurrentVersion == version.Version {
		return ctrl.Result{}, nil
	}

	selfDeployment, exists := getDeploymentByName(allDeps, utils.MCHOperatorName)
	if !exists {
		// Deployment doesn't exist so this is either being run locally or with unit tests
		return ctrl.Result{}, nil
	}

	// skip check if not deployed by OLM
	if !isACMManaged(selfDeployment) {
		return ctrl.Result{}, nil
	}

	subscriptionDeploy, exists := getDeploymentByName(allDeps, utils.SubscriptionOperatorName)
	if !exists {
		err := fmt.Errorf("Standalone subscription deployment not found")
		return ctrl.Result{}, err
	}

	// Check that the owning CSV version of the deployments match
	if selfDeployment.GetLabels() == nil || subscriptionDeploy.GetLabels() == nil {
		r.Log.Info("Missing labels on either MCH operator or subscription operator deployment")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}
	if selfDeployment.Labels["olm.owner"] != subscriptionDeploy.Labels["olm.owner"] {
		r.Log.Info("OLM owner labels do not match. Requeuing.", "MCH operator label", selfDeployment.Labels["olm.owner"], "Subscription operator label", subscriptionDeploy.Labels["olm.owner"])
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	// Check that the standalone subscription deployment is available
	if successfulDeploy(subscriptionDeploy) {
		return ctrl.Result{}, nil
	} else {
		r.Log.Info("Standalone subscription deployment is not running. Requeuing.")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}
}

// GetSubscriptionIPInformation retrieves the current subscriptions installplan information for status
func (r *MultiClusterHubReconciler) GetSubscription(sub *subv1alpha1.Subscription) (*unstructured.Unstructured, error) {
	mceSubscription := &subv1alpha1.Subscription{}
	err := r.Client.Get(context.Background(), types.NamespacedName{
		Name:      sub.GetName(),
		Namespace: sub.GetNamespace(),
	}, mceSubscription)
	if err != nil {
		return nil, err
	}
	unstructuredSub, err := runtime.DefaultUnstructuredConverter.ToUnstructured(mceSubscription)
	if err != nil {
		r.Log.Error(err, "Failed to unmarshal subscription")
		return nil, err
	}
	return &unstructured.Unstructured{Object: unstructuredSub}, nil
}

// ensureMCESubscription verifies resources needed for MCE are created
func (r *MultiClusterHubReconciler) ensureMCESubscription(ctx context.Context, multiClusterHub *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	mceSub, err := multiclusterengine.FindAndManageMCESubscription(ctx, r.Client)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	// Get sub config, catalogsource, and annotation overrides
	subConfig, err := r.GetSubConfig()
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	overrides, err := multiclusterengine.GetAnnotationOverrides(multiClusterHub)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	ctlSrc := types.NamespacedName{}
	// Search for catalogsource if not defined in overrides
	if overrides == nil || overrides.CatalogSource == "" {
		ctlSrc, err = multiclusterengine.GetCatalogSource(r.UncachedClient)
		if err != nil {
			r.Log.Info("Failed to find a suitable catalogsource.", "error", err)
			return ctrl.Result{RequeueAfter: 5 * time.Second}, err
		}
	}

	if mceSub == nil {
		result, err := r.ensureNamespace(multiClusterHub, multiclusterengine.Namespace())
		if result != (ctrl.Result{}) {
			return result, err
		}
		result, err = r.ensurePullSecret(multiClusterHub, multiclusterengine.Namespace().Name)
		if result != (ctrl.Result{}) {
			return result, err
		}
		result, err = r.ensureOperatorGroup(multiClusterHub, multiclusterengine.OperatorGroup())
		if result != (ctrl.Result{}) {
			return result, err
		}
		// Sub is nil so create a new one
		mceSub = multiclusterengine.NewSubscription(multiClusterHub, subConfig, overrides, utils.IsCommunityMode())
	} else if multiclusterengine.CreatedByMCH(mceSub, multiClusterHub) {
		result, err := r.ensurePullSecret(multiClusterHub, multiclusterengine.Namespace().Name)
		if result != (ctrl.Result{}) {
			return result, err
		}
		result, err = r.ensureOperatorGroup(multiClusterHub, multiclusterengine.OperatorGroup())
		if result != (ctrl.Result{}) {
			return result, err
		}
	}

	// Apply MCE sub
	calcSub := multiclusterengine.RenderSubscription(mceSub, subConfig, overrides, ctlSrc, utils.IsCommunityMode())

	force := true
	err = r.Client.Patch(ctx, calcSub, client.Apply, &client.PatchOptions{Force: &force, FieldManager: "multiclusterhub-operator"})
	if err != nil {
		r.Log.Info(fmt.Sprintf("Error applying subscription: %s", err.Error()))
		return ctrl.Result{Requeue: true}, err
	}
	return ctrl.Result{}, nil
}

func (r *MultiClusterHubReconciler) ensureMultiClusterEngine(ctx context.Context, multiClusterHub *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	// confirm subscription and reqs exist and are configured correctly
	result, err := r.ensureMCESubscription(ctx, multiClusterHub)
	if result != (ctrl.Result{}) {
		return result, err
	}

	result, err = r.ensureMultiClusterEngineCR(ctx, multiClusterHub)
	if result != (ctrl.Result{}) {
		return result, err
	}

	return ctrl.Result{}, nil
}

// waitForMCE checks that MCE is in a running state and at the expected version.
func (r *MultiClusterHubReconciler) waitForMCEReady(ctx context.Context, mceVersion string) (ctrl.Result, error) {
	// Wait for MCE to be ready
	existingMCE, err := multiclusterengine.GetManagedMCE(ctx, r.Client)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if utils.IsUnitTest() {
		return ctrl.Result{}, nil
	}

	if existingMCE.Status.CurrentVersion == "" {
		r.Log.Info(fmt.Sprintf("Multiclusterengine: %s is not yet available", existingMCE.GetName()))
		return ctrl.Result{RequeueAfter: resyncPeriod}, nil
	}

	split := strings.Split(mceVersion, ".")
	xyVersion := fmt.Sprintf("%s.%s.x", split[0], split[1])
	minimum, err := semver.NewConstraint(xyVersion)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	cv := existingMCE.Status.CurrentVersion
	mceSemVer, err := semver.NewVersion(cv)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	if !minimum.Check(mceSemVer) {
		return ctrl.Result{RequeueAfter: resyncPeriod}, e.New(fmt.Sprintf("%s did not meet version requirement %s", cv, xyVersion))
	}

	return ctrl.Result{}, nil
}

// GetCSVFromSubscription retrieves CSV status information from the related subscription for status
func (r *MultiClusterHubReconciler) GetCSVFromSubscription(sub *subv1alpha1.Subscription) (*unstructured.Unstructured, error) {
	if sub == nil {
		return nil, fmt.Errorf("Cannot find CSV from nil Subscription")
	}
	mceSubscription := &subv1alpha1.Subscription{}
	err := r.Client.Get(context.Background(), types.NamespacedName{
		Name:      sub.GetName(),
		Namespace: sub.GetNamespace(),
	}, mceSubscription)
	if err != nil {
		return nil, err
	}

	currentCSV := mceSubscription.Status.CurrentCSV
	if currentCSV == "" {
		return nil, fmt.Errorf("currentCSV is empty")
	}

	mceCSV := &subv1alpha1.ClusterServiceVersion{}
	err = r.Client.Get(context.Background(), types.NamespacedName{
		Name:      currentCSV,
		Namespace: sub.GetNamespace(),
	}, mceCSV)
	if err != nil {
		return nil, err
	}
	csv, err := runtime.DefaultUnstructuredConverter.ToUnstructured(mceCSV)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: csv}, nil
}

// isACMManaged returns whether this application is managed by OLM via an ACM subscription
func isACMManaged(deploy *appsv1.Deployment) bool {
	labels := deploy.GetLabels()
	if labels == nil {
		return false
	}
	if owner, ok := labels["olm.owner"]; ok {
		if strings.Contains(owner, "advanced-cluster-management") {
			return true
		}
	}
	return false
}

// getDeploymentByName returns the deployment with the matching name found from a list
func getDeploymentByName(allDeps []*appsv1.Deployment, desiredDeploy string) (*appsv1.Deployment, bool) {
	for i := range allDeps {
		if allDeps[i].Name == desiredDeploy {
			return allDeps[i], true
		}
	}
	return nil, false
}

func addProxyEnvVarsToDeployment(dep *appsv1.Deployment) *appsv1.Deployment {
	proxyEnvVars := []corev1.EnvVar{
		{
			Name:  "HTTP_PROXY",
			Value: os.Getenv("HTTP_PROXY"),
		},
		{
			Name:  "HTTPS_PROXY",
			Value: os.Getenv("HTTPS_PROXY"),
		},
		{
			Name:  "NO_PROXY",
			Value: os.Getenv("NO_PROXY"),
		},
	}
	for i := 0; i < len(dep.Spec.Template.Spec.Containers); i++ {
		dep.Spec.Template.Spec.Containers[i].Env = append(dep.Spec.Template.Spec.Containers[i].Env, proxyEnvVars...)
	}

	return dep
}

func addProxyEnvVarsToSub(u *unstructured.Unstructured) *unstructured.Unstructured {
	path := "spec.packageOverrides[].packageOverrides[].value.hubconfig."

	sub, err := injectMapIntoUnstructured(u.Object, path, map[string]interface{}{
		"proxyConfigs": map[string]interface{}{
			"HTTP_PROXY":  os.Getenv("HTTP_PROXY"),
			"HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
			"NO_PROXY":    os.Getenv("NO_PROXY"),
		},
	})
	if err != nil {
		// log.Info(fmt.Sprintf("Error inject proxy environmental variables into '%s' subscription: %s", u.GetName(), err.Error()))
	}
	u.Object = sub

	return u
}

// injects map into an unstructured objects map based off a given path
func injectMapIntoUnstructured(u map[string]interface{}, path string, mapToInject map[string]interface{}) (map[string]interface{}, error) {

	currentKey := strings.Split(path, ".")[0]
	isArr := false
	if strings.HasSuffix(currentKey, "[]") {
		isArr = true
	}
	currentKey = strings.TrimSuffix(currentKey, "[]")
	if i := strings.Index(path, "."); i >= 0 {
		// Determine remaining path
		path = path[(i + 1):]
		// If Arr, loop through each element of array
		if isArr {
			nextMap, ok := u[currentKey].([]map[string]interface{})
			if !ok {
				return u, fmt.Errorf("Failed to find key: %s", currentKey)
			}
			for i := 0; i < len(nextMap); i++ {
				_, err := injectMapIntoUnstructured(nextMap[i], path, mapToInject)
				if err != nil {
					return u, err
				}
			}
			return u, nil
		} else {
			nextMap, ok := u[currentKey].(map[string]interface{})
			if !ok {
				return u, fmt.Errorf("Failed to find key: %s", currentKey)
			}
			_, err := injectMapIntoUnstructured(nextMap, path, mapToInject)
			if err != nil {
				return u, err
			}
			return u, nil
		}
	} else {
		for key, val := range mapToInject {
			u[key] = val
			break
		}
		return u, nil
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// mergeErrors combines errors into a single string
func mergeErrors(errs []error) string {
	errStrings := []string{}
	for _, e := range errs {
		errStrings = append(errStrings, e.Error())
	}
	return strings.Join(errStrings, " ; ")
}

// GetSubConfig returns a SubscriptionConfig based on proxy variables and the mch operator configuration
func (r *MultiClusterHubReconciler) GetSubConfig() (*subv1alpha1.SubscriptionConfig, error) {
	found := &appsv1.Deployment{}
	mchOperatorNS, err := utils.FindNamespace()
	if err != nil {
		return nil, err
	}

	err = r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      utils.MCHOperatorName,
		Namespace: mchOperatorNS,
	}, found)
	if err != nil {
		return nil, err
	}

	proxyEnv := []corev1.EnvVar{}
	if utils.ProxyEnvVarsAreSet() {
		proxyEnv = []corev1.EnvVar{
			{
				Name:  "HTTP_PROXY",
				Value: os.Getenv("HTTP_PROXY"),
			},
			{
				Name:  "HTTPS_PROXY",
				Value: os.Getenv("HTTPS_PROXY"),
			},
			{
				Name:  "NO_PROXY",
				Value: os.Getenv("NO_PROXY"),
			},
		}
	}
	return &subv1alpha1.SubscriptionConfig{
		NodeSelector: found.Spec.Template.Spec.NodeSelector,
		Tolerations:  found.Spec.Template.Spec.Tolerations,
		Env:          proxyEnv,
	}, nil
}

func (r *MultiClusterHubReconciler) pluginIsSupported(multiClusterHub *operatorv1.MultiClusterHub) bool {
	// -0 allows for prerelease builds to pass the validation.
	// If -0 is removed, developer/rc builds will not pass this check
	log := r.Log
	constraint, err := semver.NewConstraint(">= 4.10.0-0")
	if err != nil {
		log.Error(err, "Failed to set constraint of minimum supported version for plugins")
		return false
	}

	currentOCPVersion := ""
	if hubOCPVersion, ok := os.LookupEnv("ACM_HUB_OCP_VERSION"); !ok {
		log.Info("ACM_HUB_OCP_VERSION environment variable not set")
		return false
	} else {
		currentOCPVersion = hubOCPVersion
	}

	currentVersion, err := semver.NewVersion(currentOCPVersion)
	if err != nil {
		log.Error(err, "Failed to convert hubOCPVersion of cluster to semver compatible value for comparison")
		return false
	}

	if constraint.Check(currentVersion) {
		log.Info("Dynamic plugins are supported. Adding ACM plugin to console")
		return true
	} else {
		log.Info("Dynamic plugins not supported.")
		return false
	}
}

func (r *MultiClusterHubReconciler) addPluginToConsole(multiClusterHub *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log
	console := &consolev1.Console{}
	// If trying to check this resource from the CLI run - `oc get consoles.operator.openshift.io cluster`.
	// The default `console` is not the correct resource
	err := r.Client.Get(ctx, types.NamespacedName{Name: "cluster"}, console)
	if err != nil {
		log.Info("Failed to find console: cluster")
		return ctrl.Result{Requeue: true}, err
	}

	if console.Spec.Plugins == nil {
		console.Spec.Plugins = []string{}
	}

	// Add acm to the plugins list if it is not already there
	if !utils.Contains(console.Spec.Plugins, "acm") {
		log.Info("Ready to add plugin")
		console.Spec.Plugins = append(console.Spec.Plugins, "acm")
		err = r.Client.Update(ctx, console)
		if err != nil {
			log.Info("Failed to add acm consoleplugin to console")
			return ctrl.Result{Requeue: true}, err
		} else {
			log.Info("Added acm consoleplugin to console")
		}
	}

	return ctrl.Result{}, nil
}

// removePluginFromConsoleResource ...
func (r *MultiClusterHubReconciler) removePluginFromConsole(multiClusterHub *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log
	console := &consolev1.Console{}
	// If trying to check this resource from the CLI run - `oc get consoles.operator.openshift.io cluster`.
	// The default `console` is not the correct resource
	err := r.Client.Get(ctx, types.NamespacedName{Name: "cluster"}, console)
	if err != nil {
		log.Info("Failed to find console: cluster")
		return ctrl.Result{Requeue: true}, err
	}

	// If No plugins, it is already removed
	if console.Spec.Plugins == nil {
		return ctrl.Result{}, nil
	}

	// Remove mce to the plugins list if it is not already there
	if utils.Contains(console.Spec.Plugins, "acm") {
		console.Spec.Plugins = utils.RemoveString(console.Spec.Plugins, "acm")
		err = r.Client.Update(ctx, console)
		if err != nil {
			log.Info("Failed to remove acm consoleplugin to console")
			return ctrl.Result{Requeue: true}, err
		} else {
			log.Info("Removed acm consoleplugin to console")
		}
	}

	return ctrl.Result{}, nil
}

// AssistedServiceConfigured returns true if assisted service has already been installed
// and configured in the hub namespace
func AssistedServiceConfigured(ctx context.Context, client client.Client) (bool, error) {
	agentServiceCRD := &apixv1.CustomResourceDefinition{}
	err := client.Get(ctx, types.NamespacedName{Name: "agentserviceconfigs.agent-install.openshift.io"}, agentServiceCRD)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	// CRD exists, check for instance
	list := &unstructured.UnstructuredList{}
	list.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "agent-install.openshift.io",
		Version: "v1beta1",
		Kind:    "AgentServiceConfigList",
	})
	if err := client.List(ctx, list); err != nil {
		return false, fmt.Errorf("unable to list AgentServiceConfigs: %s", err)
	}
	if len(list.Items) > 0 {
		return true, nil
	}
	return false, nil
}

// return current OCP version from clusterversion resource
// equivalent to `oc get clusterversion version -o=jsonpath='{.status.history[0].version}'`
func (r *MultiClusterHubReconciler) getClusterVersion(ctx context.Context) (string, error) {
	if utils.IsUnitTest() {
		// If unit test pass along a version, Can't set status in unit test
		return "4.9.0", nil
	}

	clusterVersion := &configv1.ClusterVersion{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "version"}, clusterVersion)
	if err != nil {
		return "", err
	}

	if len(clusterVersion.Status.History) == 0 {
		return "", e.New("Failed to detect status in clusterversion.status.history")
	}
	return clusterVersion.Status.History[0].Version, nil
}

func (r *MultiClusterHubReconciler) ensureSearchCR(m *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	ctx := context.Background()

	// If assisted installer is set up MCE needs to override the infrastructure
	// operator namespace
	searchList := &searchv2v1alpha1.SearchList{}
	err := r.Client.List(ctx, searchList, client.InNamespace(m.GetNamespace()))
	if err != nil {
		r.Log.Info(fmt.Sprintf("error locating Search CR. Error: %s", err.Error()))
		return ctrl.Result{Requeue: true}, err
	}

	if len(searchList.Items) == 0 {
		searchCR := &searchv2v1alpha1.Search{
			TypeMeta: metav1.TypeMeta{
				APIVersion: searchv2v1alpha1.GroupVersion.String(),
				Kind:       "Search",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "search-v2-operator",
				Namespace: m.Namespace,
			},
		}
		force := true
		err = r.Client.Patch(ctx, searchCR, client.Apply, &client.PatchOptions{Force: &force, FieldManager: "multiclusterhub-operator"})
		if err != nil {
			r.Log.Info(fmt.Sprintf("error creating Search CR. Error: %s", err.Error()))
			return ctrl.Result{Requeue: true}, err
		}
	}
	return ctrl.Result{}, nil

}

func (r *MultiClusterHubReconciler) ensureNoSearchCR(m *operatorv1.MultiClusterHub) (ctrl.Result, error) {
	ctx := context.Background()

	// If assisted installer is set up MCE needs to override the infrastructure
	// operator namespace
	searchList := &searchv2v1alpha1.SearchList{}
	err := r.Client.List(ctx, searchList, client.InNamespace(m.GetNamespace()))
	if err != nil {
		r.Log.Info(fmt.Sprintf("error locating Search CR. Error: %s", err.Error()))
		return ctrl.Result{Requeue: true}, err
	}

	if len(searchList.Items) != 0 {
		err = r.Client.Delete(context.TODO(), &searchList.Items[0])
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("Error deleting Search CR"))
			return ctrl.Result{Requeue: true}, err
		}

	}
	err = r.Client.List(ctx, searchList, client.InNamespace(m.GetNamespace()))
	if err != nil {
		r.Log.Info(fmt.Sprintf("error locating Search CR. Error: %s", err.Error()))
		return ctrl.Result{Requeue: true}, err
	}
	if len(searchList.Items) != 0 {
		r.Log.Info(fmt.Sprintf("Waiting for Search CR to be deleted"))
		return ctrl.Result{Requeue: true}, errors.NewBadRequest("Search CR has not been deleted")
	}
	return ctrl.Result{}, nil

}

// Checks if OCP Console is enabled and return true if so. If <OCP v4.12, always return true
// Otherwise check in the EnabledCapabilities spec for OCP console
func (r *MultiClusterHubReconciler) CheckConsole(ctx context.Context) (bool, error) {
	versionStatus := &configv1.ClusterVersion{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: "version"}, versionStatus)
	if err != nil {
		return false, err
	}
	ocpVersion, err := r.getClusterVersion(ctx)
	if err != nil {
		return false, err
	}
	if hubOCPVersion, ok := os.LookupEnv("ACM_HUB_OCP_VERSION"); ok {
		ocpVersion = hubOCPVersion
	}
	semverVersion, err := semver.NewVersion(ocpVersion)
	if err != nil {
		return false, fmt.Errorf("failed to convert ocp version to semver compatible value: %w", err)
	}
	// -0 allows for prerelease builds to pass the validation.
	// If -0 is removed, developer/rc builds will not pass this check
	//OCP Console can only be disabled in OCP 4.12+
	constraint, err := semver.NewConstraint(">= 4.12.0-0")
	if err != nil {
		return false, fmt.Errorf("failed to set ocp version constraint: %w", err)
	}
	if !constraint.Check(semverVersion) {
		return true, nil
	}
	for _, v := range versionStatus.Status.Capabilities.EnabledCapabilities {
		if v == "Console" {
			return true, nil
		}
	}
	return false, nil
}
