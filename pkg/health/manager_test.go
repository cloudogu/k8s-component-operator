package health

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

var one int32 = 1
var two int32 = 2
var three int32 = 3

func availableDeploymentList() *appsv1.DeploymentList {
	return &appsv1.DeploymentList{Items: []appsv1.Deployment{
		{Status: appsv1.DeploymentStatus{Replicas: 1, UpdatedReplicas: 1, AvailableReplicas: 1}},
		{Status: appsv1.DeploymentStatus{Replicas: 2, UpdatedReplicas: 2, AvailableReplicas: 2}},
		{Status: appsv1.DeploymentStatus{Replicas: 3, UpdatedReplicas: 3, AvailableReplicas: 3}},
	}}
}

func availableStatefulSetList() *appsv1.StatefulSetList {
	return &appsv1.StatefulSetList{Items: []appsv1.StatefulSet{
		{Spec: appsv1.StatefulSetSpec{Replicas: &one}, Status: appsv1.StatefulSetStatus{Replicas: 1, UpdatedReplicas: 1, AvailableReplicas: 1}},
		{Spec: appsv1.StatefulSetSpec{Replicas: &two}, Status: appsv1.StatefulSetStatus{Replicas: 2, UpdatedReplicas: 2, AvailableReplicas: 2}},
		{Spec: appsv1.StatefulSetSpec{Replicas: &three}, Status: appsv1.StatefulSetStatus{Replicas: 3, UpdatedReplicas: 3, AvailableReplicas: 3}},
	}}
}

func availableDaemonSetList() *appsv1.DaemonSetList {
	return &appsv1.DaemonSetList{Items: []appsv1.DaemonSet{
		{Status: appsv1.DaemonSetStatus{CurrentNumberScheduled: 1, DesiredNumberScheduled: 1, UpdatedNumberScheduled: 1, NumberAvailable: 1}},
		{Status: appsv1.DaemonSetStatus{CurrentNumberScheduled: 2, DesiredNumberScheduled: 2, UpdatedNumberScheduled: 2, NumberAvailable: 2}},
		{Status: appsv1.DaemonSetStatus{CurrentNumberScheduled: 3, DesiredNumberScheduled: 3, UpdatedNumberScheduled: 3, NumberAvailable: 3}},
	}}
}

func unavailableDeploymentList() *appsv1.DeploymentList {
	return &appsv1.DeploymentList{Items: []appsv1.Deployment{
		{Spec: appsv1.DeploymentSpec{Replicas: &one}, Status: appsv1.DeploymentStatus{Replicas: 1, UpdatedReplicas: 1, AvailableReplicas: 1}},
		{Spec: appsv1.DeploymentSpec{Replicas: &two}, Status: appsv1.DeploymentStatus{Replicas: 2, UpdatedReplicas: 2, AvailableReplicas: 2}},
		{Spec: appsv1.DeploymentSpec{Replicas: &three}, Status: appsv1.DeploymentStatus{Replicas: 3, UpdatedReplicas: 2, AvailableReplicas: 0}},
	}}
}

func unavailableStatefulSetList() *appsv1.StatefulSetList {
	return &appsv1.StatefulSetList{Items: []appsv1.StatefulSet{
		{Status: appsv1.StatefulSetStatus{Replicas: 1, UpdatedReplicas: 1, AvailableReplicas: 1}},
		{Status: appsv1.StatefulSetStatus{Replicas: 2, UpdatedReplicas: 2, AvailableReplicas: 2}},
		{Status: appsv1.StatefulSetStatus{Replicas: 3, UpdatedReplicas: 1, AvailableReplicas: 1}},
	}}
}

func unavailableDaemonSetList() *appsv1.DaemonSetList {
	return &appsv1.DaemonSetList{Items: []appsv1.DaemonSet{
		{Status: appsv1.DaemonSetStatus{CurrentNumberScheduled: 1, DesiredNumberScheduled: 1, UpdatedNumberScheduled: 1, NumberAvailable: 1}},
		{Status: appsv1.DaemonSetStatus{CurrentNumberScheduled: 2, DesiredNumberScheduled: 2, UpdatedNumberScheduled: 2, NumberAvailable: 2}},
		{Status: appsv1.DaemonSetStatus{CurrentNumberScheduled: 2, DesiredNumberScheduled: 3, UpdatedNumberScheduled: 1, NumberAvailable: 1}},
	}}
}

func Test_newManager(t *testing.T) {
	appsV1Mock := newMockAppsV1Client(t)
	componentMock := newMockComponentClient(t)
	componentV1Alpha1Mock := newMockComponentV1Alpha1Client(t)
	componentV1Alpha1Mock.EXPECT().Components(testNamespace).Return(componentMock)
	clientSetMock := newMockEcosystemClientSet(t)
	clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
	clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1Alpha1Mock)

	assert.NotEmpty(t, NewManager(testNamespace, clientSetMock))
}

