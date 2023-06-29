package helm

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/registry"
	"testing"
)

func TestOCI(t *testing.T) {
	t.Run("test mirror oci helm chart", func(t *testing.T) {
		client, err := registry.NewClient()
		require.NoError(t, err)

		harborHost := "staging-registry.cloudogu.com"
		err = client.Login(harborHost, registry.LoginOptBasicAuth("todo", "todo"))
		require.NoError(t, err)

		result, err := client.Pull(fmt.Sprintf("%s/testing/k8s-dogu-operator:0.1.0", harborHost))
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NoError(t, client.Logout(harborHost))

		stagexHost := "stagex.cloudogu.com/v2"
		err = client.Login(stagexHost, registry.LoginOptBasicAuth("todo", "todo"))
		require.NoError(t, err)

		push, err := client.Push(result.Chart.Data, fmt.Sprintf("%s/helmtest/k8s-dogu-operator:0.1.0", stagexHost), registry.PushOptStrictMode(false))
		require.NoError(t, err)
		require.NotNil(t, push)
		require.NoError(t, client.Logout(stagexHost))
	})
}
