//go:build k8s_integration
// +build k8s_integration

package controllers

import (
	"context"
	_ "embed"
	"fmt"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

type mockeryGinkgoLogger struct {
}

func (c mockeryGinkgoLogger) Logf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(GinkgoWriter, strings.ReplaceAll(strings.ReplaceAll(format, "PASS", "\nPASS"), "FAIL", "\nFAIL"), args...)
}

func (c mockeryGinkgoLogger) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(GinkgoWriter, format, args...)
}

func (c mockeryGinkgoLogger) FailNow() {
	println("fail now")
}

var _ = Describe("Dogu Upgrade Tests", func() {
	mockeryT := &mockeryGinkgoLogger{}

	Context("Handle component resource", func() {
		ctx := context.TODO()
		installComponent := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator", Namespace: namespace}, Spec: v1.ComponentSpec{
			Namespace: "test",
			Name:      "k8s-dogu-operator",
			Version:   "0.1.0",
		}}

		It("Should install component in cluster", func() {
			By("Creating component resource")
			*helmClientMock = MockHelmClient{}
			helmClientMock.EXPECT().InstallOrUpgrade(mock.Anything, mock.Anything).Return(nil)
			*recorderMock = MockEventRecorder{}
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Installation", "Starting installation...")
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Installation", "Installation successful")

			_, err := componentClient.Create(ctx, installComponent, metav1.CreateOptions{})
			Expect(err).Should(Succeed())

			By("Expect created component")
			Eventually(func() bool {
				get, err := componentClient.Get(ctx, "k8s-dogu-operator", metav1.GetOptions{})
				if err != nil {
					return false
				}

				finalizers := get.Finalizers
				if len(finalizers) == 1 && finalizers[0] == "component-finalizer" {
					return true
				}

				return false
			}, TimeoutInterval, PollingInterval).Should(BeTrue())

			Expect(helmClientMock.AssertExpectations(mockeryT)).To(BeTrue())
			Expect(recorderMock.AssertExpectations(mockeryT)).To(BeTrue())
		})

		It("Should upgrade component in cluster", func() {
			By("Updating component resource")
			*helmClientMock = MockHelmClient{}
			helmClientMock.EXPECT().ListDeployedReleases().Return([]*release.Release{{Name: installComponent.Spec.Name, Namespace: installComponent.Namespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.1.0"}}}}, nil)
			helmClientMock.EXPECT().InstallOrUpgrade(mock.Anything, mock.Anything).Return(nil)
			*recorderMock = MockEventRecorder{}
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Upgrade", "Starting upgrade...")
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Upgrade", "Upgrade successful")

			upgradeComponent, err := componentClient.Get(ctx, "k8s-dogu-operator", metav1.GetOptions{})
			Expect(err).Should(Succeed())
			upgradeComponent.Spec.Version = "0.2.0"
			_, err = componentClient.Update(ctx, upgradeComponent, metav1.UpdateOptions{})
			Expect(err).Should(Succeed())

			Eventually(func() bool {
				comp, err := componentClient.Get(ctx, "k8s-dogu-operator", metav1.GetOptions{})
				if err != nil {
					return false
				}

				startingResourceVersion, err := strconv.Atoi(upgradeComponent.ResourceVersion)
				if err != nil {
					return false
				}

				actualResourceVersion, err := strconv.Atoi(comp.ResourceVersion)
				if err != nil {
					return false
				}

				// Check the resource version because the upgrade routine actually does nothing other for what you can wait for
				// (e.g. a deployment creation).
				// The testcode reaches to fast the AssertExpectations calls before the component-operator calls the methods.
				// With a real cluster we should replace that check with the deployment rollout.
				return comp.Status.Status == "installed" && actualResourceVersion >= startingResourceVersion+2
			}, TimeoutInterval, PollingInterval).Should(BeTrue())

			Expect(helmClientMock.AssertExpectations(mockeryT)).To(BeTrue())
			Expect(recorderMock.AssertExpectations(mockeryT)).To(BeTrue())
		})

		It("Should delete component in cluster", func() {
			By("Delete component resource")
			*helmClientMock = MockHelmClient{}
			helmClientMock.EXPECT().Uninstall(mock.Anything).Return(nil)
			*recorderMock = MockEventRecorder{}
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Deinstallation", "Starting deinstallation...")
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Deinstallation", "Deinstallation successful")

			err := componentClient.Delete(ctx, "k8s-dogu-operator", metav1.DeleteOptions{})
			Expect(err).Should(Succeed())

			Eventually(func() bool {
				_, err := componentClient.Get(ctx, "k8s-dogu-operator", metav1.GetOptions{})
				return errors.IsNotFound(err)
			}, TimeoutInterval, PollingInterval).Should(BeTrue())

			Expect(recorderMock.AssertExpectations(mockeryT)).To(BeTrue())
			Expect(helmClientMock.AssertExpectations(mockeryT)).To(BeTrue())
		})
	})
})