func Test_defaultManager_componentHealthStatus(t *testing.T) {
	type args struct {
		deployments  *appsv1.DeploymentList
		statefulSets *appsv1.StatefulSetList
		daemonSets   *appsv1.DaemonSetList
		component    *v1.Component
	}
	tests := []struct {
		name string
		args args
		want v1.HealthStatus
	}{
		{
			name: "should be available if installed with no applications found",
			args: args{
				deployments:  &appsv1.DeploymentList{},
				statefulSets: &appsv1.StatefulSetList{},
				daemonSets:   &appsv1.DaemonSetList{},
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}},
			},
			want: "available",
		},
		{
			name: "should be available if not installed with no applications found",
			args: args{
				deployments:  &appsv1.DeploymentList{},
				statefulSets: &appsv1.StatefulSetList{},
				daemonSets:   &appsv1.DaemonSetList{},
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalling}},
			},
			want: "unavailable",
		},
		{
			name: "should be unavailable if at least one application is not available",
			args: args{
				deployments:  availableDeploymentList(),
				statefulSets: unavailableStatefulSetList(),
				daemonSets:   availableDaemonSetList(),
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}},
			},
			want: "unavailable",
		},
		{
			name: "should be unavailable if multiple applications are not available",
			args: args{
				deployments:  unavailableDeploymentList(),
				statefulSets: unavailableStatefulSetList(),
				daemonSets:   unavailableDaemonSetList(),
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}},
			},
			want: "unavailable",
		},
		{
			name: "should be unavailable if not installed and multiple applications are not available",
			args: args{
				deployments:  unavailableDeploymentList(),
				statefulSets: unavailableStatefulSetList(),
				daemonSets:   unavailableDaemonSetList(),
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusDeleting}},
			},
			want: "unavailable",
		},
		{
			name: "should be unavailable if not installed and all applications are available",
			args: args{
				deployments:  availableDeploymentList(),
				statefulSets: availableStatefulSetList(),
				daemonSets:   availableDaemonSetList(),
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusUpgrading}},
			},
			want: "unavailable",
		},
		{
			name: "should be available if installed and all applications are available",
			args: args{
				deployments:  availableDeploymentList(),
				statefulSets: availableStatefulSetList(),
				daemonSets:   availableDaemonSetList(),
				component:    &v1.Component{Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}},
			},
			want: "available",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultManager{}
			assert.Equal(t, tt.want, m.componentHealthStatus(testCtx, tt.args.deployments, tt.args.statefulSets, tt.args.daemonSets, tt.args.component))
		})
	}
}

func Test_defaultManager_UpdateComponentHealth(t *testing.T) {
	testComponent := v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}}
	type fields struct {
		applicationFinderFn func(t *testing.T) applicationFinder
		componentRepoFn     func(t *testing.T) componentRepo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get component",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).Return(nil, assert.AnError)
					return repo
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to get component %q", testComponentName), i)
			},
		},
		{
			name: "should fail to find applications",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					finder.EXPECT().findComponentApplications(testCtx, testComponentName, testNamespace).
						Return(nil, nil, nil, assert.AnError)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).
						Return(&testComponent, nil)
					return repo
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to find applications for component %q", testComponentName), i)
			},
		},
		{
			name: "should fail to update component health",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					finder.EXPECT().findComponentApplications(testCtx, testComponentName, testNamespace).
						Return(availableDeploymentList(), availableStatefulSetList(), availableDaemonSetList(), nil)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).
						Return(&testComponent, nil)
					repo.EXPECT().updateCondition(testCtx, &testComponent, v1.HealthStatus("available"), noVersionChange).
						Return(assert.AnError)
					return repo
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to update health status and installed version for component %q", testComponentName), i)
			},
		},
		{
			name: "should succeed to update component health",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					finder.EXPECT().findComponentApplications(testCtx, testComponentName, testNamespace).
						Return(availableDeploymentList(), availableStatefulSetList(), availableDaemonSetList(), nil)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).
						Return(&testComponent, nil)
					repo.EXPECT().updateCondition(testCtx, &testComponent, v1.HealthStatus("available"), noVersionChange).
						Return(nil)
					return repo
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultManager{
				applicationFinder: tt.fields.applicationFinderFn(t),
				componentRepo:     tt.fields.componentRepoFn(t),
			}
			tt.wantErr(t, m.UpdateComponentHealth(testCtx, testComponentName, testNamespace))
		})
	}
}

