// Code generated by mockery v2.20.0. DO NOT EDIT.

package client

import (
	action "helm.sh/helm/v3/pkg/action"
	chart "helm.sh/helm/v3/pkg/chart"

	context "context"

	mock "github.com/stretchr/testify/mock"

	release "helm.sh/helm/v3/pkg/release"
)

// mockUpgradeAction is an autogenerated mock type for the upgradeAction type
type mockUpgradeAction struct {
	mock.Mock
}

type mockUpgradeAction_Expecter struct {
	mock *mock.Mock
}

func (_m *mockUpgradeAction) EXPECT() *mockUpgradeAction_Expecter {
	return &mockUpgradeAction_Expecter{mock: &_m.Mock}
}

// raw provides a mock function with given fields:
func (_m *mockUpgradeAction) raw() *action.Upgrade {
	ret := _m.Called()

	var r0 *action.Upgrade
	if rf, ok := ret.Get(0).(func() *action.Upgrade); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*action.Upgrade)
		}
	}

	return r0
}

// mockUpgradeAction_raw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'raw'
type mockUpgradeAction_raw_Call struct {
	*mock.Call
}

// raw is a helper method to define mock.On call
func (_e *mockUpgradeAction_Expecter) raw() *mockUpgradeAction_raw_Call {
	return &mockUpgradeAction_raw_Call{Call: _e.mock.On("raw")}
}

func (_c *mockUpgradeAction_raw_Call) Run(run func()) *mockUpgradeAction_raw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockUpgradeAction_raw_Call) Return(_a0 *action.Upgrade) *mockUpgradeAction_raw_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockUpgradeAction_raw_Call) RunAndReturn(run func() *action.Upgrade) *mockUpgradeAction_raw_Call {
	_c.Call.Return(run)
	return _c
}

// upgrade provides a mock function with given fields: ctx, releaseName, _a2, values
func (_m *mockUpgradeAction) upgrade(ctx context.Context, releaseName string, _a2 *chart.Chart, values map[string]interface{}) (*release.Release, error) {
	ret := _m.Called(ctx, releaseName, _a2, values)

	var r0 *release.Release
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *chart.Chart, map[string]interface{}) (*release.Release, error)); ok {
		return rf(ctx, releaseName, _a2, values)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *chart.Chart, map[string]interface{}) *release.Release); ok {
		r0 = rf(ctx, releaseName, _a2, values)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*release.Release)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *chart.Chart, map[string]interface{}) error); ok {
		r1 = rf(ctx, releaseName, _a2, values)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockUpgradeAction_upgrade_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'upgrade'
type mockUpgradeAction_upgrade_Call struct {
	*mock.Call
}

// upgrade is a helper method to define mock.On call
//   - ctx context.Context
//   - releaseName string
//   - _a2 *chart.Chart
//   - values map[string]interface{}
func (_e *mockUpgradeAction_Expecter) upgrade(ctx interface{}, releaseName interface{}, _a2 interface{}, values interface{}) *mockUpgradeAction_upgrade_Call {
	return &mockUpgradeAction_upgrade_Call{Call: _e.mock.On("upgrade", ctx, releaseName, _a2, values)}
}

func (_c *mockUpgradeAction_upgrade_Call) Run(run func(ctx context.Context, releaseName string, _a2 *chart.Chart, values map[string]interface{})) *mockUpgradeAction_upgrade_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*chart.Chart), args[3].(map[string]interface{}))
	})
	return _c
}

func (_c *mockUpgradeAction_upgrade_Call) Return(_a0 *release.Release, _a1 error) *mockUpgradeAction_upgrade_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockUpgradeAction_upgrade_Call) RunAndReturn(run func(context.Context, string, *chart.Chart, map[string]interface{}) (*release.Release, error)) *mockUpgradeAction_upgrade_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockUpgradeAction interface {
	mock.TestingT
	Cleanup(func())
}

// newMockUpgradeAction creates a new instance of mockUpgradeAction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockUpgradeAction(t mockConstructorTestingTnewMockUpgradeAction) *mockUpgradeAction {
	mock := &mockUpgradeAction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
