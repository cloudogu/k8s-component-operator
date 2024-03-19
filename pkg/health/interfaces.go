package health

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

type RegistrableController interface {
	SetupWithManager(mgr ctrl.Manager) error
}

type ecosystemClientSet interface {
	ecosystem.ComponentEcosystemInterface
}

type componentClient interface {
	ecosystem.ComponentInterface
}

type appsV1Client interface {
	appsv1client.AppsV1Interface
}

type ComponentManager interface {
	UpdateComponentHealth(ctx context.Context, componentName string, namespace string) error
	UpdateComponentHealthWithInstalledVersion(ctx context.Context, componentName string, namespace string, version string) error
}

type applicationFinder interface {
	findComponentApplications(
		ctx context.Context,
		componentName string,
		namespace string,
	) (
		*appsv1.DeploymentList,
		*appsv1.StatefulSetList,
		*appsv1.DaemonSetList,
		error,
	)
}

type componentRepo interface {
	get(ctx context.Context, name string) (*v1.Component, error)
	updateCondition(ctx context.Context, component *v1.Component, status v1.HealthStatus, version string) error
}

// interfaces for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type componentV1Alpha1Client interface {
	ecosystem.ComponentV1Alpha1Interface
}

//nolint:unused
//goland:noinspection GoUnusedType
type deploymentClient interface {
	appsv1client.DeploymentInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type statefulSetClient interface {
	appsv1client.StatefulSetInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type daemonSetClient interface {
	appsv1client.DaemonSetInterface
}

//nolint:unused
//goland:noinspection GoUnusedType
type controllerManager interface {
	ctrl.Manager
}
