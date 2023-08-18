package helm

import (
	"context"
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
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
	helmClient   HelmClient
	helmRepoData OciRepositoryConfig
}

// NewClient create a new instance of the helm client.
func NewClient(namespace string, helmRepoSecret OciRepositoryConfig, debug bool, debugLog action.DebugLog) (*Client, error) {
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

// InstallOrUpgrade takes a helmChart and applies it.
func (c *Client) InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error {
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
