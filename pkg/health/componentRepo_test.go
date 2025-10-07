package health

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
)

func Test_defaultComponentRepo_get(t *testing.T) {
	type returnValues struct {
		component *v1.Component
		err       error
	}
	type wantValues struct {
		component *v1.Component
		err       assert.ErrorAssertionFunc
	}
	tests := []struct {
		name    string
		returns returnValues
		wants   wantValues
	}{
		{
			name: "should fail to get component",
			returns: returnValues{
				component: nil,
				err:       assert.AnError,
			},
			wants: wantValues{
				component: nil,
				err: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, assert.AnError, i) &&
						assert.ErrorContains(t, err, fmt.Sprintf("failed to get component %q", testComponentName), i)
				},
			},
		},
		{
			name: "should succeed to get component",
			returns: returnValues{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				err:       nil,
			},
			wants: wantValues{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				err:       assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := newMockComponentClient(t)
			clientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).
				Return(tt.returns.component, tt.returns.err)
			cr := &defaultComponentRepo{client: clientMock}
			got, err := cr.get(testCtx, testComponentName)
			if !tt.wants.err(t, err) {
				return
			}
			assert.Equal(t, tt.wants.component, got)
		})
	}
}

func Test_defaultComponentRepo_updateCondition(t *testing.T) {
	type mockValues struct {
		getComponent       *v1.Component
		getErr             error
		shouldUpdate       bool
		updateComponentIn  *v1.Component
		updateComponentOut *v1.Component
		updateErr          error
	}
	type args struct {
		component *v1.Component
		statusFn  func() (v1.HealthStatus, error)
		version   string
	}
	tests := []struct {
		name       string
		mockValues mockValues
		args       args
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get component",
			mockValues: mockValues{
				getComponent: nil,
				getErr:       assert.AnError,
				shouldUpdate: false,
			},
			args: args{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				statusFn: func() (v1.HealthStatus, error) {
					return "available", nil
				},
				version: noVersionChange,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to get component %q", testComponentName), i)
			},
		},
		{
			name: "should fail to get status",
			mockValues: mockValues{
				getComponent: nil,
				getErr:       nil,
				shouldUpdate: false,
			},
			args: args{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				statusFn: func() (v1.HealthStatus, error) {
					return "", fmt.Errorf("failed to get status: %w", assert.AnError)
				},
				version: noVersionChange,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get status", i)
			},
		},
		{
			name: "should fail to update component",
			mockValues: mockValues{
				getComponent:       &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				getErr:             nil,
				shouldUpdate:       true,
				updateComponentIn:  &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available"}},
				updateComponentOut: nil,
				updateErr:          assert.AnError,
			},
			args: args{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				statusFn: func() (v1.HealthStatus, error) {
					return "available", nil
				},
				version: noVersionChange,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to update component %q", testComponentName), i)
			},
		},
		{
			name: "should succeed to update component",
			mockValues: mockValues{
				getComponent:       &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{InstalledVersion: "0.2.1"}},
				getErr:             nil,
				shouldUpdate:       true,
				updateComponentIn:  &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available", InstalledVersion: "0.2.1"}},
				updateComponentOut: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available", InstalledVersion: "0.2.1"}},
				updateErr:          nil,
			},
			args: args{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				statusFn: func() (v1.HealthStatus, error) {
					return "available", nil
				},
				version: noVersionChange,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should succeed to update version of component",
			mockValues: mockValues{
				getComponent:       &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{InstalledVersion: "0.2.1"}},
				getErr:             nil,
				shouldUpdate:       true,
				updateComponentIn:  &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available", InstalledVersion: "0.3.0"}},
				updateComponentOut: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available", InstalledVersion: "0.3.0"}},
				updateErr:          nil,
			},
			args: args{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				statusFn: func() (v1.HealthStatus, error) {
					return "available", nil
				},
				version: "0.3.0",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := newMockComponentClient(t)
			clientMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).
				Return(tt.mockValues.getComponent, tt.mockValues.getErr)
			if tt.mockValues.shouldUpdate {
				clientMock.EXPECT().UpdateStatus(testCtx, tt.mockValues.updateComponentIn, metav1.UpdateOptions{}).
					Return(tt.mockValues.updateComponentOut, tt.mockValues.updateErr)
			}
			cr := &defaultComponentRepo{client: clientMock}
			tt.wantErr(t, cr.updateCondition(testCtx, tt.args.component, tt.args.statusFn, tt.args.version))
		})
	}
}

func Test_defaultComponentRepo_list(t *testing.T) {
	tests := []struct {
		name     string
		clientFn func(t *testing.T) componentClient
		want     *v1.ComponentList
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "should fail",
			clientFn: func(t *testing.T) componentClient {
				clientMock := newMockComponentClient(t)
				clientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, assert.AnError)
				return clientMock
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to list components", i)
			},
		},
		{
			name: "should succeed",
			clientFn: func(t *testing.T) componentClient {
				clientMock := newMockComponentClient(t)
				clientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(&v1.ComponentList{
					Items: []v1.Component{{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
						Spec:       v1.ComponentSpec{Name: "k8s-dogu-operator"},
					}},
				}, nil)
				return clientMock
			},
			want: &v1.ComponentList{Items: []v1.Component{{
				ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
				Spec:       v1.ComponentSpec{Name: "k8s-dogu-operator"},
			}}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &defaultComponentRepo{
				client: tt.clientFn(t),
			}
			got, err := cr.list(testCtx)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
