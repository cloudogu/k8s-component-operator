// Code generated by mockery v2.20.0. DO NOT EDIT.

package health

import (
	cache "sigs.k8s.io/controller-runtime/pkg/cache"
	client "sigs.k8s.io/controller-runtime/pkg/client"

	config "sigs.k8s.io/controller-runtime/pkg/config"

	context "context"

	healthz "sigs.k8s.io/controller-runtime/pkg/healthz"

	http "net/http"

	logr "github.com/go-logr/logr"

	manager "sigs.k8s.io/controller-runtime/pkg/manager"

	meta "k8s.io/apimachinery/pkg/api/meta"

	mock "github.com/stretchr/testify/mock"

	record "k8s.io/client-go/tools/record"

	rest "k8s.io/client-go/rest"

	runtime "k8s.io/apimachinery/pkg/runtime"

	webhook "sigs.k8s.io/controller-runtime/pkg/webhook"
)

// mockControllerManager is an autogenerated mock type for the controllerManager type
type mockControllerManager struct {
	mock.Mock
}

type mockControllerManager_Expecter struct {
	mock *mock.Mock
}

func (_m *mockControllerManager) EXPECT() *mockControllerManager_Expecter {
	return &mockControllerManager_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: _a0
func (_m *mockControllerManager) Add(_a0 manager.Runnable) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(manager.Runnable) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockControllerManager_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type mockControllerManager_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - _a0 manager.Runnable
func (_e *mockControllerManager_Expecter) Add(_a0 interface{}) *mockControllerManager_Add_Call {
	return &mockControllerManager_Add_Call{Call: _e.mock.On("Add", _a0)}
}

func (_c *mockControllerManager_Add_Call) Run(run func(_a0 manager.Runnable)) *mockControllerManager_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(manager.Runnable))
	})
	return _c
}

func (_c *mockControllerManager_Add_Call) Return(_a0 error) *mockControllerManager_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_Add_Call) RunAndReturn(run func(manager.Runnable) error) *mockControllerManager_Add_Call {
	_c.Call.Return(run)
	return _c
}

// AddHealthzCheck provides a mock function with given fields: name, check
func (_m *mockControllerManager) AddHealthzCheck(name string, check healthz.Checker) error {
	ret := _m.Called(name, check)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, healthz.Checker) error); ok {
		r0 = rf(name, check)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockControllerManager_AddHealthzCheck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddHealthzCheck'
type mockControllerManager_AddHealthzCheck_Call struct {
	*mock.Call
}

// AddHealthzCheck is a helper method to define mock.On call
//   - name string
//   - check healthz.Checker
func (_e *mockControllerManager_Expecter) AddHealthzCheck(name interface{}, check interface{}) *mockControllerManager_AddHealthzCheck_Call {
	return &mockControllerManager_AddHealthzCheck_Call{Call: _e.mock.On("AddHealthzCheck", name, check)}
}

func (_c *mockControllerManager_AddHealthzCheck_Call) Run(run func(name string, check healthz.Checker)) *mockControllerManager_AddHealthzCheck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(healthz.Checker))
	})
	return _c
}

func (_c *mockControllerManager_AddHealthzCheck_Call) Return(_a0 error) *mockControllerManager_AddHealthzCheck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_AddHealthzCheck_Call) RunAndReturn(run func(string, healthz.Checker) error) *mockControllerManager_AddHealthzCheck_Call {
	_c.Call.Return(run)
	return _c
}

// AddReadyzCheck provides a mock function with given fields: name, check
func (_m *mockControllerManager) AddReadyzCheck(name string, check healthz.Checker) error {
	ret := _m.Called(name, check)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, healthz.Checker) error); ok {
		r0 = rf(name, check)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockControllerManager_AddReadyzCheck_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddReadyzCheck'
type mockControllerManager_AddReadyzCheck_Call struct {
	*mock.Call
}

// AddReadyzCheck is a helper method to define mock.On call
//   - name string
//   - check healthz.Checker
func (_e *mockControllerManager_Expecter) AddReadyzCheck(name interface{}, check interface{}) *mockControllerManager_AddReadyzCheck_Call {
	return &mockControllerManager_AddReadyzCheck_Call{Call: _e.mock.On("AddReadyzCheck", name, check)}
}

func (_c *mockControllerManager_AddReadyzCheck_Call) Run(run func(name string, check healthz.Checker)) *mockControllerManager_AddReadyzCheck_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(healthz.Checker))
	})
	return _c
}

func (_c *mockControllerManager_AddReadyzCheck_Call) Return(_a0 error) *mockControllerManager_AddReadyzCheck_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_AddReadyzCheck_Call) RunAndReturn(run func(string, healthz.Checker) error) *mockControllerManager_AddReadyzCheck_Call {
	_c.Call.Return(run)
	return _c
}

