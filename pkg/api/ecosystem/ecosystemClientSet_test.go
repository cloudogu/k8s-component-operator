package ecosystem

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"testing"
)

func TestNewForConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		config := &rest.Config{}

		// when
		clientSet, err := NewForConfig(config)

		// then
		require.NoError(t, err)
		require.NotNil(t, clientSet)
	})
}

func TestEcoSystemV1Alpha1Client_Components(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		config := &rest.Config{}
		clientSet, err := NewForConfig(config)
		require.NoError(t, err)
		require.NotNil(t, clientSet)

		// when
		client := clientSet.Components("ecosystem")

		// then
		require.NotNil(t, client)
	})
}

func TestNewComponentClientset(t *testing.T) {
	t.Run("should create new componentClientset", func(t *testing.T) {
		// given
		config := &rest.Config{}
		clientSet := &kubernetes.Clientset{}

		// when
		client, err := NewComponentClientset(config, clientSet)

		// then
		require.NoError(t, err)
		require.NotNil(t, client)
		assert.Equal(t, clientSet, client.Clientset)
		assert.IsType(t, &V1Alpha1Client{}, client.ecosystemV1Alpha1)
	})
	t.Run("should fail to create new componentClientset for wrong config", func(t *testing.T) {
		// given
		config := &rest.Config{
			Host: "foo:/error",
		}
		clientSet := &kubernetes.Clientset{}

		// when
		client, err := NewComponentClientset(config, clientSet)

		// then
		require.Error(t, err)
		require.Nil(t, client)
		assert.ErrorContains(t, err, "host must be a URL or a host:port pair")
	})
}

func TestComponentV1Alpha1(t *testing.T) {
	t.Run("should return V1Alpha1Client", func(t *testing.T) {
		// given
		config := &rest.Config{}
		clientSet := &kubernetes.Clientset{}
		client, err := NewComponentClientset(config, clientSet)
		require.NoError(t, err)

		// when
		componentClient := client.ComponentV1Alpha1()

		// then
		require.NotNil(t, componentClient)
		assert.Equal(t, componentClient, client.ecosystemV1Alpha1)
	})
}
