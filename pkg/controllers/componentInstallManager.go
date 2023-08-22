package controllers

import (
	"context"
	"fmt"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentInstallManager is a central unit in the process of handling the installation process of a custom dogu resource.
type componentInstallManager struct {
	componentClient componentInterface
	helmClient      helmClient
	recorder        record.EventRecorder
}

// NewComponentInstallManager creates a new instance of componentInstallManager.
func NewComponentInstallManager(componentClient componentInterface, helmClient helmClient, recorder record.EventRecorder) *componentInstallManager {
	return &componentInstallManager{
		componentClient: componentClient,
		helmClient:      helmClient,
		recorder:        recorder,
	}
}

// Install installs a given Component Resource.
// nolint: contextcheck // uses a new non-inherited context to finish running helm-processes on SIGTERM
func (cim *componentInstallManager) Install(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	err := cim.helmClient.SatisfiesDependencies(ctx, component)
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
	helmCtx := context.Background()

	if err := cim.helmClient.InstallOrUpgrade(helmCtx, component.GetHelmChartSpec()); err != nil {
		return &genericRequeueableError{"failed to install chart for component " + component.Spec.Name, err}
	}

	component, err = cim.componentClient.UpdateStatusInstalled(helmCtx, component)
	if err != nil {
		return &genericRequeueableError{"failed to update status-installed for component " + component.Spec.Name, err}
	}

	logger.Info(fmt.Sprintf("Installed component %s.", component.Spec.Name))

	return nil
}
