package v1

import (
	"context"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestComponent_GetHelmChartSpec(t *testing.T) {
	type fields struct {
		TypeMeta   v1.TypeMeta
		ObjectMeta v1.ObjectMeta
		Spec       ComponentSpec
		Status     ComponentStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "should use deployNamespace if specified", fields: fields{Spec: ComponentSpec{DeployNamespace: "longhorn"}}, want: "longhorn"},
		{name: "should use regular namespace if no deployNamespace if specified", fields: fields{ObjectMeta: v1.ObjectMeta{Namespace: "ecosystem"}, Spec: ComponentSpec{DeployNamespace: ""}}, want: "ecosystem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			spec, _ := c.GetHelmChartSpec(context.Background())
			if got := spec.Namespace; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHelmChartSpec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMappedValuesYaml(t *testing.T) {
	testCtx := context.Background()
	t.Run("success with mapping", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"mainLogLevel": "debug",
				},
			},
		}

		doguOpMetaData := `apiVersion: v1
metavalues:
  mainLogLevel:
    name: Log-Level
    description: The central configuration value to set the log level for this component
    keys:
      - path: controllerManager.env.loglevel
        mapping:
          debug: trace
          info: info
          warn: warn
          error: error`

		spec := &client.ChartSpec{}
		helmChart := &chart.Chart{
			Files: []*chart.File{
				{
					Name: mappingMetadataFileName,
					Data: []byte(doguOpMetaData),
				},
			},
		}

		expectedYaml := `controllerManager:
  env:
    loglevel: trace
`

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, expectedYaml, mappedValuesYaml)
		assert.NoError(t, err)
	})
	t.Run("success without mapping", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"mainLogLevel": "debug",
				},
			},
		}

		doguOpMetaData := `apiVersion: v1
metavalues:
  mainLogLevel:
    name: Log-Level
    description: The central configuration value to set the log level for this component
    keys:
      - path: controllerManager.env.loglevel`

		spec := &client.ChartSpec{}
		helmChart := &chart.Chart{
			Files: []*chart.File{
				{
					Name: mappingMetadataFileName,
					Data: []byte(doguOpMetaData),
				},
			},
		}

		expectedYaml := `controllerManager:
  env:
    loglevel: debug
`

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, expectedYaml, mappedValuesYaml)
		assert.NoError(t, err)
	})
	t.Run("success with multiple path", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"mainLogLevel": "debug",
				},
			},
		}

		doguOpMetaData := `apiVersion: v1
metavalues:
  mainLogLevel:
    name: Log-Level
    description: The central configuration value to set the log level for this component
    keys:
      - path: controllerManager.env.loglevel
        mapping:
          debug: trace
      - path: my.awesome.logLevel
        mapping:
          debug: all
      - path: your.not.so.awesome.path`

		spec := &client.ChartSpec{}
		helmChart := &chart.Chart{
			Files: []*chart.File{
				{
					Name: mappingMetadataFileName,
					Data: []byte(doguOpMetaData),
				},
			},
		}

		expectedYaml := `controllerManager:
  env:
    loglevel: trace
my:
  awesome:
    logLevel: all
your:
  not:
    so:
      awesome:
        path: debug
`

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, expectedYaml, mappedValuesYaml)
		assert.NoError(t, err)
	})
}
