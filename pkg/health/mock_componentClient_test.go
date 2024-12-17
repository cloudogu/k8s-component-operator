// Code generated by mockery v2.20.0. DO NOT EDIT.

package health

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	types "k8s.io/apimachinery/pkg/types"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// mockComponentClient is an autogenerated mock type for the componentClient type
type mockComponentClient struct {
	mock.Mock
}

type mockComponentClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockComponentClient) EXPECT() *mockComponentClient_Expecter {
	return &mockComponentClient_Expecter{mock: &_m.Mock}
}

// AddFinalizer provides a mock function with given fields: ctx, component, finalizer
func (_m *mockComponentClient) AddFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	ret := _m.Called(ctx, component, finalizer)

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

// mockComponentClient_AddFinalizer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFinalizer'
type mockComponentClient_AddFinalizer_Call struct {
	*mock.Call
}

// AddFinalizer is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - finalizer string
func (_e *mockComponentClient_Expecter) AddFinalizer(ctx interface{}, component interface{}, finalizer interface{}) *mockComponentClient_AddFinalizer_Call {
	return &mockComponentClient_AddFinalizer_Call{Call: _e.mock.On("AddFinalizer", ctx, component, finalizer)}
}

func (_c *mockComponentClient_AddFinalizer_Call) Run(run func(ctx context.Context, component *v1.Component, finalizer string)) *mockComponentClient_AddFinalizer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(string))
	})
	return _c
}

func (_c *mockComponentClient_AddFinalizer_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_AddFinalizer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_AddFinalizer_Call) RunAndReturn(run func(context.Context, *v1.Component, string) (*v1.Component, error)) *mockComponentClient_AddFinalizer_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, component, opts
func (_m *mockComponentClient) Create(ctx context.Context, component *v1.Component, opts metav1.CreateOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, component, opts)

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

// mockComponentClient_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockComponentClient_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - opts metav1.CreateOptions
func (_e *mockComponentClient_Expecter) Create(ctx interface{}, component interface{}, opts interface{}) *mockComponentClient_Create_Call {
	return &mockComponentClient_Create_Call{Call: _e.mock.On("Create", ctx, component, opts)}
}

func (_c *mockComponentClient_Create_Call) Run(run func(ctx context.Context, component *v1.Component, opts metav1.CreateOptions)) *mockComponentClient_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(metav1.CreateOptions))
	})
	return _c
}

func (_c *mockComponentClient_Create_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_Create_Call) RunAndReturn(run func(context.Context, *v1.Component, metav1.CreateOptions) (*v1.Component, error)) *mockComponentClient_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *mockComponentClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentClient_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockComponentClient_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.DeleteOptions
func (_e *mockComponentClient_Expecter) Delete(ctx interface{}, name interface{}, opts interface{}) *mockComponentClient_Delete_Call {
	return &mockComponentClient_Delete_Call{Call: _e.mock.On("Delete", ctx, name, opts)}
}

func (_c *mockComponentClient_Delete_Call) Run(run func(ctx context.Context, name string, opts metav1.DeleteOptions)) *mockComponentClient_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.DeleteOptions))
	})
	return _c
}

func (_c *mockComponentClient_Delete_Call) Return(_a0 error) *mockComponentClient_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentClient_Delete_Call) RunAndReturn(run func(context.Context, string, metav1.DeleteOptions) error) *mockComponentClient_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *mockComponentClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentClient_DeleteCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCollection'
type mockComponentClient_DeleteCollection_Call struct {
	*mock.Call
}

// DeleteCollection is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.DeleteOptions
//   - listOpts metav1.ListOptions
func (_e *mockComponentClient_Expecter) DeleteCollection(ctx interface{}, opts interface{}, listOpts interface{}) *mockComponentClient_DeleteCollection_Call {
	return &mockComponentClient_DeleteCollection_Call{Call: _e.mock.On("DeleteCollection", ctx, opts, listOpts)}
}

