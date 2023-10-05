package client

import (
	"context"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/chart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	"k8s.io/client-go/rest"
)

var testCtx = context.TODO()

func TestNewClientFromRestConf(t *testing.T) {
	// given
	opt := &RestConfClientOptions{
		Options: &Options{
			RegistryConfig: "/tmp/.registry.conf",
			Namespace:      "default",
			Debug:          true,
			PlainHttp:      true,
		},
		RestConfig: &rest.Config{},
	}

	// when
	actual, err := NewClientFromRestConf(opt)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, actual)
}

func Test_setEnvSettings(t *testing.T) {
	// given
	settings := &cli.EnvSettings{}
	options := new(*Options)

	// when
	err := setEnvSettings(options, settings)

	// then
	require.NoError(t, err)
	assert.Equal(t, defaultRepositoryConfigPath, settings.RepositoryConfig)
	assert.Equal(t, defaultCachePath, settings.RepositoryCache)
	assert.Equal(t, defaultRepositoryConfigPath, (*options).RepositoryConfig)
	assert.Equal(t, defaultCachePath, (*options).RepositoryCache)
}

func TestHelmClient_UninstallReleaseByName(t *testing.T) {
	t.Run("should fail to uninstall release", func(t *testing.T) {
		// given
		uninstallMock := newMockUninstallAction(t)
		uninstallMock.EXPECT().uninstall("test-release").Return(nil, assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUninstall().Return(uninstallMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		err := sut.UninstallReleaseByName("test-release")

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to uninstall release \"test-release\"")
	})
	t.Run("should succeed to uninstall release", func(t *testing.T) {
		// given
		uninstallResponse := &release.UninstallReleaseResponse{
			Release: &release.Release{Name: "test-release"},
			Info:    "uninstall successful",
		}

		uninstallMock := newMockUninstallAction(t)
		uninstallMock.EXPECT().uninstall("test-release").Return(uninstallResponse, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUninstall().Return(uninstallMock)

		sut := &HelmClient{
			actions: providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "release uninstalled, response: %v", format)
				require.NotEmpty(t, v)
				assert.Same(t, uninstallResponse, v[0])
			},
		}

		// when
		err := sut.UninstallReleaseByName("test-release")

		// then
		require.NoError(t, err)
	})
}

func TestHelmClient_UninstallRelease(t *testing.T) {
	t.Run("should fail to uninstall release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			Timeout:     42,
		}
		uninstallAction := &action.Uninstall{}

		uninstallMock := newMockUninstallAction(t)
		uninstallMock.EXPECT().raw().Return(uninstallAction)
		uninstallMock.EXPECT().uninstall("test-release").Return(nil, assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUninstall().Return(uninstallMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		err := sut.UninstallRelease(spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to uninstall release \"test-release\"")
		assert.Equal(t, time.Duration(42), uninstallAction.Timeout)
	})
	t.Run("should succeed to uninstall release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			Timeout:     69,
		}
		uninstallAction := &action.Uninstall{}
		uninstallResponse := &release.UninstallReleaseResponse{
			Release: &release.Release{Name: "test-release"},
			Info:    "uninstall successful",
		}

		uninstallMock := newMockUninstallAction(t)
		uninstallMock.EXPECT().raw().Return(uninstallAction)
		uninstallMock.EXPECT().uninstall("test-release").Return(uninstallResponse, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUninstall().Return(uninstallMock)

		sut := &HelmClient{
			actions: providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "release uninstalled, response: %v", format)
				require.NotEmpty(t, v)
				assert.Same(t, uninstallResponse, v[0])
			},
		}

		// when
		err := sut.UninstallRelease(spec)

		// then
		require.NoError(t, err)
		assert.Equal(t, time.Duration(69), uninstallAction.Timeout)
	})
}

