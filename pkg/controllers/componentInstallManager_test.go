package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComponentInstallManager(t *testing.T) {
	// when
	manager := NewComponentInstallManager(nil, nil, nil)

	// then
	require.NotNil(t, manager)
}

func Test_componentInstallManager_Install(t *testing.T) {
	namespace := "ecosystem"
	component := getComponent(namespace, "k8s", "dogu-op", "0.1.0")

	t.Run("should install component", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, component.GetHelmChartSpec()).Return(nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
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

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			recorder:        mockRecorder,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check dependencies")
	})

	t.Run("failed to update installing status", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(nil, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set status installing")
	})

	t.Run("failed to add finalizer", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(nil, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
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

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
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
		mockHelmClient.EXPECT().InstallOrUpgrade(testCtx, component.GetHelmChartSpec()).Return(nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-installed for component dogu-op")
	})
}
