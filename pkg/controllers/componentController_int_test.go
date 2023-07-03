//go:build k8s_integration
// +build k8s_integration

package controllers

import (
	"context"
	_ "embed"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Dogu Upgrade Tests", func() {
	Context("Handle component resource", func() {
		ctx := context.TODO()
		installComponent := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator", Namespace: ""}, Spec: v1.ComponentSpec{
			Namespace: "test",
			Name:      "k8s-dogu-operator",
			Version:   "0.1.0",
		}}

		upgradeComponent := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator", Namespace: ""}, Spec: v1.ComponentSpec{
			Namespace: "test",
			Name:      "k8s-dogu-operator",
			Version:   "0.2.0",
		}}

		It("Should install component in cluster", func() {
			By("Creating component resource")
			helmClientMock.EXPECT().InstallOrUpgradeChart(ctx, mock.Anything, mock.Anything).Return(&release.Release{}, nil)
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Installation", "Starting installation...")
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
		})

		It("Should upgrade component in cluster", func() {
			By("Updating component resource")
			helmClientMock.EXPECT().InstallOrUpgradeChart(ctx, mock.Anything, mock.Anything).Return(&release.Release{}, nil)
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Upgrade", "Starting upgrade...")
			_, err := componentClient.Update(ctx, upgradeComponent, metav1.UpdateOptions{})
			Expect(err).Should(Succeed())

			By("expect updated component")
			Eventually(func() bool {
				get, err := componentClient.Get(ctx, "k8s-dogu-operator", metav1.GetOptions{})
				if err != nil {
					return false
				}

				if get.Spec.Version == "0.2.0" {
					return true
				}

				return false
			}, TimeoutInterval, PollingInterval).Should(BeTrue())
		})

		It("Should delete component in cluster", func() {
			By("Delete component resource")
			helmClientMock.EXPECT().UninstallReleaseByName("k8s-dogu-operator").Return(nil)
			recorderMock.EXPECT().Event(mock.Anything, "Normal", "Deinstallation", "Starting deinstallation...")

			err := componentClient.Delete(ctx, "k8s-dogu-operator", metav1.DeleteOptions{})
			Expect(err).Should(Succeed())
		})
	})
})
