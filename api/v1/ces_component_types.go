package v1

import (
	"embed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// ComponentName is used to select a component pod by name.
	ComponentName = "component.name"
	// ComponentVersion is used to select a component pod by version.
	ComponentVersion = "component.version"
)

const (
	ComponentStatusNotInstalled = ""
	ComponentStatusInstalling   = "installing"
	ComponentStatusUpgrading    = "upgrading"
	ComponentStatusDeleting     = "deleting"
	ComponentStatusInstalled    = "installed"
)

const FinalizerName = "component-finalizer"

// ComponentSpec defines the desired state of a component
type ComponentSpec struct {
	// Name of the component (e.g. official/ldap)
	Name string `json:"name,omitempty"`
	// Version of the component (e.g. 2.4.48-3)
	Version string `json:"version,omitempty"`
}

// ComponentStatus defines the observed state of a Component.
type ComponentStatus struct {
	// Status represents the state of the Dogu in the ecosystem
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
