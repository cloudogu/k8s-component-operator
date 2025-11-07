package client

import (
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	"helm.sh/helm/v3/pkg/getter"
)

// GetValuesMap returns the merged mapped out values of a chart,
// using both ValuesYamlOverwrite and ValuesOptions
func (spec *ChartSpec) GetValuesMap(p getter.Providers) (map[string]interface{}, error) {
	valuesYamlOverwrite := map[string]interface{}{}

	err := yaml.Unmarshal([]byte(spec.ValuesYamlOverwrite), &valuesYamlOverwrite)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesYamlOverwrite")
	}

	configRefValues := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(spec.ValuesConfigRefYaml), &configRefValues)
	if err != nil {
		return nil, err
	}
	commandLineOptionValues := map[string]interface{}{}
	if spec.ValuesOptions != nil {
		commandLineOptionValues, err = spec.ValuesOptions.MergeValues(p)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to Parse ValuesOptions")
		}
	}

	mappedValues := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(spec.MappedValuesYaml), &mappedValues)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse mappedValues")
	}

	result := values.MergeMaps(mappedValues, commandLineOptionValues)
	result = values.MergeMaps(valuesYamlOverwrite, result)
	result = values.MergeMaps(configRefValues, result)

	return result, nil
}

func hasSameValuesConfigured(a, b map[string]interface{}) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	for k, v := range a {
		if mv, ok := b[k]; ok {
			vmp, ok := v.(map[string]interface{})
			if ok {
				mvmp, ok := mv.(map[string]interface{})
				if !ok {
					return true
				}

				if hasSameValuesConfigured(vmp, mvmp) {
					return true
				} // else just continue loop
			} else {
				return true
			}

			return false
		}
	}

	return false
}
