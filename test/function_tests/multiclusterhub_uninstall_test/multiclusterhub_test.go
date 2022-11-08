// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package multiclusterhub_uninstall_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	utils "github.com/stolostron/multiclusterhub-operator/test/function_tests/utils"
	"k8s.io/client-go/dynamic"
)

var _ = Describe("Multiclusterhub", func() {

	It("Deleting and Validating MCH CR", func() {
		By(fmt.Sprintf("Deleting MultiClusterHub within %d minutes", utils.GetWaitInMinutes()))
		utils.DeleteIfExists(utils.DynamicKubeClient, utils.GVRMultiClusterHub, utils.MCHName, utils.MCHNamespace, true)

		Eventually(func() error {
			err := utils.ValidateDelete(utils.DynamicKubeClient)
			if err != nil {
				return err
			}
			return nil
		}, utils.GetWaitInMinutes()*60, 1).Should(BeNil())
	})

	if os.Getenv("full_test_suite") == "true" {
		It("Block MCH uninstall if certain resources exists", func() {
			By("Creating MultiClusterHub")
			utils.CreateMCHNotManaged()
			utils.ValidateMCH()
			utils.CreateDiscoveryConfig()

			By("Validating DiscoveryConfig blocks deletion")
			err := utils.DynamicKubeClient.Resource(utils.GVRMultiClusterHub).Namespace(utils.MCHNamespace).Delete(context.TODO(), utils.MCHName, metav1.DeleteOptions{})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(BeEquivalentTo("admission webhook \"multiclusterhub.validating-webhook.open-cluster-management.io\" denied the request: Cannot delete MultiClusterHub resource because DiscoveryConfig resource(s) exist"))

			utils.DeleteDiscoveryConfig()

			By("Validating MultiClusterObservability blocks deletion")

			utils.CreateObservabilityCRD()
			utils.CreateObservabilityCR()
			err = utils.DynamicKubeClient.Resource(utils.GVRMultiClusterHub).Namespace(utils.MCHNamespace).Delete(context.TODO(), utils.MCHName, metav1.DeleteOptions{})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).Should(BeEquivalentTo("admission webhook \"multiclusterhub.validating-webhook.open-cluster-management.io\" denied the request: Cannot delete MultiClusterHub resource because MultiClusterObservability resource(s) exist"))

			utils.DeleteObservabilityCR()
			utils.DeleteObservabilityCRD()
		})
		It("SAD CASE: Fail to remove a helmrelease (Left behind finalizer)", func() {
			By("Creating MultiClusterHub")
			utils.CreateMCHNotManaged()
			utils.ValidateMCH()
			AddFinalizerToHelmRelease(utils.DynamicKubeClient)
			utils.DeleteIfExists(utils.DynamicKubeClient, utils.GVRMultiClusterHub, utils.MCHName, utils.MCHNamespace, false)
			Expect(utils.ValidateDelete(utils.DynamicKubeClient)).ShouldNot(BeNil())
			utils.ValidateConditionDuringUninstall()
			Expect(utils.ValidatePhase("Uninstalling")).To(BeNil())

			Eventually(func() error {
				err := RemoveFinalizerFromHelmRelease(utils.DynamicKubeClient)
				if err != nil {
					return err
				}
				return nil
			}, utils.GetWaitInMinutes()*60, 1).Should(BeNil())
			utils.DeleteIfExists(utils.DynamicKubeClient, utils.GVRMultiClusterHub, utils.MCHName, utils.MCHNamespace, true)
			Expect(utils.ValidateDelete(utils.DynamicKubeClient)).Should(BeNil())
		})

		It("SAD CASE: Fail to remove managedcluster (Left behind finalizer)", func() {
			By("Creating MultiClusterHub")
			utils.CreateDefaultMCH()
			utils.ValidateMCH()
			AddFinalizerToManagedCluster(utils.DynamicKubeClient)
			utils.DeleteIfExists(utils.DynamicKubeClient, utils.GVRMultiClusterHub, utils.MCHName, utils.MCHNamespace, false)
			Expect(utils.ValidateDelete(utils.DynamicKubeClient)).ShouldNot(BeNil())
			utils.ValidateConditionDuringUninstall()

			Eventually(func() error {
				err := RemoveFinalizerFromManagedCluster(utils.DynamicKubeClient)
				if err != nil {
					return err
				}
				return nil
			}, utils.GetWaitInMinutes()*60, 1).Should(BeNil())
			utils.DeleteIfExists(utils.DynamicKubeClient, utils.GVRMultiClusterHub, utils.MCHName, utils.MCHNamespace, true)
			Expect(utils.ValidateDelete(utils.DynamicKubeClient)).Should(BeNil())
		})
	}
})

