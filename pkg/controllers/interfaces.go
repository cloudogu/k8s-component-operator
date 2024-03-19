package controllers

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/pkg/health"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"time"

	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
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

type healthManager interface {
	health.ComponentManager
}

// helmClient is an interface for managing components with helm.
type helmClient interface {
	// InstallOrUpgrade takes a helmChart and applies it.
	InstallOrUpgrade(ctx context.Context, chart *client.ChartSpec) error
	// Uninstall removes the helmRelease for the given name
	Uninstall(releaseName string) error
	// ListDeployedReleases returns all deployed helm releases
	ListDeployedReleases() ([]*release.Release, error)
	// GetReleaseValues returns the (optionally, all computed) values for the specified release.
	GetReleaseValues(name string, allValues bool) (map[string]interface{}, error)
	// GetReleaseVersion returns the version for the specified release (if the release exists).
	GetDeployedReleaseVersion(ctx context.Context, name string) (string, error)
	// GetChartSpecValues returns the additional values for the specified ChartSpec.
	GetChartSpecValues(chart *client.ChartSpec) (map[string]interface{}, error)
	// SatisfiesDependencies validates that all dependencies are installed in the required version. A nil error
	// indicates that all dependencies (if any) meet the requirements, so that the client may conduct an installation or
	// upgrade.
	SatisfiesDependencies(ctx context.Context, chart *client.ChartSpec) error
}

// eventRecorder embeds the record.EventRecorder interface for usage in this package.
type eventRecorder interface {
	record.EventRecorder
}

type requeueHandler interface {
	// Handle takes an error and handles the requeue process for the current component operation.
	Handle(ctx context.Context, contextMessage string, componentResource *k8sv1.Component, originalErr error, requeueStatus string) (ctrl.Result, error)
}

type componentEcosystemInterface interface {
	ecosystem.ComponentEcosystemInterface
}

type componentInterface interface {
	ecosystem.ComponentInterface
}

// requeuableError indicates that the current error requires the operator to requeue the component.
type requeuableError interface {
	error
	// GetRequeueTime returns the time to wait before the next reconciliation.
	GetRequeueTime(requeueTimeNanos time.Duration) time.Duration
}

//nolint:unused
//goland:noinspection GoUnusedType
type appsV1Interface interface {
	appsv1.AppsV1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type componentV1Alpha1Interface interface {
	ecosystem.ComponentV1Alpha1Interface
}
