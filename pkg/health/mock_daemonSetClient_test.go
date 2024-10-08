// Code generated by mockery v2.42.1. DO NOT EDIT.

package health

import (
	context "context"

	appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mock "github.com/stretchr/testify/mock"

	types "k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/client-go/applyconfigurations/apps/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// mockDaemonSetClient is an autogenerated mock type for the daemonSetClient type
type mockDaemonSetClient struct {
	mock.Mock
}

type mockDaemonSetClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDaemonSetClient) EXPECT() *mockDaemonSetClient_Expecter {
	return &mockDaemonSetClient_Expecter{mock: &_m.Mock}
}

// Apply provides a mock function with given fields: ctx, daemonSet, opts
func (_m *mockDaemonSetClient) Apply(ctx context.Context, daemonSet *v1.DaemonSetApplyConfiguration, opts metav1.ApplyOptions) (*appsv1.DaemonSet, error) {
	ret := _m.Called(ctx, daemonSet, opts)

	if len(ret) == 0 {
		panic("no return value specified for Apply")
	}

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, daemonSet, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, daemonSet, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, daemonSet, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_Apply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Apply'
type mockDaemonSetClient_Apply_Call struct {
	*mock.Call
}

// Apply is a helper method to define mock.On call
//   - ctx context.Context
//   - daemonSet *v1.DaemonSetApplyConfiguration
//   - opts metav1.ApplyOptions
func (_e *mockDaemonSetClient_Expecter) Apply(ctx interface{}, daemonSet interface{}, opts interface{}) *mockDaemonSetClient_Apply_Call {
	return &mockDaemonSetClient_Apply_Call{Call: _e.mock.On("Apply", ctx, daemonSet, opts)}
}

func (_c *mockDaemonSetClient_Apply_Call) Run(run func(ctx context.Context, daemonSet *v1.DaemonSetApplyConfiguration, opts metav1.ApplyOptions)) *mockDaemonSetClient_Apply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.DaemonSetApplyConfiguration), args[2].(metav1.ApplyOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_Apply_Call) Return(result *appsv1.DaemonSet, err error) *mockDaemonSetClient_Apply_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDaemonSetClient_Apply_Call) RunAndReturn(run func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_Apply_Call {
	_c.Call.Return(run)
	return _c
}

// ApplyStatus provides a mock function with given fields: ctx, daemonSet, opts
func (_m *mockDaemonSetClient) ApplyStatus(ctx context.Context, daemonSet *v1.DaemonSetApplyConfiguration, opts metav1.ApplyOptions) (*appsv1.DaemonSet, error) {
	ret := _m.Called(ctx, daemonSet, opts)

	if len(ret) == 0 {
		panic("no return value specified for ApplyStatus")
	}

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, daemonSet, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, daemonSet, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) error); ok {
		r1 = rf(ctx, daemonSet, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_ApplyStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyStatus'
type mockDaemonSetClient_ApplyStatus_Call struct {
	*mock.Call
}

// ApplyStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - daemonSet *v1.DaemonSetApplyConfiguration
//   - opts metav1.ApplyOptions
func (_e *mockDaemonSetClient_Expecter) ApplyStatus(ctx interface{}, daemonSet interface{}, opts interface{}) *mockDaemonSetClient_ApplyStatus_Call {
	return &mockDaemonSetClient_ApplyStatus_Call{Call: _e.mock.On("ApplyStatus", ctx, daemonSet, opts)}
}

func (_c *mockDaemonSetClient_ApplyStatus_Call) Run(run func(ctx context.Context, daemonSet *v1.DaemonSetApplyConfiguration, opts metav1.ApplyOptions)) *mockDaemonSetClient_ApplyStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.DaemonSetApplyConfiguration), args[2].(metav1.ApplyOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_ApplyStatus_Call) Return(result *appsv1.DaemonSet, err error) *mockDaemonSetClient_ApplyStatus_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDaemonSetClient_ApplyStatus_Call) RunAndReturn(run func(context.Context, *v1.DaemonSetApplyConfiguration, metav1.ApplyOptions) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_ApplyStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, daemonSet, opts
func (_m *mockDaemonSetClient) Create(ctx context.Context, daemonSet *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error) {
	ret := _m.Called(ctx, daemonSet, opts)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.DaemonSet, metav1.CreateOptions) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, daemonSet, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.DaemonSet, metav1.CreateOptions) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, daemonSet, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *appsv1.DaemonSet, metav1.CreateOptions) error); ok {
		r1 = rf(ctx, daemonSet, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockDaemonSetClient_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - daemonSet *appsv1.DaemonSet
//   - opts metav1.CreateOptions
func (_e *mockDaemonSetClient_Expecter) Create(ctx interface{}, daemonSet interface{}, opts interface{}) *mockDaemonSetClient_Create_Call {
	return &mockDaemonSetClient_Create_Call{Call: _e.mock.On("Create", ctx, daemonSet, opts)}
}

func (_c *mockDaemonSetClient_Create_Call) Run(run func(ctx context.Context, daemonSet *appsv1.DaemonSet, opts metav1.CreateOptions)) *mockDaemonSetClient_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appsv1.DaemonSet), args[2].(metav1.CreateOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_Create_Call) Return(_a0 *appsv1.DaemonSet, _a1 error) *mockDaemonSetClient_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDaemonSetClient_Create_Call) RunAndReturn(run func(context.Context, *appsv1.DaemonSet, metav1.CreateOptions) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *mockDaemonSetClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
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

// mockDaemonSetClient_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDaemonSetClient_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.DeleteOptions
func (_e *mockDaemonSetClient_Expecter) Delete(ctx interface{}, name interface{}, opts interface{}) *mockDaemonSetClient_Delete_Call {
	return &mockDaemonSetClient_Delete_Call{Call: _e.mock.On("Delete", ctx, name, opts)}
}

func (_c *mockDaemonSetClient_Delete_Call) Run(run func(ctx context.Context, name string, opts metav1.DeleteOptions)) *mockDaemonSetClient_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.DeleteOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_Delete_Call) Return(_a0 error) *mockDaemonSetClient_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDaemonSetClient_Delete_Call) RunAndReturn(run func(context.Context, string, metav1.DeleteOptions) error) *mockDaemonSetClient_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *mockDaemonSetClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
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

