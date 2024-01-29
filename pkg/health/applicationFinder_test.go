package health

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testCtx = context.Background()

const testComponentName = "exampleComponent"
const testNamespace = "ecosystem"

func Test_defaultApplicationFinder_findComponentApplications(t *testing.T) {
	type returnValues struct {
		deployments    *appsv1.DeploymentList
		deploymentErr  error
		statefulSets   *appsv1.StatefulSetList
		statefulSetErr error
		daemonSets     *appsv1.DaemonSetList
		daemonSetErr   error
	}
	type wantValues struct {
		deployments  *appsv1.DeploymentList
		statefulSets *appsv1.StatefulSetList
		daemonSets   *appsv1.DaemonSetList
		err          assert.ErrorAssertionFunc
	}
	tests := []struct {
		name    string
		returns returnValues
		wants   wantValues
	}{
		{
			name: "should fail to list deployments only",
			returns: returnValues{
				deployments:    nil,
				deploymentErr:  assert.AnError,
				statefulSets:   expectedStatefulSets(),
				statefulSetErr: nil,
				daemonSets:     expectedDaemonSets(),
				daemonSetErr:   nil,
			},
			wants: wantValues{
				deployments:  nil,
				statefulSets: nil,
				daemonSets:   nil,
				err: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, assert.AnError, i) &&
						assert.ErrorContains(t, err, "failed to list deployments for component \"exampleComponent\"", i)
				},
			},
		},
		{
			name: "should fail to list stateful sets only",
			returns: returnValues{
				deployments:    expectedDeployments(),
				deploymentErr:  nil,
				statefulSets:   nil,
				statefulSetErr: assert.AnError,
				daemonSets:     expectedDaemonSets(),
				daemonSetErr:   nil,
			},
			wants: wantValues{
				deployments:  nil,
				statefulSets: nil,
				daemonSets:   nil,
				err: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, assert.AnError, i) &&
						assert.ErrorContains(t, err, "failed to list stateful sets for component \"exampleComponent\"", i)
				},
			},
		},
		{
			name: "should fail to list daemon sets only",
			returns: returnValues{
				deployments:    expectedDeployments(),
				deploymentErr:  nil,
				statefulSets:   expectedStatefulSets(),
				statefulSetErr: nil,
				daemonSets:     nil,
				daemonSetErr:   assert.AnError,
			},
			wants: wantValues{
				deployments:  nil,
				statefulSets: nil,
				daemonSets:   nil,
				err: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, assert.AnError, i) &&
						assert.ErrorContains(t, err, "failed to list daemon sets for component \"exampleComponent\"", i)
				},
			},
		},
		{
			name: "should fail to list all applications",
			returns: returnValues{
				deployments:    nil,
				deploymentErr:  assert.AnError,
				statefulSets:   nil,
				statefulSetErr: assert.AnError,
				daemonSets:     nil,
				daemonSetErr:   assert.AnError,
			},
			wants: wantValues{
				deployments:  nil,
				statefulSets: nil,
				daemonSets:   nil,
				err: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, assert.AnError, i) &&
						assert.ErrorContains(t, err, "failed to list deployments for component \"exampleComponent\"", i) &&
						assert.ErrorContains(t, err, "failed to list stateful sets for component \"exampleComponent\"", i) &&
						assert.ErrorContains(t, err, "failed to list daemon sets for component \"exampleComponent\"", i)
				},
			},
		},
		{
			name: "should successfully list all applications",
			returns: returnValues{
				deployments:    expectedDeployments(),
				deploymentErr:  nil,
				statefulSets:   expectedStatefulSets(),
				statefulSetErr: nil,
				daemonSets:     expectedDaemonSets(),
				daemonSetErr:   nil,
			},
			wants: wantValues{
				deployments:  expectedDeployments(),
				statefulSets: expectedStatefulSets(),
				daemonSets:   expectedDaemonSets(),
				err:          assert.NoError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dplmntClient := newMockDeploymentClient(t)
			dplmntClient.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "k8s.cloudogu.com/component.name=exampleComponent"}).
				Return(tt.returns.deployments, tt.returns.deploymentErr)
			sttfllStClient := newMockStatefulSetClient(t)
			sttfllStClient.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "k8s.cloudogu.com/component.name=exampleComponent"}).
				Return(tt.returns.statefulSets, tt.returns.statefulSetErr)
			dmnStClient := newMockDaemonSetClient(t)
			dmnStClient.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "k8s.cloudogu.com/component.name=exampleComponent"}).
				Return(tt.returns.daemonSets, tt.returns.daemonSetErr)
			appsClient := newMockAppsV1Client(t)
			appsClient.EXPECT().Deployments(testNamespace).Return(dplmntClient)
			appsClient.EXPECT().StatefulSets(testNamespace).Return(sttfllStClient)
			appsClient.EXPECT().DaemonSets(testNamespace).Return(dmnStClient)
			af := &defaultApplicationFinder{
				appsClient: appsClient,
			}
			gotDeployments, gotStatefulSets, gotDaemonSets, err := af.findComponentApplications(
				testCtx,
				testComponentName,
				testNamespace,
			)
			tt.wants.err(t, err)
			assert.Equal(t, tt.wants.deployments, gotDeployments)
			assert.Equal(t, tt.wants.statefulSets, gotStatefulSets)
			assert.Equal(t, tt.wants.daemonSets, gotDaemonSets)
		})
	}
}

func expectedDeployments() *appsv1.DeploymentList {
	return &appsv1.DeploymentList{Items: []appsv1.Deployment{
		{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
	}}
}

func expectedStatefulSets() *appsv1.StatefulSetList {
	return &appsv1.StatefulSetList{Items: []appsv1.StatefulSet{
		{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
	}}
}

func expectedDaemonSets() *appsv1.DaemonSetList {
	return &appsv1.DaemonSetList{Items: []appsv1.DaemonSet{
		{ObjectMeta: metav1.ObjectMeta{Name: testComponentName, Namespace: testNamespace}},
	}}
}
