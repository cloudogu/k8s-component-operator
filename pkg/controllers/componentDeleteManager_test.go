package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"helm.sh/helm/v3/pkg/release"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sv1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
)

func TestNewComponentDeleteManager(t *testing.T) {
	t.Run("should create new componentDeleteManager", func(t *testing.T) {
		mockComponentClient := newMockComponentInterface(t)
		mockHelmClient := newMockHelmClient(t)

		manager := NewComponentDeleteManager(mockComponentClient, mockHelmClient)

		assert.NotNil(t, manager)
		assert.Equal(t, mockHelmClient, manager.helmClient)
		assert.Equal(t, mockComponentClient, manager.componentClient)
	})
}

func Test_componentDeleteManager_Delete(t *testing.T) {
	t.Run("should delete component", func(t *testing.T) {
		ctx := context.Background()
		componentName := "testComponent"
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      componentName,
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusDeleting(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().Get(ctx, component.Name, v1.GetOptions{}).Return(component, nil)
		mockComponentClient.EXPECT().RemoveFinalizer(ctx, component, k8sv1.FinalizerName).Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().Uninstall(component.Spec.Name).Return(nil)
		mockHelmClient.EXPECT().ListDeployedReleases().Return([]*release.Release{{
			Name: componentName,
		}}, nil)

		manager := &componentDeleteManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Delete(ctx, component)

		require.NoError(t, err)
	})

	t.Run("should fail to delete component on error while updating status", func(t *testing.T) {
		ctx := context.Background()
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      "testComponent",
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusDeleting(ctx, component).Return(component, assert.AnError)

		mockHelmClient := newMockHelmClient(t)

		manager := &componentDeleteManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Delete(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to update status-deleting for component testComponent:")
	})

	t.Run("should fail to delete component on error while uninstalling chart", func(t *testing.T) {
		ctx := context.Background()
		componentName := "testComponent"
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      componentName,
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusDeleting(ctx, component).Return(component, nil)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().Uninstall(component.Spec.Name).Return(assert.AnError)
		mockHelmClient.EXPECT().ListDeployedReleases().Return([]*release.Release{{
			Name: componentName,
		}}, nil)

		manager := &componentDeleteManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Delete(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to uninstall chart with name testComponent:")
	})

	t.Run("should fail to delete component on error while removing finalizer", func(t *testing.T) {
		ctx := context.Background()
		componentName := "testComponent"
		component := &k8sv1.Component{
			Spec: k8sv1.ComponentSpec{
				Namespace: "ecosystem",
				Name:      componentName,
				Version:   "1.0",
			},
			Status: k8sv1.ComponentStatus{Status: "installed"},
		}

		mockComponentClient := newMockComponentInterface(t)
		mockComponentClient.EXPECT().UpdateStatusDeleting(ctx, component).Return(component, nil)
		mockComponentClient.EXPECT().Get(ctx, component.Name, v1.GetOptions{}).Return(component, nil)
		mockComponentClient.EXPECT().RemoveFinalizer(ctx, component, k8sv1.FinalizerName).Return(component, assert.AnError)

		mockHelmClient := newMockHelmClient(t)
		mockHelmClient.EXPECT().Uninstall(component.Spec.Name).Return(nil)
		mockHelmClient.EXPECT().ListDeployedReleases().Return([]*release.Release{{
			Name: componentName,
		}}, nil)

		manager := &componentDeleteManager{
			componentClient: mockComponentClient,
			helmClient:      mockHelmClient,
		}
		err := manager.Delete(ctx, component)

		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to remove finalizer for component testComponent:")
	})
}
