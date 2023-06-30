package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentInstallManager is a central unit in the process of handling the installation process of a custom dogu resource.
type componentInstallManager struct {
	componentClient ecosystem.ComponentInterface
	helmClient      helmclient.Client
	helmRepoSecret  *config.HelmRepositoryData
}

// NewComponentInstallManager creates a new instance of componentInstallManager.
func NewComponentInstallManager(config *config.OperatorConfig, componentClient ecosystem.ComponentInterface, helmClient helmclient.Client) *componentInstallManager {
	return &componentInstallManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		helmRepoSecret:  config.HelmRepositoryData,
	}
}

// Install installs a given Component Resource.
func (cim *componentInstallManager) Install(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	component, err := cim.componentClient.UpdateStatusInstalling(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to set status installing: %w", err)
	}

	// Set the finalizer at the beginning of the installation procedure.
	// This is required because an error during installation would leave a component resource with its
	// k8s resources in the cluster. A deletion would tidy up those resources but would not start the
	// deletion procedure from the controller.
	component, err = cim.componentClient.AddFinalizer(ctx, component, k8sv1.FinalizerName)
	if err != nil {
		return fmt.Errorf("failed to add finalizer %s: %w", k8sv1.FinalizerName, err)
	}

	logger.Info("Add helm repo...")
	helmRepository := repo.Entry{
		Name:     component.Spec.Namespace,
		URL:      fmt.Sprintf("%s/%s", cim.helmRepoSecret.Endpoint, component.Spec.Namespace),
		Username: cim.helmRepoSecret.Username,
		Password: cim.helmRepoSecret.Password,
	}

	err = cim.helmClient.AddOrUpdateChartRepo(helmRepository)
	if err != nil {
		return fmt.Errorf("failed to add or update helm repository: %w", err)
	}

	logger.Info("Install helm chart...")
	_, err = cim.helmClient.InstallOrUpgradeChart(ctx, component.GetHelmChartSpec(), nil)
	if err != nil {
		return fmt.Errorf("failed to install chart: %w", err)
	}

	_, err = cim.componentClient.UpdateStatusInstalled(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to set status installed: %w", err)
	}

	logger.Info("Done...")

	return nil
}
