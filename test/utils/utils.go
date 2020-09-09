// Copyright (c) 2020 Red Hat, Inc.
package utils

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (

	// KubeClient ...
	KubeClient = NewKubeClient("", "", "")
	// DynamicKubeClient ...
	DynamicKubeClient = NewKubeClientDynamic("", "", "")

	// ImageOverridesCMBadImageName...
	ImageOverridesCMBadImageName = "bad-image-ref"

	// GVRMultiClusterHub ...
	GVRMultiClusterHub = schema.GroupVersionResource{
		Group:    "operator.open-cluster-management.io",
		Version:  "v1",
		Resource: "multiclusterhubs",
	}
	// GVRAppSub ...
	GVRAppSub = schema.GroupVersionResource{
		Group:    "apps.open-cluster-management.io",
		Version:  "v1",
		Resource: "subscriptions",
	}
	// GVRSub ...
	GVRSub = schema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}
	// GVROperatorGroup ...
	GVROperatorGroup = schema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1",
		Resource: "operatorgroups",
	}
	// GVRCSV ...
	GVRCSV = schema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "clusterserviceversions",
	}
	// GVRHelmRelease ...
	GVRHelmRelease = schema.GroupVersionResource{
		Group:    "apps.open-cluster-management.io",
		Version:  "v1",
		Resource: "helmreleases",
	}

	// GVRInstallPlan ...
	GVRInstallPlan = schema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "installplans",
	}

	// GVRDeployment ...
	GVRDeployment = schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	// GVRManagedCluster
	GVRManagedCluster = schema.GroupVersionResource{
		Group: "cluster.open-cluster-management.io",
		Version: "v1",
		Resource: "managedclusters",
	}

	// GVRKlusterletAddonConfig
	GVRKlusterletAddonConfig = schema.GroupVersionResource{
		Group: "agent.open-cluster-management.io",
		Version: "v1",
		Resource: "klusterletaddonconfigs",
	}
	// DefaultImageRegistry ...
	DefaultImageRegistry = "quay.io/open-cluster-management"
	// DefaultImagePullSecretName ...
	DefaultImagePullSecretName = "multiclusterhub-operator-pull-secret"

	// MCHName ...
	MCHName = "multiclusterhub"
	// MCHNamespace ...
	MCHNamespace = "open-cluster-management"
	// MCHPullSecretName ...
	MCHPullSecretName = os.Getenv("pullSecret")

	// MCHRepoName ...
	MCHRepoName = "multiclusterhub-repo"
	// MCHOperatorName ...
	MCHOperatorName = "multiclusterhub-operator"

	// OCMSubscriptionName ...
	OCMSubscriptionName = os.Getenv("name")

	// SubList contains the list of subscriptions to delete
	SubList = [...]string{
		OCMSubscriptionName,
		"hive-operator-alpha-community-operators-openshift-marketplace",
		"multicluster-operators-subscription-alpha-community-operators-openshift-marketplace",
	}

	// AppSubSlice ...
	AppSubSlice = [...]string{"application-chart-sub", "cert-manager-sub",
		"cert-manager-webhook-sub", "configmap-watcher-sub", "console-chart-sub",
		"grc-sub", "kui-web-terminal-sub", "management-ingress-sub",
		"rcm-sub", "search-prod-sub", "topology-sub"}

	// CSVName ...
	CSVName = "advanced-cluster-management"

	// WaitInMinutesDefault ...
	WaitInMinutesDefault = 20

)

// GetWaitInMinutes...
func GetWaitInMinutes() int {
	waitInMinutesAsString := os.Getenv("waitInMinutes")
	if waitInMinutesAsString == "" {
		return WaitInMinutesDefault
	}
	waitInMinutesAsInt, err := strconv.Atoi(waitInMinutesAsString)
	if err != nil {
		return WaitInMinutesDefault
	}
	return waitInMinutesAsInt
}

