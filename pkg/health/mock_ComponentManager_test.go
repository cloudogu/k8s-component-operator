// Code generated by mockery v2.42.1. DO NOT EDIT.

package health

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockComponentManager is an autogenerated mock type for the ComponentManager type
type MockComponentManager struct {
	mock.Mock
}

type MockComponentManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockComponentManager) EXPECT() *MockComponentManager_Expecter {
	return &MockComponentManager_Expecter{mock: &_m.Mock}
}

// UpdateComponentHealth provides a mock function with given fields: ctx, componentName, namespace
func (_m *MockComponentManager) UpdateComponentHealth(ctx context.Context, componentName string, namespace string) error {
	ret := _m.Called(ctx, componentName, namespace)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComponentHealth")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, componentName, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockComponentManager_UpdateComponentHealth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComponentHealth'
type MockComponentManager_UpdateComponentHealth_Call struct {
	*mock.Call
}

// UpdateComponentHealth is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName string
//   - namespace string
func (_e *MockComponentManager_Expecter) UpdateComponentHealth(ctx interface{}, componentName interface{}, namespace interface{}) *MockComponentManager_UpdateComponentHealth_Call {
	return &MockComponentManager_UpdateComponentHealth_Call{Call: _e.mock.On("UpdateComponentHealth", ctx, componentName, namespace)}
}

func (_c *MockComponentManager_UpdateComponentHealth_Call) Run(run func(ctx context.Context, componentName string, namespace string)) *MockComponentManager_UpdateComponentHealth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockComponentManager_UpdateComponentHealth_Call) Return(_a0 error) *MockComponentManager_UpdateComponentHealth_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockComponentManager_UpdateComponentHealth_Call) RunAndReturn(run func(context.Context, string, string) error) *MockComponentManager_UpdateComponentHealth_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateComponentHealthAll provides a mock function with given fields: ctx
func (_m *MockComponentManager) UpdateComponentHealthAll(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComponentHealthAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockComponentManager_UpdateComponentHealthAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComponentHealthAll'
type MockComponentManager_UpdateComponentHealthAll_Call struct {
	*mock.Call
}

// UpdateComponentHealthAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockComponentManager_Expecter) UpdateComponentHealthAll(ctx interface{}) *MockComponentManager_UpdateComponentHealthAll_Call {
	return &MockComponentManager_UpdateComponentHealthAll_Call{Call: _e.mock.On("UpdateComponentHealthAll", ctx)}
}

func (_c *MockComponentManager_UpdateComponentHealthAll_Call) Run(run func(ctx context.Context)) *MockComponentManager_UpdateComponentHealthAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockComponentManager_UpdateComponentHealthAll_Call) Return(_a0 error) *MockComponentManager_UpdateComponentHealthAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockComponentManager_UpdateComponentHealthAll_Call) RunAndReturn(run func(context.Context) error) *MockComponentManager_UpdateComponentHealthAll_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateComponentHealthWithInstalledVersion provides a mock function with given fields: ctx, componentName, namespace, version
func (_m *MockComponentManager) UpdateComponentHealthWithInstalledVersion(ctx context.Context, componentName string, namespace string, version string) error {
	ret := _m.Called(ctx, componentName, namespace, version)

	if len(ret) == 0 {
		panic("no return value specified for UpdateComponentHealthWithInstalledVersion")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, componentName, namespace, version)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateComponentHealthWithInstalledVersion'
type MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call struct {
	*mock.Call
}

// UpdateComponentHealthWithInstalledVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName string
//   - namespace string
//   - version string
func (_e *MockComponentManager_Expecter) UpdateComponentHealthWithInstalledVersion(ctx interface{}, componentName interface{}, namespace interface{}, version interface{}) *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call {
	return &MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call{Call: _e.mock.On("UpdateComponentHealthWithInstalledVersion", ctx, componentName, namespace, version)}
}

func (_c *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call) Run(run func(ctx context.Context, componentName string, namespace string, version string)) *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call) Return(_a0 error) *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call) RunAndReturn(run func(context.Context, string, string, string) error) *MockComponentManager_UpdateComponentHealthWithInstalledVersion_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockComponentManager creates a new instance of MockComponentManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockComponentManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockComponentManager {
	mock := &MockComponentManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
