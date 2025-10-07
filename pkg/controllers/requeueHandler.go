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

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/retry-lib/retry"
)

// componentRequeueHandler is responsible to requeue a component resource after it failed.
type componentRequeueHandler struct {
	clientSet componentEcosystemInterface
	namespace string
	recorder  record.EventRecorder
}

// NewComponentRequeueHandler creates a new component requeue handler.
func NewComponentRequeueHandler(clientSet componentEcosystemInterface, recorder record.EventRecorder, namespace string) *componentRequeueHandler {
	return &componentRequeueHandler{
		clientSet: clientSet,
		namespace: namespace,
		recorder:  recorder,
	}
}

// Handle takes an error and handles the requeue process for the current component operation.
func (d *componentRequeueHandler) Handle(ctx context.Context, contextMessage string, component *k8sv1.Component, originalErr error, requeueStatus string) (ctrl.Result, error) {
	requeueable, requeueableErr := shouldRequeue(originalErr)
	if !requeueable {
		return d.noLongerHandleRequeueing(ctx, component)
	}

	requeueTime := requeueableErr.GetRequeueTime(component.Status.RequeueTimeNanos)

	updateError := retry.OnConflict(func() error {
		compClient := d.clientSet.ComponentV1Alpha1().Components(d.namespace)

		updatedComponent, err := compClient.Get(ctx, component.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}

		updatedComponent.Status.Status = requeueStatus
		updatedComponent.Status.RequeueTimeNanos = requeueTime
		component, err = compClient.UpdateStatus(ctx, updatedComponent, metav1.UpdateOptions{})
		return err
	})
	if updateError != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update component status while requeueing: %w", updateError)
	}

	result := ctrl.Result{Requeue: true, RequeueAfter: requeueTime}
	d.fireRequeueEvent(component, result)

	log.FromContext(ctx).Info(fmt.Sprintf("%s: requeue in %s seconds because of: %s", contextMessage, requeueTime, originalErr.Error()))

	return result, nil
}

// noLongerHandleRequeueing returns values so the component will no longer be requeued. This will occur either on a
// successful reconciliation or errors which cannot be handled and thus not be requeued. The component may reset the
// requeue backoff if necessary in order to avoid a wrong backoff baseline time for future reconciliations.
func (d *componentRequeueHandler) noLongerHandleRequeueing(ctx context.Context, component *k8sv1.Component) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if component.Status.RequeueTimeNanos == 0 {
		logger.Info("Skipping backoff time reset")
		return ctrl.Result{}, nil
	}

	compClient := d.clientSet.ComponentV1Alpha1().Components(d.namespace)

	err := retry.OnConflict(func() error {
		updatedComponent, err := compClient.Get(ctx, component.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}

		logger.Info("Reset backoff time to 0")
		updatedComponent.Status.RequeueTimeNanos = 0
		_, err = compClient.UpdateStatus(ctx, updatedComponent, metav1.UpdateOptions{})
		return err
	})

	return ctrl.Result{}, err
}

func shouldRequeue(err error) (bool, requeuableError) {
	var requeueableError requeuableError
	return errors.As(err, &requeueableError), requeueableError
}

func (d *componentRequeueHandler) fireRequeueEvent(component *k8sv1.Component, result ctrl.Result) {
	d.recorder.Eventf(component, v1.EventTypeNormal, RequeueEventReason, "Falling back to component status %s: Trying again in %s.", component.Status.Status, result.RequeueAfter.String())
}