func (_c *mockComponentClient_DeleteCollection_Call) Run(run func(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions)) *mockComponentClient_DeleteCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.DeleteOptions), args[2].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockComponentClient_DeleteCollection_Call) Return(_a0 error) *mockComponentClient_DeleteCollection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentClient_DeleteCollection_Call) RunAndReturn(run func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error) *mockComponentClient_DeleteCollection_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockComponentClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, name, opts)

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

// mockComponentClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockComponentClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.GetOptions
func (_e *mockComponentClient_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockComponentClient_Get_Call {
	return &mockComponentClient_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockComponentClient_Get_Call) Run(run func(ctx context.Context, name string, opts metav1.GetOptions)) *mockComponentClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.GetOptions))
	})
	return _c
}

func (_c *mockComponentClient_Get_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_Get_Call) RunAndReturn(run func(context.Context, string, metav1.GetOptions) (*v1.Component, error)) *mockComponentClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *mockComponentClient) List(ctx context.Context, opts metav1.ListOptions) (*v1.ComponentList, error) {
	ret := _m.Called(ctx, opts)

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

// mockComponentClient_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type mockComponentClient_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockComponentClient_Expecter) List(ctx interface{}, opts interface{}) *mockComponentClient_List_Call {
	return &mockComponentClient_List_Call{Call: _e.mock.On("List", ctx, opts)}
}

func (_c *mockComponentClient_List_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockComponentClient_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockComponentClient_List_Call) Return(_a0 *v1.ComponentList, _a1 error) *mockComponentClient_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_List_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (*v1.ComponentList, error)) *mockComponentClient_List_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *mockComponentClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*v1.Component, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

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

// mockComponentClient_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type mockComponentClient_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - pt types.PatchType
//   - data []byte
//   - opts metav1.PatchOptions
//   - subresources ...string
func (_e *mockComponentClient_Expecter) Patch(ctx interface{}, name interface{}, pt interface{}, data interface{}, opts interface{}, subresources ...interface{}) *mockComponentClient_Patch_Call {
	return &mockComponentClient_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, name, pt, data, opts}, subresources...)...)}
}

func (_c *mockComponentClient_Patch_Call) Run(run func(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string)) *mockComponentClient_Patch_Call {
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

func (_c *mockComponentClient_Patch_Call) Return(result *v1.Component, err error) *mockComponentClient_Patch_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockComponentClient_Patch_Call) RunAndReturn(run func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*v1.Component, error)) *mockComponentClient_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveFinalizer provides a mock function with given fields: ctx, component, finalizer
func (_m *mockComponentClient) RemoveFinalizer(ctx context.Context, component *v1.Component, finalizer string) (*v1.Component, error) {
	ret := _m.Called(ctx, component, finalizer)

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

// mockComponentClient_RemoveFinalizer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveFinalizer'
type mockComponentClient_RemoveFinalizer_Call struct {
	*mock.Call
}

// RemoveFinalizer is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - finalizer string
func (_e *mockComponentClient_Expecter) RemoveFinalizer(ctx interface{}, component interface{}, finalizer interface{}) *mockComponentClient_RemoveFinalizer_Call {
	return &mockComponentClient_RemoveFinalizer_Call{Call: _e.mock.On("RemoveFinalizer", ctx, component, finalizer)}
}

func (_c *mockComponentClient_RemoveFinalizer_Call) Run(run func(ctx context.Context, component *v1.Component, finalizer string)) *mockComponentClient_RemoveFinalizer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(string))
	})
	return _c
}

func (_c *mockComponentClient_RemoveFinalizer_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_RemoveFinalizer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_RemoveFinalizer_Call) RunAndReturn(run func(context.Context, *v1.Component, string) (*v1.Component, error)) *mockComponentClient_RemoveFinalizer_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, component, opts
func (_m *mockComponentClient) Update(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, component, opts)

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

// mockComponentClient_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockComponentClient_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - opts metav1.UpdateOptions
func (_e *mockComponentClient_Expecter) Update(ctx interface{}, component interface{}, opts interface{}) *mockComponentClient_Update_Call {
	return &mockComponentClient_Update_Call{Call: _e.mock.On("Update", ctx, component, opts)}
}

