package helm

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("should create new client", func(t *testing.T) {
		namespace := "ecosystem"
		debug := false

		// override default controller method to retrieve a kube config
		oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
		defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
		ctrl.GetConfigOrDie = func() *rest.Config {
			return &rest.Config{}
		}

		client, err := New(namespace, debug, nil)

		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("should fail to create client for ???", func(t *testing.T) {
		// set an existing-file which is not a helm-config to ptovoke an error
		t.Setenv("HELM_REGISTRY_CONFIG", "client.go")

		namespace := "error"
		debug := false

		// override default controller method to retrieve a kube config
		oldGetConfigOrDieDelegate := ctrl.GetConfigOrDie
		defer func() { ctrl.GetConfigOrDie = oldGetConfigOrDieDelegate }()
		ctrl.GetConfigOrDie = func() *rest.Config {
			return &rest.Config{}
		}

		_, err := New(namespace, debug, nil)

		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to create helm client:")
	})
}
