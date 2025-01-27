package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
)

const (
	defaultCachePath            = "/tmp/.helmcache"
	defaultRepositoryConfigPath = "/tmp/.helmrepo"
)

const anyVersionConstraint = ">0.0.0-0"

var defaultDebugLog = func(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// NewClientFromRestConf returns a new Helm client constructed with the provided REST config options.
func NewClientFromRestConf(options *RestConfClientOptions) (Client, error) {
	settings := cli.New()

	clientGetter := NewRESTClientGetter(options.Namespace, nil, options.RestConfig)

	return newClient(options.Options, clientGetter, settings)
}

// newClient is used by both NewClientFromKubeConf and NewClientFromRestConf
// and returns a new Helm client via the provided options and REST config.
func newClient(options *Options, clientGetter genericclioptions.RESTClientGetter, settings *cli.EnvSettings) (Client, error) {
	err := setEnvSettings(options, settings)
	if err != nil {
		return nil, err
	}

	debugLog := options.DebugLog
	if debugLog == nil {
		debugLog = defaultDebugLog
	}

	if options.Output == nil {
		options.Output = os.Stdout
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(
		clientGetter,
		settings.Namespace(),
		"secret",
		debugLog,
	)
	if err != nil {
		return nil, err
	}

	registryClient, err := createRegistryClient(options, settings)
	if err != nil {
		return nil, err
	}
	actionConfig.RegistryClient = registryClient

	actionProvider := &provider{
		Configuration: actionConfig,
		plainHttp:     options.PlainHttp,
		insecureTls:   options.InsecureTls,
	}

	return &HelmClient{
		TagResolver: registryClient,
		Settings:    settings,
		actions:     actionProvider,
		DebugLog:    debugLog,
		output:      options.Output,
	}, nil
}

func createRegistryClient(options *Options, settings *cli.EnvSettings) (*registry.Client, error) {
	clientOpts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}

	var err error
	clientOpts, err = configureHttpRegistryClientOptions(options, clientOpts)
	if err != nil {
		return nil, err
	}

	return registry.NewClient(clientOpts...)
}

func configureHttpRegistryClientOptions(options *Options, clientOpts []registry.ClientOption) ([]registry.ClientOption, error) {
	if options.PlainHttp {
		clientOpts = append(clientOpts, registry.ClientOptPlainHTTP())
	}

	var httpTransport *http.Transport
	var err error
	httpTransport, err = getProxyTransportIfConfigured()
	if err != nil {
		return nil, err
	}

	httpTransport = configureTls(options, httpTransport)

	if httpTransport != nil {
		clientOpts = append(clientOpts, registry.ClientOptHTTPClient(&http.Client{Timeout: time.Second * 10, Transport: httpTransport}))
	}

	return clientOpts, nil
}

func configureTls(options *Options, transport *http.Transport) *http.Transport {
	if !options.PlainHttp && options.InsecureTls {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		if transport != nil {
			transport.TLSClientConfig = tlsConfig
		} else {
			transport = &http.Transport{TLSClientConfig: tlsConfig}
		}
	}

	return transport
}

func getProxyTransportIfConfigured() (*http.Transport, error) {
	proxyURL, found := os.LookupEnv("PROXY_URL")
	if !found || len(proxyURL) < 1 {
		return nil, nil
	}

	parsedProxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy: %w", err)
	}

	proxyFn := func(request *http.Request) (*url.URL, error) {
		return parsedProxy, nil
	}

	return &http.Transport{
		// From https://github.com/google/go-containerregistry/blob/c4dd792fa06c1f8b780ad90c8ab4f38b4eac05bd/pkg/v1/remote/options.go#L113
		DisableCompression: true,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 proxyFn,
	}, nil
}

