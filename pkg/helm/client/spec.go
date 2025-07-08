package client

import (
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"sigs.k8s.io/yaml"

	"helm.sh/helm/v3/pkg/getter"
)

// GetValuesMap returns the merged mapped out values of a chart,
// using both ValuesYaml and ValuesOptions
func (spec *ChartSpec) GetValuesMap(p getter.Providers) (map[string]interface{}, error) {
	logger := log.FromContext(nil)
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

	conflict := hasSameValuesConfigured(mergedValuesYamlOverwrite, mappedValuesYaml)
	if conflict {
		logger.Error(fmt.Errorf("conflicting values in valuesYamlOverwrite and mappedValues"), "you cannot set log mapped values via valuesYamlOverwrite and mappedValues. Configured value in mappedValues has priority")
	}

	return values.MergeMaps(mergedValuesYamlOverwrite, mappedValuesYaml), nil
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
				vsl, ok := v.([]map[string]interface{})
				if ok {
					mvsl, ok := mv.([]map[string]interface{})
					if !ok {
						return true
					}
					for _, ma := range vsl {
						for _, mb := range mvsl {
							if hasSameValuesConfigured(ma, mb) {
								return true
							}
						}
					}
				} else {
					return true
				}
			}

			return false
		}
	}

	return false
}
