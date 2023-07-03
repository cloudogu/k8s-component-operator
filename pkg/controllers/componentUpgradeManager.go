package controllers

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentUpgradeManager is a central unit in the process of handling the upgrade process of a custom component resource.
type componentUpgradeManager struct {
	componentClient ComponentClient
	helmClient      HelmClient
}

// NewComponentUpgradeManager creates a new instance of componentUpgradeManager.
func NewComponentUpgradeManager(componentClient ComponentClient, helmClient HelmClient) *componentUpgradeManager {
	return &componentUpgradeManager{
		componentClient: componentClient,
		helmClient:      helmClient,
	}
}

// Upgrade upgrades a given component resource.
func (cum *componentUpgradeManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	component, err := cum.componentClient.UpdateStatusUpgrading(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to update status-upgrading for component %s: %w", component.Spec.Name, err)
	}

	if err = cum.helmClient.InstallOrUpgrade(ctx, component); err != nil {
		return fmt.Errorf("failed to upgrade chart for component %s: %w", component.Spec.Name, err)
	}

	_, err = cum.componentClient.UpdateStatusInstalled(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to update status-installed for component %s: %w", component.Spec.Name, err)
	}

	logger.Info(fmt.Sprintf("Upgraded component %s.", component.Spec.Name))

	return nil
}
