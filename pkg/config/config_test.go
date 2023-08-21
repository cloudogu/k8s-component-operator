package config

import (
	"fmt"
	"github.com/cloudogu/k8s-component-operator/pkg/mocks/external"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOperatorConfig(t *testing.T) {
	_ = os.Unsetenv("NAMESPACE")

	expectedNamespace := "myNamespace"

	t.Run("Error on missing namespace env var", func(t *testing.T) {
		// when
		operatorConfig, err := NewOperatorConfig("0.0.0")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to read namespace: environment variable NAMESPACE must be set")
		assert.Nil(t, operatorConfig)
	})
	t.Setenv("NAMESPACE", expectedNamespace)
	t.Run("Create config successfully", func(t *testing.T) {
		// when
		operatorConfig, err := NewOperatorConfig("0.1.0")

		// then
		require.NoError(t, err)
		require.NotNil(t, operatorConfig)
		assert.Equal(t, expectedNamespace, operatorConfig.Namespace)
		assert.Equal(t, "0.1.0", operatorConfig.Version.Raw)
	})
}

func TestGetHelmRepositoryData(t *testing.T) {
	t.Setenv("RUNTIME", "local")
	t.Run("success from file", func(t *testing.T) {
		// given
		devHelmRepoDataPath = "../../k8s/helm-repository.yaml"

		// when
		result, err := GetHelmRepositoryData(nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, "http://192.168.56.3:30100", result.Endpoint)
	})

	t.Run("should throw error because the file does not exists", func(t *testing.T) {
		// given
		devHelmRepoDataPath = "testdata/ne.yaml"

		// when
		_, err := GetHelmRepositoryData(nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "could not find configuration at")
	})

	t.Run("should throw error because wrong yaml format", func(t *testing.T) {
		// given
		devHelmRepoDataPath = "testdata/helm-repository.yaml"

		// when
		_, err := GetHelmRepositoryData(nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to unmarshal configuration")
	})

	require.NoError(t, os.Unsetenv("RUNTIME"))
	t.Run("success with cluster", func(t *testing.T) {
		// given
		mockConfigMapInterface := external.NewMockConfigMapInterface(t)
		configMap := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "component-operator-helm-repository"},
			Data:       map[string]string{"endpoint": "endpoint"},
		}
		mockConfigMapInterface.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(configMap, nil)

		// when
		result, err := GetHelmRepositoryData(mockConfigMapInterface)

		// then
		require.NoError(t, err)
		assert.Equal(t, "endpoint", result.Endpoint)
	})

	t.Run("should return not found error if no secret was found", func(t *testing.T) {
		// given
		mockConfigMapInterface := external.NewMockConfigMapInterface(t)
		notFoundError := errors.NewNotFound(schema.GroupResource{}, "")
		mockConfigMapInterface.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(nil, notFoundError)

		// when
		_, err := GetHelmRepositoryData(mockConfigMapInterface)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "helm repository configMap component-operator-helm-repository not found")
	})

	t.Run("should return error on failed get", func(t *testing.T) {
		// given
		mockConfigMapInterface := external.NewMockConfigMapInterface(t)
		mockConfigMapInterface.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(nil, assert.AnError)

		// when
		_, err := GetHelmRepositoryData(mockConfigMapInterface)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get helm repository configMap")
	})
}

func TestHelmRepositoryData_GetOciEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		Endpoint string
		want     string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "success getOciEndpoint",
			Endpoint: "https://staging-registry.cloudogu.com",
			want:     "oci://staging-registry.cloudogu.com",
			wantErr:  assert.NoError,
		},
		{
			name:     "success getOciEndpoint with Path",
			Endpoint: "https://staging-registry.cloudogu.com/foo/bar",
			want:     "oci://staging-registry.cloudogu.com/foo/bar",
			wantErr:  assert.NoError,
		},
		{
			name:     "success getOciEndpoint with other protocol",
			Endpoint: "ftp://staging-registry.cloudogu.com",
			want:     "oci://staging-registry.cloudogu.com",
			wantErr:  assert.NoError,
		},
		{
			name:     "success no protocol",
			Endpoint: "staging-registry.cloudogu.com",
			want:     "oci://staging-registry.cloudogu.com",
			wantErr:  assert.NoError,
		},
		{
			name:     "error empty string",
			Endpoint: "",
			want:     "",
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hrd := &HelmRepositoryData{
				Endpoint: tt.Endpoint,
			}
			got, err := hrd.GetOciEndpoint()
			if !tt.wantErr(t, err, fmt.Sprintf("GetOciEndpoint()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetOciEndpoint()")
		})
	}
}
