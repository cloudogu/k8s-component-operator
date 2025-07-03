package client

import (
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/pkg/errors"

	"sigs.k8s.io/yaml"

	"helm.sh/helm/v3/pkg/getter"
)

// GetValuesMap returns the merged mapped out values of a chart,
// using both ValuesYaml and ValuesOptions
func (spec *ChartSpec) GetValuesMap(p getter.Providers) (map[string]interface{}, error) {
	valuesYaml := map[string]interface{}{}

	err := yaml.Unmarshal([]byte(spec.ValuesYaml), &valuesYaml)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesYaml")
	}

	valuesOptions, err := spec.ValuesOptions.MergeValues(p)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesOptions")
	}

	result := values.MergeMaps(valuesYaml, valuesOptions)

	valuesYaml2 := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(spec.ValuesYaml2), &valuesYaml2)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse ValuesYaml")
	}

	fmt.Println("==================================================adasds=>>")

	fmt.Printf("%v\n", valuesYaml2)
	fmt.Println("====>")
	fmt.Printf("%v\n", result)

	finalResult := values.MergeMaps(valuesYaml2, result)
	fmt.Println("====>")
	fmt.Printf("%v\n", finalResult)

	return finalResult, nil
}