// setEnvSettings sets the client's environment settings based on the provided client configuration.
func setEnvSettings(options *Options, settings *cli.EnvSettings) error {
	if options == nil {
		options = &Options{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
		}
	}

	// set the namespace with this ugly workaround because cli.EnvSettings.namespace is private
	// thank you helm!
	if options.Namespace != "" {
		pflags := pflag.NewFlagSet("", pflag.ContinueOnError)
		settings.AddFlags(pflags)
		err := pflags.Parse([]string{"-n", options.Namespace})
		if err != nil {
			return err
		}
	}

	if options.RepositoryConfig == "" {
		options.RepositoryConfig = defaultRepositoryConfigPath
	}

	if options.RepositoryCache == "" {
		options.RepositoryCache = defaultCachePath
	}

	settings.RepositoryCache = options.RepositoryCache
	settings.RepositoryConfig = options.RepositoryConfig
	settings.Debug = options.Debug

	if options.RegistryConfig != "" {
		settings.RegistryConfig = options.RegistryConfig
	}

	return nil
}

// InstallOrUpgradeChart installs or upgrades the provided chart and returns the corresponding release.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) InstallOrUpgradeChart(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	exists, err := c.chartExists(spec)
	if err != nil {
		return nil, err
	}

	if exists {
		return c.upgrade(ctx, spec)
	}

	return c.install(ctx, spec)
}

// InstallChart installs the provided chart and returns the corresponding release.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) InstallChart(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	return c.install(ctx, spec)
}

// UpgradeChart upgrades the provided chart and returns the corresponding release.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) UpgradeChart(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	return c.upgrade(ctx, spec)
}

// ListDeployedReleases lists all deployed releases.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) ListDeployedReleases() ([]*release.Release, error) {
	return c.listReleases(action.ListDeployed)
}

// ListReleasesByStateMask lists all releases filtered by stateMask.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) ListReleasesByStateMask(states action.ListStates) ([]*release.Release, error) {
	return c.listReleases(states)
}

// GetReleaseValues returns the (optionally, all computed) values for the specified release.
func (c *HelmClient) GetReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	return c.getReleaseValues(name, allValues)
}

// GetChartSpecValues returns the additional values for the specified ChartSpec.
func (c *HelmClient) GetChartSpecValues(spec *ChartSpec) (map[string]interface{}, error) {
	p := getter.All(c.Settings)
	additionalValuesYaml, err := spec.GetValuesMap(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get additional values.yaml-values from %s: %w", spec.ChartName, err)
	}

	return additionalValuesYaml, nil
}

// GetRelease returns a release specified by name.
func (c *HelmClient) GetRelease(name string) (*release.Release, error) {
	return c.getRelease(name)
}

// RollbackRelease implicitly rolls back a release to the last revision.
func (c *HelmClient) RollbackRelease(spec *ChartSpec) error {
	return c.rollbackRelease(spec)
}

// UninstallRelease uninstalls the provided release
func (c *HelmClient) UninstallRelease(spec *ChartSpec) error {
	return c.uninstallRelease(spec)
}

// UninstallReleaseByName uninstalls a release identified by the provided 'name'.
func (c *HelmClient) UninstallReleaseByName(name string) error {
	return c.uninstallReleaseByName(name)
}

// install installs the provided chart.
func (c *HelmClient) install(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	installAction := c.actions.newInstall()
	client := installAction.raw()
	mergeInstallOptions(spec, client)

	// NameAndChart returns either the TemplateName if set,
	// the ReleaseName if set or the generatedName as the first return value.
	releaseName, _, err := client.NameAndChart([]string{spec.ChartName})
	if err != nil {
		return nil, fmt.Errorf("failed to determine release name for chart %q: %w", spec.ChartName, err)
	}
	client.ReleaseName = releaseName

	if client.Version == "" {
		client.Version = anyVersionConstraint
	}

	helmChart, _, err := c.GetChart(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart for release %q: %w", spec.ReleaseName, err)
	}

	if helmChart.Metadata.Type != "" && helmChart.Metadata.Type != "application" {
		return nil, fmt.Errorf(
			"chart %q has an unsupported type and is not installable: %q",
			helmChart.Metadata.Name,
			helmChart.Metadata.Type,
		)
	}

	p := getter.All(c.Settings)
	values, err := spec.GetValuesMap(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get values for release %q: %w", spec.ReleaseName, err)
	}

	rel, err := installAction.install(ctx, helmChart, values)
	if err != nil {
		return nil, fmt.Errorf("failed to install release %q: %w", spec.ReleaseName, err)
	}

	c.DebugLog("release installed successfully: %s/%s-%s", rel.Name, rel.Chart.Metadata.Name, rel.Chart.Metadata.Version)

	return rel, nil
}

