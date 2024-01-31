package helm

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/util/yaml"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

type metadataValueClient struct {
	componentHelmClient componentHelmClient
}

func (mvc *metadataValueClient) GetChart(ctx context.Context, chartSpec *client.ChartSpec) (*chart.Chart, error) {
	return mvc.componentHelmClient.GetChart(ctx, chartSpec)
}

func (mvc *metadataValueClient) InstallOrUpgrade(ctx context.Context, chart *client.ChartSpec) error {
	return mvc.componentHelmClient.InstallOrUpgrade(ctx, chart)
}

func (mvc *metadataValueClient) InstallOrUpgradeWithMappedValues(ctx context.Context, chart *client.ChartSpec, mappedValues map[string]string) error {
	return fmt.Errorf("not implemented")
}

func (mvc *metadataValueClient) Uninstall(releaseName string) error {
	return mvc.componentHelmClient.Uninstall(releaseName)
}

func (mvc *metadataValueClient) ListDeployedReleases() ([]*release.Release, error) {
	return mvc.componentHelmClient.ListDeployedReleases()
}

func (mvc *metadataValueClient) GetReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	return mvc.componentHelmClient.GetReleaseValues(name, allValues)
}

func (mvc *metadataValueClient) GetChartSpecValues(chart *client.ChartSpec) (map[string]interface{}, error) {
	return mvc.componentHelmClient.GetChartSpecValues(chart)
}

func (mvc *metadataValueClient) SatisfiesDependencies(ctx context.Context, chart *client.ChartSpec) error {
	return mvc.componentHelmClient.SatisfiesDependencies(ctx, chart)
}

type helmValuesMetadata struct {
	ApiVersion string                             `yaml:"apiVersion"`
	Metadata   map[string]helmValuesMetadataEntry `yaml:"metadata"`
}

type helmValuesMetadataEntry struct {
	Name        string                  `yaml:"name"`
	Description string                  `yaml:"description"`
	Keys        []helmValuesMetadataKey `yaml:"keys"`
}

type helmValuesMetadataKey struct {
	Path string `yaml:"path"`
}

func NewMetadataValueClient(componentHelmClient componentHelmClient) *metadataValueClient {
	return &metadataValueClient{componentHelmClient: componentHelmClient}
}

// VerifyUserDefinedValues checks if values defined in the metadata of the chart are set in the valueYamlOverwrite Field of the Component.
// If yes the method returns an error.
// If no the method returns nil.
func (mvc *metadataValueClient) VerifyUserDefinedValues(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	valuesMetadata, err := mvc.getValuesMetadata(ctx, component)
	if err != nil {
		return err
	}

	if valuesMetadata == nil {
		logger.Info(fmt.Sprintf("found no metadata file %q for component %q", helmValuesMetadataFileName, component.Spec.Name))
		return nil
	}

	overwriteValues, err := mvc.componentHelmClient.GetChartSpecValues(component.GetHelmChartSpec())
	if err != nil {
		return fmt.Errorf("failed to read current overwriteValues: %w", err)
	}

	for _, metadataEntry := range valuesMetadata.Metadata {
		// Check if metadata values are already set in valuesYamlOverwrite
		for _, key := range metadataEntry.Keys {
			// key.Path is somethings like controllerManager.env.logLevel (dot-separated)
			path := key.Path
			if isMetadataPathInValuesMap(path, overwriteValues) {
				return fmt.Errorf("values contains path %s which should only be set in field mappedValues", path)
			}
		}
	}

	return nil
}

func (mvc *metadataValueClient) getValuesMetadata(ctx context.Context, component *k8sv1.Component) (*helmValuesMetadata, error) {
	helmChart, err := mvc.componentHelmClient.GetChart(ctx, component.GetHelmChartSpec())
	if err != nil {
		return nil, err
	}

	valuesMetadata, err := mvc.getValuesMetadataFromChart(ctx, helmChart)
	if err != nil {
		return nil, err
	}
	return valuesMetadata, nil
}

