package health

import (
	"context"
	"fmt"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/cloudogu/retry-lib/retry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const noVersionChange = "##NO_VERSION_CHANGE##"

type defaultComponentRepo struct {
	client componentClient
}

func (cr *defaultComponentRepo) list(ctx context.Context) (*v1.ComponentList, error) {
	list, err := cr.client.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list components: %w", err)
	}

	return list, nil
}

func (cr *defaultComponentRepo) get(ctx context.Context, name string) (*v1.Component, error) {
	component, err := cr.client.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get component %q: %w", name, err)
	}

	return component, nil
}

func (cr *defaultComponentRepo) updateCondition(ctx context.Context, component *v1.Component, statusFn func() (v1.HealthStatus, error), version string) error {
	return retry.OnConflict(func() error {
		component, err := cr.get(ctx, component.Name)
		if err != nil {
			return err
		}

		status, err := statusFn()
		if err != nil {
			return err
		}

		component.Status.Health = status
		if version != noVersionChange {
			component.Status.InstalledVersion = version
		}

		_, err = cr.client.UpdateStatus(ctx, component, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update component %q: %w", component.Name, err)
		}

		return nil
	})
}
