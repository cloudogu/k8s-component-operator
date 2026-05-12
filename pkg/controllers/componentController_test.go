package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const testNamespace = "testtestNamespace"
const testRequeueTime = time.Second

func TestNewComponentReconciler(t *testing.T) {
	// given
	coreV1Mock := newMockCoreV1Interface(t)
	coreV1Mock.EXPECT().ConfigMaps(testNamespace).Return(newMockConfigMapInterface(t))
	clientSetMock := newMockComponentEcosystemInterface(t)
	clientSetMock.EXPECT().CoreV1().Return(coreV1Mock)

	configMapRefReaderMock := newMockConfigMapRefReader(t)

	newHelmClientFunc := func() (*helm.Client, error) { return nil, nil }
	mockRecorder := newMockEventRecorder(t)

	// when
	manager := NewComponentReconciler(clientSetMock, newHelmClientFunc, mockRecorder, testNamespace, defaultHelmClientTimeoutMins, yaml.NewSerializer(), configMapRefReaderMock, testRequeueTime)

	// then
	require.NotNil(t, manager)
}

func Test_newHelmClientFunc_NewHelmClient(t *testing.T) {
	t.Run("should delegate to wrapped function", func(t *testing.T) {
		expectedClient := &helm.Client{}
		sut := newHelmClientFunc(func() (*helm.Client, error) {
			return expectedClient, nil
		})

		actual, err := sut.NewHelmClient()

		require.NoError(t, err)
		assert.Same(t, expectedClient, actual)
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		sut := newHelmClientFunc(func() (*helm.Client, error) {
			return nil, assert.AnError
		})

		actual, err := sut.NewHelmClient()

		require.Error(t, err)
		assert.Nil(t, actual)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_componentReconciler_Reconcile(t *testing.T) {
	helmNamespace := "k8s"
	t.Run("success install", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Installation", "Installation successful")

		manager := NewMockComponentManager(t)
		manager.EXPECT().Install(testCtx, component).Return(nil)
		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)
		componentManagerFactory.EXPECT().NewComponentManager(helmClient).Return(manager)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Installation failed with component dogu-op", component, nil, mock.Anything).Return(reconcile.Result{}, nil)

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return(Install, nil)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			recorder:                  mockRecorder,
			componentManagerFactory:   componentManagerFactory,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
			helmClientFactory:         helmClientFactory,
			requeueHandler:            mockRequeueHandler,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("success delete", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.DeletionTimestamp = &v1.Time{Time: time.Now()}

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Deinstallation", "Deinstallation successful")

		manager := NewMockComponentManager(t)
		manager.EXPECT().Delete(testCtx, component).Return(nil)
		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)
		componentManagerFactory.EXPECT().NewComponentManager(helmClient).Return(manager)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Deinstallation failed with component dogu-op", component, nil, mock.Anything).Return(reconcile.Result{}, nil)

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return(Delete, nil)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			recorder:                  mockRecorder,
			componentManagerFactory:   componentManagerFactory,
			helmClientFactory:         helmClientFactory,
			requeueHandler:            mockRequeueHandler,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("success upgrade", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Normal", "Upgrade", "Upgrade successful")

		manager := NewMockComponentManager(t)
		manager.EXPECT().Upgrade(testCtx, component).Return(nil)

		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)
		componentManagerFactory.EXPECT().NewComponentManager(helmClient).Return(manager)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Upgrade failed with component dogu-op", component, nil, mock.Anything).Return(reconcile.Result{}, nil)

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return(Upgrade, nil)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			recorder:                  mockRecorder,
			componentManagerFactory:   componentManagerFactory,
			helmClientFactory:         helmClientFactory,
			requeueHandler:            mockRequeueHandler,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("should fail on downgrade", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Warning", "Downgrade", "component downgrades are not allowed")

		manager := NewMockComponentManager(t)
		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)
		componentManagerFactory.EXPECT().NewComponentManager(helmClient).Return(manager)

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return(Downgrade, nil)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			recorder:                  mockRecorder,
			componentManagerFactory:   componentManagerFactory,
			helmClientFactory:         helmClientFactory,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "downgrades are not allowed")
	})

	t.Run("should ignore equal installed component", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"
		component.Spec.ValuesConfigRef = &k8sv1.Reference{}

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)
		configMapRefReaderMock := newMockConfigMapRefReader(t)

		mockRecorder := newMockEventRecorder(t)
		manager := NewMockComponentManager(t)
		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)
		componentManagerFactory.EXPECT().NewComponentManager(helmClient).Return(manager)

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return(Ignore, nil)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			recorder:                  mockRecorder,
			componentManagerFactory:   componentManagerFactory,
			helmClientFactory:         helmClientFactory,
			timeout:                   defaultHelmClientTimeoutMins,
			yamlSerializer:            yaml.NewSerializer(),
			reader:                    configMapRefReaderMock,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
		assert.Equal(t, reconcile.Result{}, result)
	})

	t.Run("should fail on component get error", func(t *testing.T) {
		// given
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(nil, assert.AnError)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		sut := ComponentReconciler{
			clientSet: clientSetMock,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, assert.AnError, err)
	})

	t.Run("should return nil if the component is not found", func(t *testing.T) {
		// given

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(nil, errors.NewNotFound(schema.GroupResource{}, ""))
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		sut := ComponentReconciler{
			clientSet: clientSetMock,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.NoError(t, err)
	})

	t.Run("should fail on creating helm client", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(nil, assert.AnError)

		sut := ComponentReconciler{
			clientSet:         clientSetMock,
			helmClientFactory: helmClientFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		result, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.Equal(t, reconcile.Result{}, result)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to create helm client")
	})

	t.Run("should fail on getting operation with invalid versions", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.Status.Status = "installed"

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return("", assert.AnError)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			helmClientFactory:         helmClientFactory,
			componentManagerFactory:   componentManagerFactory,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to evaluate required operation")
	})

	t.Run("should fail on error in operation", func(t *testing.T) {
		// given
		component := getComponent(testNamespace, helmNamespace, "", "dogu-op", "0.1.0")
		component.DeletionTimestamp = &v1.Time{Time: time.Now()}

		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().Get(testCtx, "dogu-op", v1.GetOptions{}).Return(component, nil)
		componentClientGetterMock := newMockComponentV1Alpha1Interface(t)
		componentClientGetterMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentClientGetterMock)

		mockRecorder := newMockEventRecorder(t)
		mockRecorder.EXPECT().Event(component, "Warning", "Deinstallation", "Deinstallation failed. Reason: assert.AnError general error for testing")
		mockRecorder.EXPECT().Eventf(component, "Warning", "Requeue", "Failed to requeue the %s.", "deinstallation").Return()

		manager := NewMockComponentManager(t)
		manager.EXPECT().Delete(testCtx, component).Return(assert.AnError)
		helmClient := newMockHelmClient(t)
		helmClientFactory := newMockHelmClientFactory(t)
		helmClientFactory.EXPECT().NewHelmClient().Return(helmClient, nil)
		componentManagerFactory := newMockComponentManagerFactory(t)
		componentManagerFactory.EXPECT().NewComponentManager(helmClient).Return(manager)

		mockRequeueHandler := newMockRequeueHandler(t)
		mockRequeueHandler.EXPECT().Handle(testCtx, "Deinstallation failed with component dogu-op", component, assert.AnError, mock.Anything).
			RunAndReturn(func(_ context.Context, _ string, _ *k8sv1.Component, err error, _ string) (reconcile.Result, error) {
				return reconcile.Result{}, err
			})

		mockOperationEvaluator := newMockOperationEvaluator(t)
		mockOperationEvaluator.EXPECT().EvaluateRequiredOperation(testCtx, component).Return(Delete, nil)
		mockOperationEvaluatorFactory := newMockOperationEvaluatorFactory(t)
		mockOperationEvaluatorFactory.EXPECT().NewOperationEvaluator(helmClient).Return(mockOperationEvaluator)

		sut := ComponentReconciler{
			clientSet:                 clientSetMock,
			recorder:                  mockRecorder,
			componentManagerFactory:   componentManagerFactory,
			helmClientFactory:         helmClientFactory,
			requeueHandler:            mockRequeueHandler,
			operationEvaluatorFactory: mockOperationEvaluatorFactory,
		}
		req := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: testNamespace, Name: "dogu-op"}}

		// when
		_, err := sut.Reconcile(testCtx, req)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func Test_componentReconciler_validateName(t *testing.T) {
	tests := []struct {
		name         string
		recorderFunc func(t *testing.T) record.EventRecorder
		component    *k8sv1.Component
		wantSuccess  bool
	}{
		{
			name: "should fail validation",
			recorderFunc: func(t *testing.T) record.EventRecorder {
				recorder := newMockEventRecorder(t)
				recorder.EXPECT().Eventf(mock.Anything, "Warning", "FailedNameValidation", "Component resource does not follow naming rules: The component's metadata.name '%s' must be the same as its spec.name '%s'.", "example", "invalid-example")
				return recorder
			},
			component:   &k8sv1.Component{ObjectMeta: v1.ObjectMeta{Name: "example"}, Spec: k8sv1.ComponentSpec{Name: "invalid-example"}},
			wantSuccess: false,
		},
		{
			name:         "should succeed validation",
			recorderFunc: func(t *testing.T) record.EventRecorder { return newMockEventRecorder(t) },
			component:    &k8sv1.Component{ObjectMeta: v1.ObjectMeta{Name: "example"}, Spec: k8sv1.ComponentSpec{Name: "example"}},
			wantSuccess:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ComponentReconciler{recorder: tt.recorderFunc(t)}
			success := r.validateName(tt.component)
			assert.Equal(t, tt.wantSuccess, success)
		})
	}
}

