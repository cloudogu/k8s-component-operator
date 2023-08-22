package helm

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"k8s.io/client-go/rest"

	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	ctrl "sigs.k8s.io/controller-runtime"
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

// Client wraps the HelmClient with config.HelmRepositoryData
type Client struct {
	helmClient        HelmClient
	helmRepoData      *config.HelmRepositoryData
	actionConfig      *action.Configuration
	dependencyChecker dependencyChecker
}

// NewClient create a new instance of the helm client.
func NewClient(namespace string, helmRepoData *config.HelmRepositoryData, debug bool, debugLog action.DebugLog) (*Client, error) {
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

// InstallOrUpgrade takes a component and applies the corresponding helmChart.
func (c *Client) InstallOrUpgrade(ctx context.Context, component *k8sv1.Component) error {
	endpoint, err := c.getOciEndpoint(component)
	if err != nil {
		return err
	}

	chartSpec := component.GetHelmChartSpec(endpoint)

	_, err = c.helmClient.InstallOrUpgradeChart(ctx, chartSpec, &helmclient.GenericHelmOptions{PlainHttp: c.helmRepoData.PlainHttp})
	if err != nil {
		return fmt.Errorf("error while installing/upgrading component %s: %w", component, err)
	}
	return nil
}

// SatisfiesDependencies checks if all dependencies are satisfied in terms of installation and ver
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

func (c *Client) getOciEndpoint(component *k8sv1.Component) (string, error) {
	endpoint, err := c.helmRepoData.GetOciEndpoint()
	if err != nil {
		return "", fmt.Errorf("error while getting oci endpoint for %s: %w", component.Spec.Name, err)
	}

	return endpoint, nil
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

// Uninstall removes the helmChart of the given component
func (c *Client) Uninstall(component *k8sv1.Component) error {
	if err := c.helmClient.UninstallReleaseByName(component.Spec.Name); err != nil {
		return fmt.Errorf("error while uninstalling helm-release %s: %w", component.Spec.Name, err)
	}
	return nil
}

// ListDeployedReleases returns all deployed helm releases
func (c *Client) ListDeployedReleases() ([]*release.Release, error) {
	return c.helmClient.ListDeployedReleases()
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
