package helm

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	ctrl "sigs.k8s.io/controller-runtime"
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
	helmClient   HelmClient
	helmRepoData *config.HelmRepositoryData
}

// NewClient create a new instance of the helm client.
func NewClient(namespace string, helmRepoSecret *config.HelmRepositoryData, debug bool, debugLog action.DebugLog) (*Client, error) {
	opt := &helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        namespace,
			RepositoryCache:  helmRepositoryCache,
			RepositoryConfig: helmRepositoryConfig,
			RegistryConfig:   helmRegistryConfigFile,
			Debug:            debug,
			DebugLog:         debugLog,
			Linting:          true,
		},
		RestConfig: ctrl.GetConfigOrDie(),
	}

	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create helm client: %w", err)
	}

	return &Client{helmClient: helmClient, helmRepoData: helmRepoSecret}, nil
}

// InstallOrUpgrade takes a component and applies the corresponding helmChart.
func (c *Client) InstallOrUpgrade(ctx context.Context, component *k8sv1.Component) error {
	endpoint, err := c.helmRepoData.GetOciEndpoint()
	if err != nil {
		return fmt.Errorf("error while getting oci endpoint for %s: %w", component.Spec.Name, err)
	}

	_, err = c.helmClient.InstallOrUpgradeChart(ctx, component.GetHelmChartSpec(endpoint), nil)
	if err != nil {
		return fmt.Errorf("error while installOrUpgrade chart %s: %w", component.Spec.Name, err)
	}
	return nil
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
