package ecoSystem

import (
	"context"
	"encoding/json"
	v1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_componentClient_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, "GET", request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/testcomponent", request.URL.Path)
			assert.Equal(t, http.NoBody, request.Body)

			writer.Header().Add("content-type", "application/json")
			component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "testcomponent", Namespace: "test"}}
			componentBytes, err := json.Marshal(component)
			require.NoError(t, err)
			_, err = writer.Write(componentBytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Get(context.TODO(), "testcomponent", metav1.GetOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodGet, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components", request.URL.Path)
			assert.Equal(t, http.NoBody, request.Body)

			writer.Header().Add("content-type", "application/json")
			componentList := v1.ComponentList{}
			component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "testcomponent", Namespace: "test"}}
			componentList.Items = append(componentList.Items, *component)
			componentBytes, err := json.Marshal(componentList)
			require.NoError(t, err)
			_, err = writer.Write(componentBytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.List(context.TODO(), metav1.ListOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "tocreate", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPost, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "tocreate", createdComponent.Name)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Create(context.TODO(), component, metav1.CreateOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "tocreate", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/tocreate", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "tocreate", createdComponent.Name)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Update(context.TODO(), component, metav1.UpdateOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_UpdateStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "tocreate", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/tocreate/status", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "tocreate", createdComponent.Name)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.UpdateStatus(context.TODO(), component, metav1.UpdateOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodDelete, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/testcomponent", request.URL.Path)

			writer.Header().Add("content-type", "application/json")
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		err = cClient.Delete(context.TODO(), "testcomponent", metav1.DeleteOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_DeleteCollection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodDelete, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components", request.URL.Path)
			assert.Equal(t, "labelSelector=test", request.URL.RawQuery)
			writer.Header().Add("content-type", "application/json")
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		err = cClient.DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "test"})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_Patch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPatch, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/testcomponent", request.URL.Path)
			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)
			assert.Equal(t, []byte("test"), bytes)
			result, err := json.Marshal(v1.Component{})
			require.NoError(t, err)

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(result)
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		patchData := []byte("test")

		// when
		_, err = cClient.Patch(context.TODO(), "testcomponent", types.JSONPatchType, patchData, metav1.PatchOptions{})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_Watch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, "GET", request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components", request.URL.Path)
			assert.Equal(t, http.NoBody, request.Body)
			assert.Equal(t, "labelSelector=test&watch=true", request.URL.RawQuery)

			writer.Header().Add("content-type", "application/json")
			_, err := writer.Write([]byte("egal"))
			require.NoError(t, err)
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Watch(context.TODO(), metav1.ListOptions{LabelSelector: "test"})

		// then
		require.NoError(t, err)
	})
}
