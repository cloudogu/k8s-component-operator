package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cloudogu/k8s-component-operator/pkg/adapter/kubernetes/configref"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	componentClient "github.com/cloudogu/k8s-component-lib/client"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/cloudogu/k8s-component-operator/pkg/controllers"
	"github.com/cloudogu/k8s-component-operator/pkg/health"
	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/logging"
	// +kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
	// set up the logger before the actual logger is instantiated
	// the logger will be replaced later-on with a more sophisticated instance
	operatorLog = ctrl.Log.WithName("component-operator")
	metricsAddr string
	probeAddr   string
)

var (
	// Version of the application
	Version = "0.0.0"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(k8sv1.AddToScheme(scheme))

	// +kubebuilder:scaffold:scheme
}

func main() {
	err := startOperator()
	if err != nil {
		operatorLog.Error(err, "failed to start operator")
		os.Exit(1)
	}
}

func startOperator() error {
	err := logging.ConfigureLogger()
	if err != nil {
		return err
	}

	operatorConfig, err := config.NewOperatorConfig(Version)
	if err != nil {
		return fmt.Errorf("failed to create new operator configuration: %w", err)
	}

	options := getK8sManagerOptions(operatorConfig)
	k8sManager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		return fmt.Errorf("failed to start manager: %w", err)
	}

	ctx := ctrl.SetupSignalHandler()

	err = configureManager(ctx, k8sManager, operatorConfig)
	if err != nil {
		return fmt.Errorf("failed to configure manager: %w", err)
	}

	return startK8sManager(ctx, k8sManager)
}

func configureManager(ctx context.Context, k8sManager manager.Manager, operatorConfig *config.OperatorConfig) error {
	clientSet, err := createEcosystemClientSet(k8sManager)
	if err != nil {
		return err
	}

	err = configureReconciler(ctx, k8sManager, clientSet, operatorConfig)
	if err != nil {
		return fmt.Errorf("failed to configure reconciler: %w", err)
	}

	err = addRunners(k8sManager, clientSet, operatorConfig)
	if err != nil {
		return err
	}

	// +kubebuilder:scaffold:builder
	err = addChecks(k8sManager)
	if err != nil {
		return fmt.Errorf("failed to add checks to the manager: %w", err)
	}

	return nil
}

func addRunners(k8sManager manager.Manager, clientSet componentClient.ComponentEcosystemInterface, operatorConfig *config.OperatorConfig) error {
	healthSyncIntervalHandler := health.NewSyncIntervalHandler(operatorConfig.Namespace, clientSet, operatorConfig.HealthSyncIntervalMins)
	err := k8sManager.Add(healthSyncIntervalHandler)
	if err != nil {
		return err
	}

	return nil
}

func getK8sManagerOptions(operatorConfig *config.OperatorConfig) manager.Options {
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")

	options := ctrl.Options{
		Scheme:  scheme,
		Metrics: server.Options{BindAddress: metricsAddr},
		Cache: cache.Options{ByObject: map[client.Object]cache.ByObject{
			// Restrict namespace for components only as we want to reconcile Deployments,
			// StatefulSets and DaemonSets across all namespaces.
			&k8sv1.Component{}: {Namespaces: map[string]cache.Config{
				operatorConfig.Namespace: {},
			}},
		}},
		HealthProbeBindAddress: probeAddr,
	}

	return options
}

func startK8sManager(ctx context.Context, k8sManager manager.Manager) error {
	operatorLog.Info("starting manager")

	err := k8sManager.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start manager: %w", err)
	}

	return nil
}

func configureReconciler(ctx context.Context, k8sManager manager.Manager, clientSet componentClient.ComponentEcosystemInterface, operatorConfig *config.OperatorConfig) error {
	eventRecorder := k8sManager.GetEventRecorderFor("k8s-component-operator")

	helmRepoData, err := config.GetHelmRepositoryData(ctx, clientSet.CoreV1().ConfigMaps(operatorConfig.Namespace))
	if err != nil {
		return err
	}
	operatorConfig.HelmRepositoryData = helmRepoData

	debug := config.Stage == config.StageDevelopment
	helmClient, err := helm.NewClient(operatorConfig.Namespace, operatorConfig.HelmRepositoryData, debug, logging.FormattingLoggerWithName("helm-client", ctrl.Log.Info))
	if err != nil {
		return fmt.Errorf("failed to create helm client: %w", err)
	}

	yamlSerializer := yaml.NewSerializer()
	reader := configref.NewConfigMapRefReader(clientSet.CoreV1().ConfigMaps(operatorConfig.Namespace))

	componentReconciler := controllers.NewComponentReconciler(clientSet, helmClient, eventRecorder, operatorConfig.Namespace, operatorConfig.HelmClientTimeoutMins, yamlSerializer, reader, operatorConfig.RequeueTime)
	err = componentReconciler.SetupWithManager(k8sManager)
	if err != nil {
		return fmt.Errorf("failed to setup reconciler with manager: %w", err)
	}

	healthReconcilers := health.NewController(operatorConfig.Namespace, clientSet)
	err = healthReconcilers.SetupWithManager(k8sManager)
	if err != nil {
		return fmt.Errorf("failed to setup health reconcilers with manager: %w", err)
	}

	return nil
}

func createEcosystemClientSet(k8sManager manager.Manager) (*componentClient.EcosystemClientset, error) {
	clientSet, err := kubernetes.NewForConfig(k8sManager.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	componentClientSet, err := componentClient.NewComponentClientset(k8sManager.GetConfig(), clientSet)
	if err != nil {
		return nil, fmt.Errorf("failed to create component client set: %w", err)
	}

	return componentClientSet, nil
}

func addChecks(mgr manager.Manager) error {
	err := mgr.AddHealthzCheck("healthz", healthz.Ping)
	if err != nil {
		return fmt.Errorf("failed to add healthz check: %w", err)
	}

	err = mgr.AddReadyzCheck("readyz", healthz.Ping)
	if err != nil {
		return fmt.Errorf("failed to add readyz check: %w", err)
	}

	return nil
}
