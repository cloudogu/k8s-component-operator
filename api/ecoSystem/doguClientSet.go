package ecoSystem

import (
	v1 "github.com/cloudogu/k8s-component-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// TODO create organization interface like appsv1 and merge with the doguclient
type ClientSetWithComponent struct {
	kubernetes.Clientset
	EcoSystemV1Alpha1Client
}

type EcoSystemV1Alpha1Interface interface {
	Components(namespace string) ComponentInterface
}

type EcoSystemV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*EcoSystemV1Alpha1Client, error) {
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

	return &EcoSystemV1Alpha1Client{restClient: client}, nil
}

func (c *EcoSystemV1Alpha1Client) Components(namespace string) ComponentInterface {
	return &componentClient{
		client: c.restClient,
		ns:     namespace,
	}
}
