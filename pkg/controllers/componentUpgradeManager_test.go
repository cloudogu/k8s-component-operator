package controllers

import (
	"context"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewComponentUpgradeManager(t *testing.T) {
	t.Run("should create new componentUpgradeManager", func(t *testing.T) {
		mockComponentClient := NewMockComponentClient(t)
		mockHelmClient := NewMockHelmClient(t)

		manager := NewComponentUpgradeManager(mockComponentClient, mockHelmClient)

		assert.NotNil(t, manager)
		assert.Equal(t, mockHelmClient, manager.helmClient)
		assert.Equal(t, mockComponentClient, manager.componentClient)
	})
}

func Test_componentUpgradeManager_Upgrade(t *testing.T) {
	t.Run("should upgrade component", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctx, component).Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UpdateChartRepos().Return(nil)
		mockHelmClient.EXPECT().UpgradeChart(ctx, component.GetHelmChartSpec(), (*helmclient.GenericHelmOptions)(nil)).Return(nil, nil)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.NoError(t, err)
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

		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, assert.AnError)

		mockHelmClient := NewMockHelmClient(t)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-upgrading for component testComponent:")
	})

	t.Run("should fail to upgrade component on error while updating chart-repos", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UpdateChartRepos().Return(assert.AnError)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update chart repositories:")
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

		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UpdateChartRepos().Return(nil)
		mockHelmClient.EXPECT().UpgradeChart(ctx, component.GetHelmChartSpec(), (*helmclient.GenericHelmOptions)(nil)).Return(nil, assert.AnError)

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

		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctx, component).Return(component, assert.AnError)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UpdateChartRepos().Return(nil)
		mockHelmClient.EXPECT().UpgradeChart(ctx, component.GetHelmChartSpec(), (*helmclient.GenericHelmOptions)(nil)).Return(nil, nil)

		manager := &componentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-installed for component testComponent:")
	})
}
