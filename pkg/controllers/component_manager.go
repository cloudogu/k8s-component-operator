package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

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
		installManager: NewComponentInstallManager(clientset, helmClient, recorder),
		deleteManager:  NewComponentDeleteManager(clientset, helmClient),
		upgradeManager: NewComponentUpgradeManager(clientset, helmClient, recorder),
		recorder:       recorder,
	}
}

// Install installs the given component resource.
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
