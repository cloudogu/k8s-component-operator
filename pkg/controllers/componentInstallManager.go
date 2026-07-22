package controllers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	helmRelease "helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ComponentInstallManager is a central unit in the process of handling the installation process of a custom dogu resource.
type ComponentInstallManager struct {
	componentClient componentInterface
	helmClient      helmClient
	healthManager   healthManager
	recorder        record.EventRecorder
	timeout         time.Duration
	reader          configMapRefReader
}

// NewComponentInstallManager creates a new instance of ComponentInstallManager.
func NewComponentInstallManager(componentClient componentInterface, helmClient helmClient, healthManager healthManager, recorder record.EventRecorder, timeout time.Duration, reader configMapRefReader) *ComponentInstallManager {
	return &ComponentInstallManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		healthManager:   healthManager,
		recorder:        recorder,
		timeout:         timeout,
		reader:          reader,
	}
}

// Install installs a given Component Resource.
// If no expected version is given in the component CR the latest version will be installed.
func (cim *ComponentInstallManager) Install(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	var err error
	// set the installed version in the component CR to use it for version-comparison in future upgrades
	var version string
	if component.Spec.Version == "" {
		version, err = cim.helmClient.GetLatestVersion(helm.GetHelmChartName(component))
		if err != nil {
			return &genericRequeueableError{fmt.Sprintf("failed to get latest version for component %q", component.Spec.Name), err}
		}

		component, err = cim.componentClient.UpdateExpectedComponentVersion(ctx, component.Spec.Name, version)
		if err != nil {
			return &genericRequeueableError{fmt.Sprintf("failed to update expected version for component %q", component.Spec.Name), err}
		}
	} else {
		version = component.Spec.Version
	}

	chartSpec, err := helm.GetHelmChartSpec(ctx, component, helm.HelmChartCreationOpts{
		HelmClient:     cim.helmClient,
		Timeout:        cim.timeout,
		YamlSerializer: yaml.NewSerializer(),
		Reader:         cim.reader,
	})
	if err != nil {
		return fmt.Errorf("failed to get helm chart spec: %w", err)
	}
	err = cim.helmClient.SatisfiesDependencies(ctx, chartSpec)
	if err != nil {
		cim.recorder.Eventf(component, corev1.EventTypeWarning, InstallEventReason, "Dependency check failed: %s", err.Error())
		return &genericRequeueableError{errMsg: "failed to check dependencies", err: err}
	}

	if component.Status.Status != k8sv1.ComponentStatusInstalling {
		component, err = cim.componentClient.UpdateStatusInstalling(ctx, component)
		if err != nil {
			return &genericRequeueableError{errMsg: "failed to set status installing", err: err}
		}
	}

	// Set the finalizer at the beginning of the installation procedure.
	// This is required because an error during installation would leave a component resource with its
	// k8s resources in the cluster. A deletion would tidy up those resources but would not start the
	// deletion procedure from the controller.
	if !slices.Contains(component.Finalizers, k8sv1.FinalizerName) {
		component, err = cim.componentClient.AddFinalizer(ctx, component, k8sv1.FinalizerName)
		if err != nil {
			return &genericRequeueableError{"failed to add finalizer " + k8sv1.FinalizerName, err}
		}
	}

	// create a new context that does not get canceled immediately on SIGTERM
	helmCtx := context.WithoutCancel(ctx)

	release, err := cim.helmClient.GetRelease(component.Spec.Name)

	switch {
	// install helm release if it does not exist
	case errors.Is(err, driver.ErrReleaseNotFound):
		logger.Info(fmt.Sprintf("No release found for component %q, creating helm release", component.Spec.Name))
		if err := cim.helmClient.InstallOrUpgrade(helmCtx, chartSpec); err != nil {
			return &genericRequeueableError{"failed to install chart for component " + component.Spec.Name, err}
		}
	// requeue if an error happens with the helm client
	case err != nil:
		return &genericRequeueableError{"failed to get release for component " + component.Spec.Name, err}
	// mark pending release as failed and reinstall
	case release.Info.Status.IsPending():
		err := handlePendingRelease(logger, component, helmCtx, chartSpec, cim.helmClient, cim.timeout)
		if err != nil {
			return &genericRequeueableError{"failed to handle pending helm release for component " + component.Spec.Name, err}
		}
		if err := cim.helmClient.InstallOrUpgrade(helmCtx, chartSpec); err != nil {
			return &genericRequeueableError{"failed to install chart for component " + component.Spec.Name, err}
		}
	// do nothing if the release is already deployed
	case release.Info.Status != helmRelease.StatusDeployed:
		logger.Info(fmt.Sprintf("Release found with status %q for component %q, trying to install/upgrade", release.Info.Status, component.Spec.Name))
		if err := cim.helmClient.InstallOrUpgrade(helmCtx, chartSpec); err != nil {
			return &genericRequeueableError{"failed to install chart for component " + component.Spec.Name, err}
		}
	}

	component, err = cim.componentClient.UpdateStatusInstalled(helmCtx, component)
	if err != nil {
		return &genericRequeueableError{fmt.Sprintf("failed to update status-installed for component %q", component.Spec.Name), err}
	}

	err = cim.healthManager.UpdateComponentHealthWithInstalledVersion(ctx, component.Spec.Name, component.Namespace, version)
	if err != nil {
		return fmt.Errorf("failed to update health status and installed version for component %q: %w", component.Spec.Name, err)
	}

	logger.Info(fmt.Sprintf("Installed component %s.", component.Spec.Name))

	return nil
}
