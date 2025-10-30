package controllers

import (
	"context"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/adapter/kubernetes/configref"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

// DefaultComponentManager is a central unit in the process of handling component custom resources.
// The DefaultComponentManager creates, updates and deletes components.
type DefaultComponentManager struct {
	installManager installManager
	deleteManager  deleteManager
	upgradeManager upgradeManager
	recorder       eventRecorder
}

// NewComponentManager creates a new instance of DefaultComponentManager.
func NewComponentManager(clientset componentInterface, helmClient helmClient, healthManager healthManager, recorder record.EventRecorder, timeout time.Duration, reader configref.ConfigMapRefReader) *DefaultComponentManager {
	return &DefaultComponentManager{
		installManager: NewComponentInstallManager(clientset, helmClient, healthManager, recorder, timeout, reader),
		deleteManager:  NewComponentDeleteManager(clientset, helmClient),
		upgradeManager: NewComponentUpgradeManager(clientset, helmClient, healthManager, recorder, timeout, reader),
		recorder:       recorder,
	}
}

// Install installs the given component resource.
func (m *DefaultComponentManager) Install(ctx context.Context, component *k8sv1.Component) error {
	m.recorder.Event(component, corev1.EventTypeNormal, InstallEventReason, "Starting installation...")
	return m.installManager.Install(ctx, component)
}

// Delete deletes the given component resource.
func (m *DefaultComponentManager) Delete(ctx context.Context, component *k8sv1.Component) error {
	m.recorder.Event(component, corev1.EventTypeNormal, DeinstallationEventReason, "Starting deinstallation...")
	return m.deleteManager.Delete(ctx, component)
}

// Upgrade upgrades the given component resource.
func (m *DefaultComponentManager) Upgrade(ctx context.Context, component *k8sv1.Component) error {
	m.recorder.Event(component, corev1.EventTypeNormal, UpgradeEventReason, "Starting upgrade...")
	return m.upgradeManager.Upgrade(ctx, component)
}
