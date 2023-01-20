package controllers

import (
	"context"
	k8sv1 "github.com/cloudogu/k8s-component-operator/api/v1"
	"github.com/cloudogu/k8s-component-operator/internal"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"k8s.io/client-go/tools/record"
)

// NewManager is an alias mainly used for testing the main package
var NewManager = NewComponentManager

// DoguManager is a central unit in the process of handling dogu custom resources
// The DoguManager creates, updates and deletes dogus
type DoguManager struct {
	installManager internal.InstallManager
	recorder       record.EventRecorder
}

// NewComponentManager creates a new instance of DoguManager
func NewComponentManager(operatorConfig *config.OperatorConfig) (*DoguManager, error) {
	installManager, err := NewComponentInstallManager(operatorConfig)
	if err != nil {
		return nil, err
	}

	return &DoguManager{
		installManager: installManager,
	}, nil
}

// Install installs component resource.
func (m *DoguManager) Install(ctx context.Context, component *k8sv1.Component) error {
	return m.installManager.Install(ctx, component)
}
