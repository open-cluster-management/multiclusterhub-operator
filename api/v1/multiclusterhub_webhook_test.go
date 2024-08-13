// Copyright Contributors to the Open Cluster Management project

package v1

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	multiClusterHubName = "multiclusterhub"
)

var _ = Describe("Multiclusterhub webhook", func() {

	Context("Creating a Multiclusterhub", func() {
		It("Should successfully create multiclusterhub", func() {
			By("by creating a new standalone Multiclusterhub resource", func() {
				mch := &MultiClusterHub{
					ObjectMeta: metav1.ObjectMeta{
						Name:      multiClusterHubName,
						Namespace: "default",
					},
				}
				Expect(k8sClient.Create(ctx, mch)).Should(Succeed())
			})
		})

		It("Should fail to create multiclusterhub", func() {
			By("because of DeploymentMode", func() {
				mch := &MultiClusterHub{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("%s-2", multiClusterHubName),
						Namespace: "default",
					},
				}
				Expect(k8sClient.Create(ctx, mch)).NotTo(BeNil(), "Only one MCE in standalone mode allowed")
			})
			By("because of invalid AvailabilityConfig", func() {
				mch := &MultiClusterHub{
					ObjectMeta: metav1.ObjectMeta{
						Name:        fmt.Sprintf("%s-2", multiClusterHubName),
						Namespace:   "default",
						Annotations: map[string]string{"deploymentmode": string(ModeHosted)},
					},
					Spec: MultiClusterHubSpec{
						AvailabilityConfig: "low",
					},
				}
				Expect(k8sClient.Create(ctx, mch)).NotTo(BeNil(), "Invalid availability config is not allowed")
			})
			By("because of component configuration", func() {
				mch := &MultiClusterHub{
					ObjectMeta: metav1.ObjectMeta{
						Name:        fmt.Sprintf("%s-2", multiClusterHubName),
						Namespace:   "default",
						Annotations: map[string]string{"deploymentmode": string(ModeHosted)},
					},
					Spec: MultiClusterHubSpec{
						Overrides: &Overrides{
							Components: []ComponentConfig{
								{
									Name:    "fake-component",
									Enabled: true,
								},
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, mch)).NotTo(BeNil(), "Invalid components not allowed in config")
			})
		})

		It("Should fail to update multiclusterhub", func() {
			mch := &MultiClusterHub{}

			By("because of DeploymentMode", func() {
				Expect(k8sClient.Get(ctx, types.NamespacedName{Name: multiClusterHubName, Namespace: "default"}, mch)).To(Succeed())
				mch.SetAnnotations(map[string]string{"deploymentmode": string(ModeHosted)})
				Expect(k8sClient.Update(ctx, mch)).NotTo(BeNil(), "DeploymentMode should not change")
			})

			By("because of invalid component", func() {
				Expect(k8sClient.Get(ctx, types.NamespacedName{Name: multiClusterHubName, Namespace: "default"}, mch)).To(Succeed())
				mch.Spec.Overrides = &Overrides{
					Components: []ComponentConfig{
						{
							Name:    "fake-component",
							Enabled: true,
						},
					},
				}
				Expect(k8sClient.Update(ctx, mch)).NotTo(BeNil(), "invalid components should not be permitted")
			})
		})

		It("Should delete multiclusterhub", func() {
			mch := &MultiClusterHub{}

			By("deleting", func() {
				Expect(k8sClient.Get(ctx, types.NamespacedName{Name: multiClusterHubName, Namespace: "default"}, mch)).To(Succeed())
				Expect(k8sClient.Delete(ctx, mch)).To(BeNil(), "MCH delete was blocked unexpectedly")
			})
		})

	})

})
