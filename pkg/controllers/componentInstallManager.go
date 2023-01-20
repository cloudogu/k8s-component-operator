package controllers

import (
	"context"
	"fmt"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/mittwald/go-helm-client"
	repo "helm.sh/helm/v3/pkg/repo"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const k8sDoguOperatorFieldManagerName = "k8s-component-operator"

// componentInstallManager is a central unit in the process of handling the installation process of a custom dogu resource.
type componentInstallManager struct {
	namespace string
}

// NewComponentInstallManager creates a new instance of componentInstallManager.
func NewComponentInstallManager(config *config.OperatorConfig) (*componentInstallManager, error) {
	return &componentInstallManager{
		namespace: config.Namespace,
	}, nil
}

// Install installs a given Component Resource.
func (cim *componentInstallManager) Install(ctx context.Context, component *k8sv1.Component) error {
	logger := log.FromContext(ctx)

	logger.Info("Creating helm client...")
	opt := &helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        cim.namespace, // Change this to the namespace you wish the client to operate in.
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
			Linting:          true, // Change this to false if you don't want linting.
			DebugLog:         logger.Info,
		},
		RestConfig: ctrl.GetConfigOrDie(),
	}

	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		return fmt.Errorf("failed to create helm client: %w", err)
	}

	logger.Info("Add helm repo...")
	devRepo := repo.Entry{
		Name:                  "k8s",
		URL:                   "http://chartmuseum.ecosystem.svc.cluster.local:8080/k8s",
		InsecureSkipTLSverify: true,
		PassCredentialsAll:    false,
	}
	err = helmClient.AddOrUpdateChartRepo(devRepo)
	if err != nil {
		return fmt.Errorf("failed to add helm repository: %w", err)
	}

	err = helmClient.UpdateChartRepos()
	if err != nil {
		return fmt.Errorf("failed to update chart repositories: %w", err)
	}

	logger.Info("Install helm chart...")
	chartSpec := &helmclient.ChartSpec{
		ReleaseName: component.Spec.Name,
		ChartName:   fmt.Sprintf("k8s/%s", component.Spec.Name),
		Namespace:   cim.namespace,
		ValuesYaml:  "",
		Version:     component.Spec.Version,
	}
	_, err = helmClient.InstallOrUpgradeChart(ctx, chartSpec, nil)
	if err != nil {
		return fmt.Errorf("failed to install chart: %w", err)
	}

	logger.Info("Done...")

	return nil
}
