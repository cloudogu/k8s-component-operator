// Code generated by mockery v2.42.1. DO NOT EDIT.

package client

import (
	action "helm.sh/helm/v3/pkg/action"
	chart "helm.sh/helm/v3/pkg/chart"

	context "context"

	mock "github.com/stretchr/testify/mock"

	release "helm.sh/helm/v3/pkg/release"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

type MockClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockClient) EXPECT() *MockClient_Expecter {
	return &MockClient_Expecter{mock: &_m.Mock}
}

// GetChart provides a mock function with given fields: spec
func (_m *MockClient) GetChart(spec *ChartSpec) (*chart.Chart, string, error) {
	ret := _m.Called(spec)

	if len(ret) == 0 {
		panic("no return value specified for GetChart")
	}

	var r0 *chart.Chart
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(*ChartSpec) (*chart.Chart, string, error)); ok {
		return rf(spec)
	}
	if rf, ok := ret.Get(0).(func(*ChartSpec) *chart.Chart); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chart.Chart)
		}
	}

	if rf, ok := ret.Get(1).(func(*ChartSpec) string); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(*ChartSpec) error); ok {
		r2 = rf(spec)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockClient_GetChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChart'
type MockClient_GetChart_Call struct {
	*mock.Call
}

// GetChart is a helper method to define mock.On call
//   - spec *ChartSpec
func (_e *MockClient_Expecter) GetChart(spec interface{}) *MockClient_GetChart_Call {
	return &MockClient_GetChart_Call{Call: _e.mock.On("GetChart", spec)}
}

func (_c *MockClient_GetChart_Call) Run(run func(spec *ChartSpec)) *MockClient_GetChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_GetChart_Call) Return(_a0 *chart.Chart, _a1 string, _a2 error) *MockClient_GetChart_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockClient_GetChart_Call) RunAndReturn(run func(*ChartSpec) (*chart.Chart, string, error)) *MockClient_GetChart_Call {
	_c.Call.Return(run)
	return _c
}

// GetChartSpecValues provides a mock function with given fields: spec
func (_m *MockClient) GetChartSpecValues(spec *ChartSpec) (map[string]interface{}, error) {
	ret := _m.Called(spec)

	if len(ret) == 0 {
		panic("no return value specified for GetChartSpecValues")
	}

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(*ChartSpec) (map[string]interface{}, error)); ok {
		return rf(spec)
	}
	if rf, ok := ret.Get(0).(func(*ChartSpec) map[string]interface{}); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(*ChartSpec) error); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_GetChartSpecValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChartSpecValues'
type MockClient_GetChartSpecValues_Call struct {
	*mock.Call
}

// GetChartSpecValues is a helper method to define mock.On call
//   - spec *ChartSpec
func (_e *MockClient_Expecter) GetChartSpecValues(spec interface{}) *MockClient_GetChartSpecValues_Call {
	return &MockClient_GetChartSpecValues_Call{Call: _e.mock.On("GetChartSpecValues", spec)}
}

func (_c *MockClient_GetChartSpecValues_Call) Run(run func(spec *ChartSpec)) *MockClient_GetChartSpecValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_GetChartSpecValues_Call) Return(_a0 map[string]interface{}, _a1 error) *MockClient_GetChartSpecValues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_GetChartSpecValues_Call) RunAndReturn(run func(*ChartSpec) (map[string]interface{}, error)) *MockClient_GetChartSpecValues_Call {
	_c.Call.Return(run)
	return _c
}

// GetRelease provides a mock function with given fields: name
func (_m *MockClient) GetRelease(name string) (*release.Release, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetRelease")
	}

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

// MockClient_GetRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRelease'
type MockClient_GetRelease_Call struct {
	*mock.Call
}

// GetRelease is a helper method to define mock.On call
//   - name string
func (_e *MockClient_Expecter) GetRelease(name interface{}) *MockClient_GetRelease_Call {
	return &MockClient_GetRelease_Call{Call: _e.mock.On("GetRelease", name)}
}

func (_c *MockClient_GetRelease_Call) Run(run func(name string)) *MockClient_GetRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockClient_GetRelease_Call) Return(_a0 *release.Release, _a1 error) *MockClient_GetRelease_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_GetRelease_Call) RunAndReturn(run func(string) (*release.Release, error)) *MockClient_GetRelease_Call {
	_c.Call.Return(run)
	return _c
}

