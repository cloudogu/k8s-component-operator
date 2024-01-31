package helm

import (
	_ "embed"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	helmChart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"testing"
)

//go:embed testdata/component-values-metadata.yaml
var simpleMetadata []byte

func Test_buildMapFromYamlPath(t *testing.T) {
	type args struct {
		path  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr error
	}{
		{
			name:    "should return empty map for empty path",
			args:    args{path: "", value: ""},
			want:    map[string]interface{}{},
			wantErr: nil,
		},
		{
			name:    "should return error on empty path and non empty value",
			args:    args{path: "", value: "not-empty"},
			want:    map[string]interface{}{},
			wantErr: fmt.Errorf("could not build map from yaml path for an empty path and value \"not-empty\""),
		},
		{
			name:    "success with on element in path",
			args:    args{path: "log-level", value: "info"},
			want:    map[string]interface{}{"log-level": "info"},
			wantErr: nil,
		},
		{
			name:    "success with two elements in path",
			args:    args{path: "controllerManager.log-level", value: "info"},
			want:    map[string]interface{}{"controllerManager": map[string]interface{}{"log-level": "info"}},
			wantErr: nil,
		},
		{
			name:    "success with three elements in path",
			args:    args{path: "env.controllerManager.log-level", value: "info"},
			want:    map[string]interface{}{"env": map[string]interface{}{"controllerManager": map[string]interface{}{"log-level": "info"}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := buildMapFromYamlPath(tt.args.path, tt.args.value)
			assert.Equalf(t, tt.want, path, "buildMapFromYamlPath(%v, %v)", tt.args.path, tt.args.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_metadataValueClient_IsValuesChanged(t *testing.T) {
	t.Run("should return false on same values in valuesYamlOverwrite", func(t *testing.T) {
		// given
		deployedRelease := &release.Release{Name: "operator"}
		deployedValues := map[string]interface{}{"loglevel": "info"}
		yamlBytes, err := yaml.Marshal(deployedValues)
		require.NoError(t, err)
		component := &k8sv1.Component{Spec: k8sv1.ComponentSpec{ValuesYamlOverwrite: string(yamlBytes)}}

		helmClientMock := newMockComponentHelmClient(t)
		helmClientMock.EXPECT().GetReleaseValues("operator", false).Return(deployedValues, nil)
		chartSpec := component.GetHelmChartSpec()
		overwriteValuesMap := deployedValues
		helmClientMock.EXPECT().GetChartSpecValues(chartSpec).Return(overwriteValuesMap, nil)
		chart := &helmChart.Chart{}
		helmClientMock.EXPECT().GetChart(testCtx, chartSpec).Return(chart, nil)

		sut := metadataValueClient{componentHelmClient: helmClientMock}

		// when
		changed, err := sut.IsValuesChanged(testCtx, deployedRelease, component)

		// then
		require.NoError(t, err)
		assert.False(t, changed)
	})

	t.Run("should return true if values changed in valuesYamlOverwrite", func(t *testing.T) {
		// given
		deployedRelease := &release.Release{Name: "operator"}
		deployedValues := map[string]interface{}{"loglevel": "info"}
		yamlBytes, err := yaml.Marshal(deployedValues)
		require.NoError(t, err)
		component := &k8sv1.Component{Spec: k8sv1.ComponentSpec{ValuesYamlOverwrite: string(yamlBytes)}}

		helmClientMock := newMockComponentHelmClient(t)
		helmClientMock.EXPECT().GetReleaseValues("operator", false).Return(deployedValues, nil)
		chartSpec := component.GetHelmChartSpec()
		overwriteValuesMap := map[string]interface{}{"loglevel": "error"}
		helmClientMock.EXPECT().GetChartSpecValues(chartSpec).Return(overwriteValuesMap, nil)
		chart := &helmChart.Chart{}
		helmClientMock.EXPECT().GetChart(testCtx, chartSpec).Return(chart, nil)

		sut := metadataValueClient{componentHelmClient: helmClientMock}

		// when
		changed, err := sut.IsValuesChanged(testCtx, deployedRelease, component)

		// then
		require.NoError(t, err)
		assert.True(t, changed)
	})

	t.Run("should return false on same values in valuesYamlOverwrite and mappedValues", func(t *testing.T) {
		// given
		deployedRelease := &release.Release{Name: "operator"}
		deployedValues := map[string]interface{}{"loglevel": "info", "stage": "development"}
		yamlBytes, err := yaml.Marshal(deployedValues)
		require.NoError(t, err)
		component := &k8sv1.Component{Spec: k8sv1.ComponentSpec{ValuesYamlOverwrite: string(yamlBytes), MappedValues: map[string]string{"mode": "development"}}}

		helmClientMock := newMockComponentHelmClient(t)
		helmClientMock.EXPECT().GetReleaseValues("operator", false).Return(deployedValues, nil)
		chartSpec := component.GetHelmChartSpec()
		overwriteValuesMap := map[string]interface{}{"loglevel": "info"}
		helmClientMock.EXPECT().GetChartSpecValues(chartSpec).Return(overwriteValuesMap, nil)
		chart := &helmChart.Chart{Metadata: &helmChart.Metadata{Name: "operator"}, Files: []*helmChart.File{{Name: "component-values-metadata.yaml", Data: simpleMetadata}}}

		helmClientMock.EXPECT().GetChart(testCtx, chartSpec).Return(chart, nil)

		sut := metadataValueClient{componentHelmClient: helmClientMock}

		// when
		changed, err := sut.IsValuesChanged(testCtx, deployedRelease, component)

		// then
		require.NoError(t, err)
		assert.False(t, changed)
	})
}
