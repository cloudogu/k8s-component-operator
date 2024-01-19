package health

import (
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"testing"
)

func Test_healthChangeEventFilter_Create(t *testing.T) {
	tests := []struct {
		name  string
		event event.CreateEvent
		want  bool
	}{
		{
			name:  "should be false if component label is not set",
			event: event.CreateEvent{Object: &appsv1.Deployment{}},
			want:  false,
		},
		{
			name: "should be true if component label is set",
			event: event.CreateEvent{Object: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
			}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &healthChangeEventFilter{}
			assert.Equal(t, tt.want, h.Create(tt.event))
		})
	}
}

func Test_healthChangeEventFilter_Delete(t *testing.T) {
	tests := []struct {
		name  string
		event event.DeleteEvent
		want  bool
	}{
		{
			name:  "should be false if component label is not set",
			event: event.DeleteEvent{Object: &appsv1.Deployment{}},
			want:  false,
		},
		{
			name: "should be true if component label is set",
			event: event.DeleteEvent{Object: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
			}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &healthChangeEventFilter{}
			assert.Equal(t, tt.want, h.Delete(tt.event))
		})
	}
}

func Test_healthChangeEventFilter_Generic(t *testing.T) {
	assert.False(t, (&healthChangeEventFilter{}).Generic(event.GenericEvent{}))
}

func Test_healthChangeEventFilter_Update(t *testing.T) {
	var one int32 = 1
	tests := []struct {
		name  string
		event event.UpdateEvent
		want  bool
	}{
		{
			name: "should be false if new object does not have component label",
			event: event.UpdateEvent{
				ObjectOld: nil,
				ObjectNew: &appsv1.Deployment{},
			},
			want: false,
		},
		{
			name: "should be false if object is not deployment, stateful set or daemon set",
			event: event.UpdateEvent{
				ObjectOld: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
				ObjectNew: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
			},
			want: false,
		},

		// Deployments
		{
			name: "should be false if type assertion fails for deployment",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
				ObjectNew: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
			},
			want: false,
		},
		{
			name: "should be false if spec.replicas, status.replicas, status.updatedReplicas and status.availableReplicas did not change for deployment",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: false,
		},
		{
			name: "should be true if spec.replicas changed for deployment",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: nil},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.replicas changed for deployment",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          0,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.updatedReplicas changed for deployment",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   0,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.availableReplicas changed for deployment",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 0,
					},
				},
				ObjectNew: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},

		// StatefulSets
		{
			name: "should be false if type assertion fails for stateful set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
				ObjectNew: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
			},
			want: false,
		},
		{
			name: "should be false if spec.replicas, status.replicas, status.updatedReplicas and status.availableReplicas did not change for stateful set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: false,
		},
		{
			name: "should be true if spec.replicas changed for stateful set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: nil},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.replicas changed for stateful set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          0,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.updatedReplicas changed for stateful set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   0,
						AvailableReplicas: 1,
					},
				},
				ObjectNew: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.availableReplicas changed for stateful set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 0,
					},
				},
				ObjectNew: &appsv1.StatefulSet{
					Spec: appsv1.StatefulSetSpec{Replicas: &one},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.StatefulSetStatus{
						Replicas:          1,
						UpdatedReplicas:   1,
						AvailableReplicas: 1,
					},
				},
			},
			want: true,
		},

		// DaemonSets
		{
			name: "should be false if type assertion fails for daemon set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
				ObjectNew: &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
				}},
			},
			want: false,
		},
		{
			name: "should be false if status.desiredNumberScheduled, status.currentNumberScheduled, status.updatedNumberScheduled and status.numberAvailable did not change for daemon set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
				ObjectNew: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
			},
			want: false,
		},
		{
			name: "should be true if status.desiredNumberScheduled changed for daemon set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 0,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
				ObjectNew: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.currentNumberScheduled changed for daemon set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 0,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
				ObjectNew: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.updatedNumberScheduled changed for daemon set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 0,
						NumberAvailable:        1,
					},
				},
				ObjectNew: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
			},
			want: true,
		},
		{
			name: "should be true if status.numberAvailable changed for daemon set",
			event: event.UpdateEvent{
				ObjectOld: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        0,
					},
				},
				ObjectNew: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{v1.ComponentNameLabelKey: testComponentName},
					},
					Status: appsv1.DaemonSetStatus{
						DesiredNumberScheduled: 1,
						CurrentNumberScheduled: 1,
						UpdatedNumberScheduled: 1,
						NumberAvailable:        1,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &healthChangeEventFilter{}
			assert.Equal(t, tt.want, h.Update(tt.event))
		})
	}
}
