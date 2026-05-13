package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
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

type newHelmClientFunc func() (*helm.Client, error)

func (f newHelmClientFunc) NewHelmClient() (helmClient, error) {
	return f()
}

// ComponentManager abstracts the simple component operations in a k8s CES.
type ComponentManager interface {
	installManager
	deleteManager
	upgradeManager
}

// ComponentReconciler watches every Component object in the cluster and handles them accordingly.
type ComponentReconciler struct {
	clientSet                 componentEcosystemInterface
	recorder                  record.EventRecorder
	componentManagerFactory   componentManagerFactory
	helmClientFactory         helmClientFactory
	operationEvaluatorFactory operationEvaluatorFactory
	requeueHandler            requeueHandler
	namespace                 string
	timeout                   time.Duration
	yamlSerializer            yaml.Serializer
	reader                    configMapRefReader
	configMapInterface        configMapInterface
}

func NewComponentReconciler(clientSet componentEcosystemInterface, newHelmClient newHelmClientFunc, recorder record.EventRecorder, namespace string, timeout time.Duration, yamlSerializer yaml.Serializer, reader configMapRefReader, requeueTime time.Duration) *ComponentReconciler {
	componentRequeueHandler := NewComponentRequeueHandler(clientSet, recorder, namespace, requeueTime)

	return &ComponentReconciler{
		clientSet: clientSet,
		recorder:  recorder,
		componentManagerFactory: &defaultComponentManagerFactory{
			namespace: namespace,
			clientSet: clientSet,
			recorder:  recorder,
			timeout:   timeout,
		},
		helmClientFactory: newHelmClient,
		operationEvaluatorFactory: &defaultOperationEvaluatorFactory{
			recorder:       recorder,
			timeout:        timeout,
			yamlSerializer: yamlSerializer,
			reader:         reader,
		},
		requeueHandler:     componentRequeueHandler,
		namespace:          namespace,
		yamlSerializer:     yamlSerializer,
		reader:             reader,
		timeout:            timeout,
		configMapInterface: clientSet.CoreV1().ConfigMaps(namespace),
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

	hc, err := r.helmClientFactory.NewHelmClient()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create helm client: %w", err)
	}

	operationEvaluator := r.operationEvaluatorFactory.NewOperationEvaluator(hc)
	operation, err := operationEvaluator.EvaluateRequiredOperation(ctx, component)
	if err != nil {
		return requeueWithError(fmt.Errorf("failed to evaluate required operation: %w", err))
	}
	logger.Info(fmt.Sprintf("Required operation is %s", operation))

	componentManager := r.componentManagerFactory.NewComponentManager(hc)

	switch operation {
	case Install:
		return r.performInstallOperation(ctx, component, componentManager)
	case Delete:
		return r.performDeleteOperation(ctx, component, componentManager)
	case Upgrade:
		return r.performUpgradeOperation(ctx, component, componentManager)
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

func (r *ComponentReconciler) performInstallOperation(ctx context.Context, component *k8sv1.Component, componentManager ComponentManager) (ctrl.Result, error) {
	return r.performOperation(ctx, component, InstallEventReason, k8sv1.ComponentStatusTryToInstall, componentManager.Install)
}

func (r *ComponentReconciler) performUpgradeOperation(ctx context.Context, component *k8sv1.Component, componentManager ComponentManager) (ctrl.Result, error) {
	return r.performOperation(ctx, component, UpgradeEventReason, k8sv1.ComponentStatusTryToUpgrade, componentManager.Upgrade)
}

func (r *ComponentReconciler) performDeleteOperation(ctx context.Context, component *k8sv1.Component, componentManager ComponentManager) (ctrl.Result, error) {
	return r.performOperation(ctx, component, DeinstallationEventReason, k8sv1.ComponentStatusTryToDelete, componentManager.Delete)
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

// SetupWithManager sets up the controller with the Manager.
func (r *ComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	controllerOptions := mgr.GetControllerOptions()
	options := controller.TypedOptions[reconcile.Request]{
		SkipNameValidation: controllerOptions.SkipNameValidation,
		RecoverPanic:       controllerOptions.RecoverPanic,
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithOptions(options).
		For(&k8sv1.Component{}).
		WatchesRawSource(r.getConfigMapKind(mgr)).
		Complete(r)
}

func (r *ComponentReconciler) getComponentRequest(ctx context.Context, cm *corev1.ConfigMap) []reconcile.Request {
	list, err := r.clientSet.ComponentV1Alpha1().Components(r.namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil
	}
	var componentRequest []reconcile.Request
	for _, component := range list.Items {
		if component.Spec.ValuesConfigRef != nil && component.Spec.ValuesConfigRef.Name == cm.Name {
			componentRequest = append(componentRequest, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      component.Name,
					Namespace: r.namespace,
				},
			})
		}
		labels := cm.Labels
		for key, label := range labels {
			if key == "k8s.cloudogu.com/component.config" && label == component.Name {
				componentRequest = append(componentRequest, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      component.Name,
						Namespace: r.namespace,
					},
				})
			}
		}
	}
	return componentRequest
}

func (r *ComponentReconciler) getConfigMapKind(mgr ctrl.Manager) source.TypedSyncingSource[reconcile.Request] {
	return source.TypedKind(
		mgr.GetCache(),
		&corev1.ConfigMap{},
		handler.TypedEnqueueRequestsFromMapFunc(func(ctx context.Context, cm *corev1.ConfigMap) []reconcile.Request {
			return r.getComponentRequest(ctx, cm)
		}),
	)
}
