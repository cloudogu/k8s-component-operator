package health

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
)

func Test_statefulSetReconciler_SetupWithManager(t *testing.T) {
	t.Run("should fail", func(t *testing.T) {
		// given
		sut := &statefulSetReconciler{}

		// when
		err := sut.SetupWithManager(nil)

		// then
		require.Error(t, err)
		assert.ErrorContains(t, err, "must provide a non-nil Manager")
	})
	t.Run("should succeed", func(t *testing.T) {
		// given
		ctrlManMock := newMockControllerManager(t)
		ctrlManMock.EXPECT().GetControllerOptions().Return(config.Controller{})
		ctrlManMock.EXPECT().GetScheme().Return(createScheme(t))
		logger := log.FromContext(testCtx)
		ctrlManMock.EXPECT().GetLogger().Return(logger)
		ctrlManMock.EXPECT().Add(mock.Anything).Return(nil)
		ctrlManMock.EXPECT().GetCache().Return(nil)

		sut := &statefulSetReconciler{}

		// when
		err := sut.SetupWithManager(ctrlManMock)

		// then
		require.NoError(t, err)
	})
}

func createScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()
	gv, err := schema.ParseGroupVersion("apps/v1")
	assert.NoError(t, err)

	scheme.AddKnownTypes(gv, &appsv1.Deployment{})
	scheme.AddKnownTypes(gv, &appsv1.StatefulSet{})
	scheme.AddKnownTypes(gv, &appsv1.DaemonSet{})
	return scheme
}

var notFoundErr = errors.NewNotFound(schema.GroupResource{}, testComponentName)

func Test_statefulSetReconciler_Reconcile(t *testing.T) {
	type fields struct {
		clientSetFn func(t *testing.T) ecosystemClientSet
		managerFn   func(t *testing.T) ComponentManager
	}
	tests := []struct {
		name    string
		fields  fields
		want    reconcile.Result
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should exit early if application is not found",
			fields: fields{
				clientSetFn: func(t *testing.T) ecosystemClientSet {
					statefulSetMock := newMockStatefulSetClient(t)
					statefulSetMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).
						Return(nil, notFoundErr)
					appsV1Mock := newMockAppsV1Client(t)
					appsV1Mock.EXPECT().StatefulSets(testNamespace).Return(statefulSetMock)
					clientSetMock := newMockEcosystemClientSet(t)
					clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
					return clientSetMock
				},
				managerFn: func(t *testing.T) ComponentManager {
					return NewMockComponentManager(t)
				},
			},
			want:    reconcile.Result{},
			wantErr: assert.NoError,
		},
		{
			name: "should fail to get application",
			fields: fields{
				clientSetFn: func(t *testing.T) ecosystemClientSet {
					statefulSetMock := newMockStatefulSetClient(t)
					statefulSetMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).
						Return(nil, assert.AnError)
					appsV1Mock := newMockAppsV1Client(t)
					appsV1Mock.EXPECT().StatefulSets(testNamespace).Return(statefulSetMock)
					clientSetMock := newMockEcosystemClientSet(t)
					clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
					return clientSetMock
				},
				managerFn: func(t *testing.T) ComponentManager {
					return NewMockComponentManager(t)
				},
			},
			want: reconcile.Result{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to get stateful set \"%s/%s\"", testNamespace, testComponentName), i)
			},
		},
		{
			name: "should fail to update component health",
			fields: fields{
				clientSetFn: func(t *testing.T) ecosystemClientSet {
					statefulSetMock := newMockStatefulSetClient(t)
					statefulSetMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).
						Return(&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"other_key":              "other_value",
								v1.ComponentNameLabelKey: testComponentName,
							},
							Name:      testComponentName,
							Namespace: testNamespace,
						}}, nil)
					appsV1Mock := newMockAppsV1Client(t)
					appsV1Mock.EXPECT().StatefulSets(testNamespace).Return(statefulSetMock)
					clientSetMock := newMockEcosystemClientSet(t)
					clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
					return clientSetMock
				},
				managerFn: func(t *testing.T) ComponentManager {
					manager := NewMockComponentManager(t)
					manager.EXPECT().UpdateComponentHealth(testCtx, testComponentName, testNamespace).
						Return(assert.AnError)
					return manager
				},
			},
			want: reconcile.Result{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, fmt.Sprintf("failed to update component health for stateful set \"%s/%s\"", testNamespace, testComponentName), i)
			},
		},
		{
			name: "should succeed to update component health",
			fields: fields{
				clientSetFn: func(t *testing.T) ecosystemClientSet {
					statefulSetMock := newMockStatefulSetClient(t)
					statefulSetMock.EXPECT().Get(testCtx, testComponentName, metav1.GetOptions{}).
						Return(&appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"other_key":              "other_value",
								v1.ComponentNameLabelKey: testComponentName,
							},
							Name:      testComponentName,
							Namespace: testNamespace,
						}}, nil)
					appsV1Mock := newMockAppsV1Client(t)
					appsV1Mock.EXPECT().StatefulSets(testNamespace).Return(statefulSetMock)
					clientSetMock := newMockEcosystemClientSet(t)
					clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
					return clientSetMock
				},
				managerFn: func(t *testing.T) ComponentManager {
					manager := NewMockComponentManager(t)
					manager.EXPECT().UpdateComponentHealth(testCtx, testComponentName, testNamespace).
						Return(nil)
					return manager
				},
			},
			want:    reconcile.Result{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssr := &statefulSetReconciler{
				clientSet: tt.fields.clientSetFn(t),
				manager:   tt.fields.managerFn(t),
			}
			got, err := ssr.Reconcile(testCtx, ctrl.Request{NamespacedName: types.NamespacedName{Name: testComponentName, Namespace: testNamespace}})
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
