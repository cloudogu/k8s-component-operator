package controllers

import (
	"context"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const k8sDoguOperatorFieldManagerName = "k8s-component-operator"

// componentInstallManager is a central unit in the process of handling the installation process of a custom dogu resource.
type componentInstallManager struct {
	client           client.Client
	doguRegistryData config.DoguRegistryData
}

// NewComponentInstallManager creates a new instance of componentInstallManager.
func NewComponentInstallManager(client client.Client, config *config.OperatorConfig) (*componentInstallManager, error) {
	return &componentInstallManager{
		client:           client,
		doguRegistryData: config.DoguRegistry,
	}, nil
}

// Install installs a given Component Resource.
func (m *componentInstallManager) Install(ctx context.Context, component *k8sv1.Component) error {
	_ = log.FromContext(ctx)

	return nil
}
