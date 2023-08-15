package controllers

import (
	"context"

	"helm.sh/helm/v3/pkg/release"

	"k8s.io/client-go/tools/record"

	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
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
	// SatisfiesDependencies validates that all dependencies are installed in the required version. A nil error
	// indicates that all dependencies (if any) meet the requirements, so that the client may conduct an installation or
	// upgrade.
	SatisfiesDependencies(ctx context.Context, component *k8sv1.Component) error
}

// ComponentClient embeds the ecosystem.ComponentInterface interface for usage in this package.
type ComponentClient interface {
	ecosystem.ComponentInterface
}

// EventRecorder embeds the record.EventRecorder interface for usage in this package.
type EventRecorder interface {
	record.EventRecorder
}
