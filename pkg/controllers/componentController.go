package controllers

import (
	"context"
	"fmt"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-component-operator/pkg/api/ecosystem"
	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	helmclient "github.com/mittwald/go-helm-client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type operation string

const (
	Install = operation("Install")
	Upgrade = operation("Upgrade")
	Delete  = operation("Delete")
	Ignore  = operation("Ignore")
)

// ComponentManager abstracts the simple component operations in a k8s CES.
type ComponentManager interface {
	InstallManager
	DeleteManager
	UpgradeManager
}

// componentReconciler watches every Component object in the cluster and handles them accordingly.
type componentReconciler struct {
	componentClient  ecosystem.ComponentInterface
	recorder         record.EventRecorder
	componentManager ComponentManager
	helmClient       helmclient.Client
}

// NewComponentReconciler creates a new component reconciler.
func NewComponentReconciler(clientset ecosystem.ComponentInterface, helmClient helmclient.Client, recorder record.EventRecorder, config *config.OperatorConfig) *componentReconciler {
	return &componentReconciler{
		recorder:         recorder,
		componentManager: NewComponentManager(config, clientset, helmClient),
		helmClient:       helmClient,
	}
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *componentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile this crd")

	component, err := r.componentClient.Get(ctx, req.Name, v1.GetOptions{})
	if err != nil {
		logger.Info(fmt.Sprintf("failed to get component %+v: %s", req, err))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info(fmt.Sprintf("Component %+v has been found", req))

	operation, err := r.evaluateRequiredOperation(ctx, component)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to evaluate required operation: %w", err)
	}
	logger.Info(fmt.Sprintf("Required operation is %s", operation))

	switch operation {
	case Install:
		return ctrl.Result{}, r.componentManager.Install(ctx, component)
	case Delete:
		return ctrl.Result{}, r.componentManager.Delete(ctx, component)
	case Upgrade:
		return ctrl.Result{}, r.componentManager.Upgrade(ctx, component)
	case Ignore:
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, nil
	}
}

func (r *componentReconciler) evaluateRequiredOperation(ctx context.Context, component *k8sv1.Component) (operation, error) {
	logger := log.FromContext(ctx)
	if component.DeletionTimestamp != nil && !component.DeletionTimestamp.IsZero() {
		return Delete, nil
	}

	switch component.Status.Status {
	case k8sv1.ComponentStatusNotInstalled:
		return Install, nil
	case k8sv1.ComponentStatusInstalled:
		upgrade, err := r.checkUpgradeAbility(component)
		if err != nil {
			return "", err
		}

		if upgrade {
			return Upgrade, nil
		}
		return Ignore, nil
	case k8sv1.ComponentStatusInstalling:
		return Ignore, nil
	case k8sv1.ComponentStatusDeleting:
		return Ignore, nil
	default:
		logger.Info(fmt.Sprintf("Found unknown operation for component status: %s", component.Status.Status))
		return Ignore, nil
	}
}

func (r *componentReconciler) checkUpgradeAbility(component *k8sv1.Component) (bool, error) {
	deployedReleases, err := r.helmClient.ListDeployedReleases()
	if err != nil {
		return false, fmt.Errorf("failed to get deployed helm releases: %w", err)
	}

	for _, release := range deployedReleases {
		// This will allow a namespace switch e. g. k8s/dogu-operator -> k8s-testing/dogu-operator.
		if release.Name == component.Spec.Name && release.Namespace == component.Namespace {
			chart := release.Chart
			deployedVersion, err := core.ParseVersion(chart.AppVersion())
			if err != nil {
				return false, fmt.Errorf("failed to parse app version %s from helm chart %s: %w", chart.AppVersion(), chart.Name(), err)
			}

			componentVersion, err := core.ParseVersion(component.Spec.Version)
			if err != nil {
				return false, fmt.Errorf("failed to parse component version %s from %s: %w", component.Spec.Version, component.Spec.Name, err)
			}

			if deployedVersion.IsOlderThan(componentVersion) {
				return true, nil
			}

			if deployedVersion.IsNewerThan(componentVersion) {
				return false, fmt.Errorf("downgrades are not allowed")
			}

			break
		}
	}

	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *componentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1.Component{}).
		Complete(r)
}
