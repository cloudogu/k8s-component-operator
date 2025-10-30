package configref

import (
	"context"
	"fmt"

	v2 "github.com/cloudogu/k8s-component-lib/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapRefReader struct {
	configMapClient configMapClient
}

func NewConfigMapRefReader(configMapClient configMapClient) *ConfigMapRefReader {
	return &ConfigMapRefReader{
		configMapClient: configMapClient,
	}
}

func (reader *ConfigMapRefReader) GetValues(ctx context.Context, configMapReference *v2.Reference) (string, error) {
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
