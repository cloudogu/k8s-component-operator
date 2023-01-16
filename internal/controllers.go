package internal

import (
	"context"
	v1 "github.com/cloudogu/k8s-component-operator/api/v1"
)

// InstallManager includes functionality to install components in the cluster.
type InstallManager interface {
	// Install installs a component resource.
	Install(ctx context.Context, component *v1.Component) error
}

// ComponentManager abstracts the simple component operations in a k8s CES.
type ComponentManager interface {
	InstallManager
}
