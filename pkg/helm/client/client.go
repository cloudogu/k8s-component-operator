package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

var storage = repo.File{}

const (
	defaultCachePath            = "/tmp/.helmcache"
	defaultRepositoryConfigPath = "/tmp/.helmrepo"
)

// NewClientFromRestConf returns a new Helm client constructed with the provided REST config options.
func NewClientFromRestConf(options *RestConfClientOptions) (Client, error) {
	settings := cli.New()

	clientGetter := NewRESTClientGetter(options.Namespace, nil, options.RestConfig)

	return newClient(options.Options, clientGetter, settings)
}

// newClient is used by both NewClientFromKubeConf and NewClientFromRestConf
// and returns a new Helm client via the provided options and REST config.
func newClient(options *Options, clientGetter genericclioptions.RESTClientGetter, settings *cli.EnvSettings) (Client, error) {
	err := setEnvSettings(&options, settings)
	if err != nil {
		return nil, err
	}

	debugLog := options.DebugLog
	if debugLog == nil {
		debugLog = func(format string, v ...interface{}) {
			log.Printf(format, v...)
		}
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

	clientOpts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}

	if options.PlainHttp {
		clientOpts = append(clientOpts, registry.ClientOptPlainHTTP())
	}

	registryClient, err := registry.NewClient(clientOpts...)
	if err != nil {
		return nil, err
	}
	actionConfig.RegistryClient = registryClient

	return &HelmClient{
		TagResolver:  registryClient,
		Settings:     settings,
		Providers:    getter.All(settings),
		storage:      &storage,
		ActionConfig: actionConfig,
		plainHttp:    options.PlainHttp,
		DebugLog:     debugLog,
		output:       options.Output,
	}, nil
}

// setEnvSettings sets the client's environment settings based on the provided client configuration.
func setEnvSettings(ppOptions **Options, settings *cli.EnvSettings) error {
	if *ppOptions == nil {
		*ppOptions = &Options{
			RepositoryConfig: defaultRepositoryConfigPath,
			RepositoryCache:  defaultCachePath,
			Linting:          true,
		}
	}

	options := *ppOptions

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
func (c *HelmClient) InstallOrUpgradeChart(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error) {
	exists, err := c.chartExists(spec)
	if err != nil {
		return nil, err
	}

	if exists {
		return c.upgrade(ctx, spec, opts)
	}

	return c.install(ctx, spec, opts)
}

// InstallChart installs the provided chart and returns the corresponding release.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) InstallChart(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error) {
	return c.install(ctx, spec, opts)
}

