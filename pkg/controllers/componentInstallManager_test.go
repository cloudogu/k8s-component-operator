package controllers

import (
	"context"
	"testing"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctxWithoutCancel = context.WithoutCancel(testCtx)
var defaultHelmClientTimeoutMins = 15 * time.Minute
var defaultRequeueTime = 3 * time.Second

func TestNewComponentInstallManager(t *testing.T) {
	// when
	manager := NewComponentInstallManager(nil, nil, nil, nil, defaultHelmClientTimeoutMins, nil)

	// then
	require.NotNil(t, manager)
}

func Test_componentInstallManager_Install(t *testing.T) {
	namespace := "ecosystem"
	component := getComponent(namespace, "k8s", "", "dogu-op", "0.1.0")
	component.Spec.ValuesConfigRef = &k8sv1.Reference{}

	t.Run("should install component", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(nil, driver.ErrReleaseNotFound)
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, mock.Anything).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(assert.AnError)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Eventf(component, "Warning", "Installation", "Dependency check failed: %s", assert.AnError.Error()).Return()

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			recorder:        mockRecorder,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		chartSpec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, chartSpec).Return(assert.AnError)
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusUnknown},
		}
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(rel, nil)
		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to install chart")
	})

	t.Run("failed to install release", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)
		mockHelmClient := newMockHelmClient(t)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		chartSpec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, chartSpec).Return(assert.AnError)
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(nil, driver.ErrReleaseNotFound)
		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to install chart for component")
	})

	t.Run("failed to get release", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)
		mockHelmClient := newMockHelmClient(t)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		mockHelmClient.EXPECT().GetRelease(component.Name).Return(nil, assert.AnError)
		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.IsType(t, err, &genericRequeueableError{})
		assert.ErrorContains(t, err, "failed to get release for component")
	})

	t.Run("install succeeds with pending release that is marked failed and reinstalled", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, component).Return(component, nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(nil)

		mockHelmClient := newMockHelmClient(t)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})

		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)

		pendingRel := &release.Release{
			Info: &release.Info{Status: release.StatusPendingInstall},
		}
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(pendingRel, nil).Once()
		mockHelmClient.EXPECT().MarkReleaseAsFailed(component.Name, "failing pending release before reinstall").Return(nil)
		nonPendingRel := &release.Release{
			Info: &release.Info{Status: release.StatusFailed},
		}
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(nonPendingRel, nil).Once()
		mockHelmClient.EXPECT().InstallOrUpgrade(mock.Anything, spec).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			healthManager:   mockHealthManager,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}

		// when
		err := sut.Install(testCtx, component)

		// then
		require.NoError(t, err)
	})

	t.Run("failed set status installed", func(t *testing.T) {
		// given
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, component).Return(component, assert.AnError)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, component, "component-finalizer").Return(component, nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		chartSpec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, chartSpec).Return(nil)
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusUnknown},
		}
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(rel, nil)
		mockHealthManager := newMockHealthManager(t)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)

		mockHelmClient := newMockHelmClient(t)
		spec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		chartSpec, _ := helm.GetHelmChartSpec(testCtx, component, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, chartSpec).Return(nil)

		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, "0.1.0").Return(assert.AnError)
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusUnknown},
		}
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(rel, nil)
		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		componentWithVersion.Spec.ValuesConfigRef = &k8sv1.Reference{}
		componentWithoutVersion.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(testCtx, componentWithVersion).Return(componentWithVersion, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(ctxWithoutCancel, componentWithVersion).Return(componentWithVersion, nil)
		mockComponentClient.EXPECT().AddFinalizer(testCtx, componentWithVersion, "component-finalizer").Return(componentWithVersion, nil)
		mockComponentClient.EXPECT().UpdateExpectedComponentVersion(testCtx, componentWithVersion.Name, componentWithVersion.Spec.Version).Return(componentWithVersion, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetLatestVersion("k8s/dogu-op").Return("4.8.3", nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)
		configMapRefReaderMock.EXPECT().GetValues(testCtx, &k8sv1.Reference{}).Return("", nil)
		spec, _ := helm.GetHelmChartSpec(testCtx, componentWithVersion, helm.HelmChartCreationOpts{
			HelmClient:     mockHelmClient,
			YamlSerializer: yaml.NewSerializer(),
			Timeout:        defaultHelmClientTimeoutMins,
			Reader:         configMapRefReaderMock,
		})
		mockHelmClient.EXPECT().SatisfiesDependencies(testCtx, spec).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgrade(ctxWithoutCancel, spec).Return(nil)
		rel := &release.Release{
			Info: &release.Info{Status: release.StatusUnknown},
		}
		mockHelmClient.EXPECT().GetRelease(component.Name).Return(rel, nil)
		mockHealthManager := newMockHealthManager(t)
		mockHealthManager.EXPECT().UpdateComponentHealthWithInstalledVersion(testCtx, component.Spec.Name, namespace, componentWithVersion.Spec.Version).Return(nil)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			healthManager:   mockHealthManager,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
		}

		// when
		err := sut.Install(testCtx, componentWithoutVersion)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to update version of component on error while getting the latest version", func(t *testing.T) {
		// given
		componentWithoutVersion := getComponent(namespace, "k8s", "", "dogu-op", "")
		componentWithoutVersion.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockComponentClient := newMockComponentInterface(t)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetLatestVersion("k8s/dogu-op").Return("", assert.AnError)
		configMapRefReaderMock := newMockConfigMapRefReader(t)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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
		componentWithoutVersion.Spec.ValuesConfigRef = &k8sv1.Reference{}
		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateExpectedComponentVersion(testCtx, component.Name, "4.8.3").Return(componentWithoutVersion, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().GetLatestVersion("k8s/dogu-op").Return("4.8.3", nil)
		configMapRefReaderMock := newMockConfigMapRefReader(t)

		sut := ComponentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			timeout:         defaultHelmClientTimeoutMins,
			reader:          configMapRefReaderMock,
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

func TestComponentInstallManager_handlePendingRelease_MarkFailedError(t *testing.T) {
	// given
	logger := logr.Discard()
	component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")

	mockHelmClient := newMockHelmClient(t)
	mockHelmClient.EXPECT().
		MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
		Return(assert.AnError)

	sut := &ComponentInstallManager{
		helmClient: mockHelmClient,
		// keep timeout reasonably high; it won't be reached because we fail earlier
		timeout: 10 * time.Second,
	}

	helmCtx := context.Background()
	chartSpec := &client.ChartSpec{}

	// when
	err := sut.handlePendingRelease(logger, component, helmCtx, chartSpec)

	// then
	require.Error(t, err)
	assert.IsType(t, &genericRequeueableError{}, err)
	assert.ErrorContains(t, err, "failed to mark release as failed")
}

func TestComponentInstallManager_handlePendingRelease_TimeoutWhileWaiting(t *testing.T) {
	// given
	logger := logr.Discard()
	component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")

	mockHelmClient := newMockHelmClient(t)
	mockHelmClient.EXPECT().
		MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
		Return(nil)

	// timeout very small so that the context is done before the first poll
	sut := &ComponentInstallManager{
		helmClient: mockHelmClient,
		timeout:    1 * time.Nanosecond,
	}

	helmCtx := context.Background()
	chartSpec := &client.ChartSpec{}

	// when
	err := sut.handlePendingRelease(logger, component, helmCtx, chartSpec)

	// then
	require.Error(t, err)
	assert.IsType(t, &genericRequeueableError{}, err)
	assert.ErrorContains(t, err, "timed out waiting for release status update after marking as failed")
}

func TestComponentInstallManager_handlePendingRelease_GetReleaseError(t *testing.T) {
	// given
	logger := logr.Discard()
	component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")

	mockHelmClient := newMockHelmClient(t)
	mockHelmClient.EXPECT().
		MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
		Return(nil)

	// give a timeout that is larger than the polling interval so that the poll actually happens
	sut := &ComponentInstallManager{
		helmClient: mockHelmClient,
		timeout:    2 * time.Second,
	}

	mockHelmClient.EXPECT().
		GetRelease(component.Spec.Name).
		Return(nil, assert.AnError)

	helmCtx := context.Background()
	chartSpec := &client.ChartSpec{}

	// when
	err := sut.handlePendingRelease(logger, component, helmCtx, chartSpec)

	// then
	require.Error(t, err)
	assert.IsType(t, &genericRequeueableError{}, err)
	assert.ErrorContains(t, err, "failed to get release while waiting for status update")
}

func TestComponentInstallManager_handlePendingRelease_InstallOrUpgradeError(t *testing.T) {
	// given
	logger := logr.Discard()
	component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")

	mockHelmClient := newMockHelmClient(t)
	mockHelmClient.EXPECT().
		MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall").
		Return(nil)

	// enough timeout so that at least one poll happens
	sut := &ComponentInstallManager{
		helmClient: mockHelmClient,
		timeout:    10 * time.Second,
	}

	pendingRel := &release.Release{
		Info: &release.Info{Status: release.StatusPendingInstall},
	}
	nonPendingRel := &release.Release{
		Info: &release.Info{Status: release.StatusFailed},
	}

	// first call: still pending, second call: not pending anymore
	mockHelmClient.EXPECT().
		GetRelease(component.Spec.Name).
		Return(pendingRel, nil).
		Once()
	mockHelmClient.EXPECT().
		GetRelease(component.Spec.Name).
		Return(nonPendingRel, nil).
		Once()

	helmCtx := context.Background()
	chartSpec := &client.ChartSpec{}

	mockHelmClient.EXPECT().
		InstallOrUpgrade(helmCtx, chartSpec).
		Return(assert.AnError)

	// when
	err := sut.handlePendingRelease(logger, component, helmCtx, chartSpec)

	// then
	require.Error(t, err)
	assert.IsType(t, &genericRequeueableError{}, err)
	assert.ErrorContains(t, err, "failed to install chart for component ")
}
