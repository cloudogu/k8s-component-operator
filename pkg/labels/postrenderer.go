package labels

import (
	"bytes"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"maps"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"

	yamlutil "github.com/cloudogu/k8s-component-operator/pkg/yaml"
)

const deploymentKind = "apps/v1/Deployment"
const statefulSetKind = "apps/v1/StatefulSet"
const daemonSetKind = "apps/v1/DaemonSet"
const jobKind = "batch/v1/Job"
const cronJobKind = "batch/v1/CronJob"

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
		documentBytes := c.documentSplitter.Bytes()
		if len(documentBytes) == 0 {
			continue
		}

		obj, _, err := c.unstructuredSerializer.Decode(documentBytes, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to parse yaml resources: %w", err)
		}

		unstructuredMap, err := c.unstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource to unstructured object: %w", err)
		}

		k8sObject := &unstructured.Unstructured{Object: unstructuredMap}
		kind := fmt.Sprintf("%s/%s", k8sObject.GetAPIVersion(), k8sObject.GetKind())

		switch kind {
		case deploymentKind:
			k8sObject, err = addLabelsToStructured(c, unstructuredMap, &appsv1.Deployment{}, func(a *appsv1.Deployment) objectWithLabels { return &a.Spec.Template })
			if err != nil {
				return nil, fmt.Errorf("failed to add labels to Deployment: %w", err)
			}
		case statefulSetKind:
			k8sObject, err = addLabelsToStructured(c, unstructuredMap, &appsv1.StatefulSet{}, func(a *appsv1.StatefulSet) objectWithLabels { return &a.Spec.Template })
			if err != nil {
				return nil, fmt.Errorf("failed to add labels to StatefulSet: %w", err)
			}
		case daemonSetKind:
			k8sObject, err = addLabelsToStructured(c, unstructuredMap, &appsv1.DaemonSet{}, func(a *appsv1.DaemonSet) objectWithLabels { return &a.Spec.Template })
			if err != nil {
				return nil, fmt.Errorf("failed to add labels to DaemonSet: %w", err)
			}
		case jobKind:
			k8sObject, err = addLabelsToStructured(c, unstructuredMap, &batchv1.Job{}, func(a *batchv1.Job) objectWithLabels { return &a.Spec.Template })
			if err != nil {
				return nil, fmt.Errorf("failed to add labels to Job: %w", err)
			}
		case cronJobKind:
			k8sObject, err = addLabelsToStructured(c, unstructuredMap, &batchv1.CronJob{}, func(a *batchv1.CronJob) objectWithLabels { return &a.Spec.JobTemplate.Spec.Template })
			if err != nil {
				return nil, fmt.Errorf("failed to add labels to CronJob: %w", err)
			}
		}

		addLabels(k8sObject, c.labels)

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

func addLabelsToStructured[k any](c *PostRenderer, unstructuredMap map[string]interface{}, obj k, getTemplate func(k) objectWithLabels) (*unstructured.Unstructured, error) {
	if err := c.unstructuredConverter.FromUnstructured(unstructuredMap, obj); err != nil {
		return nil, fmt.Errorf("failed to convert resource to structured object: %w", err)
	}

	addLabels(getTemplate(obj), c.labels)

	newUnstructuredMap, err := c.unstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to convert resource to unstructured object: %w", err)
	}

	return &unstructured.Unstructured{Object: newUnstructuredMap}, nil
}

type objectWithLabels interface {
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
}

func addLabels(obj objectWithLabels, labels map[string]string) {
	originalLabels := obj.GetLabels()
	mergedLabels := make(map[string]string, len(labels)+len(originalLabels))
	maps.Copy(mergedLabels, originalLabels)
	maps.Copy(mergedLabels, labels)
	obj.SetLabels(mergedLabels)
}