// Elected provides a mock function with given fields:
func (_m *mockControllerManager) Elected() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// mockControllerManager_Elected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Elected'
type mockControllerManager_Elected_Call struct {
	*mock.Call
}

// Elected is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) Elected() *mockControllerManager_Elected_Call {
	return &mockControllerManager_Elected_Call{Call: _e.mock.On("Elected")}
}

func (_c *mockControllerManager_Elected_Call) Run(run func()) *mockControllerManager_Elected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_Elected_Call) Return(_a0 <-chan struct{}) *mockControllerManager_Elected_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_Elected_Call) RunAndReturn(run func() <-chan struct{}) *mockControllerManager_Elected_Call {
	_c.Call.Return(run)
	return _c
}

// GetAPIReader provides a mock function with given fields:
func (_m *mockControllerManager) GetAPIReader() client.Reader {
	ret := _m.Called()

	var r0 client.Reader
	if rf, ok := ret.Get(0).(func() client.Reader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Reader)
		}
	}

	return r0
}

// mockControllerManager_GetAPIReader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAPIReader'
type mockControllerManager_GetAPIReader_Call struct {
	*mock.Call
}

// GetAPIReader is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetAPIReader() *mockControllerManager_GetAPIReader_Call {
	return &mockControllerManager_GetAPIReader_Call{Call: _e.mock.On("GetAPIReader")}
}

func (_c *mockControllerManager_GetAPIReader_Call) Run(run func()) *mockControllerManager_GetAPIReader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetAPIReader_Call) Return(_a0 client.Reader) *mockControllerManager_GetAPIReader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetAPIReader_Call) RunAndReturn(run func() client.Reader) *mockControllerManager_GetAPIReader_Call {
	_c.Call.Return(run)
	return _c
}

// GetCache provides a mock function with given fields:
func (_m *mockControllerManager) GetCache() cache.Cache {
	ret := _m.Called()

	var r0 cache.Cache
	if rf, ok := ret.Get(0).(func() cache.Cache); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(cache.Cache)
		}
	}

	return r0
}

// mockControllerManager_GetCache_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCache'
type mockControllerManager_GetCache_Call struct {
	*mock.Call
}

// GetCache is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetCache() *mockControllerManager_GetCache_Call {
	return &mockControllerManager_GetCache_Call{Call: _e.mock.On("GetCache")}
}

func (_c *mockControllerManager_GetCache_Call) Run(run func()) *mockControllerManager_GetCache_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetCache_Call) Return(_a0 cache.Cache) *mockControllerManager_GetCache_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetCache_Call) RunAndReturn(run func() cache.Cache) *mockControllerManager_GetCache_Call {
	_c.Call.Return(run)
	return _c
}

// GetClient provides a mock function with given fields:
func (_m *mockControllerManager) GetClient() client.Client {
	ret := _m.Called()

	var r0 client.Client
	if rf, ok := ret.Get(0).(func() client.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.Client)
		}
	}

	return r0
}

// mockControllerManager_GetClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetClient'
type mockControllerManager_GetClient_Call struct {
	*mock.Call
}

// GetClient is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetClient() *mockControllerManager_GetClient_Call {
	return &mockControllerManager_GetClient_Call{Call: _e.mock.On("GetClient")}
}

func (_c *mockControllerManager_GetClient_Call) Run(run func()) *mockControllerManager_GetClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetClient_Call) Return(_a0 client.Client) *mockControllerManager_GetClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetClient_Call) RunAndReturn(run func() client.Client) *mockControllerManager_GetClient_Call {
	_c.Call.Return(run)
	return _c
}

// GetConfig provides a mock function with given fields:
func (_m *mockControllerManager) GetConfig() *rest.Config {
	ret := _m.Called()

	var r0 *rest.Config
	if rf, ok := ret.Get(0).(func() *rest.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Config)
		}
	}

	return r0
}

// mockControllerManager_GetConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConfig'
type mockControllerManager_GetConfig_Call struct {
	*mock.Call
}

// GetConfig is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetConfig() *mockControllerManager_GetConfig_Call {
	return &mockControllerManager_GetConfig_Call{Call: _e.mock.On("GetConfig")}
}

func (_c *mockControllerManager_GetConfig_Call) Run(run func()) *mockControllerManager_GetConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetConfig_Call) Return(_a0 *rest.Config) *mockControllerManager_GetConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetConfig_Call) RunAndReturn(run func() *rest.Config) *mockControllerManager_GetConfig_Call {
	_c.Call.Return(run)
	return _c
}

// GetControllerOptions provides a mock function with given fields:
func (_m *mockControllerManager) GetControllerOptions() config.Controller {
	ret := _m.Called()

	var r0 config.Controller
	if rf, ok := ret.Get(0).(func() config.Controller); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(config.Controller)
	}

	return r0
}

