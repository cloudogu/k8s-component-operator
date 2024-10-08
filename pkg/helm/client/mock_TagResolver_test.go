// Code generated by mockery v2.42.1. DO NOT EDIT.

package client

import mock "github.com/stretchr/testify/mock"

// MockTagResolver is an autogenerated mock type for the TagResolver type
type MockTagResolver struct {
	mock.Mock
}

type MockTagResolver_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTagResolver) EXPECT() *MockTagResolver_Expecter {
	return &MockTagResolver_Expecter{mock: &_m.Mock}
}

// Tags provides a mock function with given fields: ref
func (_m *MockTagResolver) Tags(ref string) ([]string, error) {
	ret := _m.Called(ref)

	if len(ret) == 0 {
		panic("no return value specified for Tags")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]string, error)); ok {
		return rf(ref)
	}
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(ref)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockTagResolver_Tags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Tags'
type MockTagResolver_Tags_Call struct {
	*mock.Call
}

// Tags is a helper method to define mock.On call
//   - ref string
func (_e *MockTagResolver_Expecter) Tags(ref interface{}) *MockTagResolver_Tags_Call {
	return &MockTagResolver_Tags_Call{Call: _e.mock.On("Tags", ref)}
}

func (_c *MockTagResolver_Tags_Call) Run(run func(ref string)) *MockTagResolver_Tags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockTagResolver_Tags_Call) Return(_a0 []string, _a1 error) *MockTagResolver_Tags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockTagResolver_Tags_Call) RunAndReturn(run func(string) ([]string, error)) *MockTagResolver_Tags_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTagResolver creates a new instance of MockTagResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTagResolver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTagResolver {
	mock := &MockTagResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
