package health

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type StartupHandler struct {
	manager ComponentManager
}

func NewStartupHandler(namespace string, clientSet ecosystemClientSet) *StartupHandler {
	return &StartupHandler{manager: NewManager(namespace, clientSet)}
}

func (s *StartupHandler) Start(ctx context.Context) error {
	log.FromContext(ctx).
		WithName("health startup handler").
		Info("updating health of all components on startup")
	err := s.manager.UpdateComponentHealthAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to update health of components on startup: %w", err)
	}

	return nil
}
