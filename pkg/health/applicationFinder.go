package health

import (
	"context"
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

type defaultApplicationFinder struct {
	appsClient appsV1Client
}

func (af *defaultApplicationFinder) findComponentApplications(ctx context.Context, componentName string, namespace string) (*appsv1.DeploymentList, *appsv1.StatefulSetList, *appsv1.DaemonSetList, error) {
	componentLabelSelector := labels.ValidatedSetSelector{
		v1.ComponentNameLabelKey: componentName,
	}.String()

	var errs []error
	deploymentList, err := af.appsClient.Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: componentLabelSelector})
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list deployments for component %q: %w", componentName, err))
	}

	statefulSetList, err := af.appsClient.StatefulSets(namespace).List(ctx, metav1.ListOptions{LabelSelector: componentLabelSelector})
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list stateful sets for component %q: %w", componentName, err))
	}

	daemonSetList, err := af.appsClient.DaemonSets(namespace).List(ctx, metav1.ListOptions{LabelSelector: componentLabelSelector})
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to list daemon sets for component %q: %w", componentName, err))
	}

	if len(errs) > 0 {
		return nil, nil, nil, errors.Join(errs...)
	}

	return deploymentList, statefulSetList, daemonSetList, nil
}