// GetReleaseValues provides a mock function with given fields: name, allValues
func (_m *MockClient) GetReleaseValues(name string, allValues bool) (map[string]interface{}, error) {
	ret := _m.Called(name, allValues)

	if len(ret) == 0 {
		panic("no return value specified for GetReleaseValues")
	}

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

// MockClient_GetReleaseValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReleaseValues'
type MockClient_GetReleaseValues_Call struct {
	*mock.Call
}

// GetReleaseValues is a helper method to define mock.On call
//   - name string
//   - allValues bool
func (_e *MockClient_Expecter) GetReleaseValues(name interface{}, allValues interface{}) *MockClient_GetReleaseValues_Call {
	return &MockClient_GetReleaseValues_Call{Call: _e.mock.On("GetReleaseValues", name, allValues)}
}

func (_c *MockClient_GetReleaseValues_Call) Run(run func(name string, allValues bool)) *MockClient_GetReleaseValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(bool))
	})
	return _c
}

func (_c *MockClient_GetReleaseValues_Call) Return(_a0 map[string]interface{}, _a1 error) *MockClient_GetReleaseValues_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_GetReleaseValues_Call) RunAndReturn(run func(string, bool) (map[string]interface{}, error)) *MockClient_GetReleaseValues_Call {
	_c.Call.Return(run)
	return _c
}

// InstallChart provides a mock function with given fields: ctx, spec
func (_m *MockClient) InstallChart(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	ret := _m.Called(ctx, spec)

	if len(ret) == 0 {
		panic("no return value specified for InstallChart")
	}

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ChartSpec) (*release.Release, error)); ok {
		return rf(ctx, spec)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ChartSpec) *release.Release); ok {
		r0 = rf(ctx, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ChartSpec) error); ok {
		r1 = rf(ctx, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_InstallChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstallChart'
type MockClient_InstallChart_Call struct {
	*mock.Call
}

// InstallChart is a helper method to define mock.On call
//   - ctx context.Context
//   - spec *ChartSpec
func (_e *MockClient_Expecter) InstallChart(ctx interface{}, spec interface{}) *MockClient_InstallChart_Call {
	return &MockClient_InstallChart_Call{Call: _e.mock.On("InstallChart", ctx, spec)}
}

func (_c *MockClient_InstallChart_Call) Run(run func(ctx context.Context, spec *ChartSpec)) *MockClient_InstallChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_InstallChart_Call) Return(_a0 *release.Release, _a1 error) *MockClient_InstallChart_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_InstallChart_Call) RunAndReturn(run func(context.Context, *ChartSpec) (*release.Release, error)) *MockClient_InstallChart_Call {
	_c.Call.Return(run)
	return _c
}

// InstallOrUpgradeChart provides a mock function with given fields: ctx, spec
func (_m *MockClient) InstallOrUpgradeChart(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	ret := _m.Called(ctx, spec)

	if len(ret) == 0 {
		panic("no return value specified for InstallOrUpgradeChart")
	}

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ChartSpec) (*release.Release, error)); ok {
		return rf(ctx, spec)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ChartSpec) *release.Release); ok {
		r0 = rf(ctx, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ChartSpec) error); ok {
		r1 = rf(ctx, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_InstallOrUpgradeChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstallOrUpgradeChart'
type MockClient_InstallOrUpgradeChart_Call struct {
	*mock.Call
}

// InstallOrUpgradeChart is a helper method to define mock.On call
//   - ctx context.Context
//   - spec *ChartSpec
func (_e *MockClient_Expecter) InstallOrUpgradeChart(ctx interface{}, spec interface{}) *MockClient_InstallOrUpgradeChart_Call {
	return &MockClient_InstallOrUpgradeChart_Call{Call: _e.mock.On("InstallOrUpgradeChart", ctx, spec)}
}

func (_c *MockClient_InstallOrUpgradeChart_Call) Run(run func(ctx context.Context, spec *ChartSpec)) *MockClient_InstallOrUpgradeChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_InstallOrUpgradeChart_Call) Return(_a0 *release.Release, _a1 error) *MockClient_InstallOrUpgradeChart_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_InstallOrUpgradeChart_Call) RunAndReturn(run func(context.Context, *ChartSpec) (*release.Release, error)) *MockClient_InstallOrUpgradeChart_Call {
	_c.Call.Return(run)
	return _c
}

