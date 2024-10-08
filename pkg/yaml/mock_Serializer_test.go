// Code generated by mockery v2.42.1. DO NOT EDIT.

package yaml

import (
	mock "github.com/stretchr/testify/mock"
	sigs_k8s_ioyaml "sigs.k8s.io/yaml"
)

// MockSerializer is an autogenerated mock type for the Serializer type
type MockSerializer struct {
	mock.Mock
}

type MockSerializer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSerializer) EXPECT() *MockSerializer_Expecter {
	return &MockSerializer_Expecter{mock: &_m.Mock}
}

// Marshal provides a mock function with given fields: o
func (_m *MockSerializer) Marshal(o interface{}) ([]byte, error) {
	ret := _m.Called(o)

	if len(ret) == 0 {
		panic("no return value specified for Marshal")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}) ([]byte, error)); ok {
		return rf(o)
	}
	if rf, ok := ret.Get(0).(func(interface{}) []byte); ok {
		r0 = rf(o)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(o)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSerializer_Marshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Marshal'
type MockSerializer_Marshal_Call struct {
	*mock.Call
}

// Marshal is a helper method to define mock.On call
//   - o interface{}
func (_e *MockSerializer_Expecter) Marshal(o interface{}) *MockSerializer_Marshal_Call {
	return &MockSerializer_Marshal_Call{Call: _e.mock.On("Marshal", o)}
}

func (_c *MockSerializer_Marshal_Call) Run(run func(o interface{})) *MockSerializer_Marshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockSerializer_Marshal_Call) Return(_a0 []byte, _a1 error) *MockSerializer_Marshal_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSerializer_Marshal_Call) RunAndReturn(run func(interface{}) ([]byte, error)) *MockSerializer_Marshal_Call {
	_c.Call.Return(run)
	return _c
}

// Unmarshal provides a mock function with given fields: y, o, opts
func (_m *MockSerializer) Unmarshal(y []byte, o interface{}, opts ...sigs_k8s_ioyaml.JSONOpt) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, y, o)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Unmarshal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, interface{}, ...sigs_k8s_ioyaml.JSONOpt) error); ok {
		r0 = rf(y, o, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSerializer_Unmarshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unmarshal'
type MockSerializer_Unmarshal_Call struct {
	*mock.Call
}

// Unmarshal is a helper method to define mock.On call
//   - y []byte
//   - o interface{}
//   - opts ...sigs_k8s_ioyaml.JSONOpt
func (_e *MockSerializer_Expecter) Unmarshal(y interface{}, o interface{}, opts ...interface{}) *MockSerializer_Unmarshal_Call {
	return &MockSerializer_Unmarshal_Call{Call: _e.mock.On("Unmarshal",
		append([]interface{}{y, o}, opts...)...)}
}

func (_c *MockSerializer_Unmarshal_Call) Run(run func(y []byte, o interface{}, opts ...sigs_k8s_ioyaml.JSONOpt)) *MockSerializer_Unmarshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]sigs_k8s_ioyaml.JSONOpt, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(sigs_k8s_ioyaml.JSONOpt)
			}
		}
		run(args[0].([]byte), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockSerializer_Unmarshal_Call) Return(_a0 error) *MockSerializer_Unmarshal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSerializer_Unmarshal_Call) RunAndReturn(run func([]byte, interface{}, ...sigs_k8s_ioyaml.JSONOpt) error) *MockSerializer_Unmarshal_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSerializer creates a new instance of MockSerializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSerializer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSerializer {
	mock := &MockSerializer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
