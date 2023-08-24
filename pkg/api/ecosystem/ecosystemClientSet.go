package ecosystem

import (
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ComponentEcosystemInterface interface {
	kubernetes.Interface
	ComponentV1Alpha1() ComponentV1Alpha1Interface
}

type ComponentV1Alpha1Interface interface {
	Components(namespace string) ComponentInterface
}

// EcosystemClientset wraps the regular clientset with the ecosystemV1Alpha1 client.
type EcosystemClientset struct {
	*kubernetes.Clientset
	ecosystemV1Alpha1 *V1Alpha1Client
}

// NewComponentClientset creates a new instance of the component client.
func NewComponentClientset(config *rest.Config, clientset *kubernetes.Clientset) (*EcosystemClientset, error) {
	componentClientset, err := NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &EcosystemClientset{
		Clientset:         clientset,
		ecosystemV1Alpha1: componentClientset,
	}, nil
}

// ComponentV1Alpha1 returns the ecosystemV1Aplha1 client.
func (cswc *EcosystemClientset) ComponentV1Alpha1() ComponentV1Alpha1Interface {
	return cswc.ecosystemV1Alpha1
}

// V1Alpha1Client wraps the rest.Interface to use as a restClient for the component client.
type V1Alpha1Client struct {
	restClient rest.Interface
}

// NewForConfig creates a new V1Alpha1Client for a given rest.Config.
func NewForConfig(c *rest.Config) (*V1Alpha1Client, error) {
	config := *c
	gv := schema.GroupVersion{Group: v1.GroupVersion.Group, Version: v1.GroupVersion.Version}
	config.ContentConfig.GroupVersion = &gv
	config.APIPath = "/apis"

	s := scheme.Scheme
	err := v1.AddToScheme(s)
	if err != nil {
		return nil, err
	}

	metav1.AddToGroupVersion(s, gv)
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &V1Alpha1Client{restClient: client}, nil
}

// Components takes a namespace and returns a new component client.
func (c *V1Alpha1Client) Components(namespace string) ComponentInterface {
	return &componentClient{
		client: c.restClient,
		ns:     namespace,
	}
}
