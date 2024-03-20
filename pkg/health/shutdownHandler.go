package health

import (
	"context"
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

const componentOperatorName = "k8s-component-operator"

type ShutdownHandler struct {
	repo componentRepo
}

func NewShutdownHandler(client componentClient) *ShutdownHandler {
	return &ShutdownHandler{repo: &defaultComponentRepo{client: client}}
}

func (s *ShutdownHandler) Start(ctx context.Context) error {
	logger := log.FromContext(ctx).WithName("health shutdown handler")
	logger.Info("health shutdown handler registered, waiting for shutdown")
	<-ctx.Done()
	logger.Info("shutdown detected, handling health status")

	// context is done, we need a new one
	ctx = context.WithoutCancel(ctx)
	return s.handle(ctx)
}

func (s *ShutdownHandler) handle(ctx context.Context) error {
	components, err := s.repo.list(ctx)
	if err != nil {
		return err
	}

	var errs []error
	for _, component := range components.Items {
		// set health status of other components to unknown
		var healthStatus = v1.UnknownHealthStatus

		if component.Name == componentOperatorName {
			// set component operator health status to unavailable
			healthStatus = v1.UnavailableHealthStatus
		}

		err := s.repo.updateCondition(ctx, &component, healthStatus, noVersionChange)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to set health status of %q to %q: %w", component.Name, healthStatus, err))
		}
	}
	return errors.Join(errs...)
}
