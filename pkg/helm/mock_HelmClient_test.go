// Code generated by mockery v2.20.0. DO NOT EDIT.

package helm

import (
	action "helm.sh/helm/v3/pkg/action"
	chart "helm.sh/helm/v3/pkg/chart"

	client "github.com/cloudogu/k8s-component-operator/pkg/helm/client"

	context "context"

	mock "github.com/stretchr/testify/mock"

	release "helm.sh/helm/v3/pkg/release"
)

// MockHelmClient is an autogenerated mock type for the HelmClient type
type MockHelmClient struct {
	mock.Mock
}

type MockHelmClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHelmClient) EXPECT() *MockHelmClient_Expecter {
	return &MockHelmClient_Expecter{mock: &_m.Mock}
}

// GetChart provides a mock function with given fields: spec
func (_m *MockHelmClient) GetChart(spec *client.ChartSpec) (*chart.Chart, string, error) {
	ret := _m.Called(spec)

	var r0 *chart.Chart
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(*client.ChartSpec) (*chart.Chart, string, error)); ok {
		return rf(spec)
	}
	if rf, ok := ret.Get(0).(func(*client.ChartSpec) *chart.Chart); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chart.Chart)
		}
	}

	if rf, ok := ret.Get(1).(func(*client.ChartSpec) string); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(*client.ChartSpec) error); ok {
		r2 = rf(spec)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockHelmClient_GetChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChart'
type MockHelmClient_GetChart_Call struct {
	*mock.Call
}

// GetChart is a helper method to define mock.On call
//   - spec *client.ChartSpec
func (_e *MockHelmClient_Expecter) GetChart(spec interface{}) *MockHelmClient_GetChart_Call {
	return &MockHelmClient_GetChart_Call{Call: _e.mock.On("GetChart", spec)}
}

func (_c *MockHelmClient_GetChart_Call) Run(run func(spec *client.ChartSpec)) *MockHelmClient_GetChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*client.ChartSpec))
	})
	return _c
}

func (_c *MockHelmClient_GetChart_Call) Return(_a0 *chart.Chart, _a1 string, _a2 error) *MockHelmClient_GetChart_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockHelmClient_GetChart_Call) RunAndReturn(run func(*client.ChartSpec) (*chart.Chart, string, error)) *MockHelmClient_GetChart_Call {
	_c.Call.Return(run)
	return _c
}

// GetRelease provides a mock function with given fields: name
func (_m *MockHelmClient) GetRelease(name string) (*release.Release, error) {
	ret := _m.Called(name)

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*release.Release, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *release.Release); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_GetRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRelease'
type MockHelmClient_GetRelease_Call struct {
	*mock.Call
}

// GetRelease is a helper method to define mock.On call
//   - name string
func (_e *MockHelmClient_Expecter) GetRelease(name interface{}) *MockHelmClient_GetRelease_Call {
	return &MockHelmClient_GetRelease_Call{Call: _e.mock.On("GetRelease", name)}
}

func (_c *MockHelmClient_GetRelease_Call) Run(run func(name string)) *MockHelmClient_GetRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockHelmClient_GetRelease_Call) Return(_a0 *release.Release, _a1 error) *MockHelmClient_GetRelease_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_GetRelease_Call) RunAndReturn(run func(string) (*release.Release, error)) *MockHelmClient_GetRelease_Call {
	_c.Call.Return(run)
	return _c
}

// GetReleaseValues provides a mock function with given fields: name, allValues
func (_m *MockHelmClient) GetReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	ret := _m.Called(name, allValues)

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(string, bool) (map[string]interface{}, error)); ok {
		return rf(name, allValues)
	}
	if rf, ok := ret.Get(0).(func(string, bool) map[string]interface{}); ok {
		r0 = rf(name, allValues)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(name, allValues)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_GetReleaseValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReleaseValues'
type MockHelmClient_GetReleaseValues_Call struct {
	*mock.Call
}

// GetReleaseValues is a helper method to define mock.On call
//   - name string
//   - allValues bool
func (_e *MockHelmClient_Expecter) GetReleaseValues(name interface{}, allValues interface{}) *MockHelmClient_GetReleaseValues_Call {
	return &MockHelmClient_GetReleaseValues_Call{Call: _e.mock.On("GetReleaseValues", name, allValues)}
}

func (_c *MockHelmClient_GetReleaseValues_Call) Run(run func(name string, allValues bool)) *MockHelmClient_GetReleaseValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(bool))
	})
	return _c
}

