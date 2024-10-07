// Code generated by mockery v2.42.1. DO NOT EDIT.

package controllers

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	types "k8s.io/apimachinery/pkg/types"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// mockComponentInterface is an autogenerated mock type for the componentInterface type
type mockComponentInterface struct {
	mock.Mock
}

type mockComponentInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *mockComponentInterface) EXPECT() *mockComponentInterface_Expecter {
	return &mockComponentInterface_Expecter{mock: &_m.Mock}
}

// AddFinalizer provides a mock function with given fields: ctx, component, finalizer
func (_m *mockComponentInterface) AddFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	ret := _m.Called(ctx, component, finalizer)

	if len(ret) == 0 {
		panic("no return value specified for AddFinalizer")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, string) (*v1.Component, error)); ok {
		return rf(ctx, component, finalizer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, string) *v1.Component); ok {
		r0 = rf(ctx, component, finalizer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component, string) error); ok {
		r1 = rf(ctx, component, finalizer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_AddFinalizer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFinalizer'
type mockComponentInterface_AddFinalizer_Call struct {
	*mock.Call
}

// AddFinalizer is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - finalizer string
func (_e *mockComponentInterface_Expecter) AddFinalizer(ctx interface{}, component interface{}, finalizer interface{}) *mockComponentInterface_AddFinalizer_Call {
	return &mockComponentInterface_AddFinalizer_Call{Call: _e.mock.On("AddFinalizer", ctx, component, finalizer)}
}

func (_c *mockComponentInterface_AddFinalizer_Call) Run(run func(ctx context.Context, component *v1.Component, finalizer string)) *mockComponentInterface_AddFinalizer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(string))
	})
	return _c
}

func (_c *mockComponentInterface_AddFinalizer_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_AddFinalizer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_AddFinalizer_Call) RunAndReturn(run func(context.Context, *v1.Component, string) (*v1.Component, error)) *mockComponentInterface_AddFinalizer_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, component, opts
func (_m *mockComponentInterface) Create(ctx context.Context, component *v1.Component, opts metav1.CreateOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, component, opts)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, metav1.CreateOptions) (*v1.Component, error)); ok {
		return rf(ctx, component, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, metav1.CreateOptions) *v1.Component); ok {
		r0 = rf(ctx, component, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component, metav1.CreateOptions) error); ok {
		r1 = rf(ctx, component, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockComponentInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - opts metav1.CreateOptions
func (_e *mockComponentInterface_Expecter) Create(ctx interface{}, component interface{}, opts interface{}) *mockComponentInterface_Create_Call {
	return &mockComponentInterface_Create_Call{Call: _e.mock.On("Create", ctx, component, opts)}
}

func (_c *mockComponentInterface_Create_Call) Run(run func(ctx context.Context, component *v1.Component, opts metav1.CreateOptions)) *mockComponentInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(metav1.CreateOptions))
	})
	return _c
}

func (_c *mockComponentInterface_Create_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_Create_Call) RunAndReturn(run func(context.Context, *v1.Component, metav1.CreateOptions) (*v1.Component, error)) *mockComponentInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *mockComponentInterface) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockComponentInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.DeleteOptions
func (_e *mockComponentInterface_Expecter) Delete(ctx interface{}, name interface{}, opts interface{}) *mockComponentInterface_Delete_Call {
	return &mockComponentInterface_Delete_Call{Call: _e.mock.On("Delete", ctx, name, opts)}
}

func (_c *mockComponentInterface_Delete_Call) Run(run func(ctx context.Context, name string, opts metav1.DeleteOptions)) *mockComponentInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.DeleteOptions))
	})
	return _c
}

func (_c *mockComponentInterface_Delete_Call) Return(_a0 error) *mockComponentInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInterface_Delete_Call) RunAndReturn(run func(context.Context, string, metav1.DeleteOptions) error) *mockComponentInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *mockComponentInterface) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCollection")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInterface_DeleteCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCollection'
type mockComponentInterface_DeleteCollection_Call struct {
	*mock.Call
}

// DeleteCollection is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.DeleteOptions
//   - listOpts metav1.ListOptions
func (_e *mockComponentInterface_Expecter) DeleteCollection(ctx interface{}, opts interface{}, listOpts interface{}) *mockComponentInterface_DeleteCollection_Call {
	return &mockComponentInterface_DeleteCollection_Call{Call: _e.mock.On("DeleteCollection", ctx, opts, listOpts)}
}

