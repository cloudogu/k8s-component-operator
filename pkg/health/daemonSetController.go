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

type daemonSetReconciler struct {
	clientSet ecosystemClientSet
	manager   componentManager
}

func (dsr *daemonSetReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	daemonSet, err := dsr.clientSet.
		AppsV1().DaemonSets(request.Namespace).
		Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		return finishOrRequeue(logger,
			client.IgnoreNotFound(
				fmt.Errorf("failed to get daemon set %q: %w", request.NamespacedName, err),
			),
		)
	}

	if !hasComponentLabel(daemonSet) {
		// ignore non component daemon sets
		return finishOperation()
	}
	componentName := componentName(daemonSet)
	logger.Info(fmt.Sprintf("Found daemon set %q for component %q", daemonSet.Name, componentName))

	err = dsr.manager.UpdateComponentHealth(ctx, componentName, request.Namespace)
	if err != nil {
		return finishOrRequeue(logger, fmt.Errorf("failed to update component health for daemon set %q: %w", request.NamespacedName, err))
	}

	return finishOperation()
}

func (dsr *daemonSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		Complete(dsr)
}
