package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testCtx = context.Background()

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
		assert.Equal(t, "0.1.0", operatorConfig.Version.Original())
	})
}

func TestGetHelmRepositoryData(t *testing.T) {
	t.Run("should return local developer", func(t *testing.T) {
		// given
		t.Setenv("RUNTIME", "local")
		devHelmRepoDataPath = "testdata/helm-repository.yaml"
		expected := &HelmRepositoryData{
			Endpoint:    "192.168.56.3:30100",
			Schema:      EndpointSchemaOCI,
			PlainHttp:   true,
			InsecureTLS: true,
		}

		// when
		result, err := GetHelmRepositoryData(testCtx, nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("success with cluster", func(t *testing.T) {
		// given
		t.Setenv("RUNTIME", "")
		mockConfigMapInterface := newMockConfigMapInterface(t)
		configMap := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: "component-operator-helm-repository"},
			Data:       map[string]string{"endpoint": "endpoint", "schema": "oci", "plainHttp": "false", "insecureTls": "true"},
		}
		mockConfigMapInterface.On("Get", mock.Anything, "component-operator-helm-repository", mock.Anything).Return(configMap, nil)
		expected := &HelmRepositoryData{
			Endpoint:    "endpoint",
			Schema:      EndpointSchemaOCI,
			PlainHttp:   false,
			InsecureTLS: true,
		}

		// when
		result, err := GetHelmRepositoryData(testCtx, mockConfigMapInterface)

		// then
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestNewHelmRepoDataFromCluster(t *testing.T) {
	getOpts := metav1.GetOptions{}
	t.Run("should fail on getting config map", func(t *testing.T) {
		// given
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(nil, assert.AnError)

		// when
		_, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to get helm repository configMap component-operator-helm-repository")
	})
	t.Run("should fail to parse plainHttp", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{Data: map[string]string{"plainHttp": "invalid"}}
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(configMap, nil)

		// when
		_, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse field plainHttp from configMap component-operator-helm-repository")
	})
	t.Run("should fail to parse insecureTls", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{Data: map[string]string{"insecureTls": "invalid"}}
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(configMap, nil)

		// when
		_, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to parse field insecureTls from configMap component-operator-helm-repository")
	})
	t.Run("should fail because endpoint has empty URL", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{Data: map[string]string{"endpoint": ""}}
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(configMap, nil)

		// when
		_, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "config map 'component-operator-helm-repository' failed validation: endpoint URL must not be empty")
	})
	t.Run("should fail because endpoint schema is empty", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{Data: map[string]string{"endpoint": "myEndpoint"}}
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(configMap, nil)

		// when
		_, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "config map 'component-operator-helm-repository' failed validation: endpoint uses an unsupported schema '': valid schemas are: oci")
	})
	t.Run("should fail because endpoint schema is unsupported", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{Data: map[string]string{"endpoint": "myEndpoint", "schema": "https"}}
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(configMap, nil)

		// when
		_, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "config map 'component-operator-helm-repository' failed validation: endpoint uses an unsupported schema 'https': valid schemas are: oci")
	})
	t.Run("should succeed to parse plainHttp and insecureTls and validate endpoint", func(t *testing.T) {
		// given
		configMap := &v1.ConfigMap{Data: map[string]string{"endpoint": "myEndpoint", "schema": "oci", "plainHttp": "true", "insecureTls": "true"}}
		configMapClient := newMockConfigMapInterface(t)
		configMapClient.EXPECT().Get(testCtx, "component-operator-helm-repository", getOpts).Return(configMap, nil)

		// when
		actual, err := NewHelmRepoDataFromCluster(testCtx, configMapClient)

		// then
		expected := &HelmRepositoryData{
			Endpoint:    "myEndpoint",
			Schema:      EndpointSchemaOCI,
			PlainHttp:   true,
			InsecureTLS: true,
		}
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestNewHelmRepoDataFromFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		want     *HelmRepositoryData
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "should not find file",
			filepath: "not-exist",
			want:     nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to read configuration not-exist") &&
					assert.ErrorContains(t, err, "no such file")
			},
		},
		{
			name:     "should fail to unmarshal yaml",
			filepath: "testdata/invalid-helm-repository.yaml",
			want:     nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to unmarshal configuration testdata/invalid-helm-repository.yaml")
			},
		},
		{
			name:     "should succeed",
			filepath: "testdata/helm-repository.yaml",
			want: &HelmRepositoryData{
				Endpoint:    "192.168.56.3:30100",
				Schema:      EndpointSchemaOCI,
				PlainHttp:   true,
				InsecureTLS: true,
			},
			wantErr: assert.NoError,
		},
		{
			name:     "should fail with validation",
			filepath: "testdata/invalid-endpoint-helm-repository.yaml",
			want:     nil,
			wantErr: func(t assert.TestingT, err error, _ ...interface{}) bool {
				return assert.ErrorContains(t, err, "helm repository data from file 'testdata/invalid-endpoint-helm-repository.yaml' failed validation")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHelmRepoDataFromFile(tt.filepath)
			if !tt.wantErr(t, err, fmt.Sprintf("NewHelmRepoDataFromFile(%v)", tt.filepath)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewHelmRepoDataFromFile(%v)", tt.filepath)
		})
	}
}

