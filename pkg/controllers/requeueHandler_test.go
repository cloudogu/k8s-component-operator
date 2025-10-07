package controllers

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
)

var testCtx = context.Background()

func Test_componentRequeueHandler_Handle(t *testing.T) {
	t.Run("should exit early if there is no error", func(t *testing.T) {
		// given
		sut := &componentRequeueHandler{}
		var originalErr error = nil
		component := &v1.Component{}

		// when
		actual, err := sut.Handle(testCtx, "", component, originalErr, "installing")

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{Requeue: false, RequeueAfter: 0}, actual)
	})
	t.Run("should exit early if error is not requeuable", func(t *testing.T) {
		// given
		sut := &componentRequeueHandler{}
		var originalErr = assert.AnError
		component := &v1.Component{}

		// when
		actual, err := sut.Handle(testCtx, "", component, originalErr, "installing")

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{Requeue: false, RequeueAfter: 0}, actual)
	})
	t.Run("should fail to update component status", func(t *testing.T) {
		// given
		component := createComponent("k8s-dogu-operator", "official", "1.2.3")

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, component.Name, mock.Anything).Return(component, nil)
		componentInterfaceMock.EXPECT().UpdateStatus(testCtx, component, metav1.UpdateOptions{}).Return(nil, assert.AnError)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		sut := &componentRequeueHandler{namespace: testNamespace, clientSet: clientSetMock}

		requeueErrMock := newMockRequeuableError(t)
		requeueErrMock.EXPECT().GetRequeueTime(mock.Anything).Return(30 * time.Second)

		// when
		actual, err := sut.Handle(testCtx, "", component, requeueErrMock, "upgrading")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update component status")

		assert.Equal(t, reconcile.Result{Requeue: false, RequeueAfter: 0}, actual)
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		component := createComponent("k8s-dogu-operator", "official", "1.2.3")

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, component.Name, mock.Anything).Return(component, nil)
		componentInterfaceMock.EXPECT().UpdateStatus(testCtx, component, metav1.UpdateOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		recorderMock := newMockEventRecorder(t)
		recorderMock.EXPECT().Eventf(component, "Normal", "Requeue", "Falling back to component status %s: Trying again in %s.", "upgrading", "1s")

		sut := &componentRequeueHandler{namespace: testNamespace, clientSet: clientSetMock, recorder: recorderMock}

		requeueErrMock := newMockRequeuableError(t)
		requeueErrMock.EXPECT().GetRequeueTime(mock.Anything).Return(time.Second)
		requeueErrMock.EXPECT().Error().Return("my error")

		// when
		actual, err := sut.Handle(testCtx, "", component, requeueErrMock, "upgrading")

		// then
		require.NoError(t, err)

		assert.Equal(t, reconcile.Result{Requeue: true, RequeueAfter: 1000000000}, actual)
	})
}

func createComponent(name, namespace, version string) *v1.Component {
	return &v1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.ComponentSpec{
			Namespace: namespace,
			Name:      name,
			Version:   version,
		},
	}
}

func Test_componentRequeueHandler_noLongerHandleRequeueing(t *testing.T) {
	t.Run("reset requeue time to avoid further requeueing", func(t *testing.T) {
		// given
		finishedComponent := &v1.Component{Status: v1.ComponentStatus{
			Status:           "installed",
			RequeueTimeNanos: 3000}}

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, finishedComponent.Name, mock.Anything).Return(finishedComponent, nil)
		componentInterfaceMock.EXPECT().UpdateStatus(testCtx, finishedComponent, metav1.UpdateOptions{}).Return(finishedComponent, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		sut := &componentRequeueHandler{namespace: testNamespace, clientSet: clientSetMock}

		// when
		actual, err := sut.noLongerHandleRequeueing(testCtx, finishedComponent)

		// then
		require.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, actual)
		assert.Equal(t, time.Duration(0), finishedComponent.Status.RequeueTimeNanos)
	})
}