func Test_defaultManager_UpdateComponentHealthWithVersion(t *testing.T) {
	testComponent := v1.Component{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}, Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalled}}
	type fields struct {
		applicationFinderFn func(t *testing.T) applicationFinder
		componentRepoFn     func(t *testing.T) componentRepo
	}
	testVersion := "0.2.1"
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to get component",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).Return(nil, assert.AnError)
					return repo
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to get component %q", testComponentName), i)
			},
		},
		{
			name: "should fail to find applications",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					finder.EXPECT().findComponentApplications(testCtx, testComponentName, testNamespace).
						Return(nil, nil, nil, assert.AnError)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).
						Return(&testComponent, nil)
					return repo
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to find applications for component %q", testComponentName), i)
			},
		},
		{
			name: "should fail to update component health",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					finder.EXPECT().findComponentApplications(testCtx, testComponentName, testNamespace).
						Return(availableDeploymentList(), availableStatefulSetList(), availableDaemonSetList(), nil)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).
						Return(&testComponent, nil)
					repo.EXPECT().updateCondition(testCtx, &testComponent, v1.HealthStatus("available"), testVersion).
						Return(assert.AnError)
					return repo
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to update health status and installed version for component %q", testComponentName), i)
			},
		},
		{
			name: "should succeed to update component health",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finder := newMockApplicationFinder(t)
					finder.EXPECT().findComponentApplications(testCtx, testComponentName, testNamespace).
						Return(availableDeploymentList(), availableStatefulSetList(), availableDaemonSetList(), nil)
					return finder
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repo := newMockComponentRepo(t)
					repo.EXPECT().get(testCtx, testComponentName).
						Return(&testComponent, nil)
					repo.EXPECT().updateCondition(testCtx, &testComponent, v1.HealthStatus("available"), testVersion).
						Return(nil)
					return repo
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultManager{
				applicationFinder: tt.fields.applicationFinderFn(t),
				componentRepo:     tt.fields.componentRepoFn(t),
			}
			tt.wantErr(t, m.UpdateComponentHealthWithInstalledVersion(testCtx, testComponentName, testNamespace, testVersion))
		})
	}
}

func TestDefaultManager_UpdateComponentHealthAll(t *testing.T) {
	type fields struct {
		applicationFinderFn func(t *testing.T) applicationFinder
		componentRepoFn     func(t *testing.T) componentRepo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should fail to list components",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					return newMockApplicationFinder(t)
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repoMock := newMockComponentRepo(t)
					repoMock.EXPECT().list(testCtx).Return(nil, assert.AnError)
					return repoMock
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i)
			},
		},
		{
			name: "should fail to update health for multiple components",
			fields: fields{
				applicationFinderFn: func(t *testing.T) applicationFinder {
					finderMock := newMockApplicationFinder(t)
					finderMock.EXPECT().findComponentApplications(testCtx, "k8s-dogu-operator", testNamespace).
						Return(nil, nil, nil, assert.AnError)
					finderMock.EXPECT().findComponentApplications(testCtx, "k8s-blueprint-operator", testNamespace).
						Return(&appsv1.DeploymentList{}, &appsv1.StatefulSetList{}, &appsv1.DaemonSetList{}, nil)
					finderMock.EXPECT().findComponentApplications(testCtx, "k8s-longhorn", "longhorn-system").
						Return(&appsv1.DeploymentList{}, &appsv1.StatefulSetList{}, &appsv1.DaemonSetList{}, nil)
					return finderMock
				},
				componentRepoFn: func(t *testing.T) componentRepo {
					repoMock := newMockComponentRepo(t)
					repoMock.EXPECT().list(testCtx).Return(&v1.ComponentList{Items: []v1.Component{
						{
							ObjectMeta: metav1.ObjectMeta{Name: "k8s-dogu-operator"},
							Spec: v1.ComponentSpec{
								Name:            "k8s-dogu-operator",
								DeployNamespace: testNamespace,
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
							Spec: v1.ComponentSpec{
								Name:            "k8s-blueprint-operator",
								DeployNamespace: testNamespace,
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{Name: "k8s-longhorn"},
							Spec: v1.ComponentSpec{
								Name:            "k8s-longhorn",
								DeployNamespace: "longhorn-system",
							},
						},
					}}, nil)
					repoMock.EXPECT().updateCondition(testCtx,
						&v1.Component{
							ObjectMeta: metav1.ObjectMeta{Name: "k8s-blueprint-operator"},
							Spec: v1.ComponentSpec{
								Name:            "k8s-blueprint-operator",
								DeployNamespace: testNamespace,
							},
						}, v1.UnavailableHealthStatus, noVersionChange).Return(assert.AnError)
					repoMock.EXPECT().updateCondition(testCtx,
						&v1.Component{
							ObjectMeta: metav1.ObjectMeta{Name: "k8s-longhorn"},
							Spec: v1.ComponentSpec{
								Name:            "k8s-longhorn",
								DeployNamespace: "longhorn-system",
							},
						}, v1.UnavailableHealthStatus, noVersionChange).Return(nil)
					return repoMock
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to find applications for component \"k8s-dogu-operator\"", i) &&
					assert.ErrorContains(t, err, "failed to update health status and installed version for component \"k8s-blueprint-operator\"", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultManager{
				applicationFinder: tt.fields.applicationFinderFn(t),
				componentRepo:     tt.fields.componentRepoFn(t),
			}
			tt.wantErr(t, m.UpdateComponentHealthAll(testCtx))
		})
	}
}