// AddFinalizerToHelmRelease ...
func AddFinalizerToHelmRelease(clientHubDynamic dynamic.Interface) error {
	By("Adding a test finalizer to a helmrelease")

	appSubLink := clientHubDynamic.Resource(utils.GVRAppSub).Namespace(utils.MCHNamespace)
	appSub, err := appSubLink.Get(context.TODO(), "console-chart-sub", metav1.GetOptions{})
	Expect(err).Should(BeNil())

	helmReleaseName := fmt.Sprintf("%s-%s", strings.Replace(appSub.GetName(), "-sub", "", 1), appSub.GetUID()[0:5])

	helmReleaseLink := clientHubDynamic.Resource(utils.GVRHelmRelease).Namespace(utils.MCHNamespace)
	helmRelease, err := helmReleaseLink.Get(context.TODO(), helmReleaseName, metav1.GetOptions{})
	Expect(err).Should(BeNil())

	finalizers := []string{"test-finalizer"}

	helmRelease.SetFinalizers(finalizers)
	_, err = helmReleaseLink.Update(context.TODO(), helmRelease, metav1.UpdateOptions{})
	Expect(err).Should(BeNil())

	return nil
}

// AddFinalizerToManagedCluster ...
func AddFinalizerToManagedCluster(clientHubDynamic dynamic.Interface) error {
	By("Adding a test finalizer to managed cluster")

	mc, err := clientHubDynamic.Resource(utils.GVRManagedCluster).Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	Expect(err).Should(BeNil())

	finalizers := []string{"test-finalizer"}

	mc.SetFinalizers(finalizers)
	_, err = clientHubDynamic.Resource(utils.GVRManagedCluster).Update(context.TODO(), mc, metav1.UpdateOptions{})
	Expect(err).Should(BeNil())

	return nil
}

// RemoveFinalizerFromHelmRelease ...
func RemoveFinalizerFromHelmRelease(clientHubDynamic dynamic.Interface) error {
	By("Removing test finalizer from helmrelease")

	labelSelector := fmt.Sprintf("installer.name=%s, installer.namespace=%s", utils.MCHName, utils.MCHNamespace)
	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
		Limit:         100,
	}

	helmReleaseLink := clientHubDynamic.Resource(utils.GVRHelmRelease).Namespace(utils.MCHNamespace)
	helmReleases, err := helmReleaseLink.List(context.TODO(), listOptions)
	if err != nil {
		return err
	}

	if len(helmReleases.Items) == 0 {
		return fmt.Errorf("No helmreleases found")
	}
	helmRelease := helmReleases.Items[0]
	helmRelease.SetFinalizers([]string{})

	_, err = helmReleaseLink.Update(context.TODO(), &helmRelease, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// RemoveFinalizerFromManagedCluster ...
func RemoveFinalizerFromManagedCluster(clientHubDynamic dynamic.Interface) error {
	By("Removing test finalizer from managed cluster")

	mc, err := clientHubDynamic.Resource(utils.GVRManagedCluster).Get(context.TODO(), "local-cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}
	mc.SetFinalizers([]string{})

	_, err = clientHubDynamic.Resource(utils.GVRManagedCluster).Update(context.TODO(), mc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
