package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	"k8s.io/client-go/rest"
)

var testCtx = context.TODO()

func Test_defaultDebugLog(t *testing.T) {
	// given
	defer func() { log.Default().SetOutput(os.Stderr) }()
	buf := new(strings.Builder)
	log.Default().SetOutput(buf)

	// when
	defaultDebugLog("test %s %d %v", "1", 2, 3)

	// then
	assert.Contains(t, buf.String(), "test 1 2 3")
}

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
	t.Run("should lazy initialize options object with two field", func(t *testing.T) {
		// given
		settings := &cli.EnvSettings{}
		options := new(Options)

		// when
		err := setEnvSettings(options, settings)

		// then
		require.NoError(t, err)
		assert.Equal(t, defaultRepositoryConfigPath, settings.RepositoryConfig)
		assert.Equal(t, defaultCachePath, settings.RepositoryCache)
		assert.Equal(t, defaultRepositoryConfigPath, (*options).RepositoryConfig)
		assert.Equal(t, defaultCachePath, (*options).RepositoryCache)
	})
	t.Run("should not initialize existing options object", func(t *testing.T) {
		// given
		settings := &cli.EnvSettings{}
		options := &Options{
			Namespace:        "anamespace",
			RepositoryConfig: "arepoconfig",
			RepositoryCache:  "arepocache",
			Debug:            true,
			DebugLog:         func(format string, v ...interface{}) {},
			RegistryConfig:   "aregconfig",
			Output:           nil,
			PlainHttp:        true,
			InsecureTls:      true,
		}

		// when
		err := setEnvSettings(options, settings)

		// then
		require.NoError(t, err)
		assert.Equal(t, "arepoconfig", settings.RepositoryConfig)
		assert.Equal(t, "arepocache", settings.RepositoryCache)
		assert.Equal(t, true, settings.Debug)
		assert.Equal(t, "aregconfig", settings.RegistryConfig)
	})
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

