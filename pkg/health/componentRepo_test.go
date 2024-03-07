package health

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
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

func Test_defaultComponentRepo_updateHealthStatus(t *testing.T) {
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
		status    v1.HealthStatus
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
				status:    "available",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to get component %q", testComponentName), i)
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
				status:    "available",
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to update component %q", testComponentName), i)
			},
		},
		{
			name: "should succeed to update component",
			mockValues: mockValues{
				getComponent:       &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				getErr:             nil,
				shouldUpdate:       true,
				updateComponentIn:  &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available"}},
				updateComponentOut: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Health: "available"}},
				updateErr:          nil,
			},
			args: args{
				component: &v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
				status:    "available",
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
			tt.wantErr(t, cr.updateCondition(testCtx, tt.args.component, tt.args.status, ""))
		})
	}
}
