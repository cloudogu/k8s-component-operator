package labels

import (
	"bytes"
	"fmt"
	"maps"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	yamlutil "github.com/cloudogu/k8s-component-operator/pkg/yaml"
)

type PostRenderer struct {
	documentSplitter       documentSplitter
	unstructuredSerializer unstructuredSerializer
	unstructuredConverter  unstructuredConverter
	serializer             genericYamlSerializer
	labels                 map[string]string
}

func NewPostRenderer(labels map[string]string) *PostRenderer {
	return &PostRenderer{
		documentSplitter:       yamlutil.NewDocumentSplitter(),
		unstructuredSerializer: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
		unstructuredConverter:  runtime.DefaultUnstructuredConverter,
		serializer:             yamlutil.NewSerializer(),
		labels:                 labels,
	}
}

func (c *PostRenderer) Run(renderedManifests *bytes.Buffer) (modifiedManifests *bytes.Buffer, err error) {
	modifiedManifests = new(bytes.Buffer)

	c.documentSplitter.WithReader(renderedManifests)
	for c.documentSplitter.Next() {
		obj, _, err := c.unstructuredSerializer.Decode(c.documentSplitter.Bytes(), nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to parse yaml resources: %w", err)
		}

		unstructuredMap, err := c.unstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource to unstructured object: %w", err)
		}

		k8sObject := &unstructured.Unstructured{Object: unstructuredMap}

		originalLabels := k8sObject.GetLabels()
		mergedLabels := make(map[string]string, len(c.labels)+len(originalLabels))
		maps.Copy(mergedLabels, originalLabels)
		maps.Copy(mergedLabels, c.labels)
		k8sObject.SetLabels(mergedLabels)

		yamlBytes, err := c.serializer.Marshal(k8sObject)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize resources back to yaml: %w", err)
		}

		modifiedManifests.Write(yamlBytes)
		modifiedManifests.WriteString("\n---\n")
	}
	if err = c.documentSplitter.Err(); err != nil {
		return nil, fmt.Errorf("failed to split yaml document: %w", err)
	}

	return modifiedManifests, nil
}
