package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/go-errors/errors"
	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"

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
	timeout         time.Duration
	reader          configMapRefReader
}

// NewComponentUpgradeManager creates a new instance of ComponentUpgradeManager.
func NewComponentUpgradeManager(componentClient componentInterface, helmClient helmClient, healthManager healthManager, recorder record.EventRecorder, timeout time.Duration, reader configMapRefReader) *ComponentUpgradeManager {
	return &ComponentUpgradeManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		healthManager:   healthManager,
		recorder:        recorder,
		timeout:         timeout,
		reader:          reader,
	}
}

// Upgrade upgrades a given component resource.
func (cupm *ComponentUpgradeManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	// set the installed version in the component CR to use it for version-comparison in future upgrades
	version, component, err := cupm.updateComponentVersion(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to get component version: %w", err)
	}

	chartSpec, err := helm.GetHelmChartSpec(ctx, component, helm.HelmChartCreationOpts{
		HelmClient:     cupm.helmClient,
		Timeout:        cupm.timeout,
		YamlSerializer: yaml.NewSerializer(),
		Reader:         cupm.reader,
	})
	if err != nil {
		return fmt.Errorf("failed to get helm chart spec: %w", err)
	}

	err = cupm.helmClient.SatisfiesDependencies(ctx, chartSpec)
	if err != nil {
		cupm.recorder.Eventf(component, corev1.EventTypeWarning, UpgradeEventReason, "Dependency check failed: %s", err.Error())
		return &genericRequeueableError{errMsg: "failed to check dependencies", err: err}
	}

	if component.Status.Status != k8sv1.ComponentStatusUpgrading {
		component, err = cupm.componentClient.UpdateStatusUpgrading(ctx, component)
		if err != nil {
			return &genericRequeueableError{errMsg: fmt.Sprintf("failed to update status-upgrading for component %s", component.Spec.Name), err: err}
		}
	}

	logger.Info("Upgrade helm chart...")

	// create a new context that does not get canceled immediately on SIGTERM
	// this allows self-upgrades
	helmCtx := context.WithoutCancel(ctx)

	rel, err := cupm.helmClient.GetRelease(component.Spec.Name)

	if err := cupm.handleHelmRelease(helmCtx, component, chartSpec, rel, err); err != nil {
		return err
	}

	component, err = cupm.componentClient.UpdateStatusInstalled(helmCtx, component)
	if err != nil {
		return &genericRequeueableError{errMsg: fmt.Sprintf("failed to update status-installed for component %s", component.Spec.Name), err: err}
	}

	err = cupm.healthManager.UpdateComponentHealthWithInstalledVersion(helmCtx, component.Spec.Name, component.Namespace, version)
	if err != nil {
		return fmt.Errorf("failed to update health status for component %q: %w", component.Spec.Name, err)
	}

	logger.Info(fmt.Sprintf("Upgraded component %s.", component.Spec.Name))

	return nil
}

// updateComponentVersion updates the component version in the component CR with the latest version
func (cupm *ComponentUpgradeManager) updateComponentVersion(ctx context.Context, component *k8sv1.Component) (string, *k8sv1.Component, error) {
	var version string
	if component.Spec.Version == "" {
		version, err := cupm.helmClient.GetLatestVersion(helm.GetHelmChartName(component))
		if err != nil {
			return "", nil, &genericRequeueableError{fmt.Sprintf("failed to get latest version for component %q", component.Spec.Name), err}
		}

		component, err = cupm.componentClient.UpdateExpectedComponentVersion(ctx, component.Spec.Name, version)
		if err != nil {
			return "", nil, &genericRequeueableError{fmt.Sprintf("failed to update expected version for component %q", component.Spec.Name), err}
		}
	} else {
		version = component.Spec.Version
	}
	return version, component, nil
}

// handleHelmRelease encapsulates the switch-case logic deciding how to deal with the helm release.
func (cupm *ComponentUpgradeManager) handleHelmRelease(
	ctx context.Context,
	component *k8sv1.Component,
	chartSpec *client.ChartSpec,
	rel *release.Release,
	err error,
) error {
	logger := log.FromContext(ctx)

	switch {
	// install helm release if it does not exist
	case errors.Is(err, driver.ErrReleaseNotFound):
		logger.Info(fmt.Sprintf("No release found for component %q, creating helm release", component.Spec.Name))
		if err := cupm.helmClient.InstallOrUpgrade(ctx, chartSpec); err != nil {
			return &genericRequeueableError{errMsg: fmt.Sprintf("failed to upgrade chart for component %s", component.Spec.Name), err: err}
		}
	// requeue if an error happens with the helm client
	case err != nil:
		return &genericRequeueableError{"failed to get release for component " + component.Spec.Name, err}
	// mark pending release as failed before reinstall
	case rel.Info.Status.IsPending():
		err := cupm.handlePendingRelease(logger, component, ctx, chartSpec)
		if err != nil {
			return &genericRequeueableError{errMsg: fmt.Sprintf("failed to handle pending helm release for component %s", component.Spec.Name), err: err}
		}
	// upgrade release in all other cases
	default:
		if err := cupm.helmClient.InstallOrUpgrade(ctx, chartSpec); err != nil {
			return &genericRequeueableError{errMsg: fmt.Sprintf("failed to upgrade chart for component %s", component.Spec.Name), err: err}
		}
	}

	return nil
}

// handlePendingRelease sets the pending release as failed, waits for it to update, and installs the component
func (cupm *ComponentUpgradeManager) handlePendingRelease(logger logr.Logger, component *k8sv1.Component, helmCtx context.Context, chartSpec *client.ChartSpec) error {
	logger.Info(fmt.Sprintf("marking pending release for component %q as failed before reinstall", component.Spec.Name))
	err := cupm.helmClient.MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall")
	if err != nil {
		return &genericRequeueableError{"failed to mark release as failed", err}
	}
	waitCtx, cancel := context.WithTimeout(helmCtx, cupm.timeout)
	defer cancel()

	done := false
	for !done {
		select {
		case <-waitCtx.Done():
			return &genericRequeueableError{
				"timed out waiting for release status update after marking as failed",
				waitCtx.Err(),
			}
		case <-time.After(2 * time.Second):
			updatedRel, getErr := cupm.helmClient.GetRelease(component.Spec.Name)
			if getErr != nil {
				return &genericRequeueableError{
					"failed to get release while waiting for status update",
					getErr,
				}
			}

			if !updatedRel.Info.Status.IsPending() {
				logger.Info(fmt.Sprintf("release status for component %q updated to %q", component.Spec.Name, updatedRel.Info.Status))
				done = true
			}
		}
	}
	if err := cupm.helmClient.InstallOrUpgrade(helmCtx, chartSpec); err != nil {
		return &genericRequeueableError{"failed to install chart for component " + component.Spec.Name, err}
	}

	return nil
}
