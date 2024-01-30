package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/health"
	"reflect"
	"strings"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	semver "github.com/Masterminds/semver/v3"

	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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
	// RequeueEventReason The name of the requeue event
	RequeueEventReason = "Requeue"
	// FailedNameValidationEventReason The name of the event to validate spec.name and metadata.name of a component.
	FailedNameValidationEventReason = "FailedNameValidation"
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
	installManager
	deleteManager
	upgradeManager
}

// ComponentReconciler watches every Component object in the cluster and handles them accordingly.
type ComponentReconciler struct {
	clientSet        componentEcosystemInterface
	recorder         record.EventRecorder
	componentManager ComponentManager
	helmClient       helmClient
	requeueHandler   requeueHandler
	namespace        string
}

// NewComponentReconciler creates a new component reconciler.
func NewComponentReconciler(clientSet componentEcosystemInterface, helmClient helmClient, recorder record.EventRecorder, namespace string) *ComponentReconciler {
	componentRequeueHandler := NewComponentRequeueHandler(clientSet, recorder, namespace)
	return &ComponentReconciler{
		clientSet: clientSet,
		recorder:  recorder,
		componentManager: NewComponentManager(
			clientSet.ComponentV1Alpha1().Components(namespace),
			helmClient,
			health.NewManager(namespace, clientSet),
			recorder,
		),
		helmClient:     helmClient,
		requeueHandler: componentRequeueHandler,
		namespace:      namespace,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile this component", "component", req.Name)
	component, err := r.clientSet.ComponentV1Alpha1().Components(req.Namespace).Get(ctx, req.Name, v1.GetOptions{})

	if err != nil {
		logger.Info(fmt.Sprintf("failed to get component %+v: %s", req, err))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info(fmt.Sprintf("Component %+v has been found", req))

	success := r.validateName(component)
	if !success {
		return finishOperation()
	}

	operation, err := r.evaluateRequiredOperation(ctx, component)
	if err != nil {
		return requeueWithError(fmt.Errorf("failed to evaluate required operation: %w", err))
	}
	logger.Info(fmt.Sprintf("Required operation is %s", operation))

	switch operation {
	case Install:
		return r.performInstallOperation(ctx, component)
	case Delete:
		return r.performDeleteOperation(ctx, component)
	case Upgrade:
		return r.performUpgradeOperation(ctx, component)
	case Downgrade:
		return r.performDowngradeOperation(component)
	case Ignore:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, nil
	}
}

func (r *ComponentReconciler) validateName(component *k8sv1.Component) (success bool) {
	if component.ObjectMeta.Name != component.Spec.Name {
		r.recorder.Eventf(component, corev1.EventTypeWarning, FailedNameValidationEventReason, "Component resource does not follow naming rules: The component's metadata.name '%s' must be the same as its spec.name '%s'.", component.ObjectMeta.Name, component.Spec.Name)
		return false
	}

	return true
}

func (r *ComponentReconciler) performInstallOperation(ctx context.Context, component *k8sv1.Component) (ctrl.Result, error) {
	return r.performOperation(ctx, component, InstallEventReason, k8sv1.ComponentStatusNotInstalled, r.componentManager.Install)
}

func (r *ComponentReconciler) performUpgradeOperation(ctx context.Context, component *k8sv1.Component) (ctrl.Result, error) {
	return r.performOperation(ctx, component, UpgradeEventReason, k8sv1.ComponentStatusInstalled, r.componentManager.Upgrade)
}

func (r *ComponentReconciler) performDeleteOperation(ctx context.Context, component *k8sv1.Component) (ctrl.Result, error) {
	return r.performOperation(ctx, component, DeinstallationEventReason, k8sv1.ComponentStatusInstalled, r.componentManager.Delete)
}

func (r *ComponentReconciler) performDowngradeOperation(component *k8sv1.Component) (ctrl.Result, error) {
	r.recorder.Event(component, corev1.EventTypeWarning, DowngradeEventReason, "component downgrades are not allowed")
	return ctrl.Result{}, fmt.Errorf("downgrades are not allowed")
}

// performOperation executes the given operationFn and requeues if necessary.
// When requeueing, the sourceComponentStatus is set as the components' status.
func (r *ComponentReconciler) performOperation(
	ctx context.Context,
	component *k8sv1.Component,
	eventReason string,
	requeueStatus string,
	operationFn func(context.Context, *k8sv1.Component) error,
) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	operationError := operationFn(ctx, component)
	contextMessageOnError := fmt.Sprintf("%s failed with component %s", eventReason, component.Name)
	eventType := corev1.EventTypeNormal
	message := fmt.Sprintf("%s successful", eventReason)
	if operationError != nil {
		eventType = corev1.EventTypeWarning
		printError := strings.ReplaceAll(operationError.Error(), "\n", "")
		message = fmt.Sprintf("%s failed. Reason: %s", eventReason, printError)
		logger.Error(operationError, message)
	}

	// on self-upgrade of the component-operator this event might not get send, because the operator is already shutting down
	r.recorder.Event(component, eventType, eventReason, message)

	result, handleErr := r.requeueHandler.Handle(ctx, contextMessageOnError, component, operationError, requeueStatus)
	if handleErr != nil {
		r.recorder.Eventf(component, corev1.EventTypeWarning, RequeueEventReason,
			"Failed to requeue the %s.", strings.ToLower(eventReason))
		return requeueWithError(fmt.Errorf("failed to handle requeue: %w", handleErr))
	}

	return requeueOrFinishOperation(result)
}

// requeueWithError is a syntax sugar function to express that every non-nil error will result in a requeue
// operation.
//
// Use requeueOrFinishOperation() if the reconciler should requeue the operation because of the result instead of an
// error.
// Use finishOperation() if the reconciler should not requeue the operation.
func requeueWithError(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

// requeueOrFinishOperation is a syntax sugar function to express that the there is no error to handle but the result
// controls whether the current operation should be finished or requeued.
//
// Use requeueWithError() if the reconciler should requeue the operation because of a non-nil error.
// Use finishOperation() if the reconciler should not requeue the operation.
func requeueOrFinishOperation(result ctrl.Result) (ctrl.Result, error) {
	return result, nil
}

// finishOperation is a syntax sugar function to express that the current operation should be finished and not be
// requeued. This can happen if the operation was successful or even if an unhandleable error occurred which prevents
// requeueing.
//
// Use requeueOrFinishOperation() or requeueWithError() if the reconciler should requeue the operation.
func finishOperation() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *ComponentReconciler) evaluateRequiredOperation(ctx context.Context, component *k8sv1.Component) (operation, error) {
	logger := log.FromContext(ctx)
	if component.DeletionTimestamp != nil && !component.DeletionTimestamp.IsZero() {
		return Delete, nil
	}

	switch component.Status.Status {
	case k8sv1.ComponentStatusNotInstalled:
		return Install, nil
	case k8sv1.ComponentStatusInstalled:
		operation, err := r.getChangeOperation(ctx, component)
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

func (r *ComponentReconciler) getChangeOperation(ctx context.Context, component *k8sv1.Component) (operation, error) {
	logger := log.FromContext(ctx)

	deployedReleases, err := r.helmClient.ListDeployedReleases()
	if err != nil {
		return "", fmt.Errorf("failed to get deployed helm releases: %w", err)
	}

	for _, deployedRelease := range deployedReleases {

		isComponentToBeChanged := deployedRelease.Name == component.Spec.Name
		targetNamespace := component.Spec.DeployNamespace

		if targetNamespace == "" {
			targetNamespace = component.Namespace
		}

		existsReleaseInTargetNamespace := deployedRelease.Namespace == targetNamespace

		if isComponentToBeChanged {
			logger.Info("Found existing release for reconciled component",
				"releaseNamespace", deployedRelease.Namespace, "targetNamespace", targetNamespace)
			if existsReleaseInTargetNamespace {
				return r.getChangeOperationForRelease(component, deployedRelease)
			}
		}
	}

	return Ignore, nil
}

func (r *ComponentReconciler) isValuesChanged(deployedRelease *release.Release, component *k8sv1.Component) (bool, error) {
	deployedValues, err := r.helmClient.GetReleaseValues(deployedRelease.Name, false)
	if err != nil {
		return false, fmt.Errorf("failed to get values.yaml from release %s: %w", deployedRelease.Name, err)
	}

	// TODO Check changes for mappedValues. Merge them here. Maybe extend chartSpec and reuse them later.
	chartSpecValues, err := r.helmClient.GetChartSpecValues(component.GetHelmChartSpec())
	if err != nil {
		return false, fmt.Errorf("failed to get values.yaml from component %s: %w", component.GetHelmChartSpec().ChartName, err)
	}

	// if no additional values are set, the maps will look like this:
	// deployedValues=map[string]interface {}(nil)                                                                                                                                        │
	// chartSpecValues=map[string]interface {}{}
	// this is treated as a difference by DeepEqual, so we have to handle this edge case manually
	if len(deployedValues) == 0 && len(chartSpecValues) == 0 {
		return false, nil
	}

	return !reflect.DeepEqual(deployedValues, chartSpecValues), nil
}

func (r *ComponentReconciler) getChangeOperationForRelease(component *k8sv1.Component, release *release.Release) (operation, error) {
	chart := release.Chart
	deployedAppVersion, err := semver.NewVersion(chart.AppVersion())
	if err != nil {
		return "", fmt.Errorf("failed to parse app version %s from helm chart %s: %w", chart.AppVersion(), chart.Name(), err)
	}

	componentVersion, err := semver.NewVersion(component.Spec.Version)
	if err != nil {
		return "", fmt.Errorf("failed to parse component version %s from %s: %w", component.Spec.Version, component.Spec.Name, err)
	}

	if deployedAppVersion.LessThan(componentVersion) {
		return Upgrade, nil
	}

	if deployedAppVersion.GreaterThan(componentVersion) {
		return Downgrade, nil
	}

	isValuesChanged, err := r.isValuesChanged(release, component)
	if err != nil {
		return "", fmt.Errorf("failed to compare Values.yaml files of component %s: %w", component.Name, err)
	}
	if isValuesChanged {
		return Upgrade, nil
	}

	return Ignore, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		For(&k8sv1.Component{}).
		Complete(r)
}
