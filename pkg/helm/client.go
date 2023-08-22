package helm

import (
	"context"
	"fmt"
	"strings"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	ctrl "sigs.k8s.io/controller-runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	helmRepositoryCache    = "/tmp/.helmcache"
	helmRepositoryConfig   = "/tmp/.helmrepo"
	helmRegistryConfigFile = "/tmp/.helmregistry/config.json"
)

// HelmClient embeds the helmclient.Client interface for usage in this package.
type HelmClient interface {
	helmclient.Client
}

// OciRepositoryConfig can get an OCI-Endpoint for a helm-repository.
type OciRepositoryConfig interface {
	GetOciEndpoint() (string, error)
}

// Client wraps the HelmClient with config.HelmRepositoryData
type Client struct {
	helmClient        HelmClient
	helmRepoData      OciRepositoryConfig
	actionConfig      *action.Configuration
	dependencyChecker dependencyChecker
}

// NewClient create a new instance of the helm client.
func NewClient(namespace string, helmRepoData OciRepositoryConfig, debug bool, debugLog action.DebugLog) (*Client, error) {
	opt := &helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        namespace,
			RepositoryCache:  helmRepositoryCache,
			RepositoryConfig: helmRepositoryConfig,
			RegistryConfig:   helmRegistryConfigFile,
			Debug:            debug,
			DebugLog:         debugLog,
			Linting:          true,
			PlainHttp:        helmRepoData.PlainHttp,
		},
		RestConfig: ctrl.GetConfigOrDie(),
	}

	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm client: %w", err)
	}

	actionConfig, err := createActionConfig(namespace, helmRepoData.PlainHttp, debug, debugLog, opt.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm client: %w", err)
	}

	return &Client{
		helmClient:        helmClient,
		helmRepoData:      helmRepoData,
		actionConfig:      actionConfig,
		dependencyChecker: &installedDependencyChecker{},
	}, nil
}

func createActionConfig(namespace string, plainHttp bool, debug bool, debugLog action.DebugLog, restConfig *rest.Config) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	clientGetter := helmclient.NewRESTClientGetter(namespace, nil, restConfig)
	err := actionConfig.Init(
		clientGetter,
		namespace,
		"secret",
		debugLog,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init actionConfig: %w", err)
	}

	clientOpts := []registry.ClientOption{
		registry.ClientOptDebug(debug),
		registry.ClientOptCredentialsFile(helmRegistryConfigFile),
	}

	if plainHttp {
		clientOpts = append(clientOpts, registry.ClientOptPlainHTTP())
	}

	helmRegistryClient, err := registry.NewClient(clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm registry client: %w", err)
	}

	actionConfig.RegistryClient = helmRegistryClient
	return actionConfig, nil
}

// InstallOrUpgrade takes a helmChart and applies it.
func (c *Client) InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error {
	// This helm-client currently only works with OCI-Helm-Repositories.
	// Therefore the chartName has to include the FQDN of the repository (e.g. "oci://my.repo/...")
	// If in the future non-oci-repositories need to be used, this should be done here...
	err := c.patchOciEndpoint(chart)
	if err != nil {
		return fmt.Errorf("error while patching chart '%s': %w", chart.ChartName, err)
	}

	_, err = c.helmClient.InstallOrUpgradeChart(ctx, chart, nil)
	if err != nil {
		return fmt.Errorf("error while installOrUpgrade chart %s: %w", chart.ChartName, err)
	}

	return nil
}

// SatisfiesDependencies checks if all dependencies are satisfied in terms of installation and version.
func (c *Client) SatisfiesDependencies(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)
	logger.Info("Checking if components dependencies are satisfied", "component", component.Name)

	endpoint, err := c.getOciEndpoint(component)
	if err != nil {
		return err
	}

	chartSpec := component.GetHelmChartSpec(endpoint)

	componentChart, err := c.getChart(ctx, component, chartSpec)
	if err != nil {
		return fmt.Errorf("failed to get chart for component %s: %w", component, err)
	}

	dependencies := componentChart.Metadata.Dependencies
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

func (c *Client) getChart(ctx context.Context, component *k8sv1.Component, spec *helmclient.ChartSpec) (*chart.Chart, error) {
	logger := log.FromContext(ctx)

	// TODO extract into helper method
	// We need this installAction because it sets the registryClient in ChartPathOptions which is a private field.
	install := action.NewInstall(c.actionConfig)
	install.Version = component.Spec.Version
	install.PlainHTTP = c.helmRepoData.PlainHttp

	logger.Info("Trying to get chart with options",
		"chart", spec.ChartName,
		"version", component.Spec.Version,
		"plain http", c.helmRepoData.PlainHttp)

	componentChart, _, err := c.helmClient.GetChart(spec.ChartName, &install.ChartPathOptions)
	if err != nil {
		return nil, fmt.Errorf("error while getting chart for %s:%s: %w", component.Spec.Name, component.Spec.Version, err)
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

func (c *Client) patchOciEndpoint(chart *helmclient.ChartSpec) error {
	if strings.Index(chart.ChartName, "oci://") == 0 {
		// oci protocol already present -> nothing to do
		return nil
	}

	endpoint, err := c.helmRepoData.GetOciEndpoint()
	if err != nil {
		return fmt.Errorf("error while getting oci endpoint: %w", err)
	}

	chart.ChartName = fmt.Sprintf("%s/%s", endpoint, chart.ChartName)

	return nil
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
