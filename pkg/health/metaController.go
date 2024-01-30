package health

import (
	"errors"
	"fmt"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"
)

type metaController struct {
	controllers []RegistrableController
}

func (m metaController) SetupWithManager(mgr ctrl.Manager) error {
	var errs []error
	for _, controller := range m.controllers {
		err := controller.SetupWithManager(mgr)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to setup controller %q: %w", reflect.TypeOf(controller), err))
		}
	}

	return errors.Join(errs...)
}

func NewController(namespace string, clientSet ecosystemClientSet) RegistrableController {
	manager := NewManager(namespace, clientSet)
	return &metaController{controllers: []RegistrableController{
		&deploymentReconciler{clientSet: clientSet, manager: manager},
		&statefulSetReconciler{clientSet: clientSet, manager: manager},
		&daemonSetReconciler{clientSet: clientSet, manager: manager},
	}}
}
