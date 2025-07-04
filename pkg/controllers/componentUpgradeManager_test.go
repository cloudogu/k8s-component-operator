package controllers

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

func TestNewComponentUpgradeManager(t *testing.T) {
	t.Run("should create new ComponentUpgradeManager", func(t *testing.T) {
		mockComponentClient := newMockComponentInterface(t)
		mockHelmClient := newMockHelmClient(t)

		manager := NewComponentUpgradeManager(mockComponentClient, mockHelmClient, nil, nil, defaultHelmClientTimeoutMins)

		assert.NotNil(t, manager)
		assert.Equal(t, mockHelmClient, manager.helmClient)
		assert.Equal(t, mockComponentClient, manager.componentClient)
	})
}

func Test_componentUpgradeManager_Upgrade(t *testing.T) {
	component := &k8sv1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testComponent",
			Namespace: "ecosystem",
		},
		Spec: k8sv1.ComponentSpec{
			Namespace: "ecosystem",
			Name:      "testComponent",
			Version:   "0.1.0",
		},
		Status: k8sv1.ComponentStatus{Status: "installed"},
	}

	t.Run("should upgrade component", func(t *testing.T) {
		ctx := context.Background()

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(mock.Anything, component).Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart(mock.Anything, mock.Anything).Return(&chart.Chart{}, nil)
		spec, _ := component.GetHelmChartSpec(testCtx, k8sv1.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(mock.Anything, component.Spec.Name, "ecosystem", "0.1.0").Return(nil)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			healthManager:   mockHealthManager,
			timeout:         defaultHelmClientTimeoutMins,
		}
		err := manager.Upgrade(ctx, component)

		require.NoError(t, err)
	})

	t.Run("dependency check failed", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart(mock.Anything, mock.Anything).Return(&chart.Chart{}, nil)
		spec, _ := component.GetHelmChartSpec(testCtx, k8sv1.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(assert.AnError)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, "Warning", "Upgrade", "Dependency check failed: %s", assert.AnError.Error()).Return()

		sut := ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			recorder:        mockRecorder,
			timeout:         defaultHelmClientTimeoutMins,
		}

		// when
		err := sut.Upgrade(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
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
		mockHelmClient.EXPECT().GetChart(mock.Anything, mock.Anything).Return(&chart.Chart{}, nil)
		spec, _ := component.GetHelmChartSpec(testCtx, k8sv1.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
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
		mockHelmClient.EXPECT().GetChart(mock.Anything, mock.Anything).Return(&chart.Chart{}, nil)
		spec, _ := component.GetHelmChartSpec(testCtx, k8sv1.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(assert.AnError)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to upgrade chart for component testComponent:")
	})

	t.Run("should fail to upgrade component on error while setting installed status", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testComponent",
				Namespace: "ecosystem",
			},
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(mock.Anything, component).Return(component, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart(mock.Anything, mock.Anything).Return(&chart.Chart{}, nil)
		spec, _ := component.GetHelmChartSpec(testCtx, k8sv1.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(nil)

		mockHealthManager := newMockHealthManager(t)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			healthManager:   mockHealthManager,
			timeout:         defaultHelmClientTimeoutMins,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to update status-installed for component testComponent:")
	})

	t.Run("should fail to upgrade component on error while updating health status", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testComponent",
				Namespace: "ecosystem",
			},
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "0.1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(mock.Anything, component).Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart(mock.Anything, mock.Anything).Return(&chart.Chart{}, nil)
		spec, _ := component.GetHelmChartSpec(testCtx, k8sv1.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(mock.Anything, component.Spec.Name, "ecosystem", "0.1.0").
			Return(assert.AnError)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			healthManager:   mockHealthManager,
			timeout:         defaultHelmClientTimeoutMins,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update health status for component")
	})
}