func TestHelmClient_RollbackRelease(t *testing.T) {
	t.Run("should fail to rollback release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName:   "test-release",
			Timeout:       42,
			CleanupOnFail: true,
		}
		rollbackAction := &action.Rollback{}

		rollbackMock := newMockRollbackReleaseAction(t)
		rollbackMock.EXPECT().raw().Return(rollbackAction)
		rollbackMock.EXPECT().rollbackRelease("test-release").Return(assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newRollbackRelease().Return(rollbackMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		err := sut.RollbackRelease(spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to rollback release \"test-release\"")

		assert.Equal(t, time.Duration(42), rollbackAction.Timeout)
		assert.True(t, rollbackAction.CleanupOnFail)
	})
	t.Run("should succeed to rollback release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			Timeout:     69,
		}
		rollbackAction := &action.Rollback{}

		rollbackMock := newMockRollbackReleaseAction(t)
		rollbackMock.EXPECT().raw().Return(rollbackAction)
		rollbackMock.EXPECT().rollbackRelease("test-release").Return(nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newRollbackRelease().Return(rollbackMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		err := sut.RollbackRelease(spec)

		// then
		require.NoError(t, err)
		assert.Equal(t, time.Duration(69), rollbackAction.Timeout)
		assert.False(t, rollbackAction.CleanupOnFail)
	})
}