func TestHelmRepositoryData_URL(t *testing.T) {
	actual := &HelmRepositoryData{
		Endpoint: "example.com",
		Schema:   "oci",
	}
	assert.Equal(t, "oci://example.com", actual.URL())
}

func Test_readMinuteDurationEnv(t *testing.T) {
	tests := []struct {
		name        string
		setEnvVar   bool
		envVarValue string
		want        time.Duration
		wantLogs    bool
		wantedLogs  string
		logLevel    logrus.Level
	}{
		{
			name:       "Environment variable not set",
			setEnvVar:  false,
			want:       15 * time.Minute,
			wantLogs:   true,
			wantedLogs: "failed to read HELM_CLIENT_TIMEOUT_MINS environment variable, using default value",
			logLevel:   logrus.DebugLevel,
		},
		{
			name:        "Environment variable not set correctly",
			setEnvVar:   true,
			envVarValue: "15//",
			want:        15 * time.Minute,
			wantLogs:    true,
			wantedLogs:  "failed to parse HELM_CLIENT_TIMEOUT_MINS environment variable, using default value",
			logLevel:    logrus.WarnLevel,
		},
		{
			name:        "read negative environment variable",
			setEnvVar:   true,
			envVarValue: "-20",
			want:        15 * time.Minute,
			wantLogs:    true,
			wantedLogs:  "parsed value (-20) of HELM_CLIENT_TIMEOUT_MINS is smaller than 0, using default value",
			logLevel:    logrus.WarnLevel,
		},
		{
			name:        "Successfully read environment variable",
			setEnvVar:   true,
			envVarValue: "20",
			want:        20 * time.Minute,
			wantLogs:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnvVar {
				err := os.Setenv(envHelmClientTimeoutMins, tt.envVarValue)
				require.NoError(t, err)
			}
			var result = time.Duration(0)

			var logOutput bytes.Buffer

			originalOutput := logrus.StandardLogger().Out
			originalLevel := logrus.StandardLogger().Level
			if tt.wantLogs {
				logrus.StandardLogger().SetOutput(&logOutput)
				logrus.StandardLogger().SetLevel(tt.logLevel)
			}

			result = readMinuteDurationEnv(envHelmClientTimeoutMins, defaultHelmClientTimeoutMins)

			logrus.StandardLogger().SetOutput(originalOutput)
			logrus.StandardLogger().SetLevel(originalLevel)

			logs := logOutput.String()

			assert.Equalf(t, tt.want, result, "readMinuteDurationEnv(%s, %s)", envHelmClientTimeoutMins, defaultHelmClientTimeoutMins)

			if tt.wantLogs {
				assert.Contains(t, logs, tt.wantedLogs)
			}
		})
	}
}
