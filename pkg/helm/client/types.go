package client

import (
	"io"
	"time"

	"helm.sh/helm/v3/pkg/postrender"

	"k8s.io/client-go/rest"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
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
	// PlainHttp forces the registry client to establish plain http connections. This option will override by InsecureTls by using HTTP traffic.
	PlainHttp bool
	// InsecureTls allows invalid or selfsigned certificates to be used. This option may be overridden by PlainHttp which forces HTTP traffic.
	InsecureTls bool
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
	Settings *cli.EnvSettings
	actions  actionProvider
	output   io.Writer
	DebugLog action.DebugLog
}

type HelmTemplateOptions struct {
	KubeVersion *chartutil.KubeVersion
	// APIVersions defined here will be appended to the default list helm provides
	APIVersions chartutil.VersionSet
}

// ChartSpec defines the values of a helm chart
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
	// MappedValuesYaml is the values.yaml content, but generated from component-values-metadata.yaml file
	// use string instead of map[string]interface{}
	// https://github.com/kubernetes-sigs/kubebuilder/issues/528#issuecomment-466449483
	// and https://github.com/kubernetes-sigs/controller-tools/pull/317
	// +optional
	MappedValuesYaml string `json:"mappedValuesYaml,omitempty"`
	// ValuesConfigRef is used for configuration
	// +optional
	ValuesConfigRefYaml string `json:"valuesConfigRefYaml,omitempty"`
	// Specify values similar to the cli
	// +optional
	ValuesOptions valuesOptions `json:"valuesOptions,omitempty"`
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
	// PostRenderer can be used to apply transformations to kubernetes resources
	// on installation and upgrade after rendering the templates
	// +optional
	PostRenderer postrender.PostRenderer
}