func (_c *MockHelmClient_GetReleaseValues_Call) Return(_a0 map[string]interface{}, _a1 error) *MockHelmClient_GetReleaseValues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_GetReleaseValues_Call) RunAndReturn(run func(string, bool) (map[string]interface{}, error)) *MockHelmClient_GetReleaseValues_Call {
	_c.Call.Return(run)
	return _c
}

// InstallChart provides a mock function with given fields: ctx, spec
func (_m *MockHelmClient) InstallChart(ctx context.Context, spec *client.ChartSpec) (*release.Release, error) {
	ret := _m.Called(ctx, spec)

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChartSpec) (*release.Release, error)); ok {
		return rf(ctx, spec)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChartSpec) *release.Release); ok {
		r0 = rf(ctx, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.ChartSpec) error); ok {
		r1 = rf(ctx, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_InstallChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstallChart'
type MockHelmClient_InstallChart_Call struct {
	*mock.Call
}

// InstallChart is a helper method to define mock.On call
//   - ctx context.Context
//   - spec *client.ChartSpec
func (_e *MockHelmClient_Expecter) InstallChart(ctx interface{}, spec interface{}) *MockHelmClient_InstallChart_Call {
	return &MockHelmClient_InstallChart_Call{Call: _e.mock.On("InstallChart", ctx, spec)}
}

func (_c *MockHelmClient_InstallChart_Call) Run(run func(ctx context.Context, spec *client.ChartSpec)) *MockHelmClient_InstallChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.ChartSpec))
	})
	return _c
}

func (_c *MockHelmClient_InstallChart_Call) Return(_a0 *release.Release, _a1 error) *MockHelmClient_InstallChart_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_InstallChart_Call) RunAndReturn(run func(context.Context, *client.ChartSpec) (*release.Release, error)) *MockHelmClient_InstallChart_Call {
	_c.Call.Return(run)
	return _c
}

// InstallOrUpgradeChart provides a mock function with given fields: ctx, spec
func (_m *MockHelmClient) InstallOrUpgradeChart(ctx context.Context, spec *client.ChartSpec) (*release.Release, error) {
	ret := _m.Called(ctx, spec)

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChartSpec) (*release.Release, error)); ok {
		return rf(ctx, spec)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChartSpec) *release.Release); ok {
		r0 = rf(ctx, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.ChartSpec) error); ok {
		r1 = rf(ctx, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_InstallOrUpgradeChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstallOrUpgradeChart'
type MockHelmClient_InstallOrUpgradeChart_Call struct {
	*mock.Call
}

// InstallOrUpgradeChart is a helper method to define mock.On call
//   - ctx context.Context
//   - spec *client.ChartSpec
func (_e *MockHelmClient_Expecter) InstallOrUpgradeChart(ctx interface{}, spec interface{}) *MockHelmClient_InstallOrUpgradeChart_Call {
	return &MockHelmClient_InstallOrUpgradeChart_Call{Call: _e.mock.On("InstallOrUpgradeChart", ctx, spec)}
}

func (_c *MockHelmClient_InstallOrUpgradeChart_Call) Run(run func(ctx context.Context, spec *client.ChartSpec)) *MockHelmClient_InstallOrUpgradeChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.ChartSpec))
	})
	return _c
}

func (_c *MockHelmClient_InstallOrUpgradeChart_Call) Return(_a0 *release.Release, _a1 error) *MockHelmClient_InstallOrUpgradeChart_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_InstallOrUpgradeChart_Call) RunAndReturn(run func(context.Context, *client.ChartSpec) (*release.Release, error)) *MockHelmClient_InstallOrUpgradeChart_Call {
	_c.Call.Return(run)
	return _c
}

// ListDeployedReleases provides a mock function with given fields:
func (_m *MockHelmClient) ListDeployedReleases() ([]*release.Release, error) {
	ret := _m.Called()

	var r0 []*release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*release.Release, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*release.Release); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_ListDeployedReleases_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListDeployedReleases'
type MockHelmClient_ListDeployedReleases_Call struct {
	*mock.Call
}