func TestHelmClient_GetChartSpecValues(t *testing.T) {
	t.Run("should use ValuesYaml when no ValuesOptions are set", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			ChartName:   "test-chart",
			ValuesYaml: `
testYaml:
  key1: val1
  key2: val2
`,
		}

		sut := &HelmClient{
			Settings: &cli.EnvSettings{},
		}

		// when
		actual, err := sut.GetChartSpecValues(spec)

		// then
		expected := map[string]interface{}{
			"testYaml": map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
		}
		require.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	t.Run("should use ValuesOptions when no ValuesYaml are set", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			ChartName:   "test-chart",
			ValuesOptions: &values.Options{
				StringValues: []string{"testYaml.key3=val3", "testYaml.key4=val4"},
			},
		}

		sut := &HelmClient{
			Settings: &cli.EnvSettings{},
		}

		// when
		actual, err := sut.GetChartSpecValues(spec)

		// then
		expected := map[string]interface{}{
			"testYaml": map[string]interface{}{
				"key3": "val3",
				"key4": "val4",
			},
		}
		require.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	t.Run("should merge ValuesOptions and ValuesYaml when both are set with Options having priority", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			ChartName:   "test-chart",
			ValuesOptions: &values.Options{
				StringValues: []string{"testYaml.key2=val3", "testYaml.key3=val4"},
			},
			ValuesYaml: `
testYaml:
  key1: val1
  key2: val2
`,
		}

		sut := &HelmClient{
			Settings: &cli.EnvSettings{},
		}

		// when
		actual, err := sut.GetChartSpecValues(spec)

		// then
		expected := map[string]interface{}{
			"testYaml": map[string]interface{}{
				"key1": "val1",
				"key2": "val3",
				"key3": "val4",
			},
		}
		require.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	t.Run("should return empty map when nothing set", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			ChartName:   "test-chart",
		}

		sut := &HelmClient{
			Settings: &cli.EnvSettings{},
		}

		// when
		actual, err := sut.GetChartSpecValues(spec)

		// then
		expected := map[string]interface{}{}
		require.NoError(t, err)
		assert.Equal(t, actual, expected)
	})

	t.Run("should fail if yaml is not parsable", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			ChartName:   "test-chart",
			ValuesYaml: `
NoYaml{}
`,
		}

		sut := &HelmClient{
			Settings: &cli.EnvSettings{},
		}

		// when
		actual, err := sut.GetChartSpecValues(spec)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "Failed to Parse ValuesYaml")
		assert.ErrorContains(t, err, "failed to get additional values.yaml-values from")
	})

	t.Run("should fail if values options are not parsable", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ReleaseName: "test-release",
			ChartName:   "test-chart",
			ValuesOptions: &values.Options{
				StringValues: []string{"noKeyValue{}"},
			},
		}

		sut := &HelmClient{
			Settings: &cli.EnvSettings{},
		}

		// when
		actual, err := sut.GetChartSpecValues(spec)

		// then
		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "Failed to Parse ValuesOptions")
		assert.ErrorContains(t, err, "failed to get additional values.yaml-values from")
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

	t.Run("should fail to upgrade", func(t *testing.T) {
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
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newUpgrade().Return(upgradeMock)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

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

func TestHelmClient_InstallOrUpgradeChart(t *testing.T) {
	t.Run("should fail to check if chart exists", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}

		listMock := newMockListReleasesAction(t)
		listMock.EXPECT().raw().Return(&action.List{})
		listMock.EXPECT().listReleases().Return(nil, assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newListReleases().Return(listMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actual, err := sut.InstallOrUpgradeChart(testCtx, spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to list releases")
		assert.ErrorContains(t, err, "could not check if release \"test-release\" is already installed")
		assert.Nil(t, actual)
	})
	t.Run("should find release to upgrade", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
			Namespace:   "test-namespace",
		}
		upgradeAction := &action.Upgrade{}
		envSettings := &cli.EnvSettings{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}
		releaseToUpgrade := &release.Release{
			Name:      "test-release",
			Namespace: "test-namespace",
		}
		expectedRelease := &release.Release{
			Name: "test-release",
			Chart: &chart.Chart{Metadata: &chart.Metadata{
				Name:    "test-chart",
				Version: "1.1.0",
			}},
		}

		upgradeMock := newMockUpgradeAction(t)
		upgradeMock.EXPECT().raw().Return(upgradeAction)
		upgradeMock.EXPECT().upgrade(testCtx, "test-release", mock.Anything, mock.Anything).Return(expectedRelease, nil)
		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", envSettings).Return("testdata/test-chart", nil)
		listMock := newMockListReleasesAction(t)
		listMock.EXPECT().raw().Return(&action.List{})
		listMock.EXPECT().listReleases().Return([]*release.Release{releaseToUpgrade}, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newListReleases().Return(listMock)
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
		actual, err := sut.InstallOrUpgradeChart(testCtx, spec)

		// then
		require.NoError(t, err)
		assert.Same(t, expectedRelease, actual)
	})
	t.Run("should install if no release found", func(t *testing.T) {
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
		listMock := newMockListReleasesAction(t)
		listMock.EXPECT().raw().Return(&action.List{})
		listMock.EXPECT().listReleases().Return([]*release.Release{}, nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newListReleases().Return(listMock)
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
		actual, err := sut.InstallOrUpgradeChart(testCtx, spec)

		// then
		require.NoError(t, err)
		assert.Same(t, expectedRelease, actual)
	})
}

func TestHelmClient_GetChart(t *testing.T) {
	t.Run("should fail to locate chart", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}

		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", (*cli.EnvSettings)(nil)).Return("", assert.AnError)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actualChart, chartPath, err := sut.GetChart(spec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to locate chart \"test-chart\" with version \">0.0.0-0\"")
		assert.Nil(t, actualChart)
		assert.Empty(t, chartPath)
	})
	t.Run("should fail to load chart", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}

		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", (*cli.EnvSettings)(nil)).Return("invalid-path", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			actions: providerMock,
		}

		// when
		actualChart, chartPath, err := sut.GetChart(spec)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to load chart \"test-chart\" with version \">0.0.0-0\" from path \"invalid-path\"")
		assert.Nil(t, actualChart)
		assert.Empty(t, chartPath)
	})
	t.Run("should succeed to get chart with deprecation warning", func(t *testing.T) {
		// given
		spec := &ChartSpec{
			ChartName:   "test-chart",
			ReleaseName: "test-release",
		}

		locateMock := newMockLocateChartAction(t)
		locateMock.EXPECT().locateChart("test-chart", ">0.0.0-0", (*cli.EnvSettings)(nil)).Return("testdata/deprecated-chart", nil)
		providerMock := newMockActionProvider(t)
		providerMock.EXPECT().newLocateChart().Return(locateMock)

		sut := &HelmClient{
			actions: providerMock,
			DebugLog: func(format string, v ...interface{}) {
				t.Helper()
				assert.Equal(t, "WARNING: This chart (%q) is deprecated", format)
			},
		}

		// when
		actualChart, chartPath, err := sut.GetChart(spec)

		// then
		require.NoError(t, err)
		assert.NotEmpty(t, actualChart)
		assert.Equal(t, "testdata/deprecated-chart", chartPath)
	})
}

