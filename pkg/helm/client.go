package helm

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
)

const (
	helmRepositoryCache    = "/tmp/.helmcache"
	helmRepositoryConfig   = "/tmp/.helmrepo"
	helmRegistryConfigFile = "/tmp/.helmregistry/config.json"
	ociSchemePrefix        = string(config.EndpointSchemaOCI + "://")
)

// HelmClient embeds the client.Client interface for usage in this package.
type HelmClient interface {
	client.Client
}

// Client wraps the HelmClient with config.HelmRepositoryData
type Client struct {
	helmClient        HelmClient
	helmRepoData      *config.HelmRepositoryData
	dependencyChecker dependencyChecker
}

// NewClient create a new instance of the helm client.
func NewClient(namespace string, helmRepoData *config.HelmRepositoryData, debug bool, debugLog action.DebugLog) (*Client, error) {
	opt := &client.RestConfClientOptions{
		Options: &client.Options{
			Namespace:        namespace,
			RepositoryCache:  helmRepositoryCache,
			RepositoryConfig: helmRepositoryConfig,
			RegistryConfig:   helmRegistryConfigFile,
			Debug:            debug,
			DebugLog:         debugLog,
			PlainHttp:        helmRepoData.PlainHttp,
		},
		RestConfig: ctrl.GetConfigOrDie(),
	}

	helmClient, err := client.NewClientFromRestConf(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm client: %w", err)
	}

	return &Client{
		helmClient:        helmClient,
		helmRepoData:      helmRepoData,
		dependencyChecker: &installedDependencyChecker{},
	}, nil
}

// InstallOrUpgrade takes a helmChart and applies it.
func (c *Client) InstallOrUpgrade(ctx context.Context, chart *client.ChartSpec) error {
	// This helm-client currently only works with OCI-Helm-Repositories.
	// Therefore, the chartName has to include the FQDN of the repository (e.g. "oci://my.repo/...")
	// If in the future non-oci-repositories need to be used, this should be done here...
	c.patchOciEndpoint(chart)

	if err := c.patchChartVersion(chart); err != nil {
		return fmt.Errorf("error patching chart-version for chart %s: %w", chart.ChartName, err)
	}

	_, err := c.helmClient.InstallOrUpgradeChart(ctx, chart, nil)
	if err != nil {
		return fmt.Errorf("error while installOrUpgrade chart %s: %w", chart.ChartName, err)
	}

	return nil
}

// SatisfiesDependencies checks if all dependencies are satisfied in terms of installation and version.
func (c *Client) SatisfiesDependencies(ctx context.Context, chart *client.ChartSpec) error {
	logger := log.FromContext(ctx)
	logger.Info("Checking if components dependencies are satisfied", "component", chart.ChartName)

	c.patchOciEndpoint(chart)

	if err := c.patchChartVersion(chart); err != nil {
		return fmt.Errorf("error patching chart-version for chart %s: %w", chart.ChartName, err)
	}

	componentChart, err := c.getChart(ctx, chart)
	if err != nil {
		return fmt.Errorf("failed to get chart %s: %w", chart.ChartName, err)
	}

	dependencies := getComponentDependencies(componentChart)
	deployedReleases, err := c.ListDeployedReleases()
	if err != nil {
		return fmt.Errorf("failed to list deployed releases: %w", err)
	}

	err = c.dependencyChecker.CheckSatisfied(dependencies, deployedReleases)
	if err != nil {
		return &dependencyUnsatisfiedError{err: err}
	}

	return nil
}

func (c *Client) getChart(ctx context.Context, chartSpec *client.ChartSpec) (*chart.Chart, error) {
	logger := log.FromContext(ctx)

	logger.Info("Trying to get chart with options",
		"chart", chartSpec.ChartName,
		"version", chartSpec.Version,
		"plain http", c.helmRepoData.PlainHttp)

	componentChart, _, err := c.helmClient.GetChart(chartSpec)
	if err != nil {
		return nil, fmt.Errorf("error while getting chart for %s:%s: %w", chartSpec.ChartName, chartSpec.Version, err)
	}

	return componentChart, nil
}

// Uninstall removes the helmRelease for the given name
func (c *Client) Uninstall(releaseName string) error {
	if err := c.helmClient.UninstallReleaseByName(releaseName); err != nil {
		return fmt.Errorf("error while uninstalling helm-release %s: %w", releaseName, err)
	}
	return nil
}

// ListDeployedReleases returns all deployed helm releases
func (c *Client) ListDeployedReleases() ([]*release.Release, error) {
	return c.helmClient.ListDeployedReleases()
}

func (c *Client) patchOciEndpoint(chart *client.ChartSpec) {
	if strings.HasPrefix(chart.ChartName, ociSchemePrefix) {
		return
	}

	chart.ChartName = fmt.Sprintf("%s/%s", c.helmRepoData.URL(), chart.ChartName)
}

func (c *Client) patchChartVersion(chart *client.ChartSpec) error {
	if chart.Version != "" {
		return nil
	}

	ref := strings.TrimPrefix(chart.ChartName, ociSchemePrefix)
	tags, err := c.helmClient.Tags(ref)
	if err != nil {
		return fmt.Errorf("error resolving tags for chart %s: %w", chart.ChartName, err)
	}

	//sort tags by version
	sortedTags := sortByVersionDescending(tags)

	if len(sortedTags) <= 0 {
		return fmt.Errorf("could not find any tags for chart %s", chart.ChartName)
	}

	// set version to the latest tag
	chart.Version = sortedTags[0]

	return nil
}

func sortByVersionDescending(tags []string) []string {
	versions := make([]core.Version, 0)
	for _, tag := range tags {
		v, err := core.ParseVersion(tag)
		if err == nil {
			versions = append(versions, v)
		}
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].IsNewerOrEqualThan(versions[j])
	})

	result := make([]string, len(versions))
	for i, version := range versions {
		result[i] = version.Raw
	}

	return result
}

type dependencyUnsatisfiedError struct {
	err error
}

// Error returns the string representation of the wrapped error.
func (due *dependencyUnsatisfiedError) Error() string {
	return fmt.Sprintf("one or more dependencies are not satisfied: %s", due.err.Error())
}

// Unwrap returns the root error.
func (due *dependencyUnsatisfiedError) Unwrap() error {
	return due.err
}
