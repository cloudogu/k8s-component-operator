package ecosystem

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var testCtx = context.Background()

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
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Get(testCtx, "testcomponent", metav1.GetOptions{})

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

		timeout := int64(5)

		// when
		_, err = cClient.List(testCtx, metav1.ListOptions{TimeoutSeconds: &timeout})

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
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Create(testCtx, component, metav1.CreateOptions{})

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
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.Update(testCtx, component, metav1.UpdateOptions{})

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
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.UpdateStatus(testCtx, component, metav1.UpdateOptions{})

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
		err = cClient.Delete(testCtx, "testcomponent", metav1.DeleteOptions{})

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
			assert.Equal(t, "labelSelector=test&timeout=5s&timeoutSeconds=5", request.URL.RawQuery)
			writer.Header().Add("content-type", "application/json")
			writer.WriteHeader(200)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")
		timeout := int64(5)

		// when
		err = cClient.DeleteCollection(testCtx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: "test", TimeoutSeconds: &timeout})

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
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		patchData := []byte("test")

		// when
		_, err = cClient.Patch(testCtx, "testcomponent", types.JSONPatchType, patchData, metav1.PatchOptions{})

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
			assert.Equal(t, "labelSelector=test&timeout=5s&timeoutSeconds=5&watch=true", request.URL.RawQuery)

			writer.Header().Add("content-type", "application/json")
			_, err := writer.Write([]byte("egal"))
			require.NoError(t, err)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		timeout := int64(5)

		// when
		_, err = cClient.Watch(testCtx, metav1.ListOptions{LabelSelector: "test", TimeoutSeconds: &timeout})

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_UpdateStatusInstalling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "installing", false, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusInstalling(testCtx, component)

		// then
		require.NoError(t, err)
	})

	t.Run("success with retry", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "installing", true, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusInstalling(testCtx, component)

		// then
		require.NoError(t, err)
	})

	t.Run("fail on get component", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "installing", false, true)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusInstalling(testCtx, component)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "an error on the server (\"\") has prevented the request from succeeding (get components.k8s.cloudogu.com myComponent)")
	})
}

func Test_componentClient_UpdateStatusInstalled(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "installed", false, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusInstalled(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("success with retry", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "installed", true, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusInstalled(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("fail on get component", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "installed", false, true)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusInstalled(testCtx, component)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "an error on the server (\"\") has prevented the request from succeeding (get components.k8s.cloudogu.com myComponent)")
	})
}

func Test_componentClient_UpdateStatusNotInstalled(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "", false, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusNotInstalled(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("success with retry", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "", true, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusNotInstalled(testCtx, component)

		// then
		require.NoError(t, err)
	})
}

func Test_componentClient_UpdateStatusUpgrading(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "upgrading", false, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusUpgrading(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("success with retry", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "upgrading", true, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusUpgrading(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("fail on get component", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "upgrading", false, true)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusUpgrading(testCtx, component)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "an error on the server (\"\") has prevented the request from succeeding (get components.k8s.cloudogu.com myComponent)")
	})
}

func Test_componentClient_UpdateStatusDeleting(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "deleting", false, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusDeleting(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("success with retry", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "deleting", true, false)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusDeleting(testCtx, component)

		// then
		require.NoError(t, err)
	})
	t.Run("fail on get component", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		mockClient := mockClientForStatusUpdates(t, component, "deleting", false, true)
		cClient := mockClient.Components("test")

		// when
		_, err := cClient.UpdateStatusDeleting(testCtx, component)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "an error on the server (\"\") has prevented the request from succeeding (get components.k8s.cloudogu.com myComponent)")
	})
}

func Test_componentClient_AddFinalizer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "myComponent", createdComponent.Name)
			assert.Len(t, createdComponent.Finalizers, 1)
			assert.Equal(t, "myFinalizer", createdComponent.Finalizers[0])

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.AddFinalizer(testCtx, component, "myFinalizer")

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to set finalizer on client error", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "myComponent", createdComponent.Name)
			assert.Len(t, createdComponent.Finalizers, 1)
			assert.Equal(t, "myFinalizer", createdComponent.Finalizers[0])

			writer.WriteHeader(500)
			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.AddFinalizer(testCtx, component, "myFinalizer")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to add finalizer myFinalizer to component:")
	})
}

