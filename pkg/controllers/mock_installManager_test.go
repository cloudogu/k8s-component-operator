// Code generated by mockery v2.20.0. DO NOT EDIT.

package controllers

import (
	context "context"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	mock "github.com/stretchr/testify/mock"
)

// mockInstallManager is an autogenerated mock type for the installManager type
type mockInstallManager struct {
	mock.Mock
}

type mockInstallManager_Expecter struct {
	mock *mock.Mock
}

func (_m *mockInstallManager) EXPECT() *mockInstallManager_Expecter {
	return &mockInstallManager_Expecter{mock: &_m.Mock}
}

// Install provides a mock function with given fields: ctx, component
func (_m *mockInstallManager) Install(ctx context.Context, component *v1.Component) error {
	ret := _m.Called(ctx, component)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) error); ok {
		r0 = rf(ctx, component)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockInstallManager_Install_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Install'
type mockInstallManager_Install_Call struct {
	*mock.Call
}

// Install is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockInstallManager_Expecter) Install(ctx interface{}, component interface{}) *mockInstallManager_Install_Call {
	return &mockInstallManager_Install_Call{Call: _e.mock.On("Install", ctx, component)}
}

func (_c *mockInstallManager_Install_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockInstallManager_Install_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockInstallManager_Install_Call) Return(_a0 error) *mockInstallManager_Install_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockInstallManager_Install_Call) RunAndReturn(run func(context.Context, *v1.Component) error) *mockInstallManager_Install_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockInstallManager interface {
	mock.TestingT
	Cleanup(func())
}

// newMockInstallManager creates a new instance of mockInstallManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockInstallManager(t mockConstructorTestingTnewMockInstallManager) *mockInstallManager {
	mock := &mockInstallManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
