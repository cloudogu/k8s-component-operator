package helm

import (
	"context"
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
		component := &k8sv1.Component{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: k8sv1.ComponentSpec{
				Namespace: "testing",
				Name:      "testComponent",
				Version:   "0.1.1",
			},
			Status: k8sv1.ComponentStatus{},
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, component.GetHelmChartSpec("oci://staging.cloudogu.com"), mock.Anything).Return(nil, nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, component)

		require.NoError(t, err)
	})

	t.Run("should fail to install or upgrade chart for error in helmRepoData", func(t *testing.T) {
		component := &k8sv1.Component{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: k8sv1.ComponentSpec{
				Namespace: "testing",
				Name:      "testComponent",
				Version:   "0.1.1",
			},
			Status: k8sv1.ComponentStatus{},
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, component)

		require.Error(t, err)
		assert.ErrorContains(t, err, "error while getting oci endpoint for testComponent: error creating oci-endpoint from 'staging.cloudogu.com': wrong format")
	})

	t.Run("should fail to install or upgrade chart for error in helmClient", func(t *testing.T) {
		component := &k8sv1.Component{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: k8sv1.ComponentSpec{
				Namespace: "testing",
				Name:      "testComponent",
				Version:   "0.1.1",
			},
			Status: k8sv1.ComponentStatus{},
		}
		ctx := context.TODO()

		helmRepoData := &config.HelmRepositoryData{Endpoint: "https://staging.cloudogu.com"}
		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(ctx, component.GetHelmChartSpec("oci://staging.cloudogu.com"), mock.Anything).Return(nil, assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoData: helmRepoData}

		err := client.InstallOrUpgrade(ctx, component)

		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "error while installOrUpgrade component testing/testComponent:0.1.1")
	})
}

func TestClient_Uninstall(t *testing.T) {
	t.Run("should uninstall chart", func(t *testing.T) {
		component := &k8sv1.Component{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: k8sv1.ComponentSpec{
				Namespace: "testing",
				Name:      "testComponent",
				Version:   "0.1.1",
			},
			Status: k8sv1.ComponentStatus{},
		}

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UninstallReleaseByName(component.Spec.Name).Return(nil)

		client := &Client{helmClient: mockHelmClient, helmRepoData: nil}

		err := client.Uninstall(component)

		require.NoError(t, err)
	})

	t.Run("should fail to uninstall for error in helmClient", func(t *testing.T) {
		component := &k8sv1.Component{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: k8sv1.ComponentSpec{
				Namespace: "testing",
				Name:      "testComponent",
				Version:   "0.1.1",
			},
			Status: k8sv1.ComponentStatus{},
		}

		mockHelmClient := NewMockHelmClient(t)
		mockHelmClient.EXPECT().UninstallReleaseByName(component.Spec.Name).Return(assert.AnError)

		client := &Client{helmClient: mockHelmClient, helmRepoData: nil}

		err := client.Uninstall(component)

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
