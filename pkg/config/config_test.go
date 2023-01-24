package config

import (
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
		assert.Equal(t, "helm", result.Username)
		assert.Equal(t, "helm", result.Password)
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
		secretClientMock := external.NewSecretInterface(t)
		dataMap := make(map[string][]byte)
		dataMap["username"] = []byte("username")
		dataMap["password"] = []byte("password")
		dataMap["endpoint"] = []byte("endpoint")
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "component-operator-helm-repository"},
			Data:       dataMap,
		}
		secretClientMock.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(secret, nil)

		// when
		result, err := GetHelmRepositoryData(secretClientMock)

		// then
		require.NoError(t, err)
		assert.Equal(t, "username", result.Username)
		assert.Equal(t, "password", result.Password)
		assert.Equal(t, "endpoint", result.Endpoint)
	})

	t.Run("should return not found error if no secret was found", func(t *testing.T) {
		// given
		secretClientMock := external.NewSecretInterface(t)
		notFoundError := errors.NewNotFound(schema.GroupResource{}, "")
		secretClientMock.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(nil, notFoundError)

		// when
		_, err := GetHelmRepositoryData(secretClientMock)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "helm repository secret component-operator-helm-repository not found")
	})

	t.Run("should return error on failed get", func(t *testing.T) {
		// given
		secretClientMock := external.NewSecretInterface(t)
		secretClientMock.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(nil, assert.AnError)

		// when
		_, err := GetHelmRepositoryData(secretClientMock)

		// then
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get helm repository secret")
	})
}