func Test_getProxyTransportIfConfigured(t *testing.T) {
	testProxy := "http://user:pass@host:3128"

	parsedProxy, err := url.Parse(testProxy)
	require.NoError(t, err)

	testProxyFn := func(request *http.Request) (*url.URL, error) {
		return parsedProxy, nil
	}

	expectedTransport := &http.Transport{
		// From https://github.com/google/go-containerregistry/blob/31786c6cbb82d6ec4fb8eb79cd9387905130534e/pkg/v1/remote/options.go#L87
		DisableCompression: true,
		DialContext: (&net.Dialer{
			// By default we wrap the transport in retries, so reduce the
			// default dial timeout to 5s to avoid 5x 30s of connection
			// timeouts when doing the "ping" on certain http registries.
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 testProxyFn,
	}

	tests := []struct {
		name     string
		wantErr  assert.ErrorAssertionFunc
		setEnv   func(t *testing.T)
		expectFn func(t *testing.T, got *http.Transport)
	}{
		{
			name: "return with proxy if configured",
			expectFn: func(t *testing.T, got *http.Transport) {
				assert.Equalf(t, expectedTransport.DisableCompression, got.DisableCompression, "getProxyTransportIfConfigured()")
				assert.Equalf(t, expectedTransport.ForceAttemptHTTP2, got.ForceAttemptHTTP2, "getProxyTransportIfConfigured()")
				assert.Equalf(t, expectedTransport.MaxIdleConns, got.MaxIdleConns, "getProxyTransportIfConfigured()")
				assert.Equalf(t, expectedTransport.IdleConnTimeout, got.IdleConnTimeout, "getProxyTransportIfConfigured()")
				assert.Equalf(t, expectedTransport.TLSHandshakeTimeout, got.TLSHandshakeTimeout, "getProxyTransportIfConfigured()")
				assert.Equalf(t, expectedTransport.ExpectContinueTimeout, got.ExpectContinueTimeout, "getProxyTransportIfConfigured()")

				gotProxy, err := got.Proxy(nil)
				require.NoError(t, err)

				wantProxy, err := expectedTransport.Proxy(nil)
				require.NoError(t, err)

				assert.Equal(t, wantProxy, gotProxy)
			},
			wantErr: assert.NoError,
			setEnv: func(t *testing.T) {
				t.Setenv("PROXY_URL", testProxy)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv != nil {
				tt.setEnv(t)
			}

			got, err := getProxyTransportIfConfigured()
			if !tt.wantErr(t, err, fmt.Sprintf("getProxyTransportIfConfigured()")) {
				return
			}

			tt.expectFn(t, got)
		})
	}
}

func Test_configureTls(t *testing.T) {
	type args struct {
		options   *Options
		transport *http.Transport
		expectFn  func(t *testing.T, transport *http.Transport)
	}
	tests := []struct {
		name     string
		args     args
		expectFn func(t *testing.T, transport *http.Transport)
	}{
		{
			name: "should set tls config in existing transport",
			args: args{
				options: &Options{
					PlainHttp:   false,
					InsecureTls: true,
				},
				transport: http.DefaultTransport.(*http.Transport),
			},
			expectFn: func(t *testing.T, transport *http.Transport) {
				require.NotNil(t, transport.TLSClientConfig)
				assert.Equal(t, true, transport.TLSClientConfig.InsecureSkipVerify)
			},
		},
		{
			name: "should create tls config in non existing transport",
			args: args{
				options: &Options{
					PlainHttp:   false,
					InsecureTls: true,
				},
				transport: nil,
			},
			expectFn: func(t *testing.T, transport *http.Transport) {
				require.NotNil(t, transport.TLSClientConfig)
				assert.Equal(t, true, transport.TLSClientConfig.InsecureSkipVerify)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := configureTls(tt.args.options, tt.args.transport)

			tt.expectFn(t, got)
		})
	}
}