// ListDeployedReleases provides a mock function with given fields:
func (_m *MockClient) ListDeployedReleases() ([]*release.Release, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListDeployedReleases")
	}

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

// MockClient_ListDeployedReleases_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListDeployedReleases'
type MockClient_ListDeployedReleases_Call struct {
	*mock.Call
}

// ListDeployedReleases is a helper method to define mock.On call
func (_e *MockClient_Expecter) ListDeployedReleases() *MockClient_ListDeployedReleases_Call {
	return &MockClient_ListDeployedReleases_Call{Call: _e.mock.On("ListDeployedReleases")}
}

func (_c *MockClient_ListDeployedReleases_Call) Run(run func()) *MockClient_ListDeployedReleases_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockClient_ListDeployedReleases_Call) Return(_a0 []*release.Release, _a1 error) *MockClient_ListDeployedReleases_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_ListDeployedReleases_Call) RunAndReturn(run func() ([]*release.Release, error)) *MockClient_ListDeployedReleases_Call {
	_c.Call.Return(run)
	return _c
}

// ListReleasesByStateMask provides a mock function with given fields: _a0
func (_m *MockClient) ListReleasesByStateMask(_a0 action.ListStates) ([]*release.Release, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ListReleasesByStateMask")
	}

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

// MockClient_ListReleasesByStateMask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListReleasesByStateMask'
type MockClient_ListReleasesByStateMask_Call struct {
	*mock.Call
}

// ListReleasesByStateMask is a helper method to define mock.On call
//   - _a0 action.ListStates
func (_e *MockClient_Expecter) ListReleasesByStateMask(_a0 interface{}) *MockClient_ListReleasesByStateMask_Call {
	return &MockClient_ListReleasesByStateMask_Call{Call: _e.mock.On("ListReleasesByStateMask", _a0)}
}

func (_c *MockClient_ListReleasesByStateMask_Call) Run(run func(_a0 action.ListStates)) *MockClient_ListReleasesByStateMask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(action.ListStates))
	})
	return _c
}

func (_c *MockClient_ListReleasesByStateMask_Call) Return(_a0 []*release.Release, _a1 error) *MockClient_ListReleasesByStateMask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_ListReleasesByStateMask_Call) RunAndReturn(run func(action.ListStates) ([]*release.Release, error)) *MockClient_ListReleasesByStateMask_Call {
	_c.Call.Return(run)
	return _c
}

// RollbackRelease provides a mock function with given fields: spec
func (_m *MockClient) RollbackRelease(spec *ChartSpec) error {
	ret := _m.Called(spec)

	if len(ret) == 0 {
		panic("no return value specified for RollbackRelease")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ChartSpec) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockClient_RollbackRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RollbackRelease'
type MockClient_RollbackRelease_Call struct {
	*mock.Call
}

// RollbackRelease is a helper method to define mock.On call
//   - spec *ChartSpec
func (_e *MockClient_Expecter) RollbackRelease(spec interface{}) *MockClient_RollbackRelease_Call {
	return &MockClient_RollbackRelease_Call{Call: _e.mock.On("RollbackRelease", spec)}
}

func (_c *MockClient_RollbackRelease_Call) Run(run func(spec *ChartSpec)) *MockClient_RollbackRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_RollbackRelease_Call) Return(_a0 error) *MockClient_RollbackRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockClient_RollbackRelease_Call) RunAndReturn(run func(*ChartSpec) error) *MockClient_RollbackRelease_Call {
	_c.Call.Return(run)
	return _c
}

