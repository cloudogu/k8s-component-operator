package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponentInstallManager(t *testing.T) {
	// when
	manager := NewComponentInstallManager(nil, nil, nil, nil)

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
		mockComponentClient.EXPECT().UpdateStatusInstalled(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("0.1.0", nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, component.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
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
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, component.GetHelmChartSpec()).Return(assert.AnError)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
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
		mockComponentClient.EXPECT().UpdateStatusInstalled(testCtx, component).Return(component, assert.AnError)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("0.1.0", nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, component.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to update status-installed for component dogu-op")
	})

	t.Run("failed to update component health", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("0.1.0", nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, component.GetHelmChartSpec()).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(assert.AnError)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
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
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, componentWithoutVersion).Return(componentWithoutVersion, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(testCtx, componentWithoutVersion).Return(componentWithoutVersion, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, componentWithoutVersion, "component-finalizer").Return(componentWithoutVersion, nil)

		mockComponentClient.EXPECT().Get(testCtx, componentWithoutVersion.Name, metav1.GetOptions{}).Return(componentWithoutVersion, nil)
		componentWithVersion := getComponent(namespace, "k8s", "", "dogu-op", "4.8.3")
		mockComponentClient.EXPECT().Update(testCtx, componentWithVersion, metav1.UpdateOptions{}).Return(componentWithoutVersion, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("4.8.3", nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithVersion(testCtx, component.Spec.Name, namespace, "4.8.3").Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.NoError(t, err)
		assert.Equal(t, componentWithVersion.Spec.Version, componentWithoutVersion.Spec.Version)
	})

	t.Run("should fail to update version of component on error while listing releases", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, componentWithoutVersion).Return(componentWithoutVersion, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, componentWithoutVersion, "component-finalizer").Return(componentWithoutVersion, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("", assert.AnError)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to get release version for component")
	})

	t.Run("should fail to update version of component on error getting component", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, componentWithoutVersion).Return(componentWithoutVersion, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, componentWithoutVersion, "component-finalizer").Return(componentWithoutVersion, nil)

		mockComponentClient.EXPECT().Get(testCtx, componentWithoutVersion.Name, metav1.GetOptions{}).Return(componentWithoutVersion, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("4.8.3", nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to get component \"dogu-op\" for update")
		assert.ErrorContains(t, err, "failed to update version in component \"dogu-op\"")
	})

	t.Run("should fail to update version of component on error while updating", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, componentWithoutVersion).Return(componentWithoutVersion, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, componentWithoutVersion, "component-finalizer").Return(componentWithoutVersion, nil)

		mockComponentClient.EXPECT().Get(testCtx, componentWithoutVersion.Name, metav1.GetOptions{}).Return(componentWithoutVersion, nil)
		componentWithVersion := getComponent(namespace, "k8s", "", "dogu-op", "4.8.3")
		mockComponentClient.EXPECT().Update(testCtx, componentWithVersion, metav1.UpdateOptions{}).Return(nil, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, componentWithoutVersion.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().GetReleaseVersion(testCtx, component.Spec.Name).Return("4.8.3", nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to update version in component \"dogu-op\"")
	})
}