// mockDaemonSetClient_DeleteCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCollection'
type mockDaemonSetClient_DeleteCollection_Call struct {
	*mock.Call
}

// DeleteCollection is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.DeleteOptions
//   - listOpts metav1.ListOptions
func (_e *mockDaemonSetClient_Expecter) DeleteCollection(ctx interface{}, opts interface{}, listOpts interface{}) *mockDaemonSetClient_DeleteCollection_Call {
	return &mockDaemonSetClient_DeleteCollection_Call{Call: _e.mock.On("DeleteCollection", ctx, opts, listOpts)}
}

func (_c *mockDaemonSetClient_DeleteCollection_Call) Run(run func(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions)) *mockDaemonSetClient_DeleteCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.DeleteOptions), args[2].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_DeleteCollection_Call) Return(_a0 error) *mockDaemonSetClient_DeleteCollection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDaemonSetClient_DeleteCollection_Call) RunAndReturn(run func(context.Context, metav1.DeleteOptions, metav1.ListOptions) error) *mockDaemonSetClient_DeleteCollection_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *mockDaemonSetClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDaemonSetClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts metav1.GetOptions
func (_e *mockDaemonSetClient_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *mockDaemonSetClient_Get_Call {
	return &mockDaemonSetClient_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *mockDaemonSetClient_Get_Call) Run(run func(ctx context.Context, name string, opts metav1.GetOptions)) *mockDaemonSetClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(metav1.GetOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_Get_Call) Return(_a0 *appsv1.DaemonSet, _a1 error) *mockDaemonSetClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDaemonSetClient_Get_Call) RunAndReturn(run func(context.Context, string, metav1.GetOptions) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *mockDaemonSetClient) List(ctx context.Context, opts metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *appsv1.DaemonSetList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (*appsv1.DaemonSetList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *appsv1.DaemonSetList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSetList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type mockDaemonSetClient_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockDaemonSetClient_Expecter) List(ctx interface{}, opts interface{}) *mockDaemonSetClient_List_Call {
	return &mockDaemonSetClient_List_Call{Call: _e.mock.On("List", ctx, opts)}
}

func (_c *mockDaemonSetClient_List_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockDaemonSetClient_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_List_Call) Return(_a0 *appsv1.DaemonSetList, _a1 error) *mockDaemonSetClient_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDaemonSetClient_List_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (*appsv1.DaemonSetList, error)) *mockDaemonSetClient_List_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *mockDaemonSetClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*appsv1.DaemonSet, error) {
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

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, name, pt, data, opts, subresources...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type mockDaemonSetClient_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - pt types.PatchType
//   - data []byte
//   - opts metav1.PatchOptions
//   - subresources ...string
func (_e *mockDaemonSetClient_Expecter) Patch(ctx interface{}, name interface{}, pt interface{}, data interface{}, opts interface{}, subresources ...interface{}) *mockDaemonSetClient_Patch_Call {
	return &mockDaemonSetClient_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, name, pt, data, opts}, subresources...)...)}
}

