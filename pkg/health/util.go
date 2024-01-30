package health

import (
	"github.com/go-logr/logr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

func componentName(object client.Object) string {
	return object.GetLabels()[v1.ComponentNameLabelKey]
}

func hasComponentLabel(object client.Object) bool {
	_, exists := object.GetLabels()[v1.ComponentNameLabelKey]
	return exists
}

func finishOrRequeue(logger logr.Logger, err error) (ctrl.Result, error) {
	if err != nil {
		logger.Error(err, "reconcile failed")
	}

	return ctrl.Result{}, err
}

func finishOperation() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}