// ListDeployedReleases is a helper method to define mock.On call
func (_e *MockHelmClient_Expecter) ListDeployedReleases() *MockHelmClient_ListDeployedReleases_Call {
	return &MockHelmClient_ListDeployedReleases_Call{Call: _e.mock.On("ListDeployedReleases")}
}

func (_c *MockHelmClient_ListDeployedReleases_Call) Run(run func()) *MockHelmClient_ListDeployedReleases_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockHelmClient_ListDeployedReleases_Call) Return(_a0 []*release.Release, _a1 error) *MockHelmClient_ListDeployedReleases_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_ListDeployedReleases_Call) RunAndReturn(run func() ([]*release.Release, error)) *MockHelmClient_ListDeployedReleases_Call {
	_c.Call.Return(run)
	return _c
}

// ListReleasesByStateMask provides a mock function with given fields: _a0
func (_m *MockHelmClient) ListReleasesByStateMask(_a0 action.ListStates) ([]*release.Release, error) {
	ret := _m.Called(_a0)

	var r0 []*release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(action.ListStates) ([]*release.Release, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(action.ListStates) []*release.Release); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(action.ListStates) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_ListReleasesByStateMask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReleasesByStateMask'
type MockHelmClient_ListReleasesByStateMask_Call struct {
	*mock.Call
}

// ListReleasesByStateMask is a helper method to define mock.On call
//   - _a0 action.ListStates
func (_e *MockHelmClient_Expecter) ListReleasesByStateMask(_a0 interface{}) *MockHelmClient_ListReleasesByStateMask_Call {
	return &MockHelmClient_ListReleasesByStateMask_Call{Call: _e.mock.On("ListReleasesByStateMask", _a0)}
}

func (_c *MockHelmClient_ListReleasesByStateMask_Call) Run(run func(_a0 action.ListStates)) *MockHelmClient_ListReleasesByStateMask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(action.ListStates))
	})
	return _c
}

func (_c *MockHelmClient_ListReleasesByStateMask_Call) Return(_a0 []*release.Release, _a1 error) *MockHelmClient_ListReleasesByStateMask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_ListReleasesByStateMask_Call) RunAndReturn(run func(action.ListStates) ([]*release.Release, error)) *MockHelmClient_ListReleasesByStateMask_Call {
	_c.Call.Return(run)
	return _c
}

// RollbackRelease provides a mock function with given fields: spec
func (_m *MockHelmClient) RollbackRelease(spec *client.ChartSpec) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*client.ChartSpec) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockHelmClient_RollbackRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RollbackRelease'
type MockHelmClient_RollbackRelease_Call struct {
	*mock.Call
}

// RollbackRelease is a helper method to define mock.On call
//   - spec *client.ChartSpec
func (_e *MockHelmClient_Expecter) RollbackRelease(spec interface{}) *MockHelmClient_RollbackRelease_Call {
	return &MockHelmClient_RollbackRelease_Call{Call: _e.mock.On("RollbackRelease", spec)}
}

func (_c *MockHelmClient_RollbackRelease_Call) Run(run func(spec *client.ChartSpec)) *MockHelmClient_RollbackRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*client.ChartSpec))
	})
	return _c
}

func (_c *MockHelmClient_RollbackRelease_Call) Return(_a0 error) *MockHelmClient_RollbackRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHelmClient_RollbackRelease_Call) RunAndReturn(run func(*client.ChartSpec) error) *MockHelmClient_RollbackRelease_Call {
	_c.Call.Return(run)
	return _c
}

// Tags provides a mock function with given fields: ref
func (_m *MockHelmClient) Tags(ref string) ([]string, error) {
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

// MockHelmClient_Tags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Tags'
type MockHelmClient_Tags_Call struct {
	*mock.Call
}

// Tags is a helper method to define mock.On call
//   - ref string
func (_e *MockHelmClient_Expecter) Tags(ref interface{}) *MockHelmClient_Tags_Call {
	return &MockHelmClient_Tags_Call{Call: _e.mock.On("Tags", ref)}
}

func (_c *MockHelmClient_Tags_Call) Run(run func(ref string)) *MockHelmClient_Tags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockHelmClient_Tags_Call) Return(_a0 []string, _a1 error) *MockHelmClient_Tags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_Tags_Call) RunAndReturn(run func(string) ([]string, error)) *MockHelmClient_Tags_Call {
	_c.Call.Return(run)
	return _c
}

