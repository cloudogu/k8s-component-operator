package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"testing"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const testNamespace = "testtestNamespace"

func TestNewComponentReconciler(t *testing.T) {
	// given
	componentInterfaceMock := newMockComponentInterface(t)
	componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
	componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
	appsMock := newMockAppsV1Interface(t)
	clientSetMock := newMockComponentEcosystemInterface(t)
	clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)
	clientSetMock.EXPECT().AppsV1().Return(appsMock)

	mockHelmClient := newMockHelmClient(t)
	mockRecorder := newMockEventRecorder(t)

	// when
	manager := NewComponentReconciler(clientSetMock, mockHelmClient, mockRecorder, testNamespace)

	// then
	require.NotNil(t, manager)
}

func Test_componentReconciler_Reconcile(t *testing.T) {
	helmNamespace := "k8s"
	t.Run("success install", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Installation", "Installation successful")

		manager := NewMockComponentManager(t)
		manager.EXPECT().Install(testCtx, component).Return(nil)
		helmClient := newMockHelmClient(t)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Installation failed with component dogu-op", component, nil, mock.Anything).Return(reconcile.Result{}, nil)

		sut := ComponentReconciler{
			clientSet:        clientSetMock,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
			requeueHandler:   mockRequeueHandler,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("success delete", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.DeletionTimestamp = &v1.Time{Time: time.Now()}

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Deinstallation", "Deinstallation successful")

		manager := NewMockComponentManager(t)
		manager.EXPECT().Delete(testCtx, component).Return(nil)
		helmClient := newMockHelmClient(t)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Deinstallation failed with component dogu-op", component, nil, mock.Anything).Return(reconcile.Result{}, nil)

		sut := ComponentReconciler{
			clientSet:        clientSetMock,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
			requeueHandler:   mockRequeueHandler,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("success upgrade", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Upgrade", "Upgrade successful")

		manager := NewMockComponentManager(t)
		manager.EXPECT().Upgrade(testCtx, component).Return(nil)

		helmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: testNamespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Upgrade failed with component dogu-op", component, nil, mock.Anything).Return(reconcile.Result{}, nil)

		sut := ComponentReconciler{
			clientSet:        clientSetMock,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
			requeueHandler:   mockRequeueHandler,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("should fail on downgrade", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Warning", "Downgrade", "component downgrades are not allowed")

		manager := NewMockComponentManager(t)
		helmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: testNamespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.2.0"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			clientSet:        clientSetMock,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "downgrades are not allowed")
	})

	t.Run("should ignore equal installed component", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		manager := NewMockComponentManager(t)
		helmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: testNamespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.1.0"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		helmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{}, nil)
		helmClient.EXPECT().GetChartSpecValues(component.GetHelmChartSpec()).Return(map[string]interface{}{}, nil)

		sut := ComponentReconciler{
			clientSet:        clientSetMock,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("should fail on component get error", func(t *testing.T) {
		// given
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(nil, assert.AnError)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		sut := ComponentReconciler{
			clientSet: clientSetMock,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil if the component is not found", func(t *testing.T) {
		// given

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		sut := ComponentReconciler{
			clientSet: clientSetMock,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail on getting operation with invalid versions", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		helmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: testNamespace, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "invalidsemver"}}}}
		helmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			clientSet:  clientSetMock,
			helmClient: helmClient,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to evaluate required operation")
	})

	t.Run("should fail on error in operation", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.DeletionTimestamp = &v1.Time{Time: time.Now()}

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Warning", "Deinstallation", "Deinstallation failed. Reason: assert.AnError general error for testing")
		mockRecorder.EXPECT().Eventf(component, "Warning", "Requeue", "Failed to requeue the %s.", "deinstallation").Return()

		manager := NewMockComponentManager(t)
		manager.EXPECT().Delete(testCtx, component).Return(assert.AnError)
		helmClient := newMockHelmClient(t)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Deinstallation failed with component dogu-op", component, assert.AnError, mock.Anything).
			RunAndReturn(func(_ context.Context, _ string, _ *k8sv1.Component, err error, _ string) (reconcile.Result, error) {
				return reconcile.Result{}, err
			})

		sut := ComponentReconciler{
			clientSet:        clientSetMock,
			recorder:         mockRecorder,
			componentManager: manager,
			helmClient:       helmClient,
			requeueHandler:   mockRequeueHandler,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_componentReconciler_getChangeOperation(t *testing.T) {
	t.Run("should fail on error getting helm releases", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().ListDeployedReleases().Return(nil, assert.AnError)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get deployed helm releases")
	})

	t.Run("should fail on error parsing component version", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "notvalidsemver")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse component version")
	})

	t.Run("should fail on error getting release values", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(nil, assert.AnError)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to compare Values.yaml files of component")
		assert.ErrorContains(t, err, "failed to get values.yaml from release")
	})

	t.Run("should fail on error getting component values", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(component.GetHelmChartSpec()).Return(nil, assert.AnError)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to compare Values.yaml files of component")
		assert.ErrorContains(t, err, "failed to get values.yaml from component")
	})

	t.Run("should return downgrade-operation on downgrade", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Downgrade, op)
	})

	t.Run("should return upgrade-operation on upgrade if deploy namespace is set", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "deploy-namespace", "dogu-op", "0.0.1-2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "deploy-namespace", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1-1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return error if deploy namespace is not the same as release namespace", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "deploy-namespace", "dogu-op", "0.0.1-2")
		mockHelmClient := newMockHelmClient(t)
		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, corev1.EventTypeWarning, UpgradeEventReason, "Deploy namespace mismatch (CR: %q; deployed: %q). Deploy namespace declaration is only allowed on install. Revert deploy namespace change to prevent failing upgrade.", "deploy-namespace", "ecosystem").Return()
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1-2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
			recorder:   mockRecorder,
		}

		// when
		_, err := sut.getChangeOperation(testCtx, component)

		// then
		require.Error(t, err)
	})

	t.Run("should return upgrade-operation on upgrade", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return upgrade-operation on same version, but values-yaml difference", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{"foo": "bar", "baz": "buz"}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(component.GetHelmChartSpec()).Return(map[string]interface{}{"foo": "bar", "baz": "xyz"}, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, op)
	})

	t.Run("should return ignore-operation on same version and same values-yaml values", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{"foo": "bar", "baz": "buz"}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(component.GetHelmChartSpec()).Return(map[string]interface{}{"foo": "bar", "baz": "buz"}, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation on same version and different zero-length maps", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.2")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.2"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}(nil), nil)
		mockHelmClient.EXPECT().GetChartSpecValues(component.GetHelmChartSpec()).Return(map[string]interface{}{}, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation on same version", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		mockHelmClient := newMockHelmClient(t)
		helmReleases := []*release.Release{{Name: "dogu-op", Namespace: "ecosystem", Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.1"}}}}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)
		mockHelmClient.EXPECT().GetReleaseValues("dogu-op", false).Return(map[string]interface{}{}, nil)
		mockHelmClient.EXPECT().GetChartSpecValues(component.GetHelmChartSpec()).Return(map[string]interface{}{}, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})

	t.Run("should return ignore-operation when no release is found", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.1")
		mockHelmClient := newMockHelmClient(t)
		var helmReleases []*release.Release
		mockHelmClient.EXPECT().ListDeployedReleases().Return(helmReleases, nil)

		sut := ComponentReconciler{
			helmClient: mockHelmClient,
		}

		// when
		op, err := sut.getChangeOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, op)
	})
}

func Test_componentReconciler_evaluateRequiredOperation(t *testing.T) {
	t.Run("should return ignore on status installing", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "installing"
		sut := ComponentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return ignore on status deleting", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "deleting"
		sut := ComponentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return ignore on status upgrading", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "upgrading"
		sut := ComponentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return ignore on unrecognized status", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "foobar"
		sut := ComponentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Ignore, requiredOperation)
	})

	t.Run("should return install on tryToInstall status", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.0.0")
		component.Status.Status = "tryToInstall"
		sut := ComponentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Install, requiredOperation)
	})

	t.Run("should return upgrade on tryToUpgrade status", func(t *testing.T) {
		// given
		componentName := "dogu-op"
		component := getComponent("ecosystem", "k8s", "", componentName, "0.0.1")
		component.Status.Status = "tryToUpgrade"
		helmMock := newMockHelmClient(t)
		installedReleases := []*release.Release{{Namespace: "ecosystem", Name: componentName, Chart: &chart.Chart{Metadata: &chart.Metadata{AppVersion: "0.0.0"}}}}
		helmMock.EXPECT().ListDeployedReleases().Return(installedReleases, nil)
		sut := ComponentReconciler{helmClient: helmMock}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Upgrade, requiredOperation)
	})

	t.Run("should return delete on tryToDelete status", func(t *testing.T) {
		// given
		componentName := "dogu-op"
		component := getComponent("ecosystem", "k8s", "", componentName, "0.0.0")
		component.Status.Status = "tryToDelete"
		timeNow := v1.NewTime(time.Now())
		component.DeletionTimestamp = &timeNow
		sut := ComponentReconciler{}

		// when
		requiredOperation, err := sut.evaluateRequiredOperation(testCtx, component)

		// then
		require.NoError(t, err)
		assert.Equal(t, Delete, requiredOperation)
	})
}