// mockControllerManager_GetControllerOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetControllerOptions'
type mockControllerManager_GetControllerOptions_Call struct {
	*mock.Call
}

// GetControllerOptions is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetControllerOptions() *mockControllerManager_GetControllerOptions_Call {
	return &mockControllerManager_GetControllerOptions_Call{Call: _e.mock.On("GetControllerOptions")}
}

func (_c *mockControllerManager_GetControllerOptions_Call) Run(run func()) *mockControllerManager_GetControllerOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetControllerOptions_Call) Return(_a0 config.Controller) *mockControllerManager_GetControllerOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetControllerOptions_Call) RunAndReturn(run func() config.Controller) *mockControllerManager_GetControllerOptions_Call {
	_c.Call.Return(run)
	return _c
}

// GetEventRecorderFor provides a mock function with given fields: name
func (_m *mockControllerManager) GetEventRecorderFor(name string) record.EventRecorder {
	ret := _m.Called(name)

	var r0 record.EventRecorder
	if rf, ok := ret.Get(0).(func(string) record.EventRecorder); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(record.EventRecorder)
		}
	}

	return r0
}

// mockControllerManager_GetEventRecorderFor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEventRecorderFor'
type mockControllerManager_GetEventRecorderFor_Call struct {
	*mock.Call
}

// GetEventRecorderFor is a helper method to define mock.On call
//   - name string
func (_e *mockControllerManager_Expecter) GetEventRecorderFor(name interface{}) *mockControllerManager_GetEventRecorderFor_Call {
	return &mockControllerManager_GetEventRecorderFor_Call{Call: _e.mock.On("GetEventRecorderFor", name)}
}

func (_c *mockControllerManager_GetEventRecorderFor_Call) Run(run func(name string)) *mockControllerManager_GetEventRecorderFor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockControllerManager_GetEventRecorderFor_Call) Return(_a0 record.EventRecorder) *mockControllerManager_GetEventRecorderFor_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetEventRecorderFor_Call) RunAndReturn(run func(string) record.EventRecorder) *mockControllerManager_GetEventRecorderFor_Call {
	_c.Call.Return(run)
	return _c
}

// GetFieldIndexer provides a mock function with given fields:
func (_m *mockControllerManager) GetFieldIndexer() client.FieldIndexer {
	ret := _m.Called()

	var r0 client.FieldIndexer
	if rf, ok := ret.Get(0).(func() client.FieldIndexer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(client.FieldIndexer)
		}
	}

	return r0
}

// mockControllerManager_GetFieldIndexer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFieldIndexer'
type mockControllerManager_GetFieldIndexer_Call struct {
	*mock.Call
}

// GetFieldIndexer is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetFieldIndexer() *mockControllerManager_GetFieldIndexer_Call {
	return &mockControllerManager_GetFieldIndexer_Call{Call: _e.mock.On("GetFieldIndexer")}
}

func (_c *mockControllerManager_GetFieldIndexer_Call) Run(run func()) *mockControllerManager_GetFieldIndexer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetFieldIndexer_Call) Return(_a0 client.FieldIndexer) *mockControllerManager_GetFieldIndexer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetFieldIndexer_Call) RunAndReturn(run func() client.FieldIndexer) *mockControllerManager_GetFieldIndexer_Call {
	_c.Call.Return(run)
	return _c
}

// GetHTTPClient provides a mock function with given fields:
func (_m *mockControllerManager) GetHTTPClient() *http.Client {
	ret := _m.Called()

	var r0 *http.Client
	if rf, ok := ret.Get(0).(func() *http.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Client)
		}
	}

	return r0
}

// mockControllerManager_GetHTTPClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHTTPClient'
type mockControllerManager_GetHTTPClient_Call struct {
	*mock.Call
}

// GetHTTPClient is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetHTTPClient() *mockControllerManager_GetHTTPClient_Call {
	return &mockControllerManager_GetHTTPClient_Call{Call: _e.mock.On("GetHTTPClient")}
}

func (_c *mockControllerManager_GetHTTPClient_Call) Run(run func()) *mockControllerManager_GetHTTPClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetHTTPClient_Call) Return(_a0 *http.Client) *mockControllerManager_GetHTTPClient_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetHTTPClient_Call) RunAndReturn(run func() *http.Client) *mockControllerManager_GetHTTPClient_Call {
	_c.Call.Return(run)
	return _c
}

// GetLogger provides a mock function with given fields:
func (_m *mockControllerManager) GetLogger() logr.Logger {
	ret := _m.Called()

	var r0 logr.Logger
	if rf, ok := ret.Get(0).(func() logr.Logger); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(logr.Logger)
	}

	return r0
}