func (_c *mockComponentClient_Update_Call) Run(run func(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions)) *mockComponentClient_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockComponentClient_Update_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_Update_Call) RunAndReturn(run func(context.Context, *v1.Component, metav1.UpdateOptions) (*v1.Component, error)) *mockComponentClient_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateExpectedComponentVersion provides a mock function with given fields: ctx, componentName, version
func (_m *mockComponentClient) UpdateExpectedComponentVersion(ctx context.Context, componentName string, version string) (*v1.Component, error) {
	ret := _m.Called(ctx, componentName, version)

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

// mockComponentClient_UpdateExpectedComponentVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateExpectedComponentVersion'
type mockComponentClient_UpdateExpectedComponentVersion_Call struct {
	*mock.Call
}

// UpdateExpectedComponentVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName string
//   - version string
func (_e *mockComponentClient_Expecter) UpdateExpectedComponentVersion(ctx interface{}, componentName interface{}, version interface{}) *mockComponentClient_UpdateExpectedComponentVersion_Call {
	return &mockComponentClient_UpdateExpectedComponentVersion_Call{Call: _e.mock.On("UpdateExpectedComponentVersion", ctx, componentName, version)}
}

func (_c *mockComponentClient_UpdateExpectedComponentVersion_Call) Run(run func(ctx context.Context, componentName string, version string)) *mockComponentClient_UpdateExpectedComponentVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *mockComponentClient_UpdateExpectedComponentVersion_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateExpectedComponentVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateExpectedComponentVersion_Call) RunAndReturn(run func(context.Context, string, string) (*v1.Component, error)) *mockComponentClient_UpdateExpectedComponentVersion_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, component, opts
func (_m *mockComponentClient) UpdateStatus(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions) (*v1.Component, error) {
	ret := _m.Called(ctx, component, opts)

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

// mockComponentClient_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type mockComponentClient_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
//   - opts metav1.UpdateOptions
func (_e *mockComponentClient_Expecter) UpdateStatus(ctx interface{}, component interface{}, opts interface{}) *mockComponentClient_UpdateStatus_Call {
	return &mockComponentClient_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, component, opts)}
}

func (_c *mockComponentClient_UpdateStatus_Call) Run(run func(ctx context.Context, component *v1.Component, opts metav1.UpdateOptions)) *mockComponentClient_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockComponentClient_UpdateStatus_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateStatus_Call) RunAndReturn(run func(context.Context, *v1.Component, metav1.UpdateOptions) (*v1.Component, error)) *mockComponentClient_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusDeleting provides a mock function with given fields: ctx, component
func (_m *mockComponentClient) UpdateStatusDeleting(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

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

// mockComponentClient_UpdateStatusDeleting_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusDeleting'
type mockComponentClient_UpdateStatusDeleting_Call struct {
	*mock.Call
}

// UpdateStatusDeleting is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentClient_Expecter) UpdateStatusDeleting(ctx interface{}, component interface{}) *mockComponentClient_UpdateStatusDeleting_Call {
	return &mockComponentClient_UpdateStatusDeleting_Call{Call: _e.mock.On("UpdateStatusDeleting", ctx, component)}
}

func (_c *mockComponentClient_UpdateStatusDeleting_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentClient_UpdateStatusDeleting_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentClient_UpdateStatusDeleting_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateStatusDeleting_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateStatusDeleting_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentClient_UpdateStatusDeleting_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusInstalled provides a mock function with given fields: ctx, component
func (_m *mockComponentClient) UpdateStatusInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

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

// mockComponentClient_UpdateStatusInstalled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusInstalled'
type mockComponentClient_UpdateStatusInstalled_Call struct {
	*mock.Call
}

// UpdateStatusInstalled is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentClient_Expecter) UpdateStatusInstalled(ctx interface{}, component interface{}) *mockComponentClient_UpdateStatusInstalled_Call {
	return &mockComponentClient_UpdateStatusInstalled_Call{Call: _e.mock.On("UpdateStatusInstalled", ctx, component)}
}

