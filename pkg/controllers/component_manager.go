package controllers

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
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

// HelmClient is an interface for managing components with helm.
type HelmClient interface {
	// InstallOrUpgrade takes a component and applies the corresponding helmChart.
	InstallOrUpgrade(ctx context.Context, component *k8sv1.Component) error
	// Uninstall removes the helmChart of the given component
	Uninstall(component *k8sv1.Component) error
	// ListDeployedReleases returns all deployed helm releases
	ListDeployedReleases() ([]*release.Release, error)
}

// ComponentClient embeds the ecosystem.ComponentInterface interface for usage in this package.
type ComponentClient interface {
	ecosystem.ComponentInterface
}

// EventRecorder embeds the record.EventRecorder interface for usage in this package.
type EventRecorder interface {
	record.EventRecorder
}

// componentManager is a central unit in the process of handling component custom resources.
// The componentManager creates, updates and deletes components.
type componentManager struct {
	installManager InstallManager
	deleteManager  DeleteManager
	upgradeManager UpgradeManager
	recorder       EventRecorder
}

// NewComponentManager creates a new instance of componentManager.
func NewComponentManager(clientset ecosystem.ComponentInterface, helmClient HelmClient, recorder record.EventRecorder) *componentManager {
	return &componentManager{
		installManager: NewComponentInstallManager(clientset, helmClient),
		deleteManager:  NewComponentDeleteManager(clientset, helmClient),
		upgradeManager: NewComponentUpgradeManager(clientset, helmClient),
		recorder:       recorder,
	}
}

// Install installs  the given component resource.
func (m *componentManager) Install(ctx context.Context, component *k8sv1.Component) error {
	m.recorder.Event(component, corev1.EventTypeNormal, InstallEventReason, "Starting installation...")
	return m.installManager.Install(ctx, component)
}

// Delete deletes the given component resource.
func (m *componentManager) Delete(ctx context.Context, component *k8sv1.Component) error {
	m.recorder.Event(component, corev1.EventTypeNormal, DeinstallationEventReason, "Starting deinstallation...")
	return m.deleteManager.Delete(ctx, component)
}

// Upgrade upgrades the given component resource.
func (m *componentManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	m.recorder.Event(component, corev1.EventTypeNormal, UpgradeEventReason, "Starting upgrade...")
	return m.upgradeManager.Upgrade(ctx, component)
}
