package client

import (
	"context"

	"k8s.io/client-go/rest"
)

var helmClient *HelmClient

func ExampleNewClientFromRestConf() {
	opt := &RestConfClientOptions{
		Options: &Options{
			Namespace:        "default", // Change this to the namespace you wish the client to operate in.
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
			DebugLog: func(format string, v ...interface{}) {
				// Change this to your own logger. Default is 'log.Printf(format, v...)'.
			},
		},
		RestConfig: &rest.Config{},
	}

	helmClient, err := NewClientFromRestConf(opt)
	if err != nil {
		panic(err)
	}
	_ = helmClient
}

func ExampleHelmClient_InstallOrUpgradeChart() {
	// Define the chart to be installed
	chartSpec := ChartSpec{
		ReleaseName: "etcd-operator",
		ChartName:   "stable/etcd-operator",
		Namespace:   "default",
	}

	// Install a chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	if _, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_InstallOrUpgradeChart_useChartDirectory() {
	// Use an unpacked chart directory.
	chartSpec := ChartSpec{
		ReleaseName: "etcd-operator",
		ChartName:   "/path/to/stable/etcd-operator",
		Namespace:   "default",
	}

	if _, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_InstallOrUpgradeChart_useLocalChartArchive() {
	// Use an archived chart directory.
	chartSpec := ChartSpec{
		ReleaseName: "etcd-operator",
		ChartName:   "/path/to/stable/etcd-operator.tar.gz",
		Namespace:   "default",
	}

	if _, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_InstallOrUpgradeChart_useURL() {
	// Use an archived chart directory via URL.
	chartSpec := ChartSpec{
		ReleaseName: "etcd-operator",
		ChartName:   "http://helm.whatever.com/repo/etcd-operator.tar.gz",
		Namespace:   "default",
	}

	if _, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec); err != nil {
		panic(err)
	}
}

type customRollBack struct {
	HelmClient
}

var _ RollBack = &customRollBack{}

func (c customRollBack) RollbackRelease(spec *ChartSpec) error {
	client := c.actions.newRollbackRelease()

	client.raw().Force = true

	return client.rollbackRelease(spec.ReleaseName)
}

func ExampleHelmClient_UninstallRelease() {
	// Define the released chart to be installed.
	chartSpec := ChartSpec{
		ReleaseName: "etcd-operator",
		ChartName:   "stable/etcd-operator",
		Namespace:   "default",
	}

	// Uninstall the chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	if err := helmClient.UninstallRelease(&chartSpec); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_UninstallReleaseByName() {
	// Uninstall a release by name.
	if err := helmClient.UninstallReleaseByName("etcd-operator"); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_ListDeployedReleases() {
	// List all deployed releases.
	if _, err := helmClient.ListDeployedReleases(); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_GetReleaseValues() {
	// Get the values of a deployed release.
	if _, err := helmClient.GetReleaseValues("etcd-operator", true); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_GetRelease() {
	// Get specific details of a deployed release.
	if _, err := helmClient.GetRelease("etcd-operator"); err != nil {
		panic(err)
	}
}

func ExampleHelmClient_RollbackRelease() {
	// Define the released chart to be installed
	chartSpec := ChartSpec{
		ReleaseName: "etcd-operator",
		ChartName:   "stable/etcd-operator",
		Namespace:   "default",
	}

	// Rollback to the previous version of the release.
	if err := helmClient.RollbackRelease(&chartSpec); err != nil {
		return
	}
}
