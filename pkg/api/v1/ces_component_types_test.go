package v1

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client"
	"github.com/cloudogu/k8s-component-operator/pkg/helm/client/values"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	originalyaml "sigs.k8s.io/yaml"
	"testing"
)

//go:embed testdata/prometheus-component.yaml
var prometheusComponentBytes []byte

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
	t.Run("success without mappedValues", func(t *testing.T) {
		component := &Component{}
		spec := &client.ChartSpec{}

		mockChartGetter := NewMockChartGetter(t)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, "", mappedValuesYaml)
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
	t.Run("success no matching mapping", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"someOtherKey": "value",
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

		expectedYaml := `{}
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
	t.Run("success invalid mapping value", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"mainLogLevel": "panic",
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

		expectedYaml := `{}
`

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, expectedYaml, mappedValuesYaml)
		assert.NoError(t, err)
	})
	t.Run("error getting helm chart", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"mainLogLevel": "debug",
				},
			},
		}

		spec := &client.ChartSpec{}

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(nil, assert.AnError)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, "", mappedValuesYaml)
		assert.Error(t, err)
	})
	t.Run("error parsing yaml", func(t *testing.T) {
		component := &Component{
			Spec: ComponentSpec{
				MappedValues: map[string]string{
					"mainLogLevel": "debug",
				},
			},
		}

		doguOpMetaData := `apiVersion: v1
metavalues:
  mainLogLevel
something
invalid
        here`

		spec := &client.ChartSpec{}
		helmChart := &chart.Chart{
			Files: []*chart.File{
				{
					Name: mappingMetadataFileName,
					Data: []byte(doguOpMetaData),
				},
			},
		}

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())
		assert.Equal(t, "", mappedValuesYaml)
		assert.Error(t, err)
	})
	t.Run("error unmarshaling", func(t *testing.T) {
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
          debug: trace`

		spec := &client.ChartSpec{}
		helmChart := &chart.Chart{
			Files: []*chart.File{
				{
					Name: mappingMetadataFileName,
					Data: []byte(doguOpMetaData),
				},
			},
		}

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		marshaller := &InvalidMarshaller{
			orignalMarshaler: yaml.NewSerializer(),
			failMarshal:      true,
			failUnmarshal:    false,
		}
		mappedValuesYaml, err := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, marshaller)
		assert.Equal(t, "", mappedValuesYaml)
		assert.Error(t, err)
	})
	t.Run("success with complex mapping", func(t *testing.T) {
		serializer := yaml.NewSerializer()

		var origYaml map[string]interface{}
		_ = serializer.Unmarshal(prometheusComponentBytes, &origYaml)

		fmt.Print(origYaml)

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
      - path: spec.template.spec.containers[name=auth].env[name=LOG_LEVEL].value
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

		mockChartGetter := NewMockChartGetter(t)
		mockChartGetter.EXPECT().GetChart(testCtx, spec).Return(helmChart, nil)
		mappedValuesYaml, _ := getMappedValuesYaml(testCtx, component, spec, mockChartGetter, yaml.NewSerializer())

		var mappedYaml map[string]interface{}
		_ = serializer.Unmarshal([]byte(mappedValuesYaml), &mappedYaml)

		mergedMap := values.MergeMaps(origYaml, mappedYaml)

		fmt.Print(mergedMap)

	})
}

type InvalidMarshaller struct {
	failMarshal      bool
	failUnmarshal    bool
	orignalMarshaler yaml.Serializer
}

func (s *InvalidMarshaller) Marshal(o interface{}) ([]byte, error) {
	if s.failMarshal {
		return nil, assert.AnError
	}
	return s.orignalMarshaler.Marshal(o)
}

func (s *InvalidMarshaller) Unmarshal(y []byte, o interface{}, opts ...originalyaml.JSONOpt) error {
	if s.failUnmarshal {
		return assert.AnError
	}
	return s.orignalMarshaler.Unmarshal(y, opts)
}
