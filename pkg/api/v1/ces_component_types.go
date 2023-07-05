package v1

import (
	"embed"
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// This embed provides the crd for other applications. They can import this package and use the yaml file
// for the CRD in e.g. integration tests. The file gets refreshed by copying from the kubebuilder config/crd/bases
// folder by the "generate" make target.
//
//go:embed k8s.cloudogu.com_components.yaml
var _ embed.FS

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
)

const FinalizerName = "component-finalizer"

// ComponentSpec defines the desired state of a component.
type ComponentSpec struct {
	// Namespace of the component (e.g. k8s)
	Namespace string `json:"namespace,omitempty"`
	// Name of the component (e.g. k8s-dogu-operator)
	Name string `json:"name,omitempty"`
	// Version of the component (e.g. 2.4.48-3)
	Version string `json:"version,omitempty"`
}

// ComponentStatus defines the observed state of a Component.
type ComponentStatus struct {
	// Status represents the state of the component in the ecosystem.
	Status string `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Component is the Schema for the ces component API
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

// GetHelmChartSpec returns the helm chart for the component cr without custom values.
func (c *Component) GetHelmChartSpec(repositoryEndpoint string) *helmclient.ChartSpec {
	return &helmclient.ChartSpec{
		ReleaseName: c.Spec.Name,
		ChartName:   fmt.Sprintf("%s/%s/%s", repositoryEndpoint, c.Spec.Namespace, c.Spec.Name),
		Namespace:   c.Namespace,
		Version:     c.Spec.Version,
		// Rollback to previous release on failure.
		Atomic: true,
		// This timeout preven
		Timeout: time.Second * 30,
		// True would lead the client to delete a CRD on failure which could delete all Dogus.
		CleanupOnFail: false,
	}
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