func (_c *mockComponentInterface_DeleteCollection_Call) Run(run func(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions)) *mockComponentInterface_DeleteCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.DeleteOptions), args[2].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockComponentInterface_DeleteCollection_Call) Return(_a0 error) *mockComponentInterface_DeleteCollection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInterface_DeleteCollection_Call) RunAndReturn(run func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error) *mockComponentInterface_DeleteCollection_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockComponentInterface) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*v1.Component, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *v1.Component); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockComponentInterface_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.GetOptions
func (_e *mockComponentInterface_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockComponentInterface_Get_Call {
	return &mockComponentInterface_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockComponentInterface_Get_Call) Run(run func(ctx context.Context, name string, opts metav1.GetOptions)) *mockComponentInterface_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.GetOptions))
	})
	return _c
}

func (_c *mockComponentInterface_Get_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_Get_Call) RunAndReturn(run func(context.Context, string, metav1.GetOptions) (*v1.Component, error)) *mockComponentInterface_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *mockComponentInterface) List(ctx context.Context, opts metav1.ListOptions) (*v1.ComponentList, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *v1.ComponentList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (*v1.ComponentList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *v1.ComponentList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ComponentList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type mockComponentInterface_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockComponentInterface_Expecter) List(ctx interface{}, opts interface{}) *mockComponentInterface_List_Call {
	return &mockComponentInterface_List_Call{Call: _e.mock.On("List", ctx, opts)}
}

func (_c *mockComponentInterface_List_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockComponentInterface_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockComponentInterface_List_Call) Return(_a0 *v1.ComponentList, _a1 error) *mockComponentInterface_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_List_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (*v1.ComponentList, error)) *mockComponentInterface_List_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *mockComponentInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1.Component, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*v1.Component, error)); ok {
		return rf(ctx, name, pt, data, opts, subresources...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) *v1.Component); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type mockComponentInterface_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - pt types.PatchType
//   - data []byte
//   - opts metav1.PatchOptions
//   - subresources ...string
func (_e *mockComponentInterface_Expecter) Patch(ctx interface{}, name interface{}, pt interface{}, data interface{}, opts interface{}, subresources ...interface{}) *mockComponentInterface_Patch_Call {
	return &mockComponentInterface_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, name, pt, data, opts}, subresources...)...)}
}

func (_c *mockComponentInterface_Patch_Call) Run(run func(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string)) *mockComponentInterface_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-5)
		for i, a := range args[5:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(types.PatchType), args[3].([]byte), args[4].(metav1.PatchOptions), variadicArgs...)
	})
	return _c
}

func (_c *mockComponentInterface_Patch_Call) Return(result *v1.Component, err error) *mockComponentInterface_Patch_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockComponentInterface_Patch_Call) RunAndReturn(run func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*v1.Component, error)) *mockComponentInterface_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveFinalizer provides a mock function with given fields: ctx, component, finalizer
func (_m *mockComponentInterface) RemoveFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	ret := _m.Called(ctx, component, finalizer)

	if len(ret) == 0 {
		panic("no return value specified for RemoveFinalizer")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, string) (*v1.Component, error)); ok {
		return rf(ctx, component, finalizer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, string) *v1.Component); ok {
		r0 = rf(ctx, component, finalizer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component, string) error); ok {
		r1 = rf(ctx, component, finalizer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_RemoveFinalizer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveFinalizer'
type mockComponentInterface_RemoveFinalizer_Call struct {
	*mock.Call
}

// RemoveFinalizer is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - finalizer string
func (_e *mockComponentInterface_Expecter) RemoveFinalizer(ctx interface{}, component interface{}, finalizer interface{}) *mockComponentInterface_RemoveFinalizer_Call {
	return &mockComponentInterface_RemoveFinalizer_Call{Call: _e.mock.On("RemoveFinalizer", ctx, component, finalizer)}
}

func (_c *mockComponentInterface_RemoveFinalizer_Call) Run(run func(ctx context.Context, component *v1.Component, finalizer string)) *mockComponentInterface_RemoveFinalizer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(string))
	})
	return _c
}

