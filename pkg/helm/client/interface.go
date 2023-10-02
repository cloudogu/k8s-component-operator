package client

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

// Client holds the method signatures for a Helm client.
// NOTE: This is an interface to allow for mocking in tests.
type Client interface {
	InstallOrUpgradeChart(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error)
	InstallChart(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error)
	UpgradeChart(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error)
	ListDeployedReleases() ([]*release.Release, error)
	ListReleasesByStateMask(action.ListStates) ([]*release.Release, error)
	GetRelease(name string) (*release.Release, error)
	// RollBack is an interface to abstract a rollback action.
	RollBack
	GetReleaseValues(name string, allValues bool) (map[string]interface{}, error)
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
