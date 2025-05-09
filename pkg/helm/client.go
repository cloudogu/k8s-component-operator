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
			InsecureTls:      helmRepoData.InsecureTLS,
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
	chart.ChartName = c.patchOciEndpoint(chart.ChartName)

	if chart.Version == "" {
		return fmt.Errorf("cannot install chart %q without version", chart.ChartName)
	}

	_, err := c.helmClient.InstallOrUpgradeChart(ctx, chart)
	if err != nil {
		return fmt.Errorf("error while installOrUpgrade chart %s: %w", chart.ChartName, err)
	}

	return nil
}

// SatisfiesDependencies checks if all dependencies are satisfied in terms of installation and version.
func (c *Client) SatisfiesDependencies(ctx context.Context, chart *client.ChartSpec) error {
	logger := log.FromContext(ctx)
	logger.Info("Checking if components dependencies are satisfied", "component", chart.ChartName)

	chart.ChartName = c.patchOciEndpoint(chart.ChartName)

	if chart.Version == "" {
		return fmt.Errorf("cannot install chart %q without version", chart.ChartName)
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
		"plainHTTP", c.helmRepoData.PlainHttp,
		"insecureTLS", c.helmRepoData.InsecureTLS)

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

// GetReleaseValues returns the (optionally, all computed) values for the specified release.
func (c *Client) GetReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	return c.helmClient.GetReleaseValues(name, allValues)
}

// GetDeployedReleaseVersion returns the version for the specified release (if the release exists).
func (c *Client) GetDeployedReleaseVersion(ctx context.Context, name string) (string, error) {
	logger := log.FromContext(ctx)

	deployedReleases, err := c.ListDeployedReleases()
	if err != nil {
		return "", fmt.Errorf("could not list deployed Helm releases: %w", err)
	}

	for _, k8sRelease := range deployedReleases {
		if name == k8sRelease.Name {
			version := k8sRelease.Chart.AppVersion()
			logger.Info("Found existing release for reconciled component",
				"component", name, "version", version)
			return version, nil
		}
	}

	logger.Info("could not find a deployed release for component: ", name)
	return "", nil
}

// GetChartSpecValues returns the additional values for the specified ChartSpec.
func (c *Client) GetChartSpecValues(spec *client.ChartSpec) (map[string]interface{}, error) {
	return c.helmClient.GetChartSpecValues(spec)
}

func (c *Client) patchOciEndpoint(chartName string) string {
	if strings.HasPrefix(chartName, ociSchemePrefix) {
		return chartName
	}

	return fmt.Sprintf("%s/%s", c.helmRepoData.URL(), chartName)
}

func (c *Client) GetLatestVersion(chartName string) (string, error) {
	ref := strings.TrimPrefix(c.patchOciEndpoint(chartName), ociSchemePrefix)
	tags, err := c.helmClient.Tags(ref)
	if err != nil {
		return "", fmt.Errorf("error resolving tags for chart %s: %w", chartName, err)
	}

	//sort tags by version
	sortedTags := sortByVersionDescending(tags)

	if len(sortedTags) <= 0 {
		return "", fmt.Errorf("could not find any tags for chart %s", chartName)
	}

	// set version to the latest tag
	return sortedTags[0], nil
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
