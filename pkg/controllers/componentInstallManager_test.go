package controllers

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

	t.Run("success", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgrade(context.TODO(), component).Return(nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.NoError(t, err)
	})

	t.Run("failed to update installing status", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(nil, assert.AnError)

		mockHelmClient := NewMockHelmClient(t)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set status installing")
	})

	t.Run("failed to add finalizer", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(nil, assert.AnError)

		mockHelmClient := NewMockHelmClient(t)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to add finalizer component-finalizer")
	})

	t.Run("failed to install the chart", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgrade(context.TODO(), component).Return(assert.AnError)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to install chart")
	})

	t.Run("failed set status installed", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(context.TODO(), component).Return(component, assert.AnError)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgrade(context.TODO(), component).Return(nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-installed for component dogu-op")
	})
}
