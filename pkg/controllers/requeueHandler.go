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

// Handle takes an error and handles the requeue process for the current component operation. If the error is
// requeueable, it is returned unchanged so controller-runtime re-enqueues the request through the workqueue rate
// limiter (exponential backoff). A non-requeueable error finishes the operation without requeueing.
func (d *componentRequeueHandler) Handle(ctx context.Context, contextMessage string, component *k8sv1.Component, originalErr error, requeueStatus string) (ctrl.Result, error) {
	if !shouldRequeue(originalErr) {
		return ctrl.Result{}, nil
	}

	updateError := retry.OnConflict(func() error {
		compClient := d.clientSet.ComponentV1Alpha1().Components(d.namespace)

		updatedComponent, err := compClient.Get(ctx, component.GetName(), metav1.GetOptions{})
		if err != nil {
			return err
		}

		updatedComponent.Status.Status = requeueStatus
		_, err = compClient.UpdateStatus(ctx, updatedComponent, metav1.UpdateOptions{})
		return err
	})
	if updateError != nil {
		d.recorder.Eventf(component, v1.EventTypeWarning, RequeueEventReason, "Failed to requeue component %s.", component.GetName())
		return ctrl.Result{}, fmt.Errorf("failed to update component status while requeueing: %w", updateError)
	}

	d.fireRequeueEvent(component)

	log.FromContext(ctx).Info(fmt.Sprintf("%s: requeueing because of: %s", contextMessage, originalErr.Error()))

	return ctrl.Result{}, originalErr
}

// shouldRequeue reports whether the given error is a requeueable error, i.e. an error that should trigger a
// rate-limited requeue of the component rather than finishing the operation.
func shouldRequeue(err error) bool {
	var requeueableError *genericRequeueableError
	return errors.As(err, &requeueableError)
}

func (d *componentRequeueHandler) fireRequeueEvent(component *k8sv1.Component) {
	d.recorder.Eventf(component, v1.EventTypeNormal, RequeueEventReason, "Requeueing component %s with backoff.", component.GetName())
}
