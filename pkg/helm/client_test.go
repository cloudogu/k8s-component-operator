package helm

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudogu/k8s-component-operator/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

var testCtx = context.TODO()

func TestNew(t *testing.T) {
	t.Run("should create new client", func(t *testing.T) {
		namespace := "ecosystem"

		// override default controller method to retrieve a kube config
		oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
		defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
		ctrl.GetConfigOrDie = func() *rest.Config {
			return &rest.Config{}
		}

		client, err := NewClient(namespace, &config.HelmRepositoryData{}, false, nil)

		require.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestClient_InstallOrUpgrade(t *testing.T) {
	t.Run("should install or upgrade chart", func(t *testing.T) {
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(testCtx, chartSpec, mock.Anything).Return(nil, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(testCtx, chartSpec)

		require.NoError(t, err)
	})

	t.Run("should install or upgrade chart with oci-endpoint in chart-name", func(t *testing.T) {
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "oci://some.where/testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(testCtx, chartSpec, mock.Anything).Return(nil, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(testCtx, chartSpec)

		require.NoError(t, err)
	})

	t.Run("should fail to install or upgrade chart for error in helmRepoData", func(t *testing.T) {
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: ""}
		mockHelmClient := NewMockHelmClient(t)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, chartSpec)

		require.Error(t, err)
		assert.ErrorContains(t, err, "error while patching chart 'testing/testComponent': error while getting oci endpoint: error creating oci-endpoint from '': wrong format")
	})

	t.Run("should fail to install or upgrade chart for error in helmClient", func(t *testing.T) {
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, chartSpec, mock.Anything).Return(nil, assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, chartSpec)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while installOrUpgrade chart oci://staging.cloudogu.com/testing/testComponent:")
	})
}

func TestClient_Uninstall(t *testing.T) {
	t.Run("should uninstall chart", func(t *testing.T) {
		releaseName := "testComponent"
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UninstallReleaseByName(releaseName).Return(nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: nil}

		err := client.Uninstall(releaseName)

		require.NoError(t, err)
	})

	t.Run("should fail to uninstall for error in helmClient", func(t *testing.T) {
		releaseName := "testComponent"

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UninstallReleaseByName(releaseName).Return(assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoData: nil}

		err := client.Uninstall(releaseName)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while uninstalling helm-release testComponent")
	})
}

func TestClient_ListDeployedReleases(t *testing.T) {
	t.Run("should list deployed releases", func(t *testing.T) {
		releases := []*release.Release{
			{Name: "Test Release 1"},
			{Name: "Test Release 2"},
		}

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().ListDeployedReleases().Return(releases, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: nil}

		result, err := client.ListDeployedReleases()

		require.NoError(t, err)
		assert.Equal(t, releases, result)
	})

	t.Run("should fail to list deployed releases for error in helmClient", func(t *testing.T) {
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().ListDeployedReleases().Return(nil, assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoData: nil}

		result, err := client.ListDeployedReleases()

		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_dependencyUnsatisfiedError_Unwrap(t *testing.T) {
	testErr1 := assert.AnError
	testErr2 := errors.New("test")
	inputErr := errors.Join(testErr1, testErr2)

	sut := &dependencyUnsatisfiedError{inputErr}

	// when
	actualErr := sut.Unwrap()

	// then
	require.Error(t, sut)
	require.Error(t, actualErr)
	assert.ErrorIs(t, actualErr, testErr1)
	assert.ErrorIs(t, actualErr, testErr2)
}

func Test_dependencyUnsatisfiedError_Error(t *testing.T) {
	sut := &dependencyUnsatisfiedError{assert.AnError}
	expected := "one or more dependencies are not satisfied: assert.AnError general error for testing"
	assert.Equal(t, expected, sut.Error())
}

func TestClient_SatisfiesDependencies(t *testing.T) {
	t.Run("should fail to get oci endpoint", func(t *testing.T) {
		// given
		mockRepoConfig := newMockOciRepositoryConfig(t)
		mockRepoConfig.EXPECT().GetOciEndpoint().Return("", assert.AnError)
		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}
		sut := &Client{
			helmRepoData: mockRepoConfig,
		}

		// when
		err := sut.SatisfiesDependencies(testCtx, chartSpec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while patching chart 'testComponent': error while getting oci endpoint")
	})

	t.Run("should fail to get chart", func(t *testing.T) {
		// given
		mockRepoConfig := newMockOciRepositoryConfig(t)
		mockRepoConfig.EXPECT().GetOciEndpoint().Return("oci://some.where/testing", nil)
		mockRepoConfig.EXPECT().IsPlainHttp().Return(true)

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart("oci://some.where/testing/testComponent", mock.Anything).Return(nil, "", assert.AnError)

		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		sut := &Client{
			helmClient:   mockHelmClient,
			helmRepoData: mockRepoConfig,
			actionConfig: new(action.Configuration),
		}

		// when
		err := sut.SatisfiesDependencies(testCtx, chartSpec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get chart oci://some.where/testing/testComponent: error while getting chart for oci://some.where/testing/testComponent:0.1.1")
	})

	t.Run("should fail to list deployed releases", func(t *testing.T) {
		// given
		mockRepoConfig := newMockOciRepositoryConfig(t)
		mockRepoConfig.EXPECT().GetOciEndpoint().Return("oci://some.where/testing", nil)
		mockRepoConfig.EXPECT().IsPlainHttp().Return(true)

		helmChart := &chart.Chart{Metadata: &chart.Metadata{
			Dependencies: []*chart.Dependency{{
				Name:    "k8s-etcd",
				Version: "3.*.*",
			}},
		}}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart("oci://some.where/testing/testComponent", mock.Anything).Return(helmChart, "myPath", nil)
		mockHelmClient.EXPECT().ListDeployedReleases().Return(nil, assert.AnError)

		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		sut := &Client{
			helmClient:   mockHelmClient,
			helmRepoData: mockRepoConfig,
			actionConfig: new(action.Configuration),
		}

		// when
		err := sut.SatisfiesDependencies(testCtx, chartSpec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to list deployed releases")
	})

	t.Run("should return unsatisfied error", func(t *testing.T) {
		// given
		mockRepoConfig := newMockOciRepositoryConfig(t)
		mockRepoConfig.EXPECT().GetOciEndpoint().Return("oci://some.where/testing", nil)
		mockRepoConfig.EXPECT().IsPlainHttp().Return(true)

		dependencies := []*chart.Dependency{createDependency("k8s-etcd", "3.2.1")}
		helmChart := &chart.Chart{Metadata: &chart.Metadata{
			Dependencies: dependencies,
		}}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart("oci://some.where/testing/testComponent", mock.Anything).Return(helmChart, "myPath", nil)
		var deployedReleases []*release.Release
		mockHelmClient.EXPECT().ListDeployedReleases().Return(deployedReleases, nil)

		mockDepChecker := newMockDependencyChecker(t)
		mockDepChecker.EXPECT().CheckSatisfied(dependencies, deployedReleases).Return(assert.AnError)

		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		sut := &Client{
			helmClient:        mockHelmClient,
			helmRepoData:      mockRepoConfig,
			actionConfig:      new(action.Configuration),
			dependencyChecker: mockDepChecker,
		}

		// when
		err := sut.SatisfiesDependencies(testCtx, chartSpec)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		targetErr := &dependencyUnsatisfiedError{}
		assert.ErrorAs(t, err, &targetErr)
		assert.ErrorContains(t, err, "one or more dependencies are not satisfied")
	})

	t.Run("should succeed", func(t *testing.T) {
		// given
		mockRepoConfig := newMockOciRepositoryConfig(t)
		mockRepoConfig.EXPECT().GetOciEndpoint().Return("oci://some.where/testing", nil)
		mockRepoConfig.EXPECT().IsPlainHttp().Return(true)

		dependencies := []*chart.Dependency{createDependency("k8s-etcd", "3.2.1")}
		helmChart := &chart.Chart{Metadata: &chart.Metadata{
			Dependencies: dependencies,
		}}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().GetChart("oci://some.where/testing/testComponent", mock.Anything).Return(helmChart, "myPath", nil)
		deployedReleases := []*release.Release{createRelease("k8s-etcd", "3.2.1")}
		mockHelmClient.EXPECT().ListDeployedReleases().Return(deployedReleases, nil)

		mockDepChecker := newMockDependencyChecker(t)
		mockDepChecker.EXPECT().CheckSatisfied(dependencies, deployedReleases).Return(nil)

		chartSpec := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		sut := &Client{
			helmClient:        mockHelmClient,
			helmRepoData:      mockRepoConfig,
			actionConfig:      new(action.Configuration),
			dependencyChecker: mockDepChecker,
		}

		// when
		err := sut.SatisfiesDependencies(testCtx, chartSpec)

		// then
		require.NoError(t, err)
	})
}
