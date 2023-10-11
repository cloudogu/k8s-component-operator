package v1

import (
	"embed"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
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
	// DeployNamespace is the namespace where the helm chart should be deployed in.
	// This value is optional. If it is empty the operator deploys the helm chart in the namespace where the operator is deployed.
	DeployNamespace string `json:"deployNamespace,omitempty"`
}

// ComponentStatus defines the observed state of a Component.
type ComponentStatus struct {
	// Status represents the state of the component in the ecosystem.
	Status string `json:"status"`
	// RequeueTimeNanos contains the time in nanoseconds to wait until the next requeue.
	RequeueTimeNanos time.Duration `json:"requeueTimeNanos,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:metadata:labels=app=ces;app.kubernetes.io/name=k8s-component-operator

// Component is the Schema for the ces component API
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComponentSpec   `json:"spec,omitempty"`
	Status ComponentStatus `json:"status,omitempty"`
}

// String returns a string representation of this component.
func (c *Component) String() string {
	return fmt.Sprintf("%s/%s:%s", c.Spec.Namespace, c.Spec.Name, c.Spec.Version)
}

// GetHelmChartSpec returns the helm chart for the component cr without custom values.
func (c *Component) GetHelmChartSpec() *client.ChartSpec {
	deployNamespace := ""

	if c.Spec.DeployNamespace != "" {
		deployNamespace = c.Spec.DeployNamespace
	} else {
		deployNamespace = c.Namespace
	}

	return &client.ChartSpec{
		ReleaseName: c.Spec.Name,
		ChartName:   fmt.Sprintf("%s/%s", c.Spec.Namespace, c.Spec.Name),
		Namespace:   deployNamespace,
		Version:     c.Spec.Version,
		// Rollback to previous release on failure.
		Atomic: true,
		// This timeout prevents context exceeded errors from the used k8s client from the helm library.
		Timeout: time.Second * 300,
		// True would lead the client to delete a CRD on failure which could delete all Dogus.
		CleanupOnFail: false,
		// Create non-existent namespace so that the operator can install charts in other namespaces.
		CreateNamespace: true,
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
