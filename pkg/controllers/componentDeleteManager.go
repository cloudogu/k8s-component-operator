package controllers

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentDeleteManager is a central unit in the process of handling the deletion process of a custom component resource.
type componentDeleteManager struct {
	componentClient componentInterface
	helmClient      helmClient
}

// NewComponentDeleteManager creates a new instance of componentDeleteManager.
func NewComponentDeleteManager(componentClient componentInterface, helmClient helmClient) *componentDeleteManager {
	return &componentDeleteManager{
		componentClient: componentClient,
		helmClient:      helmClient,
	}
}

// Delete deletes a given component resource.
func (cdm *componentDeleteManager) Delete(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	component, err := cdm.componentClient.UpdateStatusDeleting(ctx, component)
	if err != nil {
		return &genericRequeueableError{fmt.Sprintf("failed to update status-deleting for component %s", component.Spec.Name), err}
	}

	deployedReleases, err := cdm.helmClient.ListDeployedReleases()
	if err != nil {
		return &genericRequeueableError{fmt.Sprintf("could not list deployed Helm releases"), err}
	}

	// Check if Helm Chart is still present before uninstalling; maybe someone has already removed it manually
	for _, release := range deployedReleases {
		if component.Spec.Name == release.Name {
			// Component Helm Chart is still present and can be uninstalled
			err = cdm.helmClient.Uninstall(component.Spec.Name)
			if err != nil {
				return &genericRequeueableError{fmt.Sprintf("failed to uninstall chart with name %s", component.Spec.Name), err}
			}
			break
		}
	}

	_, err = cdm.componentClient.RemoveFinalizer(ctx, component, k8sv1.FinalizerName)
	if err != nil {
		return &genericRequeueableError{fmt.Sprintf("failed to remove finalizer for component %s", component.Spec.Name), err}
	}

	logger.Info(fmt.Sprintf("Deleted component %s.", component.Spec.Name))

	return nil
}
