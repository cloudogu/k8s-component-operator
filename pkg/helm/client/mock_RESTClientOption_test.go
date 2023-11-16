// Code generated by mockery v2.20.0. DO NOT EDIT.

package client

import (
	mock "github.com/stretchr/testify/mock"
	rest "k8s.io/client-go/rest"
)

// MockRESTClientOption is an autogenerated mock type for the RESTClientOption type
type MockRESTClientOption struct {
	mock.Mock
}

type MockRESTClientOption_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRESTClientOption) EXPECT() *MockRESTClientOption_Expecter {
	return &MockRESTClientOption_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockRESTClientOption) Execute(_a0 *rest.Config) {
	_m.Called(_a0)
}

// MockRESTClientOption_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockRESTClientOption_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 *rest.Config
func (_e *MockRESTClientOption_Expecter) Execute(_a0 interface{}) *MockRESTClientOption_Execute_Call {
	return &MockRESTClientOption_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockRESTClientOption_Execute_Call) Run(run func(_a0 *rest.Config)) *MockRESTClientOption_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*rest.Config))
	})
	return _c
}

func (_c *MockRESTClientOption_Execute_Call) Return() *MockRESTClientOption_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockRESTClientOption_Execute_Call) RunAndReturn(run func(*rest.Config)) *MockRESTClientOption_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockRESTClientOption interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRESTClientOption creates a new instance of MockRESTClientOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRESTClientOption(t mockConstructorTestingTNewMockRESTClientOption) *MockRESTClientOption {
	mock := &MockRESTClientOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