// CreateNewUnstructured creates resources by using gvr & obj, will get object after create.
func CreateNewUnstructured(
	clientHubDynamic dynamic.Interface,
	gvr schema.GroupVersionResource,
	obj *unstructured.Unstructured,
	name, namespace string,
) {
	ns := clientHubDynamic.Resource(gvr).Namespace(namespace)
	Expect(ns.Create(context.TODO(), obj, metav1.CreateOptions{})).NotTo(BeNil())
	Expect(ns.Get(context.TODO(), name, metav1.GetOptions{})).NotTo(BeNil())
}

// CreateNewConfigMap ...
func CreateNewConfigMap(cm *corev1.ConfigMap, namespace string) error {
	_, err := KubeClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

// DeleteConfigMapIfExists ...
func DeleteConfigMapIfExists(cmName, namespace string) error {
	_, err := KubeClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmName, metav1.GetOptions{})
	if err == nil {
		return KubeClient.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), cmName, metav1.DeleteOptions{})
	}
	return nil
}

// DeleteIfExists deletes resources by using gvr, name, and namespace.
// Will wait for deletion to complete by using eventually
func DeleteIfExists(clientHubDynamic dynamic.Interface, gvr schema.GroupVersionResource, name, namespace string, wait bool) {
	ns := clientHubDynamic.Resource(gvr).Namespace(namespace)
	if _, err := ns.Get(context.TODO(), name, metav1.GetOptions{}); err != nil {
		Expect(errors.IsNotFound(err)).To(Equal(true))
		return
	}
	Expect(func() error {
		// possibly already got deleted
		err := ns.Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		return nil
	}()).To(BeNil())

	By("Wait for deletion")
	Eventually(func() error {
		var err error
		_, err = ns.Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		if err == nil {
			if wait {
				return fmt.Errorf("found object %s in namespace %s after deletion", name, namespace)
			}
			return nil
		}
		return nil
	}, GetWaitInMinutes()*60, 1).Should(BeNil())
}