func (_c *mockDaemonSetClient_Patch_Call) Run(run func(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string)) *mockDaemonSetClient_Patch_Call {
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

func (_c *mockDaemonSetClient_Patch_Call) Return(result *appsv1.DaemonSet, err error) *mockDaemonSetClient_Patch_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *mockDaemonSetClient_Patch_Call) RunAndReturn(run func(context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, daemonSet, opts
func (_m *mockDaemonSetClient) Update(ctx context.Context, daemonSet *appsv1.DaemonSet, opts metav1.UpdateOptions) (*appsv1.DaemonSet, error) {
	ret := _m.Called(ctx, daemonSet, opts)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, daemonSet, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, daemonSet, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, daemonSet, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockDaemonSetClient_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - daemonSet *appsv1.DaemonSet
//   - opts metav1.UpdateOptions
func (_e *mockDaemonSetClient_Expecter) Update(ctx interface{}, daemonSet interface{}, opts interface{}) *mockDaemonSetClient_Update_Call {
	return &mockDaemonSetClient_Update_Call{Call: _e.mock.On("Update", ctx, daemonSet, opts)}
}

func (_c *mockDaemonSetClient_Update_Call) Run(run func(ctx context.Context, daemonSet *appsv1.DaemonSet, opts metav1.UpdateOptions)) *mockDaemonSetClient_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appsv1.DaemonSet), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_Update_Call) Return(_a0 *appsv1.DaemonSet, _a1 error) *mockDaemonSetClient_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDaemonSetClient_Update_Call) RunAndReturn(run func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, daemonSet, opts
func (_m *mockDaemonSetClient) UpdateStatus(ctx context.Context, daemonSet *appsv1.DaemonSet, opts metav1.UpdateOptions) (*appsv1.DaemonSet, error) {
	ret := _m.Called(ctx, daemonSet, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 *appsv1.DaemonSet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) (*appsv1.DaemonSet, error)); ok {
		return rf(ctx, daemonSet, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) *appsv1.DaemonSet); ok {
		r0 = rf(ctx, daemonSet, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*appsv1.DaemonSet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) error); ok {
		r1 = rf(ctx, daemonSet, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDaemonSetClient_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type mockDaemonSetClient_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - daemonSet *appsv1.DaemonSet
//   - opts metav1.UpdateOptions
func (_e *mockDaemonSetClient_Expecter) UpdateStatus(ctx interface{}, daemonSet interface{}, opts interface{}) *mockDaemonSetClient_UpdateStatus_Call {
	return &mockDaemonSetClient_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, daemonSet, opts)}
}

func (_c *mockDaemonSetClient_UpdateStatus_Call) Run(run func(ctx context.Context, daemonSet *appsv1.DaemonSet, opts metav1.UpdateOptions)) *mockDaemonSetClient_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appsv1.DaemonSet), args[2].(metav1.UpdateOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_UpdateStatus_Call) Return(_a0 *appsv1.DaemonSet, _a1 error) *mockDaemonSetClient_UpdateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDaemonSetClient_UpdateStatus_Call) RunAndReturn(run func(context.Context, *appsv1.DaemonSet, metav1.UpdateOptions) (*appsv1.DaemonSet, error)) *mockDaemonSetClient_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *mockDaemonSetClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
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

// mockDaemonSetClient_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type mockDaemonSetClient_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - opts metav1.ListOptions
func (_e *mockDaemonSetClient_Expecter) Watch(ctx interface{}, opts interface{}) *mockDaemonSetClient_Watch_Call {
	return &mockDaemonSetClient_Watch_Call{Call: _e.mock.On("Watch", ctx, opts)}
}

func (_c *mockDaemonSetClient_Watch_Call) Run(run func(ctx context.Context, opts metav1.ListOptions)) *mockDaemonSetClient_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(metav1.ListOptions))
	})
	return _c
}

func (_c *mockDaemonSetClient_Watch_Call) Return(_a0 watch.Interface, _a1 error) *mockDaemonSetClient_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDaemonSetClient_Watch_Call) RunAndReturn(run func(context.Context, metav1.ListOptions) (watch.Interface, error)) *mockDaemonSetClient_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDaemonSetClient creates a new instance of mockDaemonSetClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDaemonSetClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDaemonSetClient {
	mock := &mockDaemonSetClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
