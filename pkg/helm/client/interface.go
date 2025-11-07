package client

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

// Client holds the method signatures for a Helm client.
// NOTE: This is an interface to allow for mocking in tests.
type Client interface {
	InstallOrUpgradeChart(ctx context.Context, spec *ChartSpec) (*release.Release, error)
	InstallChart(ctx context.Context, spec *ChartSpec) (*release.Release, error)
	UpgradeChart(ctx context.Context, spec *ChartSpec) (*release.Release, error)
	ListDeployedReleases() ([]*release.Release, error)
	ListReleasesByStateMask(action.ListStates) ([]*release.Release, error)
	GetRelease(name string) (*release.Release, error)
	// RollBack is an interface to abstract a rollback action.
	RollBack
	GetReleaseValues(name string, allValues bool) (map[string]interface{}, error)
	GetChartSpecValues(spec *ChartSpec) (map[string]interface{}, error)
	UninstallRelease(spec *ChartSpec) error
	UninstallReleaseByName(name string) error
	GetChart(spec *ChartSpec) (*chart.Chart, string, error)
	TagResolver
}

type TagResolver interface {
	Tags(ref string) ([]string, error)
}

type RollBack interface {
	RollbackRelease(spec *ChartSpec) error
}

type actionProvider interface {
	newInstall() installAction
	newUpgrade() upgradeAction
	newLocateChart() locateChartAction
	newUninstall() uninstallAction
	newListReleases() listReleasesAction
	newGetReleaseValues() getReleaseValuesAction
	newGetRelease() getReleaseAction
	newRollbackRelease() rollbackReleaseAction
}

type installAction interface {
	install(ctx context.Context, chart *chart.Chart, values map[string]interface{}) (*release.Release, error)
	raw() *action.Install
}

type upgradeAction interface {
	upgrade(ctx context.Context, releaseName string, chart *chart.Chart, values map[string]interface{}) (*release.Release, error)
	raw() *action.Upgrade
}

type locateChartAction interface {
	locateChart(name, version string, settings *cli.EnvSettings) (chartPath string, err error)
}

type uninstallAction interface {
	uninstall(releaseName string) (*release.UninstallReleaseResponse, error)
	raw() *action.Uninstall
}

type listReleasesAction interface {
	listReleases() ([]*release.Release, error)
	raw() *action.List
}

type getReleaseValuesAction interface {
	getReleaseValues(releaseName string) (map[string]interface{}, error)
	raw() *action.GetValues
}

type getReleaseAction interface {
	getRelease(releaseName string) (*release.Release, error)
	raw() *action.Get
}

type rollbackReleaseAction interface {
	rollbackRelease(releaseName string) error
	raw() *action.Rollback
}

type valuesOptions interface {
	MergeValues(p getter.Providers) (map[string]interface{}, error)
}
