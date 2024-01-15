package health

import appsv1 "k8s.io/api/apps/v1"

type state struct {
	kind              string
	name              string
	component         string
	desiredReplicas   int32
	scheduledReplicas int32
	updatedReplicas   int32
	availableReplicas int32
}

func (s state) IsAvailable() bool {
	if s.desiredReplicas < 1 ||
		s.updatedReplicas < s.desiredReplicas ||
		s.updatedReplicas < s.scheduledReplicas ||
		s.availableReplicas < s.updatedReplicas {
		return false
	}

	return true
}

func deploymentToState(deployment appsv1.Deployment) state {
	return state{
		kind:              deployment.Kind,
		name:              deployment.Name,
		component:         componentName(&deployment),
		desiredReplicas:   defaultToOne(deployment.Spec.Replicas),
		scheduledReplicas: deployment.Status.Replicas,
		updatedReplicas:   deployment.Status.UpdatedReplicas,
		availableReplicas: deployment.Status.AvailableReplicas,
	}
}

func statefulSetToState(statefulSet appsv1.StatefulSet) state {
	return state{
		kind:              statefulSet.Kind,
		name:              statefulSet.Name,
		component:         componentName(&statefulSet),
		desiredReplicas:   defaultToOne(statefulSet.Spec.Replicas),
		scheduledReplicas: statefulSet.Status.Replicas,
		updatedReplicas:   statefulSet.Status.UpdatedReplicas,
		availableReplicas: statefulSet.Status.AvailableReplicas,
	}
}

func daemonSetToState(daemonSet appsv1.DaemonSet) state {
	return state{
		kind:              daemonSet.Kind,
		name:              daemonSet.Name,
		component:         componentName(&daemonSet),
		desiredReplicas:   daemonSet.Status.DesiredNumberScheduled,
		scheduledReplicas: daemonSet.Status.CurrentNumberScheduled,
		updatedReplicas:   daemonSet.Status.UpdatedNumberScheduled,
		availableReplicas: daemonSet.Status.NumberAvailable,
	}
}

func defaultToOne(n *int32) int32 {
	if n != nil {
		return *n
	}

	return 1
}
