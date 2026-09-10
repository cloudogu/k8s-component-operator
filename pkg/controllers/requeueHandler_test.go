package controllers

import (
	"context"
	"testing"

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

		recorderMock := newMockEventRecorder(t)
		recorderMock.EXPECT().Eventf(component, "Warning", "Requeue", "Failed to requeue component %s.", component.GetName())

		sut := &componentRequeueHandler{namespace: testNamespace, clientSet: clientSetMock, recorder: recorderMock}

		requeueErr := &genericRequeueableError{"failed to upgrade", assert.AnError}

		// when
		actual, err := sut.Handle(testCtx, "", component, requeueErr, "upgrading")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update component status")

		assert.Equal(t, reconcile.Result{Requeue: false, RequeueAfter: 0}, actual)
	})
	t.Run("should return the original error so controller-runtime requeues with backoff", func(t *testing.T) {
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
		recorderMock.EXPECT().Eventf(component, "Normal", "Requeue", "Requeueing component %s with backoff.", component.GetName())

		sut := &componentRequeueHandler{namespace: testNamespace, clientSet: clientSetMock, recorder: recorderMock}

		requeueErr := &genericRequeueableError{"my error", assert.AnError}

		// when
		actual, err := sut.Handle(testCtx, "", component, requeueErr, "upgrading")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, requeueErr)
		assert.Equal(t, "upgrading", component.Status.Status)
		assert.Equal(t, reconcile.Result{Requeue: false, RequeueAfter: 0}, actual)
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
