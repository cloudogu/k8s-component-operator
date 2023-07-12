package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentInstallManager is a central unit in the process of handling the installation process of a custom dogu resource.
type componentInstallManager struct {
	componentClient ecosystem.ComponentInterface
	helmClient      HelmClient
}

// NewComponentInstallManager creates a new instance of componentInstallManager.
func NewComponentInstallManager(componentClient ecosystem.ComponentInterface, helmClient HelmClient) *componentInstallManager {
	return &componentInstallManager{
		componentClient: componentClient,
		helmClient:      helmClient,
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

	logger.Info("Install helm chart...")

	// create a new context that does not get cancelled immediately on SIGTERM
	helmCtx, cancelHelmCtx := context.WithCancelCause(context.Background())
	defer cancelHelmCtx(nil)

	go func() {
		if err := cim.doInstall(helmCtx, component); err != nil {
			cancelHelmCtx(err)
			return
		}
		cancelHelmCtx(nil)
	}()

	// wait for install to finish
	<-helmCtx.Done()
	if err := context.Cause(helmCtx); err != context.Canceled {
		return err
	}

	logger.Info(fmt.Sprintf("Installed component %s.", component.Spec.Name))

	return nil
}

func (cim *componentInstallManager) doInstall(ctx context.Context, component *k8sv1.Component) error {
	if err := cim.helmClient.InstallOrUpgrade(ctx, component); err != nil {
		return fmt.Errorf("failed to install chart for component %s: %w", component.Spec.Name, err)
	}

	_, err := cim.componentClient.UpdateStatusInstalled(ctx, component)
	if err != nil {
		return fmt.Errorf("failed to update status-installed for component %s: %w", component.Spec.Name, err)
	}
	return nil
}