// NewKubeClient returns a kube client
func NewKubeClient(url, kubeconfig, context string) kubernetes.Interface {
	klog.V(5).Infof("Create kubeclient for url %s using kubeconfig path %s\n", url, kubeconfig)
	config, err := LoadConfig(url, kubeconfig, context)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

// NewKubeClientDynamic returns a dynamic kube client
func NewKubeClientDynamic(url, kubeconfig, context string) dynamic.Interface {
	klog.V(5).Infof(
		"Create kubeclient dynamic for url %s using kubeconfig path %s\n",
		url,
		kubeconfig,
	)
	config, err := LoadConfig(url, kubeconfig, context)
	if err != nil {
		panic(err)
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

// LoadConfig loads kubeconfig
func LoadConfig(url, kubeconfig, context string) (*rest.Config, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	klog.V(5).Infof("Kubeconfig path %s\n", kubeconfig)
	// If we have an explicit indication of where the kubernetes config lives, read that.
	if kubeconfig != "" {
		if context == "" {
			return clientcmd.BuildConfigFromFlags(url, kubeconfig)
		}
		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{
				CurrentContext: context,
			}).ClientConfig()
	}
	// If not, try the in-cluster config.
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}
	// If no in-cluster config, try the default location in the user's home directory.
	if usr, err := user.Current(); err == nil {
		klog.V(5).Infof("clientcmd.BuildConfigFromFlags for url %s using %s\n", url, filepath.Join(usr.HomeDir, ".kube", "config"))
		if c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not create a valid kubeconfig")
}

// IsOwner checks if obj is owned by owner, obj can either be unstructured or ObjectMeta
func IsOwner(owner *unstructured.Unstructured, obj interface{}) bool {
	if obj == nil {
		return false
	}
	var owners []metav1.OwnerReference
	objMeta, ok := obj.(*metav1.ObjectMeta)
	if ok {
		owners = objMeta.GetOwnerReferences()
	} else {
		if objUnstructured, ok := obj.(*unstructured.Unstructured); ok {
			owners = objUnstructured.GetOwnerReferences()
		} else {
			klog.Error("Failed to get owners")
			return false
		}
	}

	for _, ownerRef := range owners {
		if _, ok := owner.Object["metadata"]; !ok {
			klog.Error("no meta")
			continue
		}
		meta, ok := owner.Object["metadata"].(map[string]interface{})
		if !ok || meta == nil {
			klog.Error("no meta map")
			continue
		}
		name, ok := meta["name"].(string)
		if !ok || name == "" {
			klog.Error("failed to get name")
			continue
		}
		if ownerRef.Kind == owner.Object["kind"] && ownerRef.Name == name {
			return true
		}
	}
	return false
}

// CreateDefaultMCH ...
func CreateDefaultMCH() *unstructured.Unstructured {
	mch := NewMultiClusterHub(MCHName, MCHNamespace, "")
	CreateNewUnstructured(DynamicKubeClient, GVRMultiClusterHub, mch, MCHName, MCHNamespace)
	return mch
}

// CreateMCHImageOverridesAnnotation ...
func CreateMCHImageOverridesAnnotation(imageOverridesConfigmapName string) *unstructured.Unstructured {
	mch := NewMultiClusterHub(MCHName, MCHNamespace, imageOverridesConfigmapName)
	CreateNewUnstructured(DynamicKubeClient, GVRMultiClusterHub, mch, MCHName, MCHNamespace)
	return mch
}

// BrickMCHRepo modifies the multiclusterhub-repo deployment so it becomes unhealthy
func BrickMCHRepo() error {
	deploy, err := KubeClient.AppsV1().Deployments(MCHNamespace).Get(context.TODO(), MCHRepoName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Add non-existent nodeSelector so the pod isn't scheduled
	deploy.Spec.Template.Spec.NodeSelector = map[string]string{"schedule": "never"}
	deploy.Spec.Strategy = appsv1.DeploymentStrategy{Type: appsv1.RecreateDeploymentStrategyType}

	_, err = KubeClient.AppsV1().Deployments(MCHNamespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	return err
}

// FixMCHRepo deletes the multiclusterhub-repo deployment so it can be recreated by the installer
func FixMCHRepo() error {
	return KubeClient.AppsV1().Deployments(MCHNamespace).Delete(context.TODO(), MCHRepoName, metav1.DeleteOptions{})
}

// GetMCHStatus gets the mch object and parses its status
func GetMCHStatus() (map[string]interface{}, error) {
	mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	status, ok := mch.Object["status"].(map[string]interface{})
	if !ok || status == nil {
		return nil, fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
	}
	return status, nil
}

// findPhase reports whether the hub status has the desired phase
func findPhase(status map[string]interface{}, wantPhase string) error {
	if _, ok := status["phase"]; !ok {
		return fmt.Errorf("MCH status has no 'phase' field")
	}
	if phase := status["phase"]; phase != wantPhase {
		return fmt.Errorf("MCH phase equals %s, expected %s", phase, wantPhase)
	}
	return nil
}

// waitForRepoUnavailable waits for the multiclusterhub-repo to go unready, with timeout
func waitForUnavailableRepo(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		deploy, err := KubeClient.AppsV1().Deployments(MCHNamespace).Get(context.TODO(), MCHRepoName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if deploy.Status.ReadyReplicas == 0 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("Repo failed to become unready after %s", timeout)
}

// ValidateMCHDegraded validates the install operator responds appropriately when the install components
// go into a degraded state after a successful install
func ValidateMCHDegraded() error {
	By("Validating MultiClusterHub Degraded")
	By("- Wait for degraded deployment")
	if err := waitForUnavailableRepo(time.Duration(45) * time.Second); err != nil {
		return err
	}

	status, err := GetMCHStatus()
	if err != nil {
		return err
	}

	By("- Ensuring MCH is in 'pending' phase")
	if err := findPhase(status, "Pending"); err != nil {
		return err
	}

	By("- Ensuring hub condition shows installation as incomplete")
	if err := FindCondition(status, "Complete", "False"); err != nil {
		return err
	}

	return nil
}

// ValidateDelete ...
func ValidateDelete(clientHubDynamic dynamic.Interface) error {
	By("Validating MCH has been successfully uninstalled.")

	labelSelector := fmt.Sprintf("installer.name=%s, installer.namespace=%s", MCHName, MCHNamespace)
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         100,
	}

	appSubLink := clientHubDynamic.Resource(GVRAppSub)
	appSubs, err := appSubLink.List(context.TODO(), listOptions)
	Expect(err).Should(BeNil())

	helmReleaseLink := clientHubDynamic.Resource(GVRHelmRelease)
	helmReleases, err := helmReleaseLink.List(context.TODO(), listOptions)
	Expect(err).Should(BeNil())

	By("- Ensuring Application Subscriptions have terminated")
	if len(appSubs.Items) != 0 {
		return fmt.Errorf("%d appsubs left to be uninstalled", len(appSubs.Items))
	}

	By("- Ensuring HelmReleases have terminated")
	if len(helmReleases.Items) != 0 {
		By(fmt.Sprintf("%d helmreleases left to be uninstalled", len(helmReleases.Items)))
		return fmt.Errorf("%d helmreleases left to be uninstalled", len(helmReleases.Items))
	}

	By("- Ensuring MCH Repo deployment has been terminated")
	deploymentLink := clientHubDynamic.Resource(GVRDeployment).Namespace(MCHNamespace)
	_, err = deploymentLink.Get(context.TODO(), "multiclusterhub-repo", metav1.GetOptions{})
	Expect(err).ShouldNot(BeNil())

	By("- Ensuring MCH image manifest configmap is terminated")
	labelSelector = fmt.Sprintf("ocm-configmap-type=%s", "image-manifest")
	listOptions = metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         100,
	}

	Eventually(func() error {
		configmaps, err := KubeClient.CoreV1().ConfigMaps(MCHNamespace).List(context.TODO(), listOptions)
		Expect(err).Should(BeNil())
		if len(configmaps.Items) != 0 {
			return fmt.Errorf("Expecting configmaps to terminate")
		}
		return nil
	}, GetWaitInMinutes()*60, 1).Should(BeNil())

	return nil
}

// FindCondition reports whether a hub condition of type 't' exists and matches the status 's'
func FindCondition(status map[string]interface{}, t string, s string) error {
	conditions, ok := status["conditions"].([]interface{})
	if !ok || conditions == nil {
		return fmt.Errorf("no hubConditions found")
	}
	for i := range conditions {
		condition := conditions[i]
		if condition.(map[string]interface{})["type"].(string) == t {
			if got := condition.(map[string]interface{})["status"].(string); got == s {
				return nil
			} else {
				return fmt.Errorf("hubCondition `%s` status equals '%s', expected '%s'", t, got, s)
			}
		}
	}
	return fmt.Errorf("MCH does not have a hubcondition with type '%s'", t)
}

// ValidateMCHUnsuccessful ...
func ValidateMCHUnsuccessful() error {
	By("Validating MultiClusterHub Unsuccessful")
	By(fmt.Sprintf("- Waiting %d minutes", GetWaitInMinutes()), func() {
		time.Sleep(time.Duration(GetWaitInMinutes()) * time.Minute)
	})

	By("- Ensuring MCH is in 'pending' phase")
	When("MCH Status should be `Pending`", func() {
		Eventually(func() error {
			mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
			Expect(err).To(BeNil())
			status, ok := mch.Object["status"].(map[string]interface{})
			if !ok || status == nil {
				return fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
			}
			if _, ok := status["phase"]; !ok {
				return fmt.Errorf("MultiClusterHub: %s status has no 'phase' field", mch.GetName())
			}
			if status["phase"] != "Pending" {
				return fmt.Errorf("MultiClusterHub: %s with phase %s is not in pending phase", mch.GetName(), status["phase"])
			}
			return nil
		}, 1, 1).Should(BeNil())
	})

	When("MCH Condition 'type' should be `Progressing` and 'status' should be 'true", func() {
		Eventually(func() error {
			mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
			Expect(err).To(BeNil())
			status := mch.Object["status"].(map[string]interface{})
			return FindCondition(status, "Progressing", "True")
		}, 1, 1).Should(BeNil())
	})

	return nil
}

// ValidateMCH validates MCH CR is running successfully
func ValidateMCH() error {
	By("Validating MultiClusterHub")

	By(fmt.Sprintf("- Ensuring MCH is in 'running' phase within %d minutes", GetWaitInMinutes()))
	When(fmt.Sprintf("Wait for MultiClusterHub to be in running phase (Will take up to %d minutes)", GetWaitInMinutes()), func() {
		Eventually(func() error {
			mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
			Expect(err).To(BeNil())
			status, ok := mch.Object["status"].(map[string]interface{})
			if !ok || status == nil {
				return fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
			}
			if _, ok := status["phase"]; !ok {
				return fmt.Errorf("MultiClusterHub: %s status has no 'phase' field", mch.GetName())
			}
			if status["phase"] != "Running" {
				return fmt.Errorf("MultiClusterHub: %s with phase %s is not in running phase", mch.GetName(), status["phase"])
			}
			return nil
		}, GetWaitInMinutes()*60, 1).Should(BeNil())
	})

	By("- Ensuring MCH Repo Is available")
	var deploy *appsv1.Deployment
	deploy, err := KubeClient.AppsV1().Deployments(MCHNamespace).Get(context.TODO(), MCHRepoName, metav1.GetOptions{})
	Expect(err).Should(BeNil())
	mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
	Expect(err).To(BeNil())
	Expect(deploy.Status.AvailableReplicas).ShouldNot(Equal(0))
	Expect(IsOwner(mch, &deploy.ObjectMeta)).To(Equal(true))

	By("- Ensuring components have status 'true' when MCH is in 'running' phase")
	mch, err = DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
	Expect(err).To(BeNil())
	status := mch.Object["status"].(map[string]interface{})
	if status["phase"] == "Running" {
		components, ok := mch.Object["status"].(map[string]interface{})["components"]
		if !ok || components == nil {
			return fmt.Errorf("MultiClusterHub: %s has no 'Components' map despite reporting 'running'", mch.GetName())
		}
		for k, v := range components.(map[string]interface{}) {
			compStatus := v.(map[string]interface{})["status"].(string)
			if compStatus != "True" {
				return fmt.Errorf("Component: %s does not have status of 'true'", k)
			}
		}
	}

	By("- Ensuring condition has status 'true' and type 'complete' when MCH is in 'running' phase")
	When("Component statuses should be true", func() {
		Eventually(func() error {
			mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
			Expect(err).To(BeNil())
			status := mch.Object["status"].(map[string]interface{})
			return FindCondition(status, "Complete", "True")
		}, 1, 1).Should(BeNil())
	})

	By("- Checking Appsubs")
	unstructuredAppSubs := listByGVR(DynamicKubeClient, GVRAppSub, MCHNamespace, 1, len(AppSubSlice))
	for _, appsub := range unstructuredAppSubs.Items {
		if _, ok := appsub.Object["status"]; !ok {
			return fmt.Errorf("Appsub: %s has no 'status' field", appsub.GetName())
		}
		status, ok := appsub.Object["status"].(map[string]interface{})
		if !ok || status == nil {
			return fmt.Errorf("Appsub: %s has no 'status' map", appsub.GetName())
		}
		klog.V(5).Infof("Checking Appsub - %s", appsub.GetName())
		Expect(status["message"]).To(Equal("Active"))
		Expect(status["phase"]).To(Equal("Subscribed"))
	}

	By("- Checking HelmReleases")
	unstructuredHelmReleases := listByGVR(DynamicKubeClient, GVRHelmRelease, MCHNamespace, 1, len(AppSubSlice))
	for _, helmRelease := range unstructuredHelmReleases.Items {
		klog.V(5).Infof("Checking HelmRelease - %s", helmRelease.GetName())

		status, ok := helmRelease.Object["status"].(map[string]interface{})
		if !ok || status == nil {
			return fmt.Errorf("HelmRelease: %s has no 'status' map", helmRelease.GetName())
		}

		conditions, ok := status["deployedRelease"].(map[string]interface{})
		if !ok || conditions == nil {
			return fmt.Errorf("HelmRelease: %s has no 'deployedRelease' interface", helmRelease.GetName())
		}
	}

	By("- Checking Imported Hub Cluster")
	err = ValidateManagedCluster(true)
	Expect(err).Should(BeNil())


	currentVersion, err := GetCurrentVersionFromMCH()
	Expect(err).Should(BeNil())
	v, err := semver.NewVersion(currentVersion)
	Expect(err).Should(BeNil())
	c, err := semver.NewConstraint(">= 2.1.0")
	Expect(err).Should(BeNil())
	if c.Check(v) {
		By("- Ensuring image manifest configmap is created")
		_, err = KubeClient.CoreV1().ConfigMaps(MCHNamespace).Get(context.TODO(), fmt.Sprintf("mch-image-manifest-%s", currentVersion), metav1.GetOptions{})
		Expect(err).Should(BeNil())
	}
	return nil
}

// ValidateMCHStatusExist check if mch status exists
func ValidateMCHStatusExist() error {
	Eventually(func() error {
		mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
		Expect(err).To(BeNil())
		status, ok := mch.Object["status"].(map[string]interface{})
		if !ok || status == nil {
			return fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
		}
		return nil
	}, GetWaitInMinutes()*60, 1).Should(BeNil())
	return nil
}

// ValidateComponentStatusExist check if Component statuses exist immediately when MCH is created
func ValidateComponentStatusExist() error {
	Eventually(func() error {
		mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
		Expect(err).To(BeNil())
		status, ok := mch.Object["status"].(map[string]interface{})
		if !ok || status == nil {
			return fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
		}
		if components, ok := status["components"]; !ok || components == nil {
			return fmt.Errorf("MultiClusterHub: %s has no 'Components' map in status", mch.GetName())
		} else {
			for k, v := range components.(map[string]interface{}) {
				if _, ok := v.(map[string]interface{})["status"].(string); !ok {
					return fmt.Errorf("Component: %s status does not exist", k)
				}
			}
		}
		return nil
	}, GetWaitInMinutes()*60, 1).Should(BeNil())
	return nil
}

// ValidateHubStatusExist checks if hub statuses exist immediately when MCH is created
func ValidateHubStatusExist() error {
	Eventually(func() error {
		mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
		Expect(err).To(BeNil())
		status, ok := mch.Object["status"].(map[string]interface{})
		if !ok || status == nil {
			return fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
		}
		return FindCondition(status, "Progressing", "True")
	}, GetWaitInMinutes()*60, 1).Should(BeNil())
	return nil
}

//ValidateConditionDuringUninstall check if condition is terminating during uninstall of MCH
func ValidateConditionDuringUninstall() error {
	By("- Checking HubCondition type")
	Eventually(func() error {
		mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
		Expect(err).To(BeNil())
		status := mch.Object["status"].(map[string]interface{})
		return FindCondition(status, "Terminating", "True")
	}, GetWaitInMinutes()*60, 1).Should(BeNil())
	return nil
}

// ValidateStatusesExist Confirms existence of both overall MCH and Component statuses immediately after MCH creation
func ValidateStatusesExist() error {
	By("Validating Statuses exist")

	By("- Ensuring MCH Status exists")
	if err := ValidateMCHStatusExist(); err != nil {
		return err
	}
	By("- Ensuring Component Status exist")
	if err := ValidateComponentStatusExist(); err != nil {
		return err
	}
	By("- Ensuring Hub Status exist")
	if err := ValidateHubStatusExist(); err != nil {
		return err
	}
	return nil
}

//ValidateImportHubResources confirms the existence of 3 resources that are created when importing hub as managed cluster
func ValidateImportHubResourcesExist(expected bool) error {
    //check created namespace exists
	_, nsErr := KubeClient.CoreV1().Namespaces().Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	//check created ManagedCluster exists
	mc, mcErr := DynamicKubeClient.Resource(GVRManagedCluster).Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	//check created KlusterletAddonConfig
	kac, kacErr := DynamicKubeClient.Resource(GVRKlusterletAddonConfig).Namespace("local-cluster").Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	if (expected) {
		if (mc != nil) {
			if (nsErr != nil || mcErr != nil || kacErr !=nil) {
				return fmt.Errorf("not all local-cluster resources created")
			}
			return nil
		} else {
			return fmt.Errorf("local-cluster resources exist")
		}
	} else {
		if (mc != nil || kac != nil) {
			return fmt.Errorf("local-cluster resources exist")
		}
		return nil
	}
}

// ValidateManagedCluster
func ValidateManagedCluster(importResourcesShouldExist bool) error {
	By("- Checking imported hub resources exist or not")
	By("- Confirming Necessary Resources")
	mc, _ := DynamicKubeClient.Resource(GVRManagedCluster).Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	if err := validateManagedClusterOwnerRef(mc); err != nil {
		return fmt.Errorf("Owner ref is not mch")
	}
	if err := ValidateImportHubResourcesExist(importResourcesShouldExist); err != nil {
		return fmt.Errorf("Resources are as they shouldn't")
	}
	if importResourcesShouldExist {
		if val := validateManagedClusterConditions(); val !=nil {
			return fmt.Errorf("cluster conditions")
		}
		return nil
	}
	return nil
}

// validateManagedClusterConditions 
func validateManagedClusterConditions() error {
	By("- Checking ManagedClusterConditions type true")
	mc, _ := DynamicKubeClient.Resource(GVRManagedCluster).Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	status, ok :=  mc.Object["status"].(map[string]interface{})
	if ok {
		joinErr := FindCondition(status, "ManagedClusterJoined", "True")
		avaiErr := FindCondition(status, "ManagedClusterConditionAvailable", "True")
		accpErr := FindCondition(status, "HubAcceptedManagedCluster", "True")
		if (joinErr != nil || avaiErr != nil || accpErr !=nil) {
			return fmt.Errorf("managedcluster conditions not all true")
		}
		return nil
	} else {
		return fmt.Errorf("no status")
	}
}

// validateManagedClusterOwnerRef helper func to validateManagedCluster
func validateManagedClusterOwnerRef(mc *unstructured.Unstructured) error {
	if (mc != nil) {
		name := mc.Object["metadata"].(map[string]interface{})["ownerReferences"].([]interface{})[0].(map[string]interface{})["name"]
		if name != MCHName {
			return fmt.Errorf("owner ref does not match mch name")
		}
		return nil
	}
	return nil
}

// ToggleDisableHubSelfManagement toggles the value of spec.disableHubSelfManagement from true to false or false to true
func ToggleDisableHubSelfManagement(disableHubSelfImport bool) error {
	mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
	Expect(err).To(BeNil())
	disableHubSelfManagementString := "disableHubSelfManagement"
	mch.Object["spec"].(map[string]interface{})[disableHubSelfManagementString] = disableHubSelfImport
	mch, err = DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Update(context.TODO(), mch, metav1.UpdateOptions{})
	Expect(err).To(BeNil())
	mch, err = DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
	if disableHubSelfManagement := mch.Object["spec"].(map[string]interface{})[disableHubSelfManagementString].(bool); disableHubSelfManagement != disableHubSelfImport {
		return fmt.Errorf("Spec was not updated")
	}
	return nil
}

// listByGVR keeps polling to get the object for timeout seconds
func listByGVR(clientHubDynamic dynamic.Interface, gvr schema.GroupVersionResource, namespace string, timeout int, expectedTotal int) *unstructured.UnstructuredList {
	if timeout < 1 {
		timeout = 1
	}
	var obj *unstructured.UnstructuredList

	Eventually(func() error {
		var err error
		namespace := clientHubDynamic.Resource(gvr).Namespace(namespace)

		// labelSelector := fmt.Sprintf("installer.name=%s, installer.namespace=%s", MCHName, MCHNamespace)
		// listOptions := metav1.ListOptions{
		// 	LabelSelector: labelSelector,
		// 	Limit:         100,
		// }

		obj, err = namespace.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		if len(obj.Items) < expectedTotal {
			return fmt.Errorf("Not all resources created in time. %d/%d appsubs found.", len(obj.Items), expectedTotal)
		}
		return nil
	}, timeout, 1).Should(BeNil())
	return obj
}

// GetSubscriptionSpec Returns Install Plan Mode
func GetSubscriptionSpec() map[string]interface{} {
	if os.Getenv("TEST_MODE") == "update" {
		return map[string]interface{}{
			"sourceNamespace":     os.Getenv("sourceNamespace"),
			"source":              os.Getenv("source"),
			"channel":             os.Getenv("channel"),
			"installPlanApproval": "Manual",
			"name":                os.Getenv("name"),
			"startingCSV":         fmt.Sprintf("advanced-cluster-management.v%s", os.Getenv("startVersion")),
		}
	}
	return map[string]interface{}{
		"sourceNamespace":     os.Getenv("sourceNamespace"),
		"source":              os.Getenv("source"),
		"channel":             os.Getenv("channel"),
		"installPlanApproval": "Automatic",
		"name":                os.Getenv("name"),
	}
}

// GetInstallPlanNameFromSub ...
func GetInstallPlanNameFromSub(sub *unstructured.Unstructured) (string, error) {
	if _, ok := sub.Object["status"]; !ok {
		return "", fmt.Errorf("Sub: %s has no 'status' field", sub.GetName())
	}
	status, ok := sub.Object["status"].(map[string]interface{})
	if !ok || status == nil {
		return "", fmt.Errorf("Sub: %s has no 'status' map", sub.GetName())
	}
	installplan, ok := status["installplan"].(map[string]interface{})
	if !ok || status == nil {
		return "", fmt.Errorf("Sub: %s has no 'installplan' map", sub.GetName())
	}

	return installplan["name"].(string), nil
}

// MarkInstallPlanAsApproved ...
func MarkInstallPlanAsApproved(ip *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	spec, ok := ip.Object["spec"].(map[string]interface{})
	if !ok || spec == nil {
		return nil, fmt.Errorf("Installplan: %s has no 'spec' map", ip.GetName())
	}
	spec["approved"] = true
	return ip, nil
}

// ShouldSkipSubscription skips subscription operations if set as true
func ShouldSkipSubscription() bool {
	skipSubscription := os.Getenv("skipSubscription")
	if skipSubscription == "true" {
		return true
	}
	return false
}

func GetCurrentVersionFromMCH() (string, error) {
	mch, err := DynamicKubeClient.Resource(GVRMultiClusterHub).Namespace(MCHNamespace).Get(context.TODO(), MCHName, metav1.GetOptions{})
	Expect(err).To(BeNil())
	status, ok := mch.Object["status"].(map[string]interface{})
	if !ok || status == nil {
		return "", fmt.Errorf("MultiClusterHub: %s has no 'status' map", mch.GetName())
	}
	version, ok := status["currentVersion"]
	if !ok {
		return "", fmt.Errorf("MultiClusterHub: %s status has no 'currentVersion' field", mch.GetName())
	}
	return version.(string), nil
}