func (_c *mockComponentClient_UpdateStatusInstalled_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentClient_UpdateStatusInstalled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentClient_UpdateStatusInstalled_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateStatusInstalled_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateStatusInstalled_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentClient_UpdateStatusInstalled_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusInstalling provides a mock function with given fields: ctx, component
func (_m *mockComponentClient) UpdateStatusInstalling(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

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

// mockComponentClient_UpdateStatusInstalling_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusInstalling'
type mockComponentClient_UpdateStatusInstalling_Call struct {
	*mock.Call
}

// UpdateStatusInstalling is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentClient_Expecter) UpdateStatusInstalling(ctx interface{}, component interface{}) *mockComponentClient_UpdateStatusInstalling_Call {
	return &mockComponentClient_UpdateStatusInstalling_Call{Call: _e.mock.On("UpdateStatusInstalling", ctx, component)}
}

func (_c *mockComponentClient_UpdateStatusInstalling_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentClient_UpdateStatusInstalling_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentClient_UpdateStatusInstalling_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateStatusInstalling_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateStatusInstalling_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentClient_UpdateStatusInstalling_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusNotInstalled provides a mock function with given fields: ctx, component
func (_m *mockComponentClient) UpdateStatusNotInstalled(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

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

// mockComponentClient_UpdateStatusNotInstalled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusNotInstalled'
type mockComponentClient_UpdateStatusNotInstalled_Call struct {
	*mock.Call
}

// UpdateStatusNotInstalled is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentClient_Expecter) UpdateStatusNotInstalled(ctx interface{}, component interface{}) *mockComponentClient_UpdateStatusNotInstalled_Call {
	return &mockComponentClient_UpdateStatusNotInstalled_Call{Call: _e.mock.On("UpdateStatusNotInstalled", ctx, component)}
}

func (_c *mockComponentClient_UpdateStatusNotInstalled_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentClient_UpdateStatusNotInstalled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentClient_UpdateStatusNotInstalled_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateStatusNotInstalled_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateStatusNotInstalled_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentClient_UpdateStatusNotInstalled_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusUpgrading provides a mock function with given fields: ctx, component
func (_m *mockComponentClient) UpdateStatusUpgrading(ctx context.Context, component *v1.Component) (*v1.Component, error) {
	ret := _m.Called(ctx, component)

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

// mockComponentClient_UpdateStatusUpgrading_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusUpgrading'
type mockComponentClient_UpdateStatusUpgrading_Call struct {
	*mock.Call
}

// UpdateStatusUpgrading is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *mockComponentClient_Expecter) UpdateStatusUpgrading(ctx interface{}, component interface{}) *mockComponentClient_UpdateStatusUpgrading_Call {
	return &mockComponentClient_UpdateStatusUpgrading_Call{Call: _e.mock.On("UpdateStatusUpgrading", ctx, component)}
}

func (_c *mockComponentClient_UpdateStatusUpgrading_Call) Run(run func(ctx context.Context, component *v1.Component)) *mockComponentClient_UpdateStatusUpgrading_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *mockComponentClient_UpdateStatusUpgrading_Call) Return(_a0 *v1.Component, _a1 error) *mockComponentClient_UpdateStatusUpgrading_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_UpdateStatusUpgrading_Call) RunAndReturn(run func(context.Context, *v1.Component) (*v1.Component, error)) *mockComponentClient_UpdateStatusUpgrading_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *mockComponentClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

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

// mockComponentClient_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type mockComponentClient_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockComponentClient_Expecter) Watch(ctx interface{}, opts interface{}) *mockComponentClient_Watch_Call {
	return &mockComponentClient_Watch_Call{Call: _e.mock.On("Watch", ctx, opts)}
}

func (_c *mockComponentClient_Watch_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockComponentClient_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockComponentClient_Watch_Call) Return(_a0 watch.Interface, _a1 error) *mockComponentClient_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentClient_Watch_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (watch.Interface, error)) *mockComponentClient_Watch_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockComponentClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockComponentClient creates a new instance of mockComponentClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockComponentClient(t mockConstructorTestingTnewMockComponentClient) *mockComponentClient {
	mock := &mockComponentClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
