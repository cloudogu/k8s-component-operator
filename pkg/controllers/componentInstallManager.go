package controllers

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

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
}

// NewComponentInstallManager creates a new instance of ComponentInstallManager.
func NewComponentInstallManager(componentClient componentInterface, helmClient helmClient, healthManager healthManager, recorder record.EventRecorder) *ComponentInstallManager {
	return &ComponentInstallManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		healthManager:   healthManager,
		recorder:        recorder,
	}
}

// Install installs a given Component Resource.
// If no expected version is given in the component CR the latest version will be installed.
func (cim *ComponentInstallManager) Install(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	err := cim.helmClient.SatisfiesDependencies(ctx, component.GetHelmChartSpec())
	if err != nil {
		cim.recorder.Eventf(component, corev1.EventTypeWarning, InstallEventReason, "Dependency check failed: %s", err.Error())
		return &genericRequeueableError{errMsg: "failed to check dependencies", err: err}
	}

	component, err = cim.componentClient.UpdateStatusInstalling(ctx, component)
	if err != nil {
		return &genericRequeueableError{errMsg: "failed to set status installing", err: err}
	}

	// Set the finalizer at the beginning of the installation procedure.
	// This is required because an error during installation would leave a component resource with its
	// k8s resources in the cluster. A deletion would tidy up those resources but would not start the
	// deletion procedure from the controller.
	component, err = cim.componentClient.AddFinalizer(ctx, component, k8sv1.FinalizerName)
	if err != nil {
		return &genericRequeueableError{"failed to add finalizer " + k8sv1.FinalizerName, err}
	}

	logger.Info("Install helm chart...")

	// create a new context that does not get canceled immediately on SIGTERM
	helmCtx := context.WithoutCancel(ctx)

	// if no version is set in the component CR, this will find the latest version and install it
	// TODO: if no specific version is set, find the latest version here, not in the client.
	//  A client should have no business logic (but we can bundle it with e.g. `GetLatestChartVersion()`)
	//  Then we dont have to load the installed version again with a helm call.
	if err := cim.helmClient.InstallOrUpgrade(helmCtx, component.GetHelmChartSpec()); err != nil {
		return &genericRequeueableError{"failed to install chart for component " + component.Spec.Name, err}
	}

	// set the installed version in the component CR to use it for version-comparison in future upgrades
	installedVersion, err := cim.helmClient.GetDeployedReleaseVersion(ctx, component.Spec.Name)
	if err != nil {
		return &genericRequeueableError{fmt.Sprintf("failed to get release version for component %q", component.Spec.Name), err}
	}
	if component.Spec.Version == "" {
		component, err = cim.componentClient.UpdateExpectedComponentVersion(helmCtx, component.Spec.Name, installedVersion)
		if err != nil {
			return &genericRequeueableError{fmt.Sprintf("failed to update expected version for component %q", component.Spec.Name), err}
		}
	}
	component, err = cim.componentClient.UpdateStatusInstalled(helmCtx, component)
	if err != nil {
		return &genericRequeueableError{fmt.Sprintf("failed to update status-installed for component %q", component.Spec.Name), err}
	}

	err = cim.healthManager.UpdateComponentHealthWithInstalledVersion(ctx, component.Spec.Name, component.Namespace, installedVersion)
	if err != nil {
		return fmt.Errorf("failed to update health status and installed version for component %q: %w", component.Spec.Name, err)
	}

	logger.Info(fmt.Sprintf("Installed component %s.", component.Spec.Name))

	return nil
}
