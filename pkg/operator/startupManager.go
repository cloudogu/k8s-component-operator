package operator

import (
	"context"
	"fmt"
	"slices"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// The startup manager does tasks on the operator startup. It is not registered as a runnable task in the k8s manager
// because it needs to be finished before the reconciliation loop is started.

func NewStartupManager(helmClient helmClient, componentClient componentClient) *StartupManager {
	return &StartupManager{
		helmClient:      helmClient,
		componentClient: componentClient,
	}
}

type StartupManager struct {
	helmClient      helmClient
	componentClient componentClient
}

// Start is run before the reconciliation loop is started. It runs the following tasks:
// * Sets the helm releases of all components that are currently installing or upgrading to failed.
func (s *StartupManager) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	err := s.setInstallingComponentsToFailed(ctx)
	if err != nil {
		logger.Info("failed to set installing components to failed: %w", err)
	}
	return nil
}

// setInstallingComponentsToFailed sets the helm releases of all components that are currently installing or upgrading to failed.
//
// After the operator startup some components will still be in status installing|upgrading|tryToInstall|tryToUpgrade and will be
// reconciled. If the helm release was not reset, the install/upgrade operation will fail because helm thinks that another operation
// is still in progress.
func (s *StartupManager) setInstallingComponentsToFailed(ctx context.Context) error {
	var err error
	resettableStatuses := []string{v1.ComponentStatusInstalling, v1.ComponentStatusTryToInstall, v1.ComponentStatusUpgrading, v1.ComponentStatusTryToUpgrade, v1.ComponentStatusDeleting, v1.ComponentStatusTryToDelete}
	components, err := s.componentClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list components: %w", err)
	}
	for _, component := range components.Items {
		if slices.Contains(resettableStatuses, component.Status.Status) {
			if err := s.helmClient.MarkReleaseAsFailed(component.Name, "setting unrecoverable release to failed for the next reconciliation"); err != nil {
				err = fmt.Errorf("failed to mark helm release as failed for component %s: %w", component.Name, err)
			}
		}
	}
	return nil
}