func TestHelmClient_GetRelease(t *testing.T) {
	t.Run("should fail to get release", func(t *testing.T) {
		// given
		getReleaseMock := newMockGetReleaseAction(t)
		getReleaseMock.EXPECT().getRelease("test-release").Return(nil, assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newGetRelease().Return(getReleaseMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.GetRelease("test-release")

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get release \"test-release\"")
	})
	t.Run("should succeed to get release", func(t *testing.T) {
		// given
		expectedRelease := release.Release{Name: "test-release"}
		getReleaseMock := newMockGetReleaseAction(t)
		getReleaseMock.EXPECT().getRelease("test-release").Return(&expectedRelease, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newGetRelease().Return(getReleaseMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.GetRelease("test-release")

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedRelease, *actual)
	})
}

func TestHelmClient_GetReleaseValues(t *testing.T) {
	t.Run("should fail to get release values", func(t *testing.T) {
		// given
		getValuesAction := &action.GetValues{}
		getValuesMock := newMockGetReleaseValuesAction(t)
		getValuesMock.EXPECT().raw().Return(getValuesAction)
		getValuesMock.EXPECT().getReleaseValues("test-release").Return(nil, assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newGetReleaseValues().Return(getValuesMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.GetReleaseValues("test-release", true)

		// then
		require.Error(t, err)

		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get values of release \"test-release\"")

		assert.Nil(t, actual)
		assert.True(t, getValuesAction.AllValues)
	})
	t.Run("should succeed to get release values", func(t *testing.T) {
		// given
		getValuesAction := &action.GetValues{}
		expectedValues := map[string]interface{}{"myKey": "myValue"}
		getValuesMock := newMockGetReleaseValuesAction(t)
		getValuesMock.EXPECT().raw().Return(getValuesAction)
		getValuesMock.EXPECT().getReleaseValues("test-release").Return(expectedValues, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newGetReleaseValues().Return(getValuesMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.GetReleaseValues("test-release", false)

		// then
		require.NoError(t, err)

		assert.Equal(t, expectedValues, actual)
		assert.False(t, getValuesAction.AllValues)
	})
}

func TestHelmClient_ListDeployedReleases(t *testing.T) {
	t.Run("should fail to list deployed releases", func(t *testing.T) {
		// given
		listAction := &action.List{}
		listMock := newMockListReleasesAction(t)
		listMock.EXPECT().raw().Return(listAction)
		listMock.EXPECT().listReleases().Return(nil, assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newListReleases().Return(listMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.ListDeployedReleases()

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to list releases")

		assert.Nil(t, actual)
		assert.Equal(t, action.ListDeployed, listAction.StateMask)
	})
	t.Run("should succeed to list deployed releases", func(t *testing.T) {
		// given
		listAction := &action.List{}
		expectedReleases := []*release.Release{{Name: "test-release"}}
		listMock := newMockListReleasesAction(t)
		listMock.EXPECT().raw().Return(listAction)
		listMock.EXPECT().listReleases().Return(expectedReleases, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newListReleases().Return(listMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.ListDeployedReleases()

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedReleases, actual)
		assert.Equal(t, action.ListDeployed, listAction.StateMask)
	})
}

func TestHelmClient_ListReleasesByStateMask(t *testing.T) {
	t.Run("should succeed to list failed releases", func(t *testing.T) {
		// given
		listAction := &action.List{}
		expectedReleases := []*release.Release{{Name: "test-release"}}
		listMock := newMockListReleasesAction(t)
		listMock.EXPECT().raw().Return(listAction)
		listMock.EXPECT().listReleases().Return(expectedReleases, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newListReleases().Return(listMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.ListReleasesByStateMask(action.ListFailed)

		// then
		require.NoError(t, err)
		assert.Equal(t, expectedReleases, actual)
		assert.Equal(t, action.ListFailed, listAction.StateMask)
	})
}

func TestHelmClient_InstallChart(t *testing.T) {
	t.Run("should fail on empty release name", func(t *testing.T) {
		// given
		spec := &ChartSpec{ChartName: "test-chart", ReleaseName: ""}
		installAction := &action.Install{}

		installMock := newMockInstallAction(t)
		installMock.EXPECT().raw().Return(installAction)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newInstall().Return(installMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.InstallChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to determine release name for chart \"test-chart\"")
		assert.Nil(t, actual)
	})
	t.Run("should fail to get chart", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		installAction := &action.Install{}

		installMock := newMockInstallAction(t)
		installMock.EXPECT().raw().Return(installAction)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", (*cli.EnvSettings)(nil)).Return("", assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newInstall().Return(installMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.InstallChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get chart for release \"test-release\"")
		assert.ErrorContains(t, err, "failed to locate chart \"test-chart\" with version \">0.0.0-0\"")
		assert.Nil(t, actual)
	})
	t.Run("should fail because chart has unsupported type", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		installAction := &action.Install{}

		installMock := newMockInstallAction(t)
		installMock.EXPECT().raw().Return(installAction)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", (*cli.EnvSettings)(nil)).Return("testdata/invalid-type-test-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newInstall().Return(installMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.InstallChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "chart \"test-chart\" has an unsupported type and is not installable: \"library\"")
		assert.Nil(t, actual)
	})
	t.Run("should fail to get values", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
			ValuesYaml:  "invalid YAML",
		}
		installAction := &action.Install{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}

		installMock := newMockInstallAction(t)
		installMock.EXPECT().raw().Return(installAction)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newInstall().Return(installMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
		}

		// when
		actual, err := sut.InstallChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get values for release \"test-release\"")
		assert.Nil(t, actual)
	})
	t.Run("should fail to install release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		installAction := &action.Install{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}

		installMock := newMockInstallAction(t)
		installMock.EXPECT().raw().Return(installAction)
		installMock.EXPECT().install(testCtx, mock.Anything, mock.Anything).Return(nil, assert.AnError)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newInstall().Return(installMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
		}

		// when
		actual, err := sut.InstallChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to install release \"test-release\"")
		assert.Nil(t, actual)
	})
	t.Run("should succeed to install release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		installAction := &action.Install{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}
		expectedRelease := &release.Release{
			Name: "test-release",
			Chart: &chart.Chart{Metadata: &chart.Metadata{
				Name:    "test-chart",
				Version: "1.0.0",
			}},
		}

		installMock := newMockInstallAction(t)
		installMock.EXPECT().raw().Return(installAction)
		installMock.EXPECT().install(testCtx, mock.Anything, mock.Anything).Return(expectedRelease, nil)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newInstall().Return(installMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "release installed successfully: %s/%s-%s", format)
			},
		}

		// when
		actual, err := sut.InstallChart(testCtx, spec)

		// then
		require.NoError(t, err)
		assert.Same(t, expectedRelease, actual)
	})
}

