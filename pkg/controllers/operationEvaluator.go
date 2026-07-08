package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Masterminds/semver/v3"
	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/helm"
	"github.com/cloudogu/k8s-component-operator/pkg/yaml"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type defaultOperationEvaluatorFactory struct {
	recorder       record.EventRecorder
	timeout        time.Duration
	yamlSerializer yaml.Serializer
	reader         configMapRefReader
}

func (f *defaultOperationEvaluatorFactory) NewOperationEvaluator(helmClient helmClient) operationEvaluator {
	return &defaultOperationEvaluator{
		helmClient:     helmClient,
		recorder:       f.recorder,
		timeout:        f.timeout,
		yamlSerializer: f.yamlSerializer,
		reader:         f.reader,
	}
}

type defaultOperationEvaluator struct {
	helmClient     helmClient
	recorder       record.EventRecorder
	timeout        time.Duration
	yamlSerializer yaml.Serializer
	reader         configMapRefReader
}

func (e *defaultOperationEvaluator) EvaluateRequiredOperation(ctx context.Context, component *k8sv1.Component) (operation, error) {
	logger := log.FromContext(ctx)
	if component.DeletionTimestamp != nil && !component.DeletionTimestamp.IsZero() {
		return Delete, nil
	}

	switch component.Status.Status {
	case k8sv1.ComponentStatusNotInstalled, k8sv1.ComponentStatusTryToInstall, k8sv1.ComponentStatusInstalling:
		return Install, nil
	case k8sv1.ComponentStatusInstalled, k8sv1.ComponentStatusTryToUpgrade, k8sv1.ComponentStatusTryToDelete:
		return e.getChangeOperation(ctx, component)
	case k8sv1.ComponentStatusDeleting:
		return Delete, nil
	case k8sv1.ComponentStatusUpgrading:
		return Upgrade, nil
	default:
		logger.Info(fmt.Sprintf("Found unknown operation for component status: %s", component.Status.Status))
		return Ignore, nil
	}
}

func (e *defaultOperationEvaluator) getChangeOperation(ctx context.Context, component *k8sv1.Component) (operation, error) {
	logger := log.FromContext(ctx)

	deployedReleases, err := e.helmClient.ListDeployedReleases()
	if err != nil {
		return "", fmt.Errorf("failed to get deployed helm releases: %w", err)
	}

	for _, deployedRelease := range deployedReleases {
		isComponentToBeChanged := deployedRelease.Name == component.Spec.Name
		targetNamespace := component.Spec.DeployNamespace

		if targetNamespace == "" {
			targetNamespace = component.Namespace
		}

		existsReleaseInTargetNamespace := deployedRelease.Namespace == targetNamespace

		if isComponentToBeChanged {
			logger.Info("Found existing release for reconciled component",
				"releaseNamespace", deployedRelease.Namespace, "targetNamespace", targetNamespace)
			if existsReleaseInTargetNamespace {
				return e.getChangeOperationForRelease(ctx, component, deployedRelease)
			} else {
				e.recorder.Eventf(component, corev1.EventTypeWarning, UpgradeEventReason, "Deploy namespace mismatch (CR: %q; deployed: %q). Deploy namespace declaration is only allowed on install. Revert deploy namespace change to prevent failing upgrade.", targetNamespace, deployedRelease.Namespace)
				return "", fmt.Errorf("component does not exist in target namespace (%q), but in namespace %q", targetNamespace, deployedRelease.Namespace)
			}
		}
	}

	return Ignore, nil
}

func (e *defaultOperationEvaluator) isValuesChanged(ctx context.Context, deployedRelease *release.Release, component *k8sv1.Component) (bool, error) {
	deployedValues, err := e.helmClient.GetReleaseValues(deployedRelease.Name, false)
	if err != nil {
		return false, fmt.Errorf("failed to get values.yaml from release %s: %w", deployedRelease.Name, err)
	}

	chartSpec, err := helm.GetHelmChartSpec(ctx, component, helm.HelmChartCreationOpts{
		HelmClient:     e.helmClient,
		Timeout:        e.timeout,
		YamlSerializer: e.yamlSerializer,
		Reader:         e.reader,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get helm chart spec: %w", err)
	}

	chartSpecValues, err := e.helmClient.GetChartSpecValues(chartSpec)
	if err != nil {
		return false, fmt.Errorf("failed to get values.yaml from component %s: %w", chartSpec.ChartName, err)
	}

	// if no additional values are set, the maps will look like this:
	// deployedValues=map[string]interface {}(nil)
	// chartSpecValues=map[string]interface {}{}
	// this is treated as a difference by DeepEqual, so we have to handle this edge case manually
	if len(deployedValues) == 0 && len(chartSpecValues) == 0 {
		return false, nil
	}

	return !reflect.DeepEqual(deployedValues, chartSpecValues), nil
}

func (e *defaultOperationEvaluator) getChangeOperationForRelease(ctx context.Context, component *k8sv1.Component, release *release.Release) (operation, error) {
	chart := release.Chart
	deployedAppVersion, err := semver.NewVersion(chart.AppVersion())
	if err != nil {
		return "", fmt.Errorf("failed to parse app version %s from helm chart %s: %w", chart.AppVersion(), chart.Name(), err)
	}

	componentVersion, err := semver.NewVersion(component.Spec.Version)
	if err != nil {
		return "", fmt.Errorf("failed to parse component version %s from %s: %w", component.Spec.Version, component.Spec.Name, err)
	}

	if deployedAppVersion.LessThan(componentVersion) {
		return Upgrade, nil
	}

	if deployedAppVersion.GreaterThan(componentVersion) {
		return Downgrade, nil
	}

	isValuesChanged, err := e.isValuesChanged(ctx, release, component)
	if err != nil {
		return "", fmt.Errorf("failed to compare Values.yaml files of component %s: %w", component.Name, err)
	}
	if isValuesChanged {
		return Upgrade, nil
	}

	return Ignore, nil
}
