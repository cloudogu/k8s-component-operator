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
Copied from https://github.com/helm/helm/blob/eea2f27babb0fddd9fb1907f4d8531c8f5c73c66/pkg/cli/values/options.go
Changes:
- Add generator comments
- Export MergeMaps
*/

package values

import (
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"

	"sigs.k8s.io/yaml"

	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/strvals"
)

// Options captures the different ways to specify values
type Options struct {
	ValueFiles   []string // -f/--values
	StringValues []string // --set-string
	Values       []string // --set
	FileValues   []string // --set-file
	JSONValues   []string // --set-json
}

// MergeValues merges values from files specified via -f/--values and directly
// via --set-json, --set, --set-string, or --set-file, marshaling them to YAML
func (opts *Options) MergeValues(p getter.Providers) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range opts.ValueFiles {
		currentMap := map[string]interface{}{}

		bytes, err := readFile(filePath, p)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s", filePath)
		}
		// Merge with the previous map
		base = MergeMaps(base, currentMap)
	}

	// User specified a value via --set-json
	for _, value := range opts.JSONValues {
		if err := strvals.ParseJSON(value, base); err != nil {
			return nil, errors.Errorf("failed parsing --set-json data %s", value)
		}
	}

	// User specified a value via --set
	for _, value := range opts.Values {
		if err := strvals.ParseInto(value, base); err != nil {
			return nil, errors.Wrap(err, "failed parsing --set data")
		}
	}

	// User specified a value via --set-string
	for _, value := range opts.StringValues {
		if err := strvals.ParseIntoString(value, base); err != nil {
			return nil, errors.Wrap(err, "failed parsing --set-string data")
		}
	}

	// User specified a value via --set-file
	for _, value := range opts.FileValues {
		reader := func(rs []rune) (interface{}, error) {
			bytes, err := readFile(string(rs), p)
			if err != nil {
				return nil, err
			}
			return string(bytes), err
		}
		if err := strvals.ParseIntoFile(value, base, reader); err != nil {
			return nil, errors.Wrap(err, "failed parsing --set-file data")
		}
	}

	return base, nil
}

// Merge maps recursively, with smart list merging by `name`
func MergeMaps(a, b map[string]interface{}) map[string]interface{} {
	for key, bVal := range b {
		if aVal, exists := a[key]; exists {
			switch aValTyped := aVal.(type) {
			case map[string]interface{}:
				if bValTyped, ok := bVal.(map[string]interface{}); ok {
					a[key] = MergeMaps(aValTyped, bValTyped)
				} else {
					a[key] = bVal // overwrite non-map
				}
			case []interface{}:
				if bValTyped, ok := bVal.([]interface{}); ok {
					a[key] = mergeListsByName(aValTyped, bValTyped)
				} else {
					a[key] = bVal // overwrite non-list
				}
			default:
				a[key] = bVal // primitive overwrite
			}
		} else {
			a[key] = bVal // new key
		}
	}
	return a
}

// Merge two lists of maps by matching `name` key
func mergeListsByName(aList, bList []interface{}) []interface{} {
	result := make([]interface{}, 0)
	used := map[string]bool{}

	// Index A list by name
	aIndex := map[string]map[string]interface{}{}
	for _, item := range aList {
		if m, ok := item.(map[string]interface{}); ok {
			if name, ok := m["name"].(string); ok {
				aIndex[name] = m
			}
		}
	}

	// Merge B list into A
	for _, bItem := range bList {
		if bMap, ok := bItem.(map[string]interface{}); ok {
			if name, ok := bMap["name"].(string); ok {
				if aItem, exists := aIndex[name]; exists {
					// Merge matching items
					merged := MergeMaps(copyMap(aItem), bMap)
					result = append(result, merged)
					used[name] = true
					continue
				}
			}
		}
		result = append(result, bItem)
	}

	// Add remaining A items that were not merged
	for name, aItem := range aIndex {
		if !used[name] {
			result = append(result, aItem)
		}
	}

	return result
}

// Deep copy map
func copyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// readFile load a file from stdin, the local directory, or a remote file with a url.
func readFile(filePath string, p getter.Providers) ([]byte, error) {
	if strings.TrimSpace(filePath) == "-" {
		return io.ReadAll(os.Stdin)
	}
	u, err := url.Parse(filePath)
	if err != nil {
		return nil, err
	}

	// FIXME: maybe someone handle other protocols like ftp.
	g, err := p.ByScheme(u.Scheme)
	if err != nil {
		return os.ReadFile(filePath)
	}
	data, err := g.Get(filePath, getter.WithURL(filePath))
	if err != nil {
		return nil, err
	}
	return data.Bytes(), err
}