func Test_componentClient_RemoveFinalizer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		controllerutil.AddFinalizer(component, "finalizer1")
		controllerutil.AddFinalizer(component, "finalizer2")

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "myComponent", createdComponent.Name)
			assert.Len(t, createdComponent.Finalizers, 1)
			assert.Equal(t, "finalizer2", createdComponent.Finalizers[0])

			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.RemoveFinalizer(testCtx, component, "finalizer1")

		// then
		require.NoError(t, err)
	})

	t.Run("should fail to set finalizer on client error", func(t *testing.T) {
		// given
		component := &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"}}
		controllerutil.AddFinalizer(component, "finalizer1")
		controllerutil.AddFinalizer(component, "finalizer2")

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, http.MethodPut, request.Method)
			assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)

			bytes, err := io.ReadAll(request.Body)
			require.NoError(t, err)

			createdComponent := &v1.Component{}
			require.NoError(t, json.Unmarshal(bytes, createdComponent))
			assert.Equal(t, "myComponent", createdComponent.Name)
			assert.Len(t, createdComponent.Finalizers, 1)
			assert.Equal(t, "finalizer1", createdComponent.Finalizers[0])

			writer.WriteHeader(500)
			writer.Header().Add("content-type", "application/json")
			_, err = writer.Write(bytes)
			require.NoError(t, err)
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		// when
		_, err = cClient.RemoveFinalizer(testCtx, component, "finalizer2")

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "failed to remove finalizer finalizer2 from component")
	})
}

func mockClientForStatusUpdates(t *testing.T, expectedComponent *v1.Component, expectedStatus string, withRetry bool, failOnGetComponent bool) *V1Alpha1Client {

	failGetComponent := func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, fmt.Sprintf("/apis/k8s.cloudogu.com/v1/namespaces/test/components/%s", expectedComponent.Name), request.URL.Path)

		writer.WriteHeader(500)
	}

	assertGetComponentRequest := func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		assert.Equal(t, fmt.Sprintf("/apis/k8s.cloudogu.com/v1/namespaces/test/components/%s", expectedComponent.Name), request.URL.Path)

		componentJson, err := json.Marshal(expectedComponent)
		require.NoError(t, err)

		writer.Header().Add("content-type", "application/json")
		_, err = writer.Write(componentJson)
		require.NoError(t, err)
	}

	assertUpdateStatusRequest := func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPut, request.Method)
		assert.Equal(t, fmt.Sprintf("/apis/k8s.cloudogu.com/v1/namespaces/test/components/%s/status", expectedComponent.Name), request.URL.Path)

		bytes, err := io.ReadAll(request.Body)
		require.NoError(t, err)

		createdComponent := &v1.Component{}
		require.NoError(t, json.Unmarshal(bytes, createdComponent))
		assert.Equal(t, expectedComponent.Name, createdComponent.Name)
		assert.Equal(t, expectedStatus, createdComponent.Status.Status)

		writer.Header().Add("content-type", "application/json")
		_, err = writer.Write(bytes)
		require.NoError(t, err)
	}

	conflictUpdateStatusRequest := func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPut, request.Method)
		assert.Equal(t, fmt.Sprintf("/apis/k8s.cloudogu.com/v1/namespaces/test/components/%s/status", expectedComponent.Name), request.URL.Path)

		writer.WriteHeader(409)
	}

	var requestAssertions []func(writer http.ResponseWriter, request *http.Request)

	if failOnGetComponent {
		requestAssertions = []func(writer http.ResponseWriter, request *http.Request){
			failGetComponent,
		}
	} else if withRetry {
		requestAssertions = []func(writer http.ResponseWriter, request *http.Request){
			assertGetComponentRequest,
			conflictUpdateStatusRequest,
			assertGetComponentRequest,
			assertUpdateStatusRequest,
		}
	} else {
		requestAssertions = []func(writer http.ResponseWriter, request *http.Request){
			assertGetComponentRequest,
			assertUpdateStatusRequest,
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assertRequestFunc := requestAssertions[0]
		requestAssertions = requestAssertions[1:]

		assertRequestFunc(writer, request)
	}))

	config := rest.Config{
		Host: server.URL,
	}
	client, err := NewForConfig(&config)
	require.NoError(t, err)
	return client
}

