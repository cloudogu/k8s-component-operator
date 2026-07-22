package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
)

func TestNewComponentUpgradeManager(t *testing.T) {
	t.Run("should create new ComponentUpgradeManager", func(t *testing.T) {
		mockComponentClient := newMockComponentInterface(t)
		mockHelmClient := newMockHelmClient(t)

		manager := NewComponentUpgradeManager(mockComponentClient, mockHelmClient, nil, nil, defaultHelmClientTimeoutMins, nil)

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
			Namespace:       "ecosystem",
			Name:            "testComponent",
			Version:         "0.1.0",
			ValuesConfigRef: &k8sv1.Reference{},
		},
		Status: k8sv1.ComponentStatus{Status: "installed"},
	}

	t.Run("should upgrade component", func(t *testing.T) {
		ctx := context.Background()

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(mock.Anything, component).Return(component, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusDeployed},
		}
		mockHelmClient.EXPECT().GetRelease("testComponent").Return(rel, nil)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(mock.Anything, component.Spec.Name, "ecosystem", "0.1.0").Return(nil)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			healthManager:   mockHealthManager,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}
		err := manager.Upgrade(ctx, component)

		require.NoError(t, err)
	})

	t.Run("dependency check failed", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(assert.AnError)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, "Warning", "Upgrade", "Dependency check failed: %s", assert.AnError.Error()).Return()

		sut := ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			recorder:        mockRecorder,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
				Namespace:       "ecosystem",
				Name:            "testComponent",
				Version:         "1.0",
				ValuesConfigRef: &k8sv1.Reference{},
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, assert.AnError)

		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
				Namespace:       "ecosystem",
				Name:            "testComponent",
				Version:         "1.0",
				ValuesConfigRef: &k8sv1.Reference{},
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)

		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusDeployed},
		}
		mockHelmClient.EXPECT().GetRelease("testComponent").Return(rel, nil)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(assert.AnError)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to upgrade chart for component testComponent:")
	})

	t.Run("should fail to get release", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace:       "ecosystem",
				Name:            "testComponent",
				Version:         "1.0",
				ValuesConfigRef: &k8sv1.Reference{},
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)

		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})

		mockHelmClient.EXPECT().GetRelease("testComponent").Return(nil, assert.AnError)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to get release for component")
	})

	t.Run("should fail while installing non-existing component", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace:       "ecosystem",
				Name:            "testComponent",
				Version:         "1.0",
				ValuesConfigRef: &k8sv1.Reference{},
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)

		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().GetRelease("testComponent").Return(nil, driver.ErrReleaseNotFound)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(assert.AnError)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}
		err := manager.Upgrade(ctx, component)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to upgrade chart for component")
	})

	t.Run("should fail to upgrade component on error while setting installed status", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testComponent",
				Namespace: "ecosystem",
			},
			Spec: k8sv1.ComponentSpec{
				Namespace:       "ecosystem",
				Name:            "testComponent",
				Version:         "1.0",
				ValuesConfigRef: &k8sv1.Reference{},
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(mock.Anything, component).Return(component, assert.AnError)

		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusDeployed},
		}
		mockHelmClient.EXPECT().GetRelease("testComponent").Return(rel, nil)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(nil)

		mockHealthManager := newMockHealthManager(t)

		manager := &ComponentUpgradeManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			healthManager:   mockHealthManager,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
				Namespace:       "ecosystem",
				Name:            "testComponent",
				Version:         "0.1.0",
				ValuesConfigRef: &k8sv1.Reference{},
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusUpgrading(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(mock.Anything, component).Return(component, nil)

		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			Timeout:        defaultHelmClientTimeoutMins,
			YamlSerializer: yaml.NewSerializer(),
			Reader:         configMapRefReaderMock,
		})
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusDeployed},
		}
		mockHelmClient.EXPECT().GetRelease("testComponent").Return(rel, nil)
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
			reader:          configMapRefReaderMock,
		}
		err := manager.Upgrade(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update health status for component")
	})
}

func TestComponentUpgradeManager_handlePendingRelease(t *testing.T) {
	component := &k8sv1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testComponent",
			Namespace: "ecosystem",
		},
		Spec: k8sv1.ComponentSpec{
			Namespace: "ecosystem",
			Name:      "testComponent",
			Version:   "1.0.0",
		},
	}

	t.Run("fails when MarkReleaseAsFailed returns error", func(t *testing.T) {
		// given
		logger := logr.Discard()
		mockHelmClient := newMockHelmClient(t)

		mockHelmClient.EXPECT().
			MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
			Return(assert.AnError)

		helmCtx := context.Background()
		chartSpec := &client.ChartSpec{}

		// when
		err := handlePendingRelease(logger, component, helmCtx, chartSpec, mockHelmClient, 10*time.Second)

		// then
		require.Error(t, err)
		assert.IsType(t, &genericRequeueableError{}, err)
		assert.ErrorContains(t, err, "failed to mark release as failed")
	})

	t.Run("fails on timeout while waiting for status update", func(t *testing.T) {
		// given
		logger := logr.Discard()
		mockHelmClient := newMockHelmClient(t)

		mockHelmClient.EXPECT().
			MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
			Return(nil)

		helmCtx := context.Background()
		chartSpec := &client.ChartSpec{}

		// when
		err := handlePendingRelease(logger, component, helmCtx, chartSpec, mockHelmClient, 1*time.Nanosecond)

		// then
		require.Error(t, err)
		assert.IsType(t, &genericRequeueableError{}, err)
		assert.ErrorContains(t, err, "timed out waiting for release status update after marking as failed")
	})

	t.Run("fails when GetRelease returns error while waiting", func(t *testing.T) {
		// given
		logger := logr.Discard()
		mockHelmClient := newMockHelmClient(t)

		mockHelmClient.EXPECT().
			MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
			Return(nil)

		mockHelmClient.EXPECT().
			GetRelease(component.Spec.Name).
			Return(nil, assert.AnError)

		helmCtx := context.Background()
		chartSpec := &client.ChartSpec{}

		// when
		err := handlePendingRelease(logger, component, helmCtx, chartSpec, mockHelmClient, 3*time.Second)

		// then
		require.Error(t, err)
		assert.IsType(t, &genericRequeueableError{}, err)
		assert.ErrorContains(t, err, "failed to get release while waiting for status update")
	})

}
