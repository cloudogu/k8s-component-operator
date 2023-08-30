// Code generated by mockery v2.20.0. DO NOT EDIT.

package helm

import mock "github.com/stretchr/testify/mock"

// mockTagResolver is an autogenerated mock type for the tagResolver type
type mockTagResolver struct {
	mock.Mock
}

type mockTagResolver_Expecter struct {
	mock *mock.Mock
}

func (_m *mockTagResolver) EXPECT() *mockTagResolver_Expecter {
	return &mockTagResolver_Expecter{mock: &_m.Mock}
}

// Tags provides a mock function with given fields: ref
func (_m *mockTagResolver) Tags(ref string) ([]string, error) {
	ret := _m.Called(ref)

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

// mockTagResolver_Tags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Tags'
type mockTagResolver_Tags_Call struct {
	*mock.Call
}

// Tags is a helper method to define mock.On call
//   - ref string
func (_e *mockTagResolver_Expecter) Tags(ref interface{}) *mockTagResolver_Tags_Call {
	return &mockTagResolver_Tags_Call{Call: _e.mock.On("Tags", ref)}
}

func (_c *mockTagResolver_Tags_Call) Run(run func(ref string)) *mockTagResolver_Tags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockTagResolver_Tags_Call) Return(_a0 []string, _a1 error) *mockTagResolver_Tags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockTagResolver_Tags_Call) RunAndReturn(run func(string) ([]string, error)) *mockTagResolver_Tags_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockTagResolver interface {
	mock.TestingT
	Cleanup(func())
}

// newMockTagResolver creates a new instance of mockTagResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockTagResolver(t mockConstructorTestingTnewMockTagResolver) *mockTagResolver {
	mock := &mockTagResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
