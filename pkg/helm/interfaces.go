package helm

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

// helmClient is an interface for managing components with helm.
type componentHelmClient interface {
	// InstallOrUpgrade takes a helmChart and applies it.
	InstallOrUpgrade(ctx context.Context, chart *client.ChartSpec) error
	// Uninstall removes the helmRelease for the given name
	Uninstall(releaseName string) error
	// ListDeployedReleases returns all deployed helm releases
	ListDeployedReleases() ([]*release.Release, error)
	// GetReleaseValues returns the (optionally, all computed) values for the specified release.
	GetReleaseValues(name string, allValues bool) (map[string]interface{}, error)
	// GetChartSpecValues returns the additional values for the specified ChartSpec.
	GetChartSpecValues(chart *client.ChartSpec) (map[string]interface{}, error)
	// SatisfiesDependencies validates that all dependencies are installed in the required version. A nil error
	// indicates that all dependencies (if any) meet the requirements, so that the client may conduct an installation or
	// upgrade.
	SatisfiesDependencies(ctx context.Context, chart *client.ChartSpec) error

	GetChart(ctx context.Context, chartSpec *client.ChartSpec) (*chart.Chart, error)
}
