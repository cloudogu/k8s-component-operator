package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
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
	operatorLog          = ctrl.Log.WithName("component-operator")
	metricsAddr          string
	enableLeaderElection bool
	probeAddr            string
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

	err = configureManager(k8sManager, operatorConfig)
	if err != nil {
		return fmt.Errorf("failed to configure manager: %w", err)
	}

	return startK8sManager(k8sManager)
}

func configureManager(k8sManager manager.Manager, operatorConfig *config.OperatorConfig) error {
	err := configureReconciler(k8sManager, operatorConfig)
	if err != nil {
		return fmt.Errorf("failed to configure reconciler: %w", err)
	}

	// +kubebuilder:scaffold:builder
	err = addChecks(k8sManager)
	if err != nil {
		return fmt.Errorf("failed to add checks to the manager: %w", err)
	}

	return nil
}

func getK8sManagerOptions(operatorConfig *config.OperatorConfig) manager.Options {
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

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
		WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "951e217a.cloudogu.com",
	}

	return options
}

func startK8sManager(k8sManager manager.Manager) error {
	operatorLog.Info("starting manager")
	err := k8sManager.Start(ctrl.SetupSignalHandler())
	if err != nil {
		return fmt.Errorf("failed to start manager: %w", err)
	}

	return nil
}

func configureReconciler(k8sManager manager.Manager, operatorConfig *config.OperatorConfig) error {
	eventRecorder := k8sManager.GetEventRecorderFor("k8s-component-operator")

	clientSet, err := kubernetes.NewForConfig(k8sManager.GetConfig())
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	ctx := context.Background()
	helmRepoData, err := config.GetHelmRepositoryData(ctx, clientSet.CoreV1().ConfigMaps(operatorConfig.Namespace))
	if err != nil {
		return err
	}
	operatorConfig.HelmRepositoryData = helmRepoData

	componentClientSet, err := ecosystem.NewComponentClientset(k8sManager.GetConfig(), clientSet)
	if err != nil {
		return fmt.Errorf("failed to create component client set: %w", err)
	}

	debug := config.Stage == config.StageDevelopment
	helmClient, err := helm.NewClient(operatorConfig.Namespace, operatorConfig.HelmRepositoryData, debug, logging.FormattingLoggerWithName("helm-client", ctrl.Log.Info))
	if err != nil {
		return fmt.Errorf("failed to create helm client: %w", err)
	}

	componentReconciler := controllers.NewComponentReconciler(componentClientSet, helmClient, eventRecorder, operatorConfig.Namespace)
	err = componentReconciler.SetupWithManager(k8sManager)
	if err != nil {
		return fmt.Errorf("failed to setup reconciler with manager: %w", err)
	}

	healthReconcilers := health.NewController(operatorConfig.Namespace, componentClientSet)
	err = healthReconcilers.SetupWithManager(k8sManager)
	if err != nil {
		return fmt.Errorf("failed to setup health reconcilers with manager: %w", err)
	}

	return nil
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
