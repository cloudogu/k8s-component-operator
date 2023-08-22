package controllers

import (
	"context"
	"fmt"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentUpgradeManager is a central unit in the process of handling the upgrade process of a custom component resource.
type componentUpgradeManager struct {
	componentClient componentInterface
	helmClient      helmClient
	recorder        record.EventRecorder
}

// NewComponentUpgradeManager creates a new instance of componentUpgradeManager.
func NewComponentUpgradeManager(componentClient componentInterface, helmClient helmClient, recorder record.EventRecorder) *componentUpgradeManager {
	return &componentUpgradeManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		recorder:        recorder,
	}
}

// Upgrade upgrades a given component resource.
// nolint: contextcheck // uses a new non-inherited context to finish running helm-processes on SIGTERM
func (cum *componentUpgradeManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	err := cum.helmClient.SatisfiesDependencies(ctx, component.GetHelmChartSpec())
	if err != nil {
		cum.recorder.Eventf(component, corev1.EventTypeWarning, UpgradeEventReason, "Dependency check failed: %s", err.Error())
		return &genericRequeueableError{errMsg: "failed to check dependencies", err: err}
	}

	component, err = cum.componentClient.UpdateStatusUpgrading(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to update status-upgrading for component %s: %w", component.Spec.Name, err)
	}

	logger.Info("Upgrade helm chart...")

	// create a new context that does not get canceled immediately on SIGTERM
	helmCtx := context.Background()

	if err := cum.helmClient.InstallOrUpgrade(helmCtx, component.GetHelmChartSpec()); err != nil {
		return fmt.Errorf("failed to upgrade chart for component %s: %w", component.Spec.Name, err)
	}

	component, err = cum.componentClient.UpdateStatusInstalled(helmCtx, component)
	if err != nil {
		return fmt.Errorf("failed to update status-installed for component %s: %w", component.Spec.Name, err)
	}

	logger.Info(fmt.Sprintf("Upgraded component %s.", component.Spec.Name))

	return nil
}
