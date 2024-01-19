package health

import (
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type healthChangeEventFilter struct{}

func (h *healthChangeEventFilter) Create(event event.CreateEvent) bool {
	return hasComponentLabel(event.Object)
}

func (h *healthChangeEventFilter) Delete(event event.DeleteEvent) bool {
	return hasComponentLabel(event.Object)
}

func (h *healthChangeEventFilter) Update(event event.UpdateEvent) bool {
	if !hasComponentLabel(event.ObjectNew) {
		return false
	}

	switch oldObj := event.ObjectOld.(type) {
	case *appsv1.Deployment:
		newObj := event.ObjectNew.(*appsv1.Deployment)
		return oldObj.Spec.Replicas != newObj.Spec.Replicas ||
			oldObj.Status.Replicas != newObj.Status.Replicas ||
			oldObj.Status.UpdatedReplicas != newObj.Status.UpdatedReplicas ||
			oldObj.Status.AvailableReplicas != newObj.Status.AvailableReplicas
	case *appsv1.StatefulSet:
		newObj := event.ObjectNew.(*appsv1.StatefulSet)
		return oldObj.Spec.Replicas != newObj.Spec.Replicas ||
			oldObj.Status.Replicas != newObj.Status.Replicas ||
			oldObj.Status.UpdatedReplicas != newObj.Status.UpdatedReplicas ||
			oldObj.Status.AvailableReplicas != newObj.Status.AvailableReplicas
	case *appsv1.DaemonSet:
		newObj := event.ObjectNew.(*appsv1.DaemonSet)
		return oldObj.Status.DesiredNumberScheduled != newObj.Status.DesiredNumberScheduled ||
			oldObj.Status.CurrentNumberScheduled != newObj.Status.CurrentNumberScheduled ||
			oldObj.Status.UpdatedNumberScheduled != newObj.Status.UpdatedNumberScheduled ||
			oldObj.Status.NumberAvailable != newObj.Status.NumberAvailable
	default:
		return false
	}
}

func (h *healthChangeEventFilter) Generic(_ event.GenericEvent) bool {
	return false
}
