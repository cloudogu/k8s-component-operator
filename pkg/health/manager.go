package health

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/util"
)

type defaultManager struct {
	applicationFinder
	componentRepo
}

func newManager(namespace string, clientSet ecosystemClientSet) *defaultManager {
	return &defaultManager{
		applicationFinder: &defaultApplicationFinder{appsClient: clientSet.AppsV1()},
		componentRepo:     &defaultComponentRepo{client: clientSet.ComponentV1Alpha1().Components(namespace)},
	}
}

func (m *defaultManager) UpdateComponentHealth(ctx context.Context, componentName string, namespace string) error {
	deploymentList, statefulSetList, daemonSetList, err := m.findComponentApplications(ctx, componentName, namespace)
	if err != nil {
		return fmt.Errorf("failed to find applications for component %q: %w", componentName, err)
	}

	healthStatus := m.componentHealthStatus(ctx, deploymentList, statefulSetList, daemonSetList)

	component, err := m.get(ctx, componentName)
	if err != nil {
		return fmt.Errorf("failed to get component %q: %w", componentName, err)
	}

	err = m.updateHealthStatus(ctx, component, healthStatus)
	if err != nil {
		return fmt.Errorf("failed to update health status for component %q: %w", componentName, err)
	}

	return nil
}

func (m *defaultManager) componentHealthStatus(ctx context.Context, deployments *appsv1.DeploymentList, statefulSets *appsv1.StatefulSetList, daemonSets *appsv1.DaemonSetList) v1.HealthStatus {
	logger := log.FromContext(ctx).WithName("componentHealthStatus")

	states := make([]state, 0, len(deployments.Items)+len(statefulSets.Items)+len(daemonSets.Items))
	states = append(states, util.Map(deployments.Items, deploymentToState)...)
	states = append(states, util.Map(statefulSets.Items, statefulSetToState)...)
	states = append(states, util.Map(daemonSets.Items, daemonSetToState)...)

	for _, state := range states {
		if !state.IsAvailable() {
			logger.Info(fmt.Sprintf("%s %q of component %q is not (yet?) available", state.kind, state.name, state.component))
		}
	}

	componentAvailable := util.Reduce(states, true, func(value state, acc bool) bool {
		return value.IsAvailable() && acc
	})

	if componentAvailable {
		return v1.AvailableHealthStatus
	}
	return v1.UnavailableHealthStatus
}
