package helm

import (
	"context"
	"fmt"
	"time"

	componentV1 "github.com/cloudogu/k8s-component-lib/api/v1"

	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/cloudogu/k8s-component-operator/pkg/labels"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"helm.sh/helm/v3/pkg/chart"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const defaultHelmClientTimeoutMins = time.Duration(15) * time.Minute
const mappingMetadataFileName = "component-values-metadata.yaml"

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

// GetHelmChartSpec returns the helm chart for the component cr without custom values.
func GetHelmChartSpec(ctx context.Context, c *componentV1.Component, opts ...HelmChartCreationOpts) (*client.ChartSpec, error) {
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
		ChartName:   GetHelmChartName(c),
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
			componentV1.ComponentNameLabelKey:    c.Spec.Name,
			componentV1.ComponentVersionLabelKey: c.Spec.Version,
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

func getMappedValuesYaml(ctx context.Context, component *componentV1.Component, spec *client.ChartSpec, helmClient ChartGetter, yamlSerializer yaml.Serializer) (string, error) {
	logger := log.FromContext(ctx)

	if len(component.Spec.MappedValues) == 0 {
		return "", nil
	}

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
			continue
		}
		for _, key := range mappings.Metavalues[k].Keys {
			if key.Mapping == nil {
				nestedYaml, e := yaml.PathToYAML(key.Path, v, yamlSerializer)
				if e != nil {
					logger.Error(fmt.Errorf("error parsing key path %s", key.Path), "")
					continue
				}
				mappingYaml = values.MergeMaps(mappingYaml, nestedYaml)
				continue
			}
			if value, ok := key.Mapping[v]; ok {
				nestedYaml, e := yaml.PathToYAML(key.Path, value, yamlSerializer)
				if e != nil {
					logger.Error(fmt.Errorf("error parsing key path %s", key.Path), "")
					continue
				}
				mappingYaml = values.MergeMaps(mappingYaml, nestedYaml)
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

func GetHelmChartName(c *componentV1.Component) string {
	return fmt.Sprintf("%s/%s", c.Spec.Namespace, c.Spec.Name)
}