// UpgradeChart upgrades the provided chart and returns the corresponding release.
// Namespace and other context is provided via the client.Options struct when instantiating a client.
func (c *HelmClient) UpgradeChart(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error) {
	return c.upgrade(ctx, spec, opts)
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
// Optionally lints the chart if the linting flag is set.
func (c *HelmClient) install(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error) {
	client := action.NewInstall(c.ActionConfig)
	mergeInstallOptions(spec, client)

	// NameAndChart returns either the TemplateName if set,
	// the ReleaseName if set or the generatedName as the first return value.
	releaseName, _, err := client.NameAndChart([]string{spec.ChartName})
	if err != nil {
		return nil, err
	}
	client.ReleaseName = releaseName

	if client.Version == "" {
		client.Version = ">0.0.0-0"
	}

	if opts != nil {
		if opts.PostRenderer != nil {
			client.PostRenderer = opts.PostRenderer
		}

		client.PlainHTTP = opts.PlainHttp
	}

	helmChart, _, err := c.GetChart(spec)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	rel, err := client.RunWithContext(ctx, helmChart, values)
	if err != nil {
		return rel, err
	}

	c.DebugLog("release installed successfully: %s/%s-%s", rel.Name, rel.Chart.Metadata.Name, rel.Chart.Metadata.Version)

	return rel, nil
}

// upgrade upgrades a chart and CRDs.
// Optionally lints the chart if the linting flag is set.
func (c *HelmClient) upgrade(ctx context.Context, spec *ChartSpec, opts *GenericHelmOptions) (*release.Release, error) {
	client := action.NewUpgrade(c.ActionConfig)
	mergeUpgradeOptions(spec, client)
	client.Install = true

	if client.Version == "" {
		client.Version = ">0.0.0-0"
	}

	if opts != nil {
		if opts.PostRenderer != nil {
			client.PostRenderer = opts.PostRenderer
		}

		client.PlainHTTP = opts.PlainHttp
	}

	helmChart, _, err := c.GetChart(spec)
	if err != nil {
		return nil, err
	}

	p := getter.All(c.Settings)
	values, err := spec.GetValuesMap(p)
	if err != nil {
		return nil, err
	}

	upgradedRelease, upgradeErr := client.RunWithContext(ctx, spec.ReleaseName, helmChart, values)
	if upgradeErr != nil {
		resultErr := upgradeErr
		if upgradedRelease == nil && opts != nil && opts.RollBack != nil {
			rollbackErr := opts.RollBack.RollbackRelease(spec)
			if rollbackErr != nil {
				resultErr = fmt.Errorf("release failed, rollback failed: release error: %w, rollback error: %v", upgradeErr, rollbackErr)
			} else {
				resultErr = fmt.Errorf("release failed, rollback succeeded: release error: %w", upgradeErr)
			}
		}
		c.DebugLog("release upgrade failed: %s", resultErr)
		return nil, resultErr
	}

	c.DebugLog("release upgraded successfully: %s/%s-%s", upgradedRelease.Name, upgradedRelease.Chart.Metadata.Name, upgradedRelease.Chart.Metadata.Version)

	return upgradedRelease, nil
}

// uninstallRelease uninstalls the provided release.
func (c *HelmClient) uninstallRelease(spec *ChartSpec) error {
	client := action.NewUninstall(c.ActionConfig)

	mergeUninstallReleaseOptions(spec, client)

	resp, err := client.Run(spec.ReleaseName)
	if err != nil {
		return err
	}

	c.DebugLog("release uninstalled, response: %v", resp)

	return nil
}

// uninstallReleaseByName uninstalls a release identified by the provided 'name'.
func (c *HelmClient) uninstallReleaseByName(name string) error {
	client := action.NewUninstall(c.ActionConfig)

	resp, err := client.Run(name)
	if err != nil {
		return err
	}

	c.DebugLog("release uninstalled, response: %v", resp)

	return nil
}

// GetChart returns a chart matching the provided chart name and options.
func (c *HelmClient) GetChart(spec *ChartSpec) (*chart.Chart, string, error) {
	install := action.NewInstall(c.ActionConfig)
	install.Version = spec.Version
	install.PlainHTTP = c.plainHttp

	chartPath, err := install.ChartPathOptions.LocateChart(spec.ChartName, c.Settings)
	if err != nil {
		return nil, "", err
	}

	helmChart, err := loader.Load(chartPath)
	if err != nil {
		return nil, "", err
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
		return false, err
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
	listClient := action.NewList(c.ActionConfig)
	listClient.StateMask = state

	return listClient.Run()
}

// getReleaseValues returns the values for the provided release 'name'.
// If allValues = true is specified, all computed values are returned.
func (c *HelmClient) getReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	getReleaseValuesClient := action.NewGetValues(c.ActionConfig)

	getReleaseValuesClient.AllValues = allValues

	return getReleaseValuesClient.Run(name)
}

// getRelease returns a release matching the provided 'name'.
func (c *HelmClient) getRelease(name string) (*release.Release, error) {
	getReleaseClient := action.NewGet(c.ActionConfig)

	return getReleaseClient.Run(name)
}

// rollbackRelease implicitly rolls back a release to the last revision.
func (c *HelmClient) rollbackRelease(spec *ChartSpec) error {
	client := action.NewRollback(c.ActionConfig)

	mergeRollbackOptions(spec, client)

	return client.Run(spec.ReleaseName)
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
}

// mergeUninstallReleaseOptions merges values of the provided chart to helm uninstall options used by the client.
func mergeUninstallReleaseOptions(chartSpec *ChartSpec, uninstallReleaseOptions *action.Uninstall) {
	uninstallReleaseOptions.Timeout = chartSpec.Timeout
}
