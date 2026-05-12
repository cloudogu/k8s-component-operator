package configref

import (
	"context"
	"fmt"

	v2 "github.com/cloudogu/k8s-component-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ConfigMapRefReader struct {
	configMapClient configMapClient
}

func NewConfigMapRefReader(configMapClient configMapClient) ConfigMapRefReader {
	return ConfigMapRefReader{
		configMapClient: configMapClient,
	}
}

func (reader ConfigMapRefReader) GetValues(ctx context.Context, configMapReference *v2.Reference) (string, error) {
	if configMapReference == nil || configMapReference.Name == "" {
		return "", nil
	}
	configMap, err := reader.configMapClient.Get(ctx, configMapReference.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	value, exists := configMap.Data[configMapReference.Key]

	if !exists {
		return "", fmt.Errorf("key %s does not exist in configmap %s", configMapReference.Key, configMapReference.Name)
	}

	return value, nil
}

func (reader ConfigMapRefReader) GetSystemValues(ctx context.Context, component *v2.Component) (string, error) {
	if component == nil {
		return "", nil
	}
	cmLabelSelector := labels.ValidatedSetSelector{
		"k8s.cloudogu.com/component.config": component.Name,
	}.String()
	configMaps, err := reader.configMapClient.List(ctx, metav1.ListOptions{LabelSelector: cmLabelSelector})
	if err != nil {
		return "", err
	}
	if len(configMaps.Items) != 1 {
		return "", nil
	}
	cm := configMaps.Items[0]
	value, exists := cm.Data["values"]

	if !exists {
		return "", fmt.Errorf("key `values` does not exist in configmap %s", cm.Name)
	}

	return value, nil
}
