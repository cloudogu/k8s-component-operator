package client

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"k8s.io/client-go/rest"
)

//go:embed testdata/kubeconfig.yaml
var kubeconfigBytes []byte

//go:embed testdata/invalid-kubeconfig.yaml
var invalidKubeconfigBytes []byte

func TestRESTClientGetter_ToRESTConfig(t *testing.T) {
	t.Run("should return rest config if not nil", func(t *testing.T) {
		// given
		restConfig := &rest.Config{}
		sut := &RESTClientGetter{restConfig: restConfig}

		// when
		actual, err := sut.ToRESTConfig()

		// then
		require.NoError(t, err)
		assert.Same(t, restConfig, actual)
	})
	t.Run("should return rest config for kubeconfig", func(t *testing.T) {
		// given
		sut := &RESTClientGetter{kubeConfig: kubeconfigBytes}

		// when
		actual, err := sut.ToRESTConfig()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRESTClientGetter_ToDiscoveryClient(t *testing.T) {
	t.Run("should fail to create rest config", func(t *testing.T) {
		// given
		sut := &RESTClientGetter{kubeConfig: invalidKubeconfigBytes}

		// when
		_, err := sut.ToDiscoveryClient()

		// then
		require.Error(t, err)
	})
	t.Run("should create discovery client", func(t *testing.T) {
		// given
		sut := &RESTClientGetter{
			kubeConfig: kubeconfigBytes,
			opts: []RESTClientOption{func(config *rest.Config) {
				config.UserAgent = "my-user-agent"
			}},
		}

		// when
		actual, err := sut.ToDiscoveryClient()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRESTClientGetter_ToRESTMapper(t *testing.T) {
	t.Run("should fail to create discovery client", func(t *testing.T) {
		// given
		sut := &RESTClientGetter{kubeConfig: invalidKubeconfigBytes}

		// when
		_, err := sut.ToRESTMapper()

		// then
		require.Error(t, err)
	})
	t.Run("should succeed to create rest mapper", func(t *testing.T) {
		// given
		sut := &RESTClientGetter{kubeConfig: kubeconfigBytes}

		// when
		actual, err := sut.ToRESTMapper()

		// then
		require.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestRESTClientGetter_ToRawKubeConfigLoader(t *testing.T) {
	t.Run("should return config loader", func(t *testing.T) {
		// given
		sut := &RESTClientGetter{namespace: "test-namespace"}

		// when
		actual := sut.ToRawKubeConfigLoader()

		// then
		assert.NotNil(t, actual)
	})
}
