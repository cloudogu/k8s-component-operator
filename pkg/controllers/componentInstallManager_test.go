package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctxWithoutCancel = context.WithoutCancel(testCtx)
var defaultHelmClientTimeoutMins = 15 * time.Minute

func TestNewComponentInstallManager(t *testing.T) {
	// when
	manager := NewComponentInstallManager(nil, nil, nil, nil, defaultHelmClientTimeoutMins)

	// then
	require.NotNil(t, manager)
}

func Test_componentInstallManager_Install(t *testing.T) {
	namespace := "ecosystem"
	component := getComponent(namespace, "k8s", "", "dogu-op", "0.1.0")

	t.Run("should install component", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, component.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.NoError(t, err)
	})

	t.Run("dependency check failed", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(assert.AnError)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, "Warning", "Installation", "Dependency check failed: %s", assert.AnError.Error()).Return()

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			recorder:        mockRecorder,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to check dependencies")
	})

	t.Run("failed to update installing status", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(nil, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to set status installing")
	})

	t.Run("failed to add finalizer", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(nil, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to add finalizer component-finalizer")
	})

	t.Run("failed to install the chart", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, component.GetHelmChartSpec()).Return(assert.AnError)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to install chart")
	})

	t.Run("failed set status installed", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, component).Return(component, assert.AnError)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, component.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to update status-installed for component \"dogu-op\"")
	})

	t.Run("failed to update component health", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, component.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(assert.AnError)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update health status and installed version for component")
	})

	t.Run("should update version of component", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		componentWithVersion := getComponent(namespace, "k8s", "", "dogu-op", "4.8.3")
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, componentWithVersion).Return(componentWithVersion, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, componentWithVersion).Return(componentWithVersion, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, componentWithVersion, "component-finalizer").Return(componentWithVersion, nil)
		mockComponentClient.EXPECT().UpdateExpectedComponentVersion(testCtx, componentWithVersion.Name, componentWithVersion.Spec.Version).Return(componentWithVersion, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetLatestVersion("k8s/dogu-op").Return("4.8.3", nil)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, componentWithVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, componentWithVersion.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, componentWithVersion.Spec.Version).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to update version of component on error while getting the latest version", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		mockComponentClient := newMockComponentInterface(t)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetLatestVersion("k8s/dogu-op").Return("", assert.AnError)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to get latest version for component \"dogu-op\"")
	})

	t.Run("should fail to update expected version of component", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateExpectedComponentVersion(testCtx, component.Name, "4.8.3").Return(componentWithoutVersion, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetLatestVersion("k8s/dogu-op").Return("4.8.3", nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to update expected version for component \"dogu-op\"")
	})
}
