package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/rest"
	"testing"
	"time"
)

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
		assert.ErrorContains(t, err, "failed to uninstall release 'test-release'")
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
		assert.ErrorContains(t, err, "failed to uninstall release 'test-release'")
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
		assert.ErrorContains(t, err, "failed to rollback release 'test-release'")

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
		assert.ErrorContains(t, err, "failed to get release 'test-release'")
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