// Tags provides a mock function with given fields: ref
func (_m *MockClient) Tags(ref string) ([]string, error) {
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

// MockClient_Tags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Tags'
type MockClient_Tags_Call struct {
	*mock.Call
}

// Tags is a helper method to define mock.On call
//   - ref string
func (_e *MockClient_Expecter) Tags(ref interface{}) *MockClient_Tags_Call {
	return &MockClient_Tags_Call{Call: _e.mock.On("Tags", ref)}
}

func (_c *MockClient_Tags_Call) Run(run func(ref string)) *MockClient_Tags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockClient_Tags_Call) Return(_a0 []string, _a1 error) *MockClient_Tags_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_Tags_Call) RunAndReturn(run func(string) ([]string, error)) *MockClient_Tags_Call {
	_c.Call.Return(run)
	return _c
}

// UninstallRelease provides a mock function with given fields: spec
func (_m *MockClient) UninstallRelease(spec *ChartSpec) error {
	ret := _m.Called(spec)

	if len(ret) == 0 {
		panic("no return value specified for UninstallRelease")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ChartSpec) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockClient_UninstallRelease_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UninstallRelease'
type MockClient_UninstallRelease_Call struct {
	*mock.Call
}

// UninstallRelease is a helper method to define mock.On call
//   - spec *ChartSpec
func (_e *MockClient_Expecter) UninstallRelease(spec interface{}) *MockClient_UninstallRelease_Call {
	return &MockClient_UninstallRelease_Call{Call: _e.mock.On("UninstallRelease", spec)}
}

func (_c *MockClient_UninstallRelease_Call) Run(run func(spec *ChartSpec)) *MockClient_UninstallRelease_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_UninstallRelease_Call) Return(_a0 error) *MockClient_UninstallRelease_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockClient_UninstallRelease_Call) RunAndReturn(run func(*ChartSpec) error) *MockClient_UninstallRelease_Call {
	_c.Call.Return(run)
	return _c
}

// UninstallReleaseByName provides a mock function with given fields: name
func (_m *MockClient) UninstallReleaseByName(name string) error {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for UninstallReleaseByName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockClient_UninstallReleaseByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UninstallReleaseByName'
type MockClient_UninstallReleaseByName_Call struct {
	*mock.Call
}

// UninstallReleaseByName is a helper method to define mock.On call
//   - name string
func (_e *MockClient_Expecter) UninstallReleaseByName(name interface{}) *MockClient_UninstallReleaseByName_Call {
	return &MockClient_UninstallReleaseByName_Call{Call: _e.mock.On("UninstallReleaseByName", name)}
}

func (_c *MockClient_UninstallReleaseByName_Call) Run(run func(name string)) *MockClient_UninstallReleaseByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockClient_UninstallReleaseByName_Call) Return(_a0 error) *MockClient_UninstallReleaseByName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockClient_UninstallReleaseByName_Call) RunAndReturn(run func(string) error) *MockClient_UninstallReleaseByName_Call {
	_c.Call.Return(run)
	return _c
}

// UpgradeChart provides a mock function with given fields: ctx, spec
func (_m *MockClient) UpgradeChart(ctx context.Context, spec *ChartSpec) (*release.Release, error) {
	ret := _m.Called(ctx, spec)

	if len(ret) == 0 {
		panic("no return value specified for UpgradeChart")
	}

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ChartSpec) (*release.Release, error)); ok {
		return rf(ctx, spec)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ChartSpec) *release.Release); ok {
		r0 = rf(ctx, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ChartSpec) error); ok {
		r1 = rf(ctx, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockClient_UpgradeChart_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpgradeChart'
type MockClient_UpgradeChart_Call struct {
	*mock.Call
}

// UpgradeChart is a helper method to define mock.On call
//   - ctx context.Context
//   - spec *ChartSpec
func (_e *MockClient_Expecter) UpgradeChart(ctx interface{}, spec interface{}) *MockClient_UpgradeChart_Call {
	return &MockClient_UpgradeChart_Call{Call: _e.mock.On("UpgradeChart", ctx, spec)}
}

func (_c *MockClient_UpgradeChart_Call) Run(run func(ctx context.Context, spec *ChartSpec)) *MockClient_UpgradeChart_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ChartSpec))
	})
	return _c
}

func (_c *MockClient_UpgradeChart_Call) Return(_a0 *release.Release, _a1 error) *MockClient_UpgradeChart_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockClient_UpgradeChart_Call) RunAndReturn(run func(context.Context, *ChartSpec) (*release.Release, error)) *MockClient_UpgradeChart_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