// mockControllerManager_GetLogger_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLogger'
type mockControllerManager_GetLogger_Call struct {
	*mock.Call
}

// GetLogger is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetLogger() *mockControllerManager_GetLogger_Call {
	return &mockControllerManager_GetLogger_Call{Call: _e.mock.On("GetLogger")}
}

func (_c *mockControllerManager_GetLogger_Call) Run(run func()) *mockControllerManager_GetLogger_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetLogger_Call) Return(_a0 logr.Logger) *mockControllerManager_GetLogger_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetLogger_Call) RunAndReturn(run func() logr.Logger) *mockControllerManager_GetLogger_Call {
	_c.Call.Return(run)
	return _c
}

// GetRESTMapper provides a mock function with given fields:
func (_m *mockControllerManager) GetRESTMapper() meta.RESTMapper {
	ret := _m.Called()

	var r0 meta.RESTMapper
	if rf, ok := ret.Get(0).(func() meta.RESTMapper); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(meta.RESTMapper)
		}
	}

	return r0
}

// mockControllerManager_GetRESTMapper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRESTMapper'
type mockControllerManager_GetRESTMapper_Call struct {
	*mock.Call
}

// GetRESTMapper is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetRESTMapper() *mockControllerManager_GetRESTMapper_Call {
	return &mockControllerManager_GetRESTMapper_Call{Call: _e.mock.On("GetRESTMapper")}
}

func (_c *mockControllerManager_GetRESTMapper_Call) Run(run func()) *mockControllerManager_GetRESTMapper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetRESTMapper_Call) Return(_a0 meta.RESTMapper) *mockControllerManager_GetRESTMapper_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetRESTMapper_Call) RunAndReturn(run func() meta.RESTMapper) *mockControllerManager_GetRESTMapper_Call {
	_c.Call.Return(run)
	return _c
}

// GetScheme provides a mock function with given fields:
func (_m *mockControllerManager) GetScheme() *runtime.Scheme {
	ret := _m.Called()

	var r0 *runtime.Scheme
	if rf, ok := ret.Get(0).(func() *runtime.Scheme); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtime.Scheme)
		}
	}

	return r0
}

// mockControllerManager_GetScheme_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetScheme'
type mockControllerManager_GetScheme_Call struct {
	*mock.Call
}

// GetScheme is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetScheme() *mockControllerManager_GetScheme_Call {
	return &mockControllerManager_GetScheme_Call{Call: _e.mock.On("GetScheme")}
}

func (_c *mockControllerManager_GetScheme_Call) Run(run func()) *mockControllerManager_GetScheme_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetScheme_Call) Return(_a0 *runtime.Scheme) *mockControllerManager_GetScheme_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetScheme_Call) RunAndReturn(run func() *runtime.Scheme) *mockControllerManager_GetScheme_Call {
	_c.Call.Return(run)
	return _c
}

// GetWebhookServer provides a mock function with given fields:
func (_m *mockControllerManager) GetWebhookServer() webhook.Server {
	ret := _m.Called()

	var r0 webhook.Server
	if rf, ok := ret.Get(0).(func() webhook.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(webhook.Server)
		}
	}

	return r0
}

// mockControllerManager_GetWebhookServer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWebhookServer'
type mockControllerManager_GetWebhookServer_Call struct {
	*mock.Call
}

// GetWebhookServer is a helper method to define mock.On call
func (_e *mockControllerManager_Expecter) GetWebhookServer() *mockControllerManager_GetWebhookServer_Call {
	return &mockControllerManager_GetWebhookServer_Call{Call: _e.mock.On("GetWebhookServer")}
}

func (_c *mockControllerManager_GetWebhookServer_Call) Run(run func()) *mockControllerManager_GetWebhookServer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockControllerManager_GetWebhookServer_Call) Return(_a0 webhook.Server) *mockControllerManager_GetWebhookServer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_GetWebhookServer_Call) RunAndReturn(run func() webhook.Server) *mockControllerManager_GetWebhookServer_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: ctx
func (_m *mockControllerManager) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockControllerManager_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type mockControllerManager_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockControllerManager_Expecter) Start(ctx interface{}) *mockControllerManager_Start_Call {
	return &mockControllerManager_Start_Call{Call: _e.mock.On("Start", ctx)}
}

func (_c *mockControllerManager_Start_Call) Run(run func(ctx context.Context)) *mockControllerManager_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockControllerManager_Start_Call) Return(_a0 error) *mockControllerManager_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockControllerManager_Start_Call) RunAndReturn(run func(context.Context) error) *mockControllerManager_Start_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockControllerManager interface {
	mock.TestingT
	Cleanup(func())
}

// newMockControllerManager creates a new instance of mockControllerManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockControllerManager(t mockConstructorTestingTnewMockControllerManager) *mockControllerManager {
	mock := &mockControllerManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
