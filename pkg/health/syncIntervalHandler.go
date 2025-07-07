package health

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

const componentOperatorName = "k8s-component-operator"

type SyncIntervalHandler struct {
	manager            ComponentManager
	repo               componentRepo
	healthSyncInterval time.Duration
}

func NewSyncIntervalHandler(namespace string, clientSet ecosystemClientSet, healthSyncInterval time.Duration) *SyncIntervalHandler {
	return &SyncIntervalHandler{
		manager:            NewManager(namespace, clientSet),
		repo:               &defaultComponentRepo{client: clientSet.ComponentV1Alpha1().Components(namespace)},
		healthSyncInterval: healthSyncInterval,
	}
}

func (s *SyncIntervalHandler) Start(ctx context.Context) error {
	logger := log.FromContext(ctx).
		WithName("health sync interval handler")
	logger.Info(fmt.Sprintf("started regularly syncing health of all components with interval %s", s.healthSyncInterval))

	ticker := time.NewTicker(s.healthSyncInterval)
	defer ticker.Stop()

	var errs []error
	for {
		select {
		case <-ctx.Done():
			logger.Info("shutdown detected, handling health status")

			// context is done, we need a new one
			shutdownCtx := context.WithoutCancel(ctx)
			err := s.handleShutdown(shutdownCtx)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to correctly handle health status during shutdown: %w", err))
			}

			if len(errs) > 0 {
				return fmt.Errorf("some errors occurred when regularly syncing health of all components: %w", errors.Join(errs...))
			}

			return nil
		case <-ticker.C:
			logger.Info("regularly syncing health of all components...")
			err := s.manager.UpdateComponentHealthAll(ctx)
			if err != nil {
				logger.Error(err, "failed to regularly sync health of all components")
				errs = append(errs, err)
			}
		}
	}
}

func (s *SyncIntervalHandler) handleShutdown(ctx context.Context) error {
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

		err := s.repo.updateCondition(ctx, &component, func() (v1.HealthStatus, error) {
			return healthStatus, nil
		}, noVersionChange)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to set health status of %q to %q: %w", component.Name, healthStatus, err))
		}
	}
	return errors.Join(errs...)
}
