package health

import (
	"context"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewShutdownHandler(t *testing.T) {
	// given
	clientMock := newMockComponentClient(t)

	// when
	actual := NewShutdownHandler(clientMock)

	// then
	assert.NotEmpty(t, actual)
	assert.NotEmpty(t, actual.repo)
}

func TestShutdownHandler_Start(t *testing.T) {
	tests := []struct {
		name    string
		repoFn  func(t *testing.T) componentRepo
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get component operator component",
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
				repoMock.EXPECT().get(mock.Anything, "k8s-component-operator").Return(nil, assert.AnError)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to get component for \"k8s-component-operator\"", i)
			},
		},
		{
			name: "should fail to update health status for component operator",
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
				repoMock.EXPECT().get(mock.Anything, "k8s-component-operator").Return(&k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, k8sv1.UnavailableHealthStatus).Return(assert.AnError)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to set health status of \"k8s-component-operator\" to \"unavailable\"", i)
			},
		},
		{
			name: "should fail to list components",
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
				repoMock.EXPECT().get(mock.Anything, "k8s-component-operator").Return(&k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, k8sv1.UnavailableHealthStatus).Return(nil)
				repoMock.EXPECT().list(mock.Anything).Return(nil, assert.AnError)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i)
			},
		},
		{
			name: "should fail to update health status for multiple components",
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
				repoMock.EXPECT().get(mock.Anything, "k8s-component-operator").Return(&k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, k8sv1.UnavailableHealthStatus).Return(nil)
				repoMock.EXPECT().list(mock.Anything).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
					},
				},
				}, nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
				}, k8sv1.UnknownHealthStatus).Return(assert.AnError)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
				}, k8sv1.UnknownHealthStatus).Return(nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
				}, k8sv1.UnknownHealthStatus).Return(assert.AnError)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to set health status of \"k8s-dogu-operator\" to \"unknown\"", i) &&
					assert.ErrorContains(t, err, "failed to set health status of \"k8s-backup-operator\" to \"unknown\"", i)
			},
		},
		{
			name: "should succeed",
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
				repoMock.EXPECT().get(mock.Anything, "k8s-component-operator").Return(&k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, k8sv1.UnavailableHealthStatus).Return(nil)
				repoMock.EXPECT().list(mock.Anything).Return(&k8sv1.ComponentList{Items: []k8sv1.Component{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
						Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
					},
				},
				}, nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
				}, k8sv1.UnknownHealthStatus).Return(nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
				}, k8sv1.UnknownHealthStatus).Return(nil)
				repoMock.EXPECT().updateHealthStatus(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
				}, k8sv1.UnknownHealthStatus).Return(nil)
				return repoMock
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShutdownHandler{
				repo: tt.repoFn(t),
			}
			ctx, cancel := context.WithCancel(testCtx)
			cancel()
			tt.wantErr(t, s.Start(ctx))
		})
	}
}
