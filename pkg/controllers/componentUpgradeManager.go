package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/mittwald/go-helm-client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentUpgradeManager is a central unit in the process of handling the upgrade process of a custom component resource.
type componentUpgradeManager struct {
	clientset       *ecosystem.EcosystemClientset
	componentClient ecosystem.ComponentInterface
	helmClient      helmclient.Client
	namespace       string
}

// NewComponentUpgradeManager creates a new instance of componentUpgradeManager.
func NewComponentUpgradeManager(config *config.OperatorConfig, clientset *ecosystem.EcosystemClientset, helmClient helmclient.Client) *componentUpgradeManager {
	return &componentUpgradeManager{
		clientset:       clientset,
		namespace:       config.Namespace,
		componentClient: clientset.EcosystemV1Alpha1().Components(config.Namespace),
		helmClient:      helmClient,
	}
}

// Upgrade upgrades a given component resource.
func (cum *componentUpgradeManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	component, err := cum.componentClient.UpdateStatusUpgrading(ctx, component)
	if err != nil {
		return err
	}

	err = cum.helmClient.UpdateChartRepos()
	if err != nil {
		return fmt.Errorf("failed to update chart repositories: %w", err)
	}

	_, err = cum.helmClient.UpgradeChart(ctx, component.GetHelmChartSpec(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade chart: %w", err)
	}

	_, err = cum.componentClient.UpdateStatusInstalled(ctx, component)
	if err != nil {
		return err
	}

	logger.Info("Done...")

	return nil
}
