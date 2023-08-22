package controllers

import (
	"context"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

// installManager includes functionality to install components in the cluster.
type installManager interface {
	// Install installs a component resource.
	Install(ctx context.Context, component *k8sv1.Component) error
}

// deleteManager includes functionality to delete components in the cluster.
type deleteManager interface {
	// Delete deletes a component resource.
	Delete(ctx context.Context, component *k8sv1.Component) error
}

// upgradeManager includes functionality to upgrade components in the cluster.
type upgradeManager interface {
	// Upgrade upgrades a component resource.
	Upgrade(ctx context.Context, component *k8sv1.Component) error
}

// helmClient is an interface for managing components with helm.
type helmClient interface {
	// InstallOrUpgrade takes a helmChart and applies it.
	InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error
	// Uninstall removes the helmRelease for the given name
	Uninstall(releaseName string) error
	// ListDeployedReleases returns all deployed helm releases
	ListDeployedReleases() ([]*release.Release, error)
	// SatisfiesDependencies validates that all dependencies are installed in the required version. A nil error
	// indicates that all dependencies (if any) meet the requirements, so that the client may conduct an installation or
	// upgrade.
	SatisfiesDependencies(ctx context.Context, chart *helmclient.ChartSpec) error
}

// eventRecorder embeds the record.EventRecorder interface for usage in this package.
type eventRecorder interface {
	record.EventRecorder
}

type requeueHandler interface {
	Handle(ctx context.Context, contextMessage string, componentResource *k8sv1.Component, originalErr error, onRequeue func()) (ctrl.Result, error)
}

type componentEcosystemInterface interface {
	ecosystem.ComponentEcosystemInterface
}

type componentInterface interface {
	ecosystem.ComponentInterface
}

// requeuableError indicates that the current error requires the operator to requeue the dogu.
type requeuableError interface {
	error
	// GetRequeueTime returns the time to wait before the next reconciliation.
	GetRequeueTime(requeueTimeNanos time.Duration) time.Duration
}
