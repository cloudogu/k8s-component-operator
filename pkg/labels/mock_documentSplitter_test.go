// Code generated by mockery v2.42.1. DO NOT EDIT.

package labels

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
	runtime "k8s.io/apimachinery/pkg/runtime"

	yaml "github.com/cloudogu/k8s-component-operator/pkg/yaml"
)

// mockDocumentSplitter is an autogenerated mock type for the documentSplitter type
type mockDocumentSplitter struct {
	mock.Mock
}

type mockDocumentSplitter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDocumentSplitter) EXPECT() *mockDocumentSplitter_Expecter {
	return &mockDocumentSplitter_Expecter{mock: &_m.Mock}
}

// Bytes provides a mock function with given fields:
func (_m *mockDocumentSplitter) Bytes() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Bytes")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// mockDocumentSplitter_Bytes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Bytes'
type mockDocumentSplitter_Bytes_Call struct {
	*mock.Call
}

// Bytes is a helper method to define mock.On call
func (_e *mockDocumentSplitter_Expecter) Bytes() *mockDocumentSplitter_Bytes_Call {
	return &mockDocumentSplitter_Bytes_Call{Call: _e.mock.On("Bytes")}
}

func (_c *mockDocumentSplitter_Bytes_Call) Run(run func()) *mockDocumentSplitter_Bytes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDocumentSplitter_Bytes_Call) Return(_a0 []byte) *mockDocumentSplitter_Bytes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDocumentSplitter_Bytes_Call) RunAndReturn(run func() []byte) *mockDocumentSplitter_Bytes_Call {
	_c.Call.Return(run)
	return _c
}

// Err provides a mock function with given fields:
func (_m *mockDocumentSplitter) Err() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Err")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDocumentSplitter_Err_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Err'
type mockDocumentSplitter_Err_Call struct {
	*mock.Call
}

// Err is a helper method to define mock.On call
func (_e *mockDocumentSplitter_Expecter) Err() *mockDocumentSplitter_Err_Call {
	return &mockDocumentSplitter_Err_Call{Call: _e.mock.On("Err")}
}

func (_c *mockDocumentSplitter_Err_Call) Run(run func()) *mockDocumentSplitter_Err_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDocumentSplitter_Err_Call) Return(_a0 error) *mockDocumentSplitter_Err_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDocumentSplitter_Err_Call) RunAndReturn(run func() error) *mockDocumentSplitter_Err_Call {
	_c.Call.Return(run)
	return _c
}

// Next provides a mock function with given fields:
func (_m *mockDocumentSplitter) Next() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Next")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// mockDocumentSplitter_Next_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Next'
type mockDocumentSplitter_Next_Call struct {
	*mock.Call
}

// Next is a helper method to define mock.On call
func (_e *mockDocumentSplitter_Expecter) Next() *mockDocumentSplitter_Next_Call {
	return &mockDocumentSplitter_Next_Call{Call: _e.mock.On("Next")}
}

func (_c *mockDocumentSplitter_Next_Call) Run(run func()) *mockDocumentSplitter_Next_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDocumentSplitter_Next_Call) Return(_a0 bool) *mockDocumentSplitter_Next_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDocumentSplitter_Next_Call) RunAndReturn(run func() bool) *mockDocumentSplitter_Next_Call {
	_c.Call.Return(run)
	return _c
}

// Object provides a mock function with given fields:
func (_m *mockDocumentSplitter) Object() runtime.Object {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Object")
	}

	var r0 runtime.Object
	if rf, ok := ret.Get(0).(func() runtime.Object); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(runtime.Object)
		}
	}

	return r0
}

// mockDocumentSplitter_Object_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Object'
type mockDocumentSplitter_Object_Call struct {
	*mock.Call
}

// Object is a helper method to define mock.On call
func (_e *mockDocumentSplitter_Expecter) Object() *mockDocumentSplitter_Object_Call {
	return &mockDocumentSplitter_Object_Call{Call: _e.mock.On("Object")}
}

func (_c *mockDocumentSplitter_Object_Call) Run(run func()) *mockDocumentSplitter_Object_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDocumentSplitter_Object_Call) Return(_a0 runtime.Object) *mockDocumentSplitter_Object_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDocumentSplitter_Object_Call) RunAndReturn(run func() runtime.Object) *mockDocumentSplitter_Object_Call {
	_c.Call.Return(run)
	return _c
}

// WithReader provides a mock function with given fields: r
func (_m *mockDocumentSplitter) WithReader(r io.Reader) yaml.DocumentSplitter {
	ret := _m.Called(r)

	if len(ret) == 0 {
		panic("no return value specified for WithReader")
	}

	var r0 yaml.DocumentSplitter
	if rf, ok := ret.Get(0).(func(io.Reader) yaml.DocumentSplitter); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(yaml.DocumentSplitter)
		}
	}

	return r0
}

// mockDocumentSplitter_WithReader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithReader'
type mockDocumentSplitter_WithReader_Call struct {
	*mock.Call
}

// WithReader is a helper method to define mock.On call
//   - r io.Reader
func (_e *mockDocumentSplitter_Expecter) WithReader(r interface{}) *mockDocumentSplitter_WithReader_Call {
	return &mockDocumentSplitter_WithReader_Call{Call: _e.mock.On("WithReader", r)}
}

func (_c *mockDocumentSplitter_WithReader_Call) Run(run func(r io.Reader)) *mockDocumentSplitter_WithReader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(io.Reader))
	})
	return _c
}

func (_c *mockDocumentSplitter_WithReader_Call) Return(_a0 yaml.DocumentSplitter) *mockDocumentSplitter_WithReader_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDocumentSplitter_WithReader_Call) RunAndReturn(run func(io.Reader) yaml.DocumentSplitter) *mockDocumentSplitter_WithReader_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDocumentSplitter creates a new instance of mockDocumentSplitter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDocumentSplitter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDocumentSplitter {
	mock := &mockDocumentSplitter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