func Test_componentClient_UpdateExpectedComponentVersion(t *testing.T) {
	t.Run("ok, no retry", func(t *testing.T) {
		// given
		componentWithoutVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "",
			},
		}
		componentWithVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "1.0.0",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Contains(t, []string{"GET", "PUT"}, request.Method)
			if request.Method == "GET" {
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)
				assert.Equal(t, http.NoBody, request.Body)

				writer.Header().Add("content-type", "application/json")
				componentBytes, err := json.Marshal(componentWithoutVersion)
				require.NoError(t, err)
				_, err = writer.Write(componentBytes)
				require.NoError(t, err)
			} else {
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)

				bytes, err := io.ReadAll(request.Body)
				require.NoError(t, err)

				createdComponent := &v1.Component{}
				require.NoError(t, json.Unmarshal(bytes, createdComponent))
				assert.Equal(t, "myComponent", createdComponent.Name)

				writer.Header().Add("content-type", "application/json")
				_, err = writer.Write(bytes)
				require.NoError(t, err)
			}
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		//then
		updatedComponent, err := cClient.UpdateExpectedComponentVersion(
			testCtx,
			componentWithoutVersion.Name,
			"1.0.0",
		)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, componentWithVersion, updatedComponent)
	})
	t.Run("ok, with retry", func(t *testing.T) {
		// given
		componentWithoutVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "",
			},
		}
		componentWithVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "1.0.0",
			},
		}
		requestPutCounter := 0
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Contains(t, []string{"GET", "PUT"}, request.Method)
			if request.Method == "GET" {
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)
				assert.Equal(t, http.NoBody, request.Body)

				writer.Header().Add("content-type", "application/json")
				componentBytes, err := json.Marshal(componentWithoutVersion)
				require.NoError(t, err)
				_, err = writer.Write(componentBytes)
				require.NoError(t, err)
			} else {
				requestPutCounter++
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)
				bytes, err := io.ReadAll(request.Body)
				createdComponent := &v1.Component{}
				require.NoError(t, json.Unmarshal(bytes, createdComponent))
				assert.Equal(t, "myComponent", createdComponent.Name)
				require.NoError(t, err)
				if requestPutCounter == 1 {
					writer.WriteHeader(http.StatusConflict)
				} else if requestPutCounter == 2 {
					writer.Header().Add("content-type", "application/json")
					_, err = writer.Write(bytes)
					require.NoError(t, err)
				}
			}
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		//then
		updatedComponent, err := cClient.UpdateExpectedComponentVersion(
			testCtx,
			componentWithoutVersion.Name,
			"1.0.0",
		)

		//assert
		assert.NoError(t, err)
		assert.Equal(t, componentWithVersion, updatedComponent)
	})
	t.Run("error getting component", func(t *testing.T) {
		// given
		componentWithoutVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Contains(t, []string{"GET"}, request.Method)
			if request.Method == "GET" {
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)
				assert.Equal(t, http.NoBody, request.Body)
				writer.WriteHeader(500)
			}
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		//then
		_, err = cClient.UpdateExpectedComponentVersion(
			testCtx,
			componentWithoutVersion.Name,
			"1.0.0",
		)

		//assert
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to update version in component \"myComponent\": failed to get component \"myComponent\" for update")
	})
	t.Run("error updating component", func(t *testing.T) {
		// given
		componentWithoutVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "",
			},
		}
		componentWithVersion := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "myComponent", Namespace: "test"},
			Spec: v1.ComponentSpec{
				Name:    "myComponent",
				Version: "1.0.0",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Contains(t, []string{"GET", "PUT"}, request.Method)
			if request.Method == "GET" {
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)
				assert.Equal(t, http.NoBody, request.Body)

				writer.Header().Add("content-type", "application/json")
				componentBytes, err := json.Marshal(componentWithoutVersion)
				require.NoError(t, err)
				_, err = writer.Write(componentBytes)
				require.NoError(t, err)
			} else {
				assert.Equal(t, "/apis/k8s.cloudogu.com/v1/namespaces/test/components/myComponent", request.URL.Path)

				bytes, err := io.ReadAll(request.Body)
				require.NoError(t, err)

				createdComponent := &v1.Component{}
				require.NoError(t, json.Unmarshal(bytes, createdComponent))
				assert.Equal(t, componentWithVersion.Spec.Name, createdComponent.Spec.Name)
				assert.Equal(t, componentWithVersion.Spec.Version, createdComponent.Spec.Version)
				writer.WriteHeader(500)
			}
		}))

		config := rest.Config{
			Host: server.URL,
		}
		client, err := NewForConfig(&config)
		require.NoError(t, err)
		cClient := client.Components("test")

		//then
		_, err = cClient.UpdateExpectedComponentVersion(
			testCtx,
			componentWithoutVersion.Name,
			"1.0.0",
		)

		//assert
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to update version in component \"myComponent\"")
	})
}
