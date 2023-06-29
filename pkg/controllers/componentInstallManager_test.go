package controllers

import (
	"context"
	"fmt"
	v1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/repo"
	"testing"
)

func TestNewComponentInstallManager(t *testing.T) {
	// when
	manager := NewComponentInstallManager(&config.OperatorConfig{}, nil, nil)

	// then
	require.NotNil(t, manager)
}

func Test_componentInstallManager_Install(t *testing.T) {
	endpoint := "test"
	username := "admin"
	password := "adminpw"
	helmSecret := &config.HelmRepositoryData{
		Endpoint: endpoint,
		Username: username,
		Password: password,
	}
	namespace := "ecosystem"
	component := getComponent(namespace, "dogu-op", "0.1.0")

	t.Run("success", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		helmRepo := getHelmRepo(component, endpoint, namespace)
		mockHelmClient.EXPECT().AddOrUpdateChartRepo(helmRepo).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(context.TODO(), component.GetHelmChartSpec(), mock.Anything).Return(nil, nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			helmRepoSecret:  helmSecret,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.NoError(t, err)
	})

	t.Run("failed to update installing status", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(nil, assert.AnError)

		mockHelmClient := NewMockHelmClient(t)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			helmRepoSecret:  helmSecret,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set status installing")
	})

	t.Run("failed to add finalizer", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(nil, assert.AnError)

		mockHelmClient := NewMockHelmClient(t)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			helmRepoSecret:  helmSecret,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to add finalizer component-finalizer")
	})

	t.Run("failed to add or update chart repository", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		helmRepo := getHelmRepo(component, endpoint, namespace)
		mockHelmClient.EXPECT().AddOrUpdateChartRepo(helmRepo).Return(assert.AnError)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			helmRepoSecret:  helmSecret,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to add or update helm repository")
	})

	t.Run("failed to install the chart", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		helmRepo := getHelmRepo(component, endpoint, namespace)
		mockHelmClient.EXPECT().AddOrUpdateChartRepo(helmRepo).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(context.TODO(), component.GetHelmChartSpec(), mock.Anything).Return(nil, assert.AnError)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			helmRepoSecret:  helmSecret,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to install chart")
	})

	t.Run("failed set status installed", func(t *testing.T) {
		// given
		mockComponentClient := NewMockComponentClient(t)
		mockComponentClient.EXPECT().UpdateStatusInstalling(context.TODO(), component).Return(component, nil)
		mockComponentClient.EXPECT().UpdateStatusInstalled(context.TODO(), component).Return(nil, assert.AnError)
		mockComponentClient.EXPECT().AddFinalizer(context.TODO(), component, "component-finalizer").Return(component, nil)

		mockHelmClient := NewMockHelmClient(t)
		helmRepo := getHelmRepo(component, endpoint, namespace)
		mockHelmClient.EXPECT().AddOrUpdateChartRepo(helmRepo).Return(nil)
		mockHelmClient.EXPECT().InstallOrUpgradeChart(context.TODO(), component.GetHelmChartSpec(), mock.Anything).Return(nil, nil)

		sut := componentInstallManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
			helmRepoSecret:  helmSecret,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to set status installed")
	})
}

func getHelmRepo(component *v1.Component, endpoint string, namespace string) repo.Entry {
	return repo.Entry{
		Name: component.Spec.Namespace,
		URL:  fmt.Sprintf("%s/%s", endpoint, namespace),
	}
}