func (_c *mockComponentInterface_RemoveFinalizer_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_RemoveFinalizer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_RemoveFinalizer_Call) RunAndReturn(run func(context.Context, *v1.Component, string) (*v1.Component, error)) *mockComponentInterface_RemoveFinalizer_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, component, opts
func (_m *mockComponentInterface) Update(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, component, opts)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, metav1.UpdateOptions) (*v1.Component, error)); ok {
		return rf(ctx, component, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, metav1.UpdateOptions) *v1.Component); ok {
		r0 = rf(ctx, component, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, component, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockComponentInterface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - opts metav1.UpdateOptions
func (_e *mockComponentInterface_Expecter) Update(ctx interface{}, component interface{}, opts interface{}) *mockComponentInterface_Update_Call {
	return &mockComponentInterface_Update_Call{Call: _e.mock.On("Update", ctx, component, opts)}
}

func (_c *mockComponentInterface_Update_Call) Run(run func(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions)) *mockComponentInterface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockComponentInterface_Update_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_Update_Call) RunAndReturn(run func(context.Context, *v1.Component, metav1.UpdateOptions) (*v1.Component, error)) *mockComponentInterface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateExpectedComponentVersion provides a mock function with given fields: ctx, componentName, version
func (_m *mockComponentInterface) UpdateExpectedComponentVersion(ctx context.Context, componentName string, version string) (*v1.Component, error) {
	ret := _m.Called(ctx, componentName, version)

	if len(ret) == 0 {
		panic("no return value specified for UpdateExpectedComponentVersion")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*v1.Component, error)); ok {
		return rf(ctx, componentName, version)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *v1.Component); ok {
		r0 = rf(ctx, componentName, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, componentName, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateExpectedComponentVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateExpectedComponentVersion'
type mockComponentInterface_UpdateExpectedComponentVersion_Call struct {
	*mock.Call
}

// UpdateExpectedComponentVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName string
//   - version string
func (_e *mockComponentInterface_Expecter) UpdateExpectedComponentVersion(ctx interface{}, componentName interface{}, version interface{}) *mockComponentInterface_UpdateExpectedComponentVersion_Call {
	return &mockComponentInterface_UpdateExpectedComponentVersion_Call{Call: _e.mock.On("UpdateExpectedComponentVersion", ctx, componentName, version)}
}

func (_c *mockComponentInterface_UpdateExpectedComponentVersion_Call) Run(run func(ctx context.Context, componentName string, version string)) *mockComponentInterface_UpdateExpectedComponentVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateExpectedComponentVersion_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateExpectedComponentVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateExpectedComponentVersion_Call) RunAndReturn(run func(context.Context, string, string) (*v1.Component, error)) *mockComponentInterface_UpdateExpectedComponentVersion_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, component, opts
func (_m *mockComponentInterface) UpdateStatus(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, component, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, metav1.UpdateOptions) (*v1.Component, error)); ok {
		return rf(ctx, component, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component, metav1.UpdateOptions) *v1.Component); ok {
		r0 = rf(ctx, component, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, component, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type mockComponentInterface_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - opts metav1.UpdateOptions
func (_e *mockComponentInterface_Expecter) UpdateStatus(ctx interface{}, component interface{}, opts interface{}) *mockComponentInterface_UpdateStatus_Call {
	return &mockComponentInterface_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, component, opts)}
}

func (_c *mockComponentInterface_UpdateStatus_Call) Run(run func(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions)) *mockComponentInterface_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateStatus_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateStatus_Call) RunAndReturn(run func(context.Context, *v1.Component, metav1.UpdateOptions) (*v1.Component, error)) *mockComponentInterface_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusDeleting provides a mock function with given fields: ctx, component
func (_m *mockComponentInterface) UpdateStatusDeleting(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatusDeleting")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) (*v1.Component, error)); ok {
		return rf(ctx, component)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) *v1.Component); ok {
		r0 = rf(ctx, component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component) error); ok {
		r1 = rf(ctx, component)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateStatusDeleting_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusDeleting'
type mockComponentInterface_UpdateStatusDeleting_Call struct {
	*mock.Call
}

// UpdateStatusDeleting is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentInterface_Expecter) UpdateStatusDeleting(ctx interface{}, component interface{}) *mockComponentInterface_UpdateStatusDeleting_Call {
	return &mockComponentInterface_UpdateStatusDeleting_Call{Call: _e.mock.On("UpdateStatusDeleting", ctx, component)}
}

func (_c *mockComponentInterface_UpdateStatusDeleting_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentInterface_UpdateStatusDeleting_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateStatusDeleting_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateStatusDeleting_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateStatusDeleting_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentInterface_UpdateStatusDeleting_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusInstalled provides a mock function with given fields: ctx, component
func (_m *mockComponentInterface) UpdateStatusInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatusInstalled")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) (*v1.Component, error)); ok {
		return rf(ctx, component)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) *v1.Component); ok {
		r0 = rf(ctx, component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component) error); ok {
		r1 = rf(ctx, component)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateStatusInstalled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusInstalled'
type mockComponentInterface_UpdateStatusInstalled_Call struct {
	*mock.Call
}

// UpdateStatusInstalled is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentInterface_Expecter) UpdateStatusInstalled(ctx interface{}, component interface{}) *mockComponentInterface_UpdateStatusInstalled_Call {
	return &mockComponentInterface_UpdateStatusInstalled_Call{Call: _e.mock.On("UpdateStatusInstalled", ctx, component)}
}

