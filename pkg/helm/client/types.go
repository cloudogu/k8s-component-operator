package client

import (
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"io"
	"time"

	"k8s.io/client-go/rest"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
)

// Type Guard asserting that HelmClient satisfies the HelmClient interface.
var _ Client = &HelmClient{}

// KubeConfClientOptions defines the options used for constructing a client via kubeconfig.
type KubeConfClientOptions struct {
	*Options
	KubeContext string
	KubeConfig  []byte
}

// RestConfClientOptions defines the options used for constructing a client via REST config.
type RestConfClientOptions struct {
	*Options
	RestConfig *rest.Config
}

// Options defines the options of a client. If Output is not set, os.Stdout will be used.
type Options struct {
	Namespace        string
	RepositoryConfig string
	RepositoryCache  string
	Debug            bool
	DebugLog         action.DebugLog
	RegistryConfig   string
	Output           io.Writer
	// PlainHttp forces the registry client to establish plain http connections.
	PlainHttp bool
}

// RESTClientOption is a function that can be used to set the RESTClientOptions of a HelmClient.
type RESTClientOption func(*rest.Config)

// RESTClientGetter defines the values of a helm REST client.
type RESTClientGetter struct {
	namespace  string
	kubeConfig []byte
	restConfig *rest.Config

	opts []RESTClientOption
}

// HelmClient Client defines the values of a helm client.
type HelmClient struct {
	TagResolver
	// Settings defines the environment settings of a client.
	Settings  *cli.EnvSettings
	Providers getter.Providers
	actions   actionProvider
	output    io.Writer
	DebugLog  action.DebugLog
}

type GenericHelmOptions struct {
	RollBack RollBack
	// PlainHttp forces use of plain HTTP for communication with the helm registry.
	PlainHttp bool
}

type HelmTemplateOptions struct {
	KubeVersion *chartutil.KubeVersion
	// APIVersions defined here will be appended to the default list helm provides
	APIVersions chartutil.VersionSet
}

// ChartSpec defines the values of a helm chart
// +kubebuilder:object:generate:=true
type ChartSpec struct {
	ReleaseName string `json:"release"`
	ChartName   string `json:"chart"`
	// Namespace where the chart release is deployed.
	// Note that client.Options.Namespace should ideally match the namespace configured here.
	Namespace string `json:"namespace"`
	// ValuesYaml is the values.yaml content.
	// use string instead of map[string]interface{}
	// https://github.com/kubernetes-sigs/kubebuilder/issues/528#issuecomment-466449483
	// and https://github.com/kubernetes-sigs/controller-tools/pull/317
	// +optional
	ValuesYaml string `json:"valuesYaml,omitempty"`
	// Specify values similar to the cli
	// +optional
	ValuesOptions values.Options `json:"valuesOptions,omitempty"`
	// Version of the chart release.
	// +optional
	Version string `json:"version,omitempty"`
	// CreateNamespace indicates whether to create the namespace if it does not exist.
	// +optional
	CreateNamespace bool `json:"createNamespace,omitempty"`
	// Timeout configures the time to wait for any individual Kubernetes operation (like Jobs for hooks).
	// +optional
	Timeout time.Duration `json:"timeout,omitempty"`
	// Atomic indicates whether to install resources atomically.
	// 'Wait' will automatically be set to true when using Atomic.
	// +optional
	Atomic bool `json:"atomic,omitempty"`
	// ResetValues indicates whether to reset the values.yaml file during installation.
	// +optional
	ResetValues bool `json:"resetValues,omitempty"`
	// ReuseValues indicates whether to reuse the values.yaml file during installation.
	// +optional
	ReuseValues bool `json:"reuseValues,omitempty"`
	// CleanupOnFail indicates whether to cleanup the release on failure.
	// +optional
	CleanupOnFail bool `json:"cleanupOnFail,omitempty"`
}
