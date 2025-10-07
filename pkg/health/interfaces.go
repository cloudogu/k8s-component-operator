package health

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-lib/client"
)

type RegistrableController interface {
	SetupWithManager(mgr ctrl.Manager) error
}

type ecosystemClientSet interface {
	client.ComponentEcosystemInterface
}

type componentClient interface {
	client.ComponentInterface
}

type appsV1Client interface {
	appsv1client.AppsV1Interface
}

type ComponentManager interface {
	UpdateComponentHealth(ctx context.Context, componentName string, namespace string) error
	UpdateComponentHealthWithInstalledVersion(ctx context.Context, componentName string, namespace string, version string) error
	UpdateComponentHealthAll(ctx context.Context) error
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
	list(ctx context.Context) (*v1.ComponentList, error)
	updateCondition(ctx context.Context, component *v1.Component, statusFn func() (v1.HealthStatus, error), version string) error
}

// interfaces for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type componentV1Alpha1Client interface {
	client.ComponentV1Alpha1Interface
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
