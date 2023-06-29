package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/mittwald/go-helm-client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// componentDeleteManager is a central unit in the process of handling the deletion process of a custom component resource.
type componentDeleteManager struct {
	componentClient ecosystem.ComponentInterface
	helmClient      helmclient.Client
	namespace       string
}

// NewComponentDeleteManager creates a new instance of componentDeleteManager.
func NewComponentDeleteManager(config *config.OperatorConfig, componentClient ecosystem.ComponentInterface, helmClient helmclient.Client) *componentDeleteManager {
	return &componentDeleteManager{
		namespace:       config.Namespace,
		componentClient: componentClient,
		helmClient:      helmClient,
	}
}

// Delete deletes a given component resource.
func (cdm *componentDeleteManager) Delete(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	component, err := cdm.componentClient.UpdateStatusDeleting(ctx, component)
	if err != nil {
		return err
	}

	err = cdm.helmClient.UninstallReleaseByName(component.Spec.Name)
	if err != nil {
		return fmt.Errorf("failed to uninstall chart: %w", err)
	}

	_, err = cdm.componentClient.RemoveFinalizer(ctx, component, k8sv1.FinalizerName)
	if err != nil {
		return err
	}

	logger.Info("Done...")

	return nil
}
