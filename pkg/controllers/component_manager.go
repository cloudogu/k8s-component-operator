package controllers

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/mittwald/go-helm-client"
	"k8s.io/client-go/tools/record"
)

// InstallManager includes functionality to install components in the cluster.
type InstallManager interface {
	// Install installs a component resource.
	Install(ctx context.Context, component *k8sv1.Component) error
}

// DeleteManager includes functionality to delete components in the cluster.
type DeleteManager interface {
	// Delete deletes a component resource.
	Delete(ctx context.Context, component *k8sv1.Component) error
}

// UpgradeManager includes functionality to upgrade components in the cluster.
type UpgradeManager interface {
	// Upgrade upgrades a component resource.
	Upgrade(ctx context.Context, component *k8sv1.Component) error
}

// HelmClient embeds the helmclient.Client interface for usage in this package.
type HelmClient interface {
	helmclient.Client
}

// ComponentClient embeds the ecosystem.ComponentInterface interface for usage in this package.
type ComponentClient interface {
	ecosystem.ComponentInterface
}

// componentManager is a central unit in the process of handling component custom resources.
// The componentManager creates, updates and deletes components.
type componentManager struct {
	installManager InstallManager
	deleteManager  DeleteManager
	upgradeManager UpgradeManager
	recorder       record.EventRecorder
}

// NewComponentManager creates a new instance of componentManager.
func NewComponentManager(operatorConfig *config.OperatorConfig, clientset ecosystem.ComponentInterface, helmClient helmclient.Client) *componentManager {
	return &componentManager{
		installManager: NewComponentInstallManager(operatorConfig, clientset, helmClient),
		deleteManager:  NewComponentDeleteManager(operatorConfig, clientset, helmClient),
		upgradeManager: NewComponentUpgradeManager(operatorConfig, clientset, helmClient),
	}
}

// Install installs  the given component resource.
func (m *componentManager) Install(ctx context.Context, component *k8sv1.Component) error {
	return m.installManager.Install(ctx, component)
}

// Delete deletes the given component resource.
func (m *componentManager) Delete(ctx context.Context, component *k8sv1.Component) error {
	return m.deleteManager.Delete(ctx, component)
}

// Upgrade upgrades the given component resource.
func (m *componentManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	return m.upgradeManager.Upgrade(ctx, component)
}
