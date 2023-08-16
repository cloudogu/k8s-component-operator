package controllers

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

// componentRequeueHandler is responsible to requeue a dogu resource after it failed.
type componentRequeueHandler struct {
	clientSet componentEcosystemInterface
	namespace string
	recorder  record.EventRecorder
}

// NewComponentRequeueHandler creates a new dogu requeue handler.
func NewComponentRequeueHandler(clientSet componentEcosystemInterface, recorder record.EventRecorder, namespace string) *componentRequeueHandler {
	return &componentRequeueHandler{
		clientSet: clientSet,
		namespace: namespace,
		recorder:  recorder,
	}
}

// Handle takes an error and handles the requeue process for the current dogu operation.
func (d *componentRequeueHandler) Handle(ctx context.Context, contextMessage string, component *k8sv1.Component, originalErr error, onRequeue func()) (ctrl.Result, error) {
	requeueable, requeueableErr := shouldRequeue(originalErr)
	if !requeueable {
		return ctrl.Result{}, nil
	}
	if onRequeue != nil {
		onRequeue()
	}

	_, updateError := d.clientSet.ComponentV1Alpha1().Components(d.namespace).UpdateStatus(ctx, component, metav1.UpdateOptions{})
	if updateError != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update component status: %w", updateError)
	}

	requeueTime := requeueableErr.GetRequeueTime(component.Status.RequeueTimeNanos)
	result := ctrl.Result{Requeue: true, RequeueAfter: requeueTime}
	d.fireRequeueEvent(component, result)

	log.FromContext(ctx).Info(fmt.Sprintf("%s: requeue in %s seconds because of: %s", contextMessage, requeueTime, originalErr.Error()))

	return result, nil
}

func shouldRequeue(err error) (bool, requeuableError) {
	var requeueableError requeuableError
	return errors.As(err, &requeueableError), requeueableError
}

func (d *componentRequeueHandler) fireRequeueEvent(component *k8sv1.Component, result ctrl.Result) {
	d.recorder.Eventf(component, v1.EventTypeNormal, RequeueEventReason, "Trying again in %s.", result.RequeueAfter.String())
}
