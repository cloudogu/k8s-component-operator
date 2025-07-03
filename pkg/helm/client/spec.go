package client

import (
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/pkg/errors"

	"sigs.k8s.io/yaml"

	"helm.sh/helm/v3/pkg/getter"
)

// GetValuesMap returns the merged mapped out values of a chart,
// using both ValuesYaml and ValuesOptions
func (spec *ChartSpec) GetValuesMap(p getter.Providers) (map[string]interface{}, error) {
	originalValuesYaml := map[string]interface{}{}

	err := yaml.Unmarshal([]byte(spec.ValuesYaml), &originalValuesYaml)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesYaml")
	}

	valuesOptions, err := spec.ValuesOptions.MergeValues(p)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesOptions")
	}

	mergedValuesYamlOverwrite := values.MergeMaps(originalValuesYaml, valuesOptions)

	mappedValuesYaml := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(spec.MappedValuesYaml), &mappedValuesYaml)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesYaml")
	}

	return values.MergeMaps(mappedValuesYaml, mergedValuesYamlOverwrite), nil
}
