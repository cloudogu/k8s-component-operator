package helm

import (
	"context"
	"errors"
	"testing"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

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
		chart := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, chart, mock.Anything).Return(nil, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, chart)

		require.NoError(t, err)
	})

	t.Run("should install or upgrade chart with oci-endpoint in chart-name", func(t *testing.T) {
		chart := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "oci://some.where/testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}

		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, chart, mock.Anything).Return(nil, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, chart)

		require.NoError(t, err)
	})

	t.Run("should fail to install or upgrade chart for error in helmRepoData", func(t *testing.T) {
		chart := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: ""}
		mockHelmClient := NewMockHelmClient(t)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, chart)

		require.Error(t, err)
		assert.ErrorContains(t, err, "error while patching chart 'testing/testComponent': error while getting oci endpoint: error creating oci-endpoint from '': wrong format")
	})

	t.Run("should fail to install or upgrade chart for error in helmClient", func(t *testing.T) {
		chart := &helmclient.ChartSpec{
			ReleaseName: "testComponent",
			ChartName:   "testing/testComponent",
			Namespace:   "testNS",
			Version:     "0.1.1",
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, chart, mock.Anything).Return(nil, assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, chart)

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