// path should be a dot-separated yaml path like: `controllerManager.env.logLevel`
func isMetadataPathInValuesMap(path string, values map[string]interface{}) bool {
	before, after, _ := strings.Cut(path, ".")

	i, ok := values[before]
	if !ok {
		return false
	}

	ii, ok := i.(map[string]interface{})
	if !ok && after == "" {
		return true
	}

	return isMetadataPathInValuesMap(after, ii)
}

func (mvc *metadataValueClient) getValuesMetadataFromChart(ctx context.Context, chart *chart.Chart) (*helmValuesMetadata, error) {
	logger := log.FromContext(ctx)
	files := chart.Files
	var metadataFile *helmValuesMetadata
	for _, file := range files {
		logger.Info(file.Name)
		if file.Name == helmValuesMetadataFileName {
			err := yaml.Unmarshal(file.Data, &metadataFile)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal %s: %w", helmValuesMetadataFileName, err)
			}
			logger.Info(fmt.Sprintf("found metadata file %q for chart %q", helmValuesMetadataFileName, chart.Metadata.Name))
			logger.Info(fmt.Sprintf("%+v", metadataFile))
			break
		}
	}

	return metadataFile, nil
}

func (mvc *metadataValueClient) IsValuesChanged(ctx context.Context, deployedRelease *release.Release, component *k8sv1.Component) (bool, error) {
	deployedValues, err := mvc.componentHelmClient.GetReleaseValues(deployedRelease.Name, false)
	if err != nil {
		return false, fmt.Errorf("failed to get values.yaml from release %s: %w", deployedRelease.Name, err)
	}

	chartSpecValues, err := mvc.componentHelmClient.GetChartSpecValues(component.GetHelmChartSpec())
	if err != nil {
		return false, fmt.Errorf("failed to get values.yaml from component %q: %w", component.GetHelmChartSpec().ChartName, err)
	}

	mappedValues, err := mvc.getMappedValues(ctx, component)
	if err != nil {
		return false, fmt.Errorf("failed to get mapped values for component %q: %w", component.Spec.Name, err)
	}
	userValues := MergeMapSlice(chartSpecValues, mappedValues)

	log.FromContext(ctx).Info("Compare actual deployed values:\n %+v", deployedValues)
	log.FromContext(ctx).Info("Compare actual user values:\n %+v", userValues)

	return !reflect.DeepEqual(deployedValues, userValues), nil
}

func (mvc *metadataValueClient) getMappedValues(ctx context.Context, component *k8sv1.Component) (map[string]interface{}, error) {
	valuesMetadata, err := mvc.getValuesMetadata(ctx, component)
	if err != nil {
		return nil, err
	}

	var valueMaps []map[string]interface{}
	for key, value := range component.Spec.MappedValues {
		metadata, ok := valuesMetadata.Metadata[key]
		if !ok {
			return nil, fmt.Errorf("failed to get metadata for key %q", key)
		}

		for _, yamlKey := range metadata.Keys {
			// TODO MAP VALUES
			path, yamlPathErr := buildMapFromYamlPath(yamlKey.Path, value)
			if yamlPathErr != nil {
				return nil, yamlPathErr
			}
			valueMaps = append(valueMaps, path)
		}
	}

	return MergeMapSlice(valueMaps...), nil
}

func buildMapFromYamlPath(path, value string) (map[string]interface{}, error) {
	ref := map[string]interface{}{}
	head := ref

	if path == "" {
		if value != "" {
			return ref, fmt.Errorf("could not build map from yaml path for an empty path and value %q", value)
		}
		return ref, nil
	}

	pathElements := strings.Split(path, ".")

	for i, element := range pathElements {
		if i == len(pathElements)-1 {
			ref[element] = value
		} else {
			child := map[string]interface{}{}
			ref[element] = child
			ref = child
		}
	}

	return head, nil
}

func MergeMapSlice(maps ...map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for _, element := range maps {
		result = values.MergeMaps(result, element)
	}

	return result
}