func (_c *mockComponentInterface_UpdateStatusInstalled_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentInterface_UpdateStatusInstalled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateStatusInstalled_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateStatusInstalled_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateStatusInstalled_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentInterface_UpdateStatusInstalled_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusInstalling provides a mock function with given fields: ctx, component
func (_m *mockComponentInterface) UpdateStatusInstalling(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatusInstalling")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) (*v1.Component, error)); ok {
		return rf(ctx, component)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) *v1.Component); ok {
		r0 = rf(ctx, component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component) error); ok {
		r1 = rf(ctx, component)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateStatusInstalling_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusInstalling'
type mockComponentInterface_UpdateStatusInstalling_Call struct {
	*mock.Call
}

// UpdateStatusInstalling is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentInterface_Expecter) UpdateStatusInstalling(ctx interface{}, component interface{}) *mockComponentInterface_UpdateStatusInstalling_Call {
	return &mockComponentInterface_UpdateStatusInstalling_Call{Call: _e.mock.On("UpdateStatusInstalling", ctx, component)}
}

func (_c *mockComponentInterface_UpdateStatusInstalling_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentInterface_UpdateStatusInstalling_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateStatusInstalling_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateStatusInstalling_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateStatusInstalling_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentInterface_UpdateStatusInstalling_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusNotInstalled provides a mock function with given fields: ctx, component
func (_m *mockComponentInterface) UpdateStatusNotInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatusNotInstalled")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) (*v1.Component, error)); ok {
		return rf(ctx, component)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) *v1.Component); ok {
		r0 = rf(ctx, component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component) error); ok {
		r1 = rf(ctx, component)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateStatusNotInstalled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusNotInstalled'
type mockComponentInterface_UpdateStatusNotInstalled_Call struct {
	*mock.Call
}

// UpdateStatusNotInstalled is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentInterface_Expecter) UpdateStatusNotInstalled(ctx interface{}, component interface{}) *mockComponentInterface_UpdateStatusNotInstalled_Call {
	return &mockComponentInterface_UpdateStatusNotInstalled_Call{Call: _e.mock.On("UpdateStatusNotInstalled", ctx, component)}
}

func (_c *mockComponentInterface_UpdateStatusNotInstalled_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentInterface_UpdateStatusNotInstalled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateStatusNotInstalled_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateStatusNotInstalled_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateStatusNotInstalled_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentInterface_UpdateStatusNotInstalled_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusUpgrading provides a mock function with given fields: ctx, component
func (_m *mockComponentInterface) UpdateStatusUpgrading(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatusUpgrading")
	}

	var r0 *v1.Component
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) (*v1.Component, error)); ok {
		return rf(ctx, component)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) *v1.Component); ok {
		r0 = rf(ctx, component)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Component)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Component) error); ok {
		r1 = rf(ctx, component)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_UpdateStatusUpgrading_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusUpgrading'
type mockComponentInterface_UpdateStatusUpgrading_Call struct {
	*mock.Call
}

// UpdateStatusUpgrading is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentInterface_Expecter) UpdateStatusUpgrading(ctx interface{}, component interface{}) *mockComponentInterface_UpdateStatusUpgrading_Call {
	return &mockComponentInterface_UpdateStatusUpgrading_Call{Call: _e.mock.On("UpdateStatusUpgrading", ctx, component)}
}

func (_c *mockComponentInterface_UpdateStatusUpgrading_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentInterface_UpdateStatusUpgrading_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentInterface_UpdateStatusUpgrading_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentInterface_UpdateStatusUpgrading_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_UpdateStatusUpgrading_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentInterface_UpdateStatusUpgrading_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *mockComponentInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Watch")
	}

	var r0 watch.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (watch.Interface, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInterface_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type mockComponentInterface_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockComponentInterface_Expecter) Watch(ctx interface{}, opts interface{}) *mockComponentInterface_Watch_Call {
	return &mockComponentInterface_Watch_Call{Call: _e.mock.On("Watch", ctx, opts)}
}

func (_c *mockComponentInterface_Watch_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockComponentInterface_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockComponentInterface_Watch_Call) Return(_a0 watch.Interface, _a1 error) *mockComponentInterface_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInterface_Watch_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (watch.Interface, error)) *mockComponentInterface_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// newMockComponentInterface creates a new instance of mockComponentInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockComponentInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockComponentInterface {
	mock := &mockComponentInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
