// Code generated by mockery v2.20.0. DO NOT EDIT.

package helm

import mock "github.com/stretchr/testify/mock"

// mockOciRepositoryConfig is an autogenerated mock type for the ociRepositoryConfig type
type mockOciRepositoryConfig struct {
	mock.Mock
}

type mockOciRepositoryConfig_Expecter struct {
	mock *mock.Mock
}

func (_m *mockOciRepositoryConfig) EXPECT() *mockOciRepositoryConfig_Expecter {
	return &mockOciRepositoryConfig_Expecter{mock: &_m.Mock}
}

// GetOciEndpoint provides a mock function with given fields:
func (_m *mockOciRepositoryConfig) GetOciEndpoint() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockOciRepositoryConfig_GetOciEndpoint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOciEndpoint'
type mockOciRepositoryConfig_GetOciEndpoint_Call struct {
	*mock.Call
}

// GetOciEndpoint is a helper method to define mock.On call
func (_e *mockOciRepositoryConfig_Expecter) GetOciEndpoint() *mockOciRepositoryConfig_GetOciEndpoint_Call {
	return &mockOciRepositoryConfig_GetOciEndpoint_Call{Call: _e.mock.On("GetOciEndpoint")}
}

func (_c *mockOciRepositoryConfig_GetOciEndpoint_Call) Run(run func()) *mockOciRepositoryConfig_GetOciEndpoint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockOciRepositoryConfig_GetOciEndpoint_Call) Return(_a0 string, _a1 error) *mockOciRepositoryConfig_GetOciEndpoint_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockOciRepositoryConfig_GetOciEndpoint_Call) RunAndReturn(run func() (string, error)) *mockOciRepositoryConfig_GetOciEndpoint_Call {
	_c.Call.Return(run)
	return _c
}

// IsPlainHttp provides a mock function with given fields:
func (_m *mockOciRepositoryConfig) IsPlainHttp() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// mockOciRepositoryConfig_IsPlainHttp_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsPlainHttp'
type mockOciRepositoryConfig_IsPlainHttp_Call struct {
	*mock.Call
}

// IsPlainHttp is a helper method to define mock.On call
func (_e *mockOciRepositoryConfig_Expecter) IsPlainHttp() *mockOciRepositoryConfig_IsPlainHttp_Call {
	return &mockOciRepositoryConfig_IsPlainHttp_Call{Call: _e.mock.On("IsPlainHttp")}
}

func (_c *mockOciRepositoryConfig_IsPlainHttp_Call) Run(run func()) *mockOciRepositoryConfig_IsPlainHttp_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockOciRepositoryConfig_IsPlainHttp_Call) Return(_a0 bool) *mockOciRepositoryConfig_IsPlainHttp_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockOciRepositoryConfig_IsPlainHttp_Call) RunAndReturn(run func() bool) *mockOciRepositoryConfig_IsPlainHttp_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockOciRepositoryConfig interface {
	mock.TestingT
	Cleanup(func())
}

// newMockOciRepositoryConfig creates a new instance of mockOciRepositoryConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockOciRepositoryConfig(t mockConstructorTestingTnewMockOciRepositoryConfig) *mockOciRepositoryConfig {
	mock := &mockOciRepositoryConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
