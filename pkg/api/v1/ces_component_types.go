package v1

import (
	"context"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"helm.sh/helm/v3/pkg/chart"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/labels"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const defaultHelmClientTimeoutMins = time.Duration(15) * time.Minute

const (
	// ComponentStatusNotInstalled represents a status for a component that is not installed
	ComponentStatusNotInstalled = ""
	// ComponentStatusInstalling represents a status for a component that is currently being installed
	ComponentStatusInstalling = "installing"
	// ComponentStatusUpgrading represents a status for a component that is currently being upgraded
	ComponentStatusUpgrading = "upgrading"
	// ComponentStatusDeleting represents a status for a component that is currently being deleted
	ComponentStatusDeleting = "deleting"
	// ComponentStatusInstalled represents a status for a component that was successfully installed
	ComponentStatusInstalled = "installed"
	// ComponentStatusTryToInstall represents a status for a component that is not installed but its install process is in requeue loop.
	ComponentStatusTryToInstall = "tryToInstall"
	// ComponentStatusTryToUpgrade represents a status for a component that is installed but its actual upgrade process is in requeue loop.
	// In this state the component can be healthy but the version in the spec is not installed.
	ComponentStatusTryToUpgrade = "tryToUpgrade"
	// ComponentStatusTryToDelete represents a status for a component that is installed but its delete process is in requeue loop.
	// In this state the component can be healthy.
	ComponentStatusTryToDelete = "tryToDelete"
)

var _ webhook.Validator = &Component{}

const FinalizerName = "component-finalizer"

const (
	ComponentNameLabelKey    = "k8s.cloudogu.com/component.name"
	ComponentVersionLabelKey = "k8s.cloudogu.com/component.version"
)

// ComponentSpec defines the desired state of a component.
type ComponentSpec struct {
	// Namespace of the component (e.g. k8s)
	Namespace string `json:"namespace,omitempty"`
	// Name of the component (e.g. k8s-dogu-operator)
	Name string `json:"name,omitempty"`
	// Desired version of the component (e.g. 2.4.48-3)
	Version string `json:"version,omitempty"`
	// Mapped values from component cr to specific component attributes, usually log levels
	MappedValues map[string]string `json:"mappedValues,omitempty"`
	// DeployNamespace is the namespace where the helm chart should be deployed in.
	// This value is optional. If it is empty the operator deploys the helm chart in the namespace where the operator is deployed.
	DeployNamespace string `json:"deployNamespace,omitempty"`
	// ValuesYamlOverwrite is a multiline-yaml string that is applied alongside the original values.yaml-file of the component.
	// It can be used to overwrite specific configurations. Lists are overwritten, maps are merged.
	// +optional
	ValuesYamlOverwrite string `json:"valuesYamlOverwrite,omitempty"`
}

type HealthStatus string

const (
	PendingHealthStatus     HealthStatus = ""
	AvailableHealthStatus   HealthStatus = "available"
	UnavailableHealthStatus HealthStatus = "unavailable"
	UnknownHealthStatus     HealthStatus = "unknown"
	mappingMetadataFileName              = "component-values-metadata.yaml"
)

// ComponentStatus defines the observed state of a Component.
type ComponentStatus struct {
	// Status represents the state of the component in the ecosystem.
	Status string `json:"status"`
	// RequeueTimeNanos contains the time in nanoseconds to wait until the next requeue.
	RequeueTimeNanos time.Duration `json:"requeueTimeNanos,omitempty"`
	// Health describes the health status of the component.
	// A component becomes 'available' if its Status is 'installed',
	// and all its deployments, stateful sets, and daemon sets are available.
	Health HealthStatus `json:"health,omitempty"`
	// Installed version of the component (e.g. 2.4.48-3)
	InstalledVersion string `json:"installedVersion,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels=app=ces;app.kubernetes.io/name=k8s-component-operator;k8s.cloudogu.com/component.name=k8s-component-operator-crd
// +kubebuilder:resource:shortName="comp"
// +kubebuilder:printcolumn:name="Spec-Version",type="string",JSONPath=".spec.version",description="The desired version of the component"
// +kubebuilder:printcolumn:name="Installed Version",type="string",JSONPath=".status.installedVersion",description="The current version of the component"
// +kubebuilder:printcolumn:name="Health",type="string",JSONPath=".status.health",description="The current health state of the component"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.status",description="The current status of the component"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of the component"
// +kubebuilder:webhook:path=/validate-k8s-cloudogu-com-component,mutating=false,failurePolicy=fail,groups=k8s.cloudgou.com,resources=Component,verbs=create;update,versions=v1,name=validatecomponent.kb.io,admissionReviewVersions=v1,sideEffects=None

// Component is the Schema for the ces component API
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

type ChartGetter interface {
	GetChart(ctx context.Context, spec *client.ChartSpec) (*chart.Chart, error)
}

type HelmChartCreationOpts struct {
	HelmClient     ChartGetter
	Timeout        time.Duration
	YamlSerializer yaml.Serializer
}

type Mapping struct {
	Path    string            `yaml:"path"`
	Mapping map[string]string `yaml:"Mapping"`
}

type MetaValue struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Keys        []Mapping
}

type MetadataMapping struct {
	ApiVersion string               `yaml:"apiVersion"`
	Metavalues map[string]MetaValue `yaml:"metavalues"`
}

// String returns a string representation of this component.
func (c *Component) String() string {
	return fmt.Sprintf("%s/%s:%s", c.Spec.Namespace, c.Spec.Name, c.Spec.Version)
}

// GetHelmChartSpec returns the helm chart for the component cr without custom values.
func (c *Component) GetHelmChartSpec(ctx context.Context, opts ...HelmChartCreationOpts) (*client.ChartSpec, error) {
	deployNamespace := ""

	if c.Spec.DeployNamespace != "" {
		deployNamespace = c.Spec.DeployNamespace
	} else {
		deployNamespace = c.Namespace
	}

	timeout := defaultHelmClientTimeoutMins
	var chartGetter ChartGetter
	var yamlSerializer yaml.Serializer
	if len(opts) > 0 {
		timeout = opts[0].Timeout
		chartGetter = opts[0].HelmClient
		yamlSerializer = opts[0].YamlSerializer
	}

	chartSpec := &client.ChartSpec{
		ReleaseName: c.Spec.Name,
		ChartName:   c.GetHelmChartName(),
		Namespace:   deployNamespace,
		Version:     c.Spec.Version,
		ValuesYaml:  c.Spec.ValuesYamlOverwrite,
		// Rollback to previous release on failure.
		Atomic: true,
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: timeout,
		// True would lead the client to delete a CRD on failure which could delete all Dogus.
		CleanupOnFail: false,
		// Create non-existent namespace so that the operator can install charts in other namespaces.
		CreateNamespace: true,
		PostRenderer: labels.NewPostRenderer(map[string]string{
			ComponentNameLabelKey:    c.Spec.Name,
			ComponentVersionLabelKey: c.Spec.Version,
		}),
	}

	if len(opts) > 0 {
		var err error
		chartSpec.MappedValuesYaml, err = getMappedValuesYaml(ctx, c, chartSpec, chartGetter, yamlSerializer)
		if err != nil {
			return nil, fmt.Errorf("failed to create mapped values: %w", err)
		}
	}

	return chartSpec, nil
}

func pathToNestedYAML(path string, value any) map[string]any {
	parts := strings.Split(path, ".")
	n := len(parts)

	result := value
	for i := n - 1; i >= 0; i-- {
		result = map[string]any{
			parts[i]: result,
		}
	}

	return result.(map[string]any)
}

func getMappedValuesYaml(ctx context.Context, component *Component, spec *client.ChartSpec, helmClient ChartGetter, yamlSerializer yaml.Serializer) (string, error) {
	logger := log.FromContext(ctx)

	hChart, err := helmClient.GetChart(ctx, spec)
	if err != nil {
		return "", fmt.Errorf("failed to get helm chart: %w", err)
	}

	var mappings MetadataMapping
	for _, file := range hChart.Files {
		logger.Info(fmt.Sprintf("Found file %s in component %s", file.Name, component.Name))
		if file.Name == mappingMetadataFileName {
			logger.Info("Serializing metadata-file...")
			err = yamlSerializer.Unmarshal(file.Data, &mappings)
			if err != nil {
				return "", fmt.Errorf("failed to parse Mapping metadata: %w", err)
			}
		}
	}

	mappingYaml := map[string]interface{}{}

	for k, v := range component.Spec.MappedValues {
		if _, ok := mappings.Metavalues[k]; !ok {
			fmt.Printf("key %s not found in metaValues\n", k)
			continue
		}
		for _, key := range mappings.Metavalues[k].Keys {
			fmt.Printf("checking key %s...\n", key)
			if key.Mapping == nil {
				mappingYaml = values.MergeMaps(mappingYaml, pathToNestedYAML(key.Path, v))
				continue
			}
			if value, ok := key.Mapping[v]; ok {
				mappingYaml = values.MergeMaps(mappingYaml, pathToNestedYAML(key.Path, value))
			} else {
				logger.Error(fmt.Errorf("no Mapping found for key %s", v), "")
			}
		}
	}

	serialized, err := yamlSerializer.Marshal(mappingYaml)
	if err != nil {
		return "", fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return string(serialized), nil
}

func (c *Component) GetHelmChartName() string {
	return fmt.Sprintf("%s/%s", c.Spec.Namespace, c.Spec.Name)
}

// +kubebuilder:object:root=true

// ComponentList contains a list of Component
type ComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Component `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Component{}, &ComponentList{})
}

func (c *Component) ValidateCreate() (admission.Warnings, error) {
	fmt.Println("===================================>")
	fmt.Println("===================================>")
	fmt.Println("===================================>")
	fmt.Println("Created new component")
	return nil, fmt.Errorf("Validate create")
}

func (c *Component) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	fmt.Println("===================================>")
	fmt.Println("===================================>")
	fmt.Println("===================================>")
	fmt.Println("Update component")
	return nil, fmt.Errorf("Validate update")
}

func (c *Component) ValidateDelete() (admission.Warnings, error) {
	fmt.Println("===================================>")
	fmt.Println("===================================>")
	fmt.Println("===================================>")
	fmt.Println("Delete component")
	return nil, nil
}

func SetupComponentValidatorForManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&Component{}).
		Complete()
}
