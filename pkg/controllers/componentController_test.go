package controllers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
	"time"
)

func TestNewComponentReconciler(t *testing.T) {
	// given
	mockComponentClient := NewMockComponentClient(t)
	mockHelmClient := NewMockHelmClient(t)
	mockRecorder := NewMockEventRecorder(t)

	// when
	manager := NewComponentReconciler(mockComponentClient, mockHelmClient, mockRecorder)

	// then
	require.NotNil(t, manager)
}

func Test_componentReconciler_Reconcile(t *testing.T) {
	namespace := "ecosystem"
	helmNamespace := "k8s"
	t.Run("success install", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		mockRecorder := NewMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Installation", "Installation successful")
		manager := NewMockComponentManager(t)
		manager.EXPECT().Install(context.TODO(), component).Return(nil)
		helmClient := NewMockHelmClient(t)
		sut := componentReconciler{
			componentClient:  mockComponentClient,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(context.TODO(), req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("success delete", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		component.DeletionTimestamp = &v1.Time{Time: time.Now()}
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		mockRecorder := NewMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Deinstallation", "Deinstallation successful")
		manager := NewMockComponentManager(t)
		manager.EXPECT().Delete(context.TODO(), component).Return(nil)
		helmClient := NewMockHelmClient(t)
		sut := componentReconciler{
			componentClient:  mockComponentClient,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(context.TODO(), req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("success upgrade", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		component.Status.Status = "installed"
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		mockRecorder := NewMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Upgrade", "Upgrade successful")
		manager := NewMockComponentManager(t)
		manager.EXPECT().Upgrade(context.TODO(), component).Return(nil)
		helmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: namespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		sut := componentReconciler{
			componentClient:  mockComponentClient,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(context.TODO(), req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("should fail on downgrade", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		component.Status.Status = "installed"
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		mockRecorder := NewMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Warning", "Downgrade", "component downgrades are not allowed")
		manager := NewMockComponentManager(t)
		helmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: namespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.2.0"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		sut := componentReconciler{
			componentClient:  mockComponentClient,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(context.TODO(), req)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "downgrades are not allowed")
	})

	t.Run("should ignore equal installed component", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		component.Status.Status = "installed"
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		mockRecorder := NewMockEventRecorder(t)
		manager := NewMockComponentManager(t)
		helmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: namespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.1.0"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		sut := componentReconciler{
			componentClient:  mockComponentClient,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(context.TODO(), req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("should fail on component get error", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(nil, assert.AnError)
		sut := componentReconciler{
			componentClient: mockComponentClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(context.TODO(), req)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil if the component is not found", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		sut := componentReconciler{
			componentClient: mockComponentClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(context.TODO(), req)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail on getting operation with invalid versions", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		component.Status.Status = "installed"
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		helmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: namespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "invalidsemver"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		sut := componentReconciler{
			componentClient: mockComponentClient,
			helmClient:      helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(context.TODO(), req)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to evaluate required operation")
	})

	t.Run("should fail on error in operation", func(t *testing.T) {
		// given
		component := getComponent(namespace, helmNamespace, "dogu-op", "0.1.0")
		component.DeletionTimestamp = &v1.Time{Time: time.Now()}
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().Get(context.TODO(), "dogu-op", v1.GetOptions{}).Return(component, nil)
		mockRecorder := NewMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Warning", "Deinstallation", "Deinstallation failed. Reason: assert.AnError general error for testing")
		manager := NewMockComponentManager(t)
		manager.EXPECT().Delete(context.TODO(), component).Return(assert.AnError)
		helmClient := NewMockHelmClient(t)
		sut := componentReconciler{
			componentClient:  mockComponentClient,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: namespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(context.TODO(), req)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_componentReconciler_getChangeOperation(t *testing.T) {
	t.Run("should fail on error getting helm releases", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.1.0")
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().ListDeployedReleases().Return(nil, assert.AnError)

		sut := componentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		_, err := sut.getChangeOperation(component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get deployed helm releases")
	})

	t.Run("should fail on error parsing component version", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "notvalidsemver")
		mockHelmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := componentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		_, err := sut.getChangeOperation(component)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse component version")
	})

	t.Run("should return downgrade-operation on downgrade", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.0")
		mockHelmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := componentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Downgrade, op)
	})

	t.Run("should return upgrade-operation on upgrade", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.2")
		mockHelmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := componentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return ignore-operation on same version", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.1")
		mockHelmClient := NewMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := componentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation when no release is found", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.1")
		mockHelmClient := NewMockHelmClient(t)
		var helmReleases []*release.Release
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := componentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})
}

func Test_componentReconciler_evaluateRequiredOperation(t *testing.T) {
	t.Run("should return ignore on status installing", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.0")
		component.Status.Status = "installing"
		sut := componentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(context.TODO(), component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return ignore on status deleting", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.0")
		component.Status.Status = "deleting"
		sut := componentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(context.TODO(), component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return ignore on status upgrading", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.0")
		component.Status.Status = "upgrading"
		sut := componentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(context.TODO(), component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return ignore on unrecognized status", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.0.0")
		component.Status.Status = "foobar"
		sut := componentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(context.TODO(), component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})
}