// upgrade upgrades a chart and CRDs.
func (c *HelmClient) upgrade(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	upgradeAction := c.actions.newUpgrade()
	client := upgradeAction.raw()
	mergeUpgradeOptions(spec, client)
	client.Install = true

	if client.Version == "" {
		client.Version = anyVersionConstraint
	}

	helmChart, _, err := c.GetChart(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to get chart for release %q: %w", spec.ReleaseName, err)
	}

	p := getter.All(c.Settings)
	values, err := spec.GetValuesMap(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get values for release %q: %w", spec.ReleaseName, err)
	}

	upgradedRelease, upgradeErr := upgradeAction.upgrade(ctx, spec.ReleaseName, helmChart, values)
	if upgradeErr != nil {
		c.DebugLog("release upgrade failed: %s", upgradeErr)
		return nil, fmt.Errorf("failed to upgrade release %q: %w", spec.ReleaseName, upgradeErr)
	}

	c.DebugLog("release upgraded successfully: %s/%s-%s", upgradedRelease.Name, upgradedRelease.Chart.Metadata.Name, upgradedRelease.Chart.Metadata.Version)

	return upgradedRelease, nil
}

// uninstallRelease uninstalls the provided release.
func (c *HelmClient) uninstallRelease(spec *ChartSpec) error {
	uninstallAction := c.actions.newUninstall()
	mergeUninstallReleaseOptions(spec, uninstallAction.raw())

	resp, err := uninstallAction.uninstall(spec.ReleaseName)
	if err != nil {
		return fmt.Errorf("failed to uninstall release %q: %w", spec.ReleaseName, err)
	}

	c.DebugLog("release uninstalled, response: %v", resp)

	return nil
}

// uninstallReleaseByName uninstalls a release identified by the provided 'name'.
func (c *HelmClient) uninstallReleaseByName(name string) error {
	uninstallAction := c.actions.newUninstall()

	resp, err := uninstallAction.uninstall(name)
	if err != nil {
		return fmt.Errorf("failed to uninstall release %q: %w", name, err)
	}

	c.DebugLog("release uninstalled, response: %v", resp)

	return nil
}

// GetChart returns a chart matching the provided chart name and options.
func (c *HelmClient) GetChart(spec *ChartSpec) (*chart.Chart, string, error) {
	locateAction := c.actions.newLocateChart()

	if spec.Version == "" {
		spec.Version = anyVersionConstraint
	}

	chartPath, err := locateAction.locateChart(spec.ChartName, spec.Version, c.Settings)
	if err != nil {
		return nil, "", fmt.Errorf("failed to locate chart %q with version %q: %w", spec.ChartName, spec.Version, err)
	}

	helmChart, err := loader.Load(chartPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load chart %q with version %q from path %q: %w", spec.ChartName, spec.Version, chartPath, err)
	}

	if helmChart.Metadata.Deprecated {
		c.DebugLog("WARNING: This chart (%q) is deprecated", helmChart.Metadata.Name)
	}

	return helmChart, chartPath, err
}