func TestHelmClient_UpgradeChart(t *testing.T) {
	t.Run("should fail to get chart", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		upgradeAction := &action.Upgrade{}

		upgradeMock := newMockUpgradeAction(t)
		upgradeMock.EXPECT().raw().Return(upgradeAction)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", (*cli.EnvSettings)(nil)).Return("", assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUpgrade().Return(upgradeMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.UpgradeChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get chart for release \"test-release\"")

		assert.Nil(t, actual)
	})
	t.Run("should fail to get values for release", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
			ValuesYaml:  "invalid YAML",
		}
		upgradeAction := &action.Upgrade{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}

		upgradeMock := newMockUpgradeAction(t)
		upgradeMock.EXPECT().raw().Return(upgradeAction)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUpgrade().Return(upgradeMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
		}

		// when
		actual, err := sut.UpgradeChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to get values for release \"test-release\"")

		assert.Nil(t, actual)
	})
	t.Run("should fail to upgrade and fail to rollback", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		upgradeAction := &action.Upgrade{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}

		upgradeMock := newMockUpgradeAction(t)
		upgradeMock.EXPECT().raw().Return(upgradeAction)
		upgradeMock.EXPECT().upgrade(testCtx, "test-release", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		rollbackMock := newMockRollbackReleaseAction(t)
		rollbackMock.EXPECT().raw().Return(&action.Rollback{})
		rollbackMock.EXPECT().rollbackRelease("test-release").Return(assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUpgrade().Return(upgradeMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)
		providerMock.EXPECT().newRollbackRelease().Return(rollbackMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "release upgrade failed: %s", format)
			},
		}

		// when
		actual, err := sut.UpgradeChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "release failed, rollback failed")
		assert.ErrorContains(t, err, "failed to upgrade release \"test-release\"")
		assert.ErrorContains(t, err, "failed to rollback release \"test-release\"")

		assert.Nil(t, actual)
	})
	t.Run("should fail to upgrade and succeed to rollback", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		upgradeAction := &action.Upgrade{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}

		upgradeMock := newMockUpgradeAction(t)
		upgradeMock.EXPECT().raw().Return(upgradeAction)
		upgradeMock.EXPECT().upgrade(testCtx, "test-release", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		rollbackMock := newMockRollbackReleaseAction(t)
		rollbackMock.EXPECT().raw().Return(&action.Rollback{})
		rollbackMock.EXPECT().rollbackRelease("test-release").Return(nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUpgrade().Return(upgradeMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)
		providerMock.EXPECT().newRollbackRelease().Return(rollbackMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "release upgrade failed: %s", format)
			},
		}

		// when
		actual, err := sut.UpgradeChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "release failed, rollback succeeded")
		assert.ErrorContains(t, err, "failed to upgrade release \"test-release\"")

		assert.Nil(t, actual)
	})
	t.Run("should succeed to upgrade", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}
		upgradeAction := &action.Upgrade{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}
		expectedRelease := &release.Release{
			Name: "test-release",
			Chart: &chart.Chart{Metadata: &chart.Metadata{
				Name:    "test-chart",
				Version: "1.0.0",
			}},
		}

		upgradeMock := newMockUpgradeAction(t)
		upgradeMock.EXPECT().raw().Return(upgradeAction)
		upgradeMock.EXPECT().upgrade(testCtx, "test-release", mock.Anything, mock.Anything).Return(expectedRelease, nil)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUpgrade().Return(upgradeMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			Settings: envSettings,
			actions:  providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "release upgraded successfully: %s/%s-%s", format)
			},
		}

		// when
		actual, err := sut.UpgradeChart(testCtx, spec)

		// then
		require.NoError(t, err)
		assert.Same(t, expectedRelease, actual)
	})
}
