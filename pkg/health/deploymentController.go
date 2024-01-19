package health

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type deploymentReconciler struct {
	clientSet ecosystemClientSet
	manager   ComponentManager
}

func (dr *deploymentReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	deployment, err := dr.clientSet.
		AppsV1().Deployments(request.Namespace).
		Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		return finishOrRequeue(logger,
			client.IgnoreNotFound(
				fmt.Errorf("failed to get deployment %q: %w", request.NamespacedName, err),
			),
		)
	}

	// we know that this object belongs to a component since healthChangeEventFilter checked that for us already
	componentName := componentName(deployment)
	logger.Info(fmt.Sprintf("Found deployment %q for component %q", deployment.Name, componentName))

	err = dr.manager.UpdateComponentHealth(ctx, componentName, request.Namespace)
	if err != nil {
		return finishOrRequeue(logger, fmt.Errorf("failed to update component health for deployment %q: %w", request.NamespacedName, err))
	}

	return finishOperation()
}

func (dr *deploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(&healthChangeEventFilter{}).
		Complete(dr)
}
