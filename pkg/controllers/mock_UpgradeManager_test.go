// Code generated by mockery v2.20.0. DO NOT EDIT.

package controllers

import (
	context "context"

	v1 "github.com/cloudogu/k8s-component-operator/api/v1"
	mock "github.com/stretchr/testify/mock"
)

// MockUpgradeManager is an autogenerated mock type for the UpgradeManager type
type MockUpgradeManager struct {
	mock.Mock
}

type MockUpgradeManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpgradeManager) EXPECT() *MockUpgradeManager_Expecter {
	return &MockUpgradeManager_Expecter{mock: &_m.Mock}
}

// Upgrade provides a mock function with given fields: ctx, component
func (_m *MockUpgradeManager) Upgrade(ctx context.Context, component *v1.Component) error {
	ret := _m.Called(ctx, component)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Component) error); ok {
		r0 = rf(ctx, component)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockUpgradeManager_Upgrade_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upgrade'
type MockUpgradeManager_Upgrade_Call struct {
	*mock.Call
}

// Upgrade is a helper method to define mock.On call
//   - ctx context.Context
//   - component *v1.Component
func (_e *MockUpgradeManager_Expecter) Upgrade(ctx interface{}, component interface{}) *MockUpgradeManager_Upgrade_Call {
	return &MockUpgradeManager_Upgrade_Call{Call: _e.mock.On("Upgrade", ctx, component)}
}

func (_c *MockUpgradeManager_Upgrade_Call) Run(run func(ctx context.Context, component *v1.Component)) *MockUpgradeManager_Upgrade_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Component))
	})
	return _c
}

func (_c *MockUpgradeManager_Upgrade_Call) Return(_a0 error) *MockUpgradeManager_Upgrade_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUpgradeManager_Upgrade_Call) RunAndReturn(run func(context.Context, *v1.Component) error) *MockUpgradeManager_Upgrade_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockUpgradeManager interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockUpgradeManager creates a new instance of MockUpgradeManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUpgradeManager(t mockConstructorTestingTNewMockUpgradeManager) *MockUpgradeManager {
	mock := &MockUpgradeManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
