package controllers

import (
	"context"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentController watches every Component object in the cluster and handles them accordingly.
type componentController struct {
	client   client.Client
	recorder record.EventRecorder
}

// NewComponentController creates a new component reconciler.
func NewComponentController(client client.Client, recorder record.EventRecorder) *componentController {
	return &componentController{
		client:   client,
		recorder: recorder,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *componentController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile this crd")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *componentController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1.Component{}).
		Complete(r)
}
