package client

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/getter"
	"sigs.k8s.io/yaml"

	yaml2 "github.com/cloudogu/k8s-component-operator/pkg/yaml"
)

//go:embed testdata/values/mappedValues.yaml
var mappedValuesBytes []byte

//go:embed testdata/values/valuesYamlOverwrite.yaml
var valuesYamlOverwriteBytes []byte

//go:embed testdata/values/configMapData.yaml
var configMapDataBytes []byte

//go:embed testdata/values/values.yaml
var valuesBytes []byte

//go:embed testdata/values/result.yaml
var resultBytes []byte

//go:embed testdata/values/valuesYamlOverwriteOverwritesConfigReferenceValues/valuesYamlOverwrite.yaml
var vyocrvValuesYamlOverwriteBytes []byte

//go:embed testdata/values/valuesYamlOverwriteOverwritesConfigReferenceValues/configMapData.yaml
var vyocrvConfigMapDataBytes []byte

//go:embed testdata/values/valuesYamlOverwriteOverwritesConfigReferenceValues/result.yaml
var vyocrvResultBytes []byte

func Test_haveSameKeyWithDifferentValues(t *testing.T) {
	tests := []struct {
		name string
		a    map[string]interface{}
		b    map[string]interface{}
		want bool
	}{
		{
			name: "two empty maps",
			a:    map[string]interface{}{},
			b:    map[string]interface{}{},
			want: false,
		},
		{
			name: "only different keys, more keys in a",
			a: map[string]interface{}{
				"a": "b",
				"c": "d",
			},
			b: map[string]interface{}{
				"e": "f",
			},
			want: false,
		},
		{
			a: map[string]interface{}{
				"e": "f",
			},
			name: "only different keys, more keys in b",
			b: map[string]interface{}{
				"a": "b",
				"c": "d",
			},
			want: false,
		},
		{
			name: "two equal maps",
			a: map[string]interface{}{
				"a": "b",
			},
			b: map[string]interface{}{
				"a": "b",
			},
			want: true,
		},
		{
			name: "two nested maps with different keys",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"e": "f",
				},
			},
			want: false,
		},
		{
			name: "two nested maps with same keys and values",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			want: true,
		},
		{
			name: "two nested maps with same keys different values",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "d",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "e",
				},
			},
			want: true,
		},
		{
			name: "string value in second map overwrites map value in first map",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": map[string]interface{}{
						"d": "e",
					},
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "e",
				},
			},
			want: true,
		},
		{
			name: "string value in first map overwrites string value in first map",
			a: map[string]interface{}{
				"a": map[string]interface{}{
					"c": "e",
				},
			},
			b: map[string]interface{}{
				"a": map[string]interface{}{
					"c": map[string]interface{}{
						"d": "e",
					},
				},
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equalf(t, test.want, hasSameValuesConfigured(test.a, test.b), "hasSameValuesConfigured(%v, %v)", test.a, test.b)
		})
	}
}

func TestChartSpec_GetValuesMap(t *testing.T) {
	type fields struct {
		MappedValuesYamlFn    func(t *testing.T) string
		ValuesOptionsFn       func(t *testing.T) valuesOptions
		ValuesConfigRefYamlFn func(t *testing.T) string
		ValuesYamlFn          func(t *testing.T) string
	}
	type args struct {
		p getter.Providers
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantFn  func(t *testing.T) map[string]interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "values yaml overwrite overwrites config reference values",
			fields: fields{
				MappedValuesYamlFn: func(t *testing.T) string {
					return ""
				},
				ValuesOptionsFn: func(t *testing.T) valuesOptions {
					mck := newMockValuesOptions(t)
					valuesYamlOverwrite := map[string]interface{}{}

					mck.EXPECT().MergeValues(mock.Anything).Return(valuesYamlOverwrite, nil)
					return mck
				},
				ValuesConfigRefYamlFn: func(t *testing.T) string {
					return string(vyocrvConfigMapDataBytes)
				},
				ValuesYamlFn: func(t *testing.T) string {
					return string(vyocrvValuesYamlOverwriteBytes)
				},
			},
			args: args{
				p: make(getter.Providers, 0),
			},
			wantErr: assert.NoError,
			wantFn: func(t *testing.T) map[string]interface{} {
				var resultYaml map[string]interface{}
				serializer := yaml2.NewSerializer()
				err := serializer.Unmarshal(vyocrvResultBytes, &resultYaml)
				require.NoError(t, err)
				return resultYaml
			},
		},
		{
			name: "mappedValues before values yaml overwrite. values yaml overwrite before configmap reference. Config map reference before values yaml",
			fields: fields{
				MappedValuesYamlFn: func(t *testing.T) string {
					return string(mappedValuesBytes)
				},
				ValuesOptionsFn: func(t *testing.T) valuesOptions {
					mck := newMockValuesOptions(t)
					valuesYamlOverwrite := map[string]interface{}{}
					err := yaml.Unmarshal(valuesBytes, &valuesYamlOverwrite)
					require.NoError(t, err)

					mck.EXPECT().MergeValues(mock.Anything).Return(valuesYamlOverwrite, nil)
					return mck
				},
				ValuesConfigRefYamlFn: func(t *testing.T) string {
					return string(configMapDataBytes)
				},
				ValuesYamlFn: func(t *testing.T) string {
					return string(valuesYamlOverwriteBytes)
				},
			},
			args: args{
				p: make(getter.Providers, 0),
			},
			wantErr: assert.NoError,
			wantFn: func(t *testing.T) map[string]interface{} {
				var resultYaml map[string]interface{}
				serializer := yaml2.NewSerializer()
				err := serializer.Unmarshal(resultBytes, &resultYaml)
				require.NoError(t, err)
				return resultYaml
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &ChartSpec{
				ValuesYamlOverwrite: tt.fields.ValuesYamlFn(t),
				MappedValuesYaml:    tt.fields.MappedValuesYamlFn(t),
				ValuesConfigRefYaml: tt.fields.ValuesConfigRefYamlFn(t),
				ValuesOptions:       tt.fields.ValuesOptionsFn(t),
			}
			got, err := spec.GetValuesMap(tt.args.p)
			if !tt.wantErr(t, err, fmt.Sprintf("GetValuesMap(%v)", tt.args.p)) {
				return
			}
			assert.Equalf(t, tt.wantFn(t), got, "GetValuesMap(%v)", tt.args.p)
		})
	}
}
