package controllers

import (
	"context"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewComponentUpgradeManager(t *testing.T) {
	t.Run("should create new componentUpgradeManager", func(t *testing.T) {
		mockComponentClient := newMockComponentInterface(t)
		mockHelmClient := newMockHelmClient(t)

		manager := NewComponentUpgradeManager(mockComponentClient, mockHelmClient, nil)

		assert.NotNil(t, manager)
		assert.Equal(t, mockHelmClient, manager.helmClient)
		assert.Equal(t, mockComponentClient, manager.componentClient)
	})
}

func Test_componentUpgradeManager_Upgrade(t *testing.T) {
	component := &k8sv1.Component{
		Spec: k8sv1.ComponentSpec{
			Namespace: "ecosystem",
			Name:      "testComponent",
			Version:   "1.0",
		},
		Status: k8sv1.ComponentStatus{Status: "installed"},
	}

	t.Run("should upgrade component", func(t *testing.T) {
		ctx := context.Background()

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctx, component).Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctx, component.GetHelmChartSpec()).Return(nil)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.NoError(t, err)
	})

	t.Run("dependency check failed", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(assert.AnError)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, "Warning", "Upgrade", "Dependency check failed: %s", assert.AnError.Error()).Return()

		sut := componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			recorder:        mockRecorder,
		}

		// when
		err := sut.Upgrade(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to check dependencies")
	})

	t.Run("should fail to upgrade component on error while setting upgrading status", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-upgrading for component testComponent:")
	})

	t.Run("should fail to upgrade component on error while upgrading chart", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctx, component.GetHelmChartSpec()).Return(assert.AnError)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to upgrade chart for component testComponent:")
	})

	t.Run("should fail to upgrade component on error while setting installed status", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctx, component).Return(component, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, component.GetHelmChartSpec()).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctx, component.GetHelmChartSpec()).Return(nil)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-installed for component testComponent:")
	})
}
