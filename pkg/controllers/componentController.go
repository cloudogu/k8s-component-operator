package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/internal"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type operation string

const (
	Install = operation("Install")
	Upgrade = operation("Upgrade")
	Delete  = operation("Delete")
	Ignore  = operation("Ignore")
)

// componentReconciler watches every Component object in the cluster and handles them accordingly.
type componentReconciler struct {
	client          *ecosystem.EcosystemClientset
	recorder        record.EventRecorder
	componentManger internal.ComponentManager
}

// NewComponentReconciler creates a new component reconciler.
func NewComponentReconciler(client *ecosystem.EcosystemClientset, recorder record.EventRecorder) *componentReconciler {
	return &componentReconciler{
		client:   client,
		recorder: recorder,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *componentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile this crd")

	component, err := r.client.EcosystemV1Alpha1().Components(req.Namespace).Get(ctx, req.Name, v1.GetOptions{})
	if err != nil {
		logger.Info(fmt.Sprintf("failed to get component %+v: %s", req, err))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info(fmt.Sprintf("Component %+v has been found", req))

	operation, err := r.evaluateRequiredOperation(ctx, component)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to evaluate required operation: %w", err)
	}

	switch operation {
	case Install:
		err = r.componentManger.Install(ctx, component)
		return ctrl.Result{}, err
	case Ignore:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, nil
	}
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
		return Ignore, nil
	case k8sv1.ComponentStatusInstalling:
		return Ignore, nil
	case k8sv1.ComponentStatusDeleting:
		return Ignore, nil
	default:
		logger.Info(fmt.Sprintf("Found unknown operation for component status: %s", component.Status.Status))
		return Ignore, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *componentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1.Component{}).
		Complete(r)
}
