//go:build k8s_integration
// +build k8s_integration

package controllers

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bombsimon/logrusr/v2"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.
var testEnv *envtest.Environment
var cancel context.CancelFunc

// Used in other integration tests
var (
	componentClient ecosystem.ComponentInterface

	helmClientMock *MockHelmClient

	recorderMock *MockEventRecorder
)

const TimeoutInterval = time.Second * 10
const PollingInterval = time.Second * 1

var oldGetConfig func() (*rest.Config, error)
var oldGetConfigOrDie func() *rest.Config

func TestAPIs(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Controller Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	// We need to ensure that the development stage flag is not passed by our makefiles to prevent the dogu operator
	// from running in the developing mode. The developing mode changes some operator behaviour. Our integration test
	// aim to test the production functionality of the operator.
	err := os.Unsetenv(config.StageEnvironmentVariable)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = os.Setenv(config.StageEnvironmentVariable, config.StageProduction)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	config.Stage = config.StageProduction

	logf.SetLogger(logrusr.New(logrus.New()))

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.TODO())

	ginkgo.By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(cfg).NotTo(gomega.BeNil())

	oldGetConfig = ctrl.GetConfig
	ctrl.GetConfig = func() (*rest.Config, error) {
		return cfg, nil
	}

	oldGetConfigOrDie = ctrl.GetConfigOrDie
	ctrl.GetConfigOrDie = func() *rest.Config {
		return cfg
	}

	err = k8sv1.AddToScheme(scheme.Scheme)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// +kubebuilder:scaffold:scheme
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	t := &testing.T{}
	helmClientMock = NewMockHelmClient(t)
	recorderMock = NewMockEventRecorder(t)

	config := &config.OperatorConfig{}

	clientSet, err := kubernetes.NewForConfig(cfg)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	componentClientSet, err := ecosystem.NewComponentClientset(k8sManager.GetConfig(), clientSet)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	componentClient = componentClientSet.EcosystemV1Alpha1().Components(config.Namespace)
	reconciler := NewComponentReconciler(componentClient, helmClientMock, recorderMock, config)

	err = reconciler.SetupWithManager(k8sManager)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	go func() {
		err = k8sManager.Start(ctx)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}()
}, 60)

var _ = ginkgo.AfterSuite(func() {
	cancel()
	ginkgo.By("tearing down the test environment")
	err := testEnv.Stop()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ctrl.GetConfig = oldGetConfig
	ctrl.GetConfigOrDie = oldGetConfigOrDie
})