func TestComponentReconciler_getComponentRequest(t *testing.T) {
	t.Run("should fail to get component list", func(t *testing.T) {
		// given
		cm := &corev1.ConfigMap{}
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		componentV1InterfaceMock := newMockComponentV1Alpha1Interface(t)
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().List(testCtx, v1.ListOptions{}).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{*component}}, assert.AnError)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1InterfaceMock)
		componentV1InterfaceMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)

		sut := ComponentReconciler{
			namespace: testNamespace,
			clientSet: clientSetMock,
		}

		// when
		requests := sut.getComponentRequest(testCtx, cm)

		// then
		assert.Nil(t, requests)
	})
	t.Run("should get empty component list", func(t *testing.T) {
		// given
		cm := &corev1.ConfigMap{}
		componentV1InterfaceMock := newMockComponentV1Alpha1Interface(t)
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().List(testCtx, v1.ListOptions{}).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{}}, nil)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1InterfaceMock)
		componentV1InterfaceMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)

		sut := ComponentReconciler{
			namespace: testNamespace,
			clientSet: clientSetMock,
		}

		// when
		requests := sut.getComponentRequest(testCtx, cm)

		// then
		assert.Empty(t, requests)
	})
	t.Run("should get component list with no config map reference", func(t *testing.T) {
		// given
		cm := &corev1.ConfigMap{}
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		componentV1InterfaceMock := newMockComponentV1Alpha1Interface(t)
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().List(testCtx, v1.ListOptions{}).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{*component}}, nil)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1InterfaceMock)
		componentV1InterfaceMock.EXPECT().Components(testNamespace).Return(componentInterfaceMock)

		sut := ComponentReconciler{
			namespace: testNamespace,
			clientSet: clientSetMock,
		}

		// when
		requests := sut.getComponentRequest(testCtx, cm)

		// then
		assert.Empty(t, requests)
	})
	t.Run("should get component list with config map reference", func(t *testing.T) {
		// given
		cm := &corev1.ConfigMap{ObjectMeta: v1.ObjectMeta{Name: "configmap"}}
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		component.Spec.ValuesConfigRef = &k8sv1.Reference{
			Name: "configmap",
			Key:  "key",
		}
		componentV1InterfaceMock := newMockComponentV1Alpha1Interface(t)
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().List(testCtx, v1.ListOptions{}).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{*component}}, nil)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1InterfaceMock)
		componentV1InterfaceMock.EXPECT().Components("ecosystem").Return(componentInterfaceMock)

		sut := ComponentReconciler{
			namespace: "ecosystem",
			clientSet: clientSetMock,
		}

		// when
		requests := sut.getComponentRequest(testCtx, cm)

		// then
		assert.NotEmpty(t, requests)
		assert.Equal(t, reconcile.Request{NamespacedName: types.NamespacedName{Name: "dogu-op", Namespace: "ecosystem"}}, requests[0])
	})
	t.Run("should get component list with default config config map", func(t *testing.T) {
		// given
		cm := &corev1.ConfigMap{ObjectMeta: v1.ObjectMeta{
			Name: "dogu-op-config",
			Labels: map[string]string{
				"k8s.cloudogu.com/component.config": "dogu-op",
			},
		}}
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		componentV1InterfaceMock := newMockComponentV1Alpha1Interface(t)
		componentInterfaceMock := newMockComponentInterface(t)
		componentInterfaceMock.EXPECT().List(testCtx, v1.ListOptions{}).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{*component}}, nil)
		clientSetMock := newMockComponentEcosystemInterface(t)
		clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1InterfaceMock)
		componentV1InterfaceMock.EXPECT().Components("ecosystem").Return(componentInterfaceMock)

		sut := ComponentReconciler{
			namespace: "ecosystem",
			clientSet: clientSetMock,
		}

		// when
		requests := sut.getComponentRequest(testCtx, cm)

		// then
		assert.NotEmpty(t, requests)
		assert.Equal(t, reconcile.Request{NamespacedName: types.NamespacedName{Name: "dogu-op", Namespace: "ecosystem"}}, requests[0])
	})
}