// UninstallRelease provides a mock function with given fields: spec
func (_m *MockHelmClient) UninstallRelease(spec *client.ChartSpec) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*client.ChartSpec) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockHelmClient_UninstallRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UninstallRelease'
type MockHelmClient_UninstallRelease_Call struct {
	*mock.Call
}

// UninstallRelease is a helper method to define mock.On call
//   - spec *client.ChartSpec
func (_e *MockHelmClient_Expecter) UninstallRelease(spec interface{}) *MockHelmClient_UninstallRelease_Call {
	return &MockHelmClient_UninstallRelease_Call{Call: _e.mock.On("UninstallRelease", spec)}
}

func (_c *MockHelmClient_UninstallRelease_Call) Run(run func(spec *client.ChartSpec)) *MockHelmClient_UninstallRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*client.ChartSpec))
	})
	return _c
}

func (_c *MockHelmClient_UninstallRelease_Call) Return(_a0 error) *MockHelmClient_UninstallRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHelmClient_UninstallRelease_Call) RunAndReturn(run func(*client.ChartSpec) error) *MockHelmClient_UninstallRelease_Call {
	_c.Call.Return(run)
	return _c
}

// UninstallReleaseByName provides a mock function with given fields: name
func (_m *MockHelmClient) UninstallReleaseByName(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockHelmClient_UninstallReleaseByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UninstallReleaseByName'
type MockHelmClient_UninstallReleaseByName_Call struct {
	*mock.Call
}

// UninstallReleaseByName is a helper method to define mock.On call
//   - name string
func (_e *MockHelmClient_Expecter) UninstallReleaseByName(name interface{}) *MockHelmClient_UninstallReleaseByName_Call {
	return &MockHelmClient_UninstallReleaseByName_Call{Call: _e.mock.On("UninstallReleaseByName", name)}
}

func (_c *MockHelmClient_UninstallReleaseByName_Call) Run(run func(name string)) *MockHelmClient_UninstallReleaseByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockHelmClient_UninstallReleaseByName_Call) Return(_a0 error) *MockHelmClient_UninstallReleaseByName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHelmClient_UninstallReleaseByName_Call) RunAndReturn(run func(string) error) *MockHelmClient_UninstallReleaseByName_Call {
	_c.Call.Return(run)
	return _c
}

// UpgradeChart provides a mock function with given fields: ctx, spec
func (_m *MockHelmClient) UpgradeChart(ctx context.Context, spec *client.ChartSpec) (*release.Release, error) {
	ret := _m.Called(ctx, spec)

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChartSpec) (*release.Release, error)); ok {
		return rf(ctx, spec)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChartSpec) *release.Release); ok {
		r0 = rf(ctx, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.ChartSpec) error); ok {
		r1 = rf(ctx, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHelmClient_UpgradeChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpgradeChart'
type MockHelmClient_UpgradeChart_Call struct {
	*mock.Call
}

// UpgradeChart is a helper method to define mock.On call
//   - ctx context.Context
//   - spec *client.ChartSpec
func (_e *MockHelmClient_Expecter) UpgradeChart(ctx interface{}, spec interface{}) *MockHelmClient_UpgradeChart_Call {
	return &MockHelmClient_UpgradeChart_Call{Call: _e.mock.On("UpgradeChart", ctx, spec)}
}

func (_c *MockHelmClient_UpgradeChart_Call) Run(run func(ctx context.Context, spec *client.ChartSpec)) *MockHelmClient_UpgradeChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.ChartSpec))
	})
	return _c
}

func (_c *MockHelmClient_UpgradeChart_Call) Return(_a0 *release.Release, _a1 error) *MockHelmClient_UpgradeChart_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHelmClient_UpgradeChart_Call) RunAndReturn(run func(context.Context, *client.ChartSpec) (*release.Release, error)) *MockHelmClient_UpgradeChart_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockHelmClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockHelmClient creates a new instance of MockHelmClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockHelmClient(t mockConstructorTestingTNewMockHelmClient) *MockHelmClient {
	mock := &MockHelmClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
