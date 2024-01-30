package controllers

import (
	"context"
	"fmt"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ComponentUpgradeManager is a central unit in the process of handling the upgrade process of a custom component resource.
type ComponentUpgradeManager struct {
	componentClient componentInterface
	helmClient      helmClient
	healthManager   healthManager
	recorder        record.EventRecorder
}

// NewComponentUpgradeManager creates a new instance of ComponentUpgradeManager.
func NewComponentUpgradeManager(componentClient componentInterface, helmClient helmClient, healthManager healthManager, recorder record.EventRecorder) *ComponentUpgradeManager {
	return &ComponentUpgradeManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		healthManager:   healthManager,
		recorder:        recorder,
	}
}

// Upgrade upgrades a given component resource.
func (cum *ComponentUpgradeManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	err := cum.helmClient.SatisfiesDependencies(ctx, component.GetHelmChartSpec())
	if err != nil {
		cum.recorder.Eventf(component, corev1.EventTypeWarning, UpgradeEventReason, "Dependency check failed: %s", err.Error())
		return &genericRequeueableError{errMsg: "failed to check dependencies", err: err}
	}

	component, err = cum.componentClient.UpdateStatusUpgrading(ctx, component)
	if err != nil {
		return &genericRequeueableError{errMsg: fmt.Sprintf("failed to update status-upgrading for component %s", component.Spec.Name), err: err}
	}

	logger.Info("Upgrade helm chart...")

	// create a new context that does not get canceled immediately on SIGTERM
	// this allows self-upgrades
	helmCtx := context.WithoutCancel(ctx)

	if err := cum.helmClient.InstallOrUpgrade(helmCtx, component.GetHelmChartSpec()); err != nil {
		return &genericRequeueableError{errMsg: fmt.Sprintf("failed to upgrade chart for component %s", component.Spec.Name), err: err}
	}

	component, err = cum.componentClient.UpdateStatusInstalled(helmCtx, component)
	if err != nil {
		return &genericRequeueableError{errMsg: fmt.Sprintf("failed to update status-installed for component %s", component.Spec.Name), err: err}
	}

	err = cum.healthManager.UpdateComponentHealth(helmCtx, component.Spec.Name, component.Namespace)
	if err != nil {
		return fmt.Errorf("failed to update health status for component %q: %w", component.Spec.Name, err)
	}

	logger.Info(fmt.Sprintf("Upgraded component %s.", component.Spec.Name))

	return nil
}
