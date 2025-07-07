package health

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewSyncIntervalHandler(t *testing.T) {
	// given
	clientSetMock := newMockEcosystemClientSet(t)
	appsV1Mock := newMockAppsV1Client(t)
	clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
	componentV1Mock := newMockComponentV1Alpha1Client(t)
	componentMock := newMockComponentClient(t)
	componentV1Mock.EXPECT().Components(testNamespace).Return(componentMock)
	clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1Mock)

	// when
	actual := NewSyncIntervalHandler(testNamespace, clientSetMock, 2*time.Minute)

	// then
	assert.NotEmpty(t, actual)
	assert.NotEmpty(t, actual.manager)
}

func TestSyncIntervalHandler_Start(t *testing.T) {
	tests := []struct {
		name      string
		managerFn func(t *testing.T) ComponentManager
		repoFn    func(t *testing.T) componentRepo
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to regularly sync health of all components",
			managerFn: func(t *testing.T) ComponentManager {
				managerMock := NewMockComponentManager(t)
				managerMock.EXPECT().UpdateComponentHealthAll(mock.Anything).Return(assert.AnError)
				return managerMock
			},
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
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
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnavailableHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "some errors occurred when regularly syncing health of all components", i)
			},
		},
		{
			name: "should fail to list components on shutdown",
			managerFn: func(t *testing.T) ComponentManager {
				managerMock := NewMockComponentManager(t)
				managerMock.EXPECT().UpdateComponentHealthAll(mock.Anything).Return(nil)
				return managerMock
			},
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
				repoMock.EXPECT().list(mock.Anything).Return(nil, assert.AnError)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "some errors occurred when regularly syncing health of all components", i) &&
					assert.ErrorContains(t, err, "failed to correctly handle health status during shutdown", i)
			},
		},
		{
			name: "should fail to update health status for multiple components on shutdown",
			managerFn: func(t *testing.T) ComponentManager {
				managerMock := NewMockComponentManager(t)
				managerMock.EXPECT().UpdateComponentHealthAll(mock.Anything).Return(nil)
				return managerMock
			},
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
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
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(assert.AnError)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnavailableHealthStatus, status)
					}).
					Return(assert.AnError)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(assert.AnError)
				return repoMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "some errors occurred when regularly syncing health of all components", i) &&
					assert.ErrorContains(t, err, "failed to correctly handle health status during shutdown", i) &&
					assert.ErrorContains(t, err, "failed to set health status of \"k8s-dogu-operator\" to \"unknown\"", i) &&
					assert.ErrorContains(t, err, "failed to set health status of \"k8s-backup-operator\" to \"unknown\"", i) &&
					assert.ErrorContains(t, err, "failed to set health status of \"k8s-component-operator\" to \"unavailable\"", i)
			},
		},
		{
			name: "should succeed",
			managerFn: func(t *testing.T) ComponentManager {
				managerMock := NewMockComponentManager(t)
				managerMock.EXPECT().UpdateComponentHealthAll(mock.Anything).Return(nil)
				return managerMock
			},
			repoFn: func(t *testing.T) componentRepo {
				repoMock := newMockComponentRepo(t)
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
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-dogu-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-component-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-component-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnavailableHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-blueprint-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				repoMock.EXPECT().updateCondition(mock.Anything, &k8sv1.Component{
					ObjectMeta: metav1.ObjectMeta{Name: "k8s-backup-operator"},
					Spec:       k8sv1.ComponentSpec{Name: "k8s-backup-operator"},
				}, mock.Anything, noVersionChange).
					Run(func(ctx context.Context, component *k8sv1.Component, statusFn func() (k8sv1.HealthStatus, error), version string) {
						status, err := statusFn()
						assert.NoError(t, err)
						assert.Equal(t, k8sv1.UnknownHealthStatus, status)
					}).
					Return(nil)
				return repoMock
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncIntervalHandler{
				manager:            tt.managerFn(t),
				repo:               tt.repoFn(t),
				healthSyncInterval: time.Millisecond,
			}
			ctx, cancelFunc := context.WithTimeout(testCtx, 10*time.Millisecond)
			defer cancelFunc()
			tt.wantErr(t, s.Start(ctx))
		})
	}
}
