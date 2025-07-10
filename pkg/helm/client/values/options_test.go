/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/*
Copied from https://github.com/helm/helm/blob/eea2f27babb0fddd9fb1907f4d8531c8f5c73c66/pkg/cli/values/options_test.go
No changes.
*/

package values

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sigs.k8s.io/yaml"
	"testing"

	"helm.sh/helm/v3/pkg/getter"
)

func TestMergeValues(t *testing.T) {
	nestedMap := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]string{
			"cool": "stuff",
		},
	}
	anotherNestedMap := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]string{
			"cool":    "things",
			"awesome": "stuff",
		},
	}
	flatMap := map[string]interface{}{
		"foo": "bar",
		"baz": "stuff",
	}
	anotherFlatMap := map[string]interface{}{
		"testing": "fun",
	}

	testMap := MergeMaps(flatMap, nestedMap)
	equal := reflect.DeepEqual(testMap, nestedMap)
	if !equal {
		t.Errorf("Expected a nested map to overwrite a flat value. Expected: %v, got %v", nestedMap, testMap)
	}

	testMap = MergeMaps(nestedMap, flatMap)
	equal = reflect.DeepEqual(testMap, flatMap)
	if !equal {
		t.Errorf("Expected a flat value to overwrite a map. Expected: %v, got %v", flatMap, testMap)
	}

	testMap = MergeMaps(nestedMap, anotherNestedMap)
	equal = reflect.DeepEqual(testMap, anotherNestedMap)
	if !equal {
		t.Errorf("Expected a nested map to overwrite another nested map. Expected: %v, got %v", anotherNestedMap, testMap)
	}

	testMap = MergeMaps(anotherFlatMap, anotherNestedMap)
	expectedMap := map[string]interface{}{
		"testing": "fun",
		"foo":     "bar",
		"baz": map[string]string{
			"cool":    "things",
			"awesome": "stuff",
		},
	}
	equal = reflect.DeepEqual(testMap, expectedMap)
	if !equal {
		t.Errorf("Expected a map with different keys to merge properly with another map. Expected: %v, got %v", expectedMap, testMap)
	}
}

func TestReadFile(t *testing.T) {
	var p getter.Providers
	filePath := "%a.txt"
	_, err := readFile(filePath, p)
	if err == nil {
		t.Errorf("Expected error when has special strings")
	}
}

func TestMergeMaps_WithLists(t *testing.T) {
	t.Run("map with complete overwrite", func(t *testing.T) {
		var mapA map[string]interface{}
		strA := `kube-prometheus-stack:
  prometheus:
    containers:
      - name: auth
        env:
          - name: LOG_LEVEL
            value: info
`
		_ = yaml.Unmarshal([]byte(strA), &mapA)

		var mapB map[string]interface{}
		strB := `kube-prometheus-stack:
  prometheus:
    containers:
      - name: auth
        env:
          - name: LOG_LEVEL
            value: debug
`
		_ = yaml.Unmarshal([]byte(strB), &mapB)

		testMap := MergeMaps(mapA, mapB)
		equal := reflect.DeepEqual(testMap, mapB)
		assert.True(t, equal)
	})
	t.Run("map with additional values in source map", func(t *testing.T) {
		var mapA map[string]interface{}
		strA := `kube-prometheus-stack:
  prometheus:
    containers:
      - name: auth
        env:
          - name: LOG_LEVEL
            value: info
      - name: test
        env:
          - name: LOG_LEVEL
            value: error
`
		_ = yaml.Unmarshal([]byte(strA), &mapA)

		var mapB map[string]interface{}
		strB := `kube-prometheus-stack:
  prometheus:
    containers:
      - name: auth
        env:
          - name: LOG_LEVEL
            value: debug
`
		_ = yaml.Unmarshal([]byte(strB), &mapB)

		var mapResult map[string]interface{}
		strResult := `kube-prometheus-stack:
  prometheus:
    containers:
      - name: auth
        env:
          - name: LOG_LEVEL
            value: debug
      - name: test
        env:
          - name: LOG_LEVEL
            value: error
`
		_ = yaml.Unmarshal([]byte(strResult), &mapResult)

		testMap := MergeMaps(mapA, mapB)
		equal := reflect.DeepEqual(testMap, mapResult)
		assert.True(t, equal)
	})
}
