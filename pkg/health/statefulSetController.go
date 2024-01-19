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

type statefulSetReconciler struct {
	clientSet ecosystemClientSet
	manager   ComponentManager
}

func (ssr *statefulSetReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx)

	statefulSet, err := ssr.clientSet.
		AppsV1().StatefulSets(request.Namespace).
		Get(ctx, request.Name, metav1.GetOptions{})
	if err != nil {
		return finishOrRequeue(logger,
			client.IgnoreNotFound(
				fmt.Errorf("failed to get stateful set %q: %w", request.NamespacedName, err),
			),
		)
	}

	// we know that this object belongs to a component since healthChangeEventFilter checked that for us already
	componentName := componentName(statefulSet)
	logger.Info(fmt.Sprintf("Found stateful set %q for component %q", statefulSet.Name, componentName))

	err = ssr.manager.UpdateComponentHealth(ctx, componentName, request.Namespace)
	if err != nil {
		return finishOrRequeue(logger, fmt.Errorf("failed to update component health for stateful set %q: %w", request.NamespacedName, err))
	}

	return finishOperation()
}

func (ssr *statefulSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.StatefulSet{}).
		WithEventFilter(&healthChangeEventFilter{}).
		Complete(ssr)
}
