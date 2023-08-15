package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

type operation string

const (
	// InstallEventReason The name of the installation event
	InstallEventReason = "Installation"
	// DeinstallationEventReason The name of the deinstallation event
	DeinstallationEventReason = "Deinstallation"
	// UpgradeEventReason The name of the upgrade event
	UpgradeEventReason = "Upgrade"
	// DowngradeEventReason The name of the downgrade event
	DowngradeEventReason = "Downgrade"
	// Install represents the install-operation
	Install = operation("Install")
	// Upgrade represents the upgrade-operation
	Upgrade = operation("Upgrade")
	// Downgrade represents the downgrade-operation. Currently not supported.
	Downgrade = operation("Downgrade")
	// Delete represents the delete-operation
	Delete = operation("Delete")
	// Ignore represents the ignore-operation
	Ignore = operation("Ignore")
)

// ComponentManager abstracts the simple component operations in a k8s CES.
type ComponentManager interface {
	InstallManager
	DeleteManager
	UpgradeManager
}

// componentReconciler watches every Component object in the cluster and handles them accordingly.
type componentReconciler struct {
	componentClient  ecosystem.ComponentInterface
	recorder         record.EventRecorder
	componentManager ComponentManager
	helmClient       HelmClient
}

// NewComponentReconciler creates a new component reconciler.
func NewComponentReconciler(componentClient ecosystem.ComponentInterface, helmClient HelmClient, recorder record.EventRecorder) *componentReconciler {
	return &componentReconciler{
		componentClient:  componentClient,
		recorder:         recorder,
		componentManager: NewComponentManager(componentClient, helmClient, recorder),
		helmClient:       helmClient,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *componentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile this component", "component", req.Name)
	component, err := r.componentClient.Get(ctx, req.Name, v1.GetOptions{})

	if err != nil {
		logger.Info(fmt.Sprintf("failed to get component %+v: %s", req, err))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info(fmt.Sprintf("Component %+v has been found", req))

	operation, err := r.evaluateRequiredOperation(ctx, component)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to evaluate required operation: %w", err)
	}
	logger.Info(fmt.Sprintf("Required operation is %s", operation))

	switch operation {
	case Install:
		return ctrl.Result{}, r.performInstallOperation(ctx, component)
	case Delete:
		return ctrl.Result{}, r.performDeleteOperation(ctx, component)
	case Upgrade:
		return ctrl.Result{}, r.performUpgradeOperation(ctx, component)
	case Downgrade:
		return ctrl.Result{}, r.performDowngradeOperation(component)
	case Ignore:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, nil
	}
}

func (r *componentReconciler) performInstallOperation(ctx context.Context, component *k8sv1.Component) error {
	return r.performOperation(ctx, component, InstallEventReason, r.componentManager.Install)
}

func (r *componentReconciler) performUpgradeOperation(ctx context.Context, component *k8sv1.Component) error {
	return r.performOperation(ctx, component, UpgradeEventReason, r.componentManager.Upgrade)
}

func (r *componentReconciler) performDeleteOperation(ctx context.Context, component *k8sv1.Component) error {
	return r.performOperation(ctx, component, DeinstallationEventReason, r.componentManager.Delete)
}

func (r *componentReconciler) performDowngradeOperation(component *k8sv1.Component) error {
	r.recorder.Event(component, corev1.EventTypeWarning, DowngradeEventReason, "component downgrades are not allowed")
	return fmt.Errorf("downgrades are not allowed")
}

func (r *componentReconciler) performOperation(ctx context.Context, component *k8sv1.Component, eventReason string, operationFn func(context.Context, *k8sv1.Component) error) error {
	err := operationFn(ctx, component)
	eventType := corev1.EventTypeNormal
	message := fmt.Sprintf("%s successful", eventReason)
	if err != nil {
		eventType = corev1.EventTypeWarning
		printError := strings.ReplaceAll(err.Error(), "\n", "")
		message = fmt.Sprintf("%s failed. Reason: %s", eventReason, printError)
	}

	// on self-upgrade of the component-operator this event might not get send, because the operator is already shutting down
	r.recorder.Event(component, eventType, eventReason, message)

	return err
}

func (r *componentReconciler) evaluateRequiredOperation(ctx context.Context, component *k8sv1.Component) (operation, error) {
	logger := log.FromContext(ctx)
	if component.DeletionTimestamp != nil && !component.DeletionTimestamp.IsZero() {
		return Delete, nil
	}

	switch component.Status.Status {
	case k8sv1.ComponentStatusNotInstalled:
		return Install, nil
	case k8sv1.ComponentStatusInstalled:
		operation, err := r.getChangeOperation(component)
		if err != nil {
			return "", err
		}

		return operation, nil
	case k8sv1.ComponentStatusInstalling:
		return Ignore, nil
	case k8sv1.ComponentStatusDeleting:
		return Ignore, nil
	case k8sv1.ComponentStatusUpgrading:
		return Ignore, nil
	default:
		logger.Info(fmt.Sprintf("Found unknown operation for component status: %s", component.Status.Status))
		return Ignore, nil
	}
}

func (r *componentReconciler) getChangeOperation(component *k8sv1.Component) (operation, error) {
	deployedReleases, err := r.helmClient.ListDeployedReleases()
	if err != nil {
		return "", fmt.Errorf("failed to get deployed helm releases: %w", err)
	}

	for _, deployedRelease := range deployedReleases {
		// This will allow a namespace switch e. g. k8s/dogu-operator -> k8s-testing/dogu-operator.
		if deployedRelease.Name == component.Spec.Name && deployedRelease.Namespace == component.Namespace {
			return getChangeOperationForRelease(component, deployedRelease)
		}
	}

	return Ignore, nil
}

func getChangeOperationForRelease(component *k8sv1.Component, release *release.Release) (operation, error) {
	chart := release.Chart
	deployedAppVersion, err := core.ParseVersion(chart.AppVersion())
	if err != nil {
		return "", fmt.Errorf("failed to parse app version %s from helm chart %s: %w", chart.AppVersion(), chart.Name(), err)
	}

	componentVersion, err := core.ParseVersion(component.Spec.Version)
	if err != nil {
		return "", fmt.Errorf("failed to parse component version %s from %s: %w", component.Spec.Version, component.Spec.Name, err)
	}

	if deployedAppVersion.IsOlderThan(componentVersion) {
		return Upgrade, nil
	}

	if deployedAppVersion.IsNewerThan(componentVersion) {
		return Downgrade, nil
	}

	return Ignore, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *componentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		For(&k8sv1.Component{}).
		Complete(r)
}