func Test_componentReconciler_validateName(t *testing.T) {
	tests := []struct {
		name         string
		recorderFunc func(t *testing.T) record.EventRecorder
		component    *k8sv1.Component
		wantSuccess  bool
	}{
		{
			name: "should fail validation",
			recorderFunc: func(t *testing.T) record.EventRecorder {
				recorder := newMockEventRecorder(t)
				recorder.EXPECT().Eventf(mock.Anything, "Warning", "FailedNameValidation", "Component resource does not follow naming rules: The component's metadata.name '%s' must be the same as its spec.name '%s'.", "example", "invalid-example")
				return recorder
			},
			component:   &k8sv1.Component{ObjectMeta: v1.ObjectMeta{Name: "example"}, Spec: k8sv1.ComponentSpec{Name: "invalid-example"}},
			wantSuccess: false,
		},
		{
			name:         "should succeed validation",
			recorderFunc: func(t *testing.T) record.EventRecorder { return newMockEventRecorder(t) },
			component:    &k8sv1.Component{ObjectMeta: v1.ObjectMeta{Name: "example"}, Spec: k8sv1.ComponentSpec{Name: "example"}},
			wantSuccess:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ComponentReconciler{recorder: tt.recorderFunc(t)}
			success := r.validateName(tt.component)
			assert.Equal(t, tt.wantSuccess, success)
		})
	}
}
