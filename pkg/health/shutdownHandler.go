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
	logger := log.FromContext(ctx).WithName("shutdown handler")
	logger.Info("waiting for shutdown")
	<-ctx.Done()
	logger.Info("shutdown detected, handling health status")

	// context is done, we need a new one
	ctx = context.WithoutCancel(ctx)
	return s.handle(ctx)
}

func (s *ShutdownHandler) handle(ctx context.Context) error {
	// set component operator health status to unavailable
	componentOperator, err := s.repo.get(ctx, componentOperatorName)
	if err != nil {
		return fmt.Errorf("failed to get component for %q: %w", componentOperatorName, err)
	}

	err = s.repo.updateHealthStatus(ctx, componentOperator, v1.UnavailableHealthStatus)
	if err != nil {
		return fmt.Errorf("failed to set health status of %q to %q: %w", componentOperatorName, v1.UnavailableHealthStatus, err)
	}

	// set health status of other components to unknown
	components, err := s.repo.list(ctx)
	if err != nil {
		return err
	}

	var errs []error
	for _, component := range components.Items {
		if component.Name == componentOperatorName {
			continue
		}

		err := s.repo.updateHealthStatus(ctx, &component, v1.UnknownHealthStatus)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to set health status of %q to %q: %w", component.Name, v1.UnknownHealthStatus, err))
		}
	}
	return errors.Join(errs...)
}