// chartExists checks whether a chart is already installed
// in a namespace or not based on the provided chart spec.
// Note that this function only considers the contained chart name and namespace.
func (c *HelmClient) chartExists(spec *ChartSpec) (bool, error) {
	releases, err := c.listReleases(action.ListAll)
	if err != nil {
		return false, fmt.Errorf("could not check if release %q is already installed: %w", spec.ReleaseName, err)
	}

	for _, r := range releases {
		if r.Name == spec.ReleaseName && r.Namespace == spec.Namespace {
			return true, nil
		}
	}

	return false, nil
}

// listReleases lists all releases that match the given state.
func (c *HelmClient) listReleases(state action.ListStates) ([]*release.Release, error) {
	listAction := c.actions.newListReleases()
	listAction.raw().StateMask = state

	releases, err := listAction.listReleases()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	return releases, nil
}

// getReleaseValues returns the values for the provided release 'name'.
// If allValues = true is specified, all computed values are returned.
func (c *HelmClient) getReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	getReleaseValuesAction := c.actions.newGetReleaseValues()
	getReleaseValuesAction.raw().AllValues = allValues

	values, err := getReleaseValuesAction.getReleaseValues(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get values of release %q: %w", name, err)
	}

	return values, nil
}

// getRelease returns a release matching the provided 'name'.
func (c *HelmClient) getRelease(name string) (*release.Release, error) {
	getReleaseAction := c.actions.newGetRelease()

	rel, err := getReleaseAction.getRelease(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get release %q: %w", name, err)
	}

	return rel, nil
}

// rollbackRelease implicitly rolls back a release to the last revision.
func (c *HelmClient) rollbackRelease(spec *ChartSpec) error {
	rollbackAction := c.actions.newRollbackRelease()
	mergeRollbackOptions(spec, rollbackAction.raw())

	err := rollbackAction.rollbackRelease(spec.ReleaseName)
	if err != nil {
		return fmt.Errorf("failed to rollback release %q: %w", spec.ReleaseName, err)
	}

	return nil
}

// mergeRollbackOptions merges values of the provided chart to helm rollback options used by the client.
func mergeRollbackOptions(chartSpec *ChartSpec, rollbackOptions *action.Rollback) {
	rollbackOptions.Timeout = chartSpec.Timeout
	rollbackOptions.CleanupOnFail = chartSpec.CleanupOnFail
}

// mergeInstallOptions merges values of the provided chart to helm install options used by the client.
func mergeInstallOptions(chartSpec *ChartSpec, installOptions *action.Install) {
	installOptions.CreateNamespace = chartSpec.CreateNamespace
	installOptions.Timeout = chartSpec.Timeout
	installOptions.Namespace = chartSpec.Namespace
	installOptions.ReleaseName = chartSpec.ReleaseName
	installOptions.Version = chartSpec.Version
	installOptions.Atomic = chartSpec.Atomic
	installOptions.PostRenderer = chartSpec.PostRenderer
}

// mergeUpgradeOptions merges values of the provided chart to helm upgrade options used by the client.
func mergeUpgradeOptions(chartSpec *ChartSpec, upgradeOptions *action.Upgrade) {
	upgradeOptions.Version = chartSpec.Version
	upgradeOptions.Namespace = chartSpec.Namespace
	upgradeOptions.Timeout = chartSpec.Timeout
	upgradeOptions.ResetValues = chartSpec.ResetValues
	upgradeOptions.ReuseValues = chartSpec.ReuseValues
	upgradeOptions.Atomic = chartSpec.Atomic
	upgradeOptions.CleanupOnFail = chartSpec.CleanupOnFail
	upgradeOptions.PostRenderer = chartSpec.PostRenderer
}

// mergeUninstallReleaseOptions merges values of the provided chart to helm uninstall options used by the client.
func mergeUninstallReleaseOptions(chartSpec *ChartSpec, uninstallReleaseOptions *action.Uninstall) {
	uninstallReleaseOptions.Timeout = chartSpec.Timeout
}
