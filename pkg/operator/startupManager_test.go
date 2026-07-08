package operator

import (
	"context"
	"testing"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testCtx = context.Background()

func TestNewStartupManager(t *testing.T) {
	t.Run("should create a new startup manager", func(t *testing.T) {
		// given
		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)

		// when
		sut := NewStartupManager(helmClientMock, componentClientMock)

		// then
		require.NotNil(t, sut)
		assert.Equal(t, helmClientMock, sut.helmClient)
		assert.Equal(t, componentClientMock, sut.componentClient)
	})
}

func TestStartupManager_Start(t *testing.T) {
	t.Run("should succeed when no components exist", func(t *testing.T) {
		// given
		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should succeed when components with non-resettable status exist", func(t *testing.T) {
		// given
		component1 := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "component1"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalled},
		}
		component2 := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "component2"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusNotInstalled},
		}

		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*component1, *component2}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark installing components as failed", func(t *testing.T) {
		// given
		installingComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "installing-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalling},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("installing-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*installingComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark upgrading components as failed", func(t *testing.T) {
		// given
		upgradingComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "upgrading-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusUpgrading},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("upgrading-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*upgradingComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark tryToInstall components as failed", func(t *testing.T) {
		// given
		tryToInstallComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "tryToInstall-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusTryToInstall},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("tryToInstall-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*tryToInstallComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark tryToUpgrade components as failed", func(t *testing.T) {
		// given
		tryToUpgradeComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "tryToUpgrade-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusTryToUpgrade},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("tryToUpgrade-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*tryToUpgradeComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark deleting components as failed", func(t *testing.T) {
		// given
		deletingComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "deleting-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusDeleting},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("deleting-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*deletingComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark tryToDelete components as failed", func(t *testing.T) {
		// given
		tryToDeleteComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "tryToDelete-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusTryToDelete},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("tryToDelete-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*tryToDeleteComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should handle multiple components with mixed statuses", func(t *testing.T) {
		// given
		installedComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "installed-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalled},
		}
		installingComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "installing-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalling},
		}
		upgradingComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "upgrading-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusUpgrading},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("installing-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)
		helmClientMock.EXPECT().MarkReleaseAsFailed("upgrading-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*installedComponent, *installingComponent, *upgradingComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil even when listing components fails", func(t *testing.T) {
		// given
		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(nil, assert.AnError)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil even when marking release as failed fails", func(t *testing.T) {
		// given
		installingComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "installing-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalling},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("installing-component", "setting unrecoverable release to failed for the next reconciliation").
			Return(assert.AnError)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*installingComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.Start(testCtx)

		// then
		require.NoError(t, err)
	})
}

func TestStartupManager_setInstallingComponentsToFailed(t *testing.T) {
	t.Run("should fail when listing components fails", func(t *testing.T) {
		// given
		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(nil, assert.AnError)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.setInstallingComponentsToFailed(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, "failed to list components")
	})

	t.Run("should succeed when no components exist", func(t *testing.T) {
		// given
		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.setInstallingComponentsToFailed(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should succeed when only non-resettable components exist", func(t *testing.T) {
		// given
		installedComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "installed-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalled},
		}
		notInstalledComponent := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "notinstalled-component"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusNotInstalled},
		}

		helmClientMock := newMockHelmClient(t)
		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*installedComponent, *notInstalledComponent}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.setInstallingComponentsToFailed(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should mark all resettable statuses as failed", func(t *testing.T) {
		// given
		components := []v1.Component{
			{ObjectMeta: metav1.ObjectMeta{Name: "installing"}, Status: v1.ComponentStatus{Status: v1.ComponentStatusInstalling}},
			{ObjectMeta: metav1.ObjectMeta{Name: "tryToInstall"}, Status: v1.ComponentStatus{Status: v1.ComponentStatusTryToInstall}},
			{ObjectMeta: metav1.ObjectMeta{Name: "upgrading"}, Status: v1.ComponentStatus{Status: v1.ComponentStatusUpgrading}},
			{ObjectMeta: metav1.ObjectMeta{Name: "tryToUpgrade"}, Status: v1.ComponentStatus{Status: v1.ComponentStatusTryToUpgrade}},
			{ObjectMeta: metav1.ObjectMeta{Name: "deleting"}, Status: v1.ComponentStatus{Status: v1.ComponentStatusDeleting}},
			{ObjectMeta: metav1.ObjectMeta{Name: "tryToDelete"}, Status: v1.ComponentStatus{Status: v1.ComponentStatusTryToDelete}},
		}

		helmClientMock := newMockHelmClient(t)
		for _, comp := range components {
			helmClientMock.EXPECT().MarkReleaseAsFailed(comp.Name, "setting unrecoverable release to failed for the next reconciliation").
				Return(nil)
		}

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: components}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.setInstallingComponentsToFailed(testCtx)

		// then
		require.NoError(t, err)
	})

	t.Run("should return nil when marking a release fails", func(t *testing.T) {
		// given
		component1 := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "component1"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusInstalling},
		}
		component2 := &v1.Component{
			ObjectMeta: metav1.ObjectMeta{Name: "component2"},
			Status:     v1.ComponentStatus{Status: v1.ComponentStatusUpgrading},
		}

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().MarkReleaseAsFailed("component1", "setting unrecoverable release to failed for the next reconciliation").
			Return(assert.AnError)
		helmClientMock.EXPECT().MarkReleaseAsFailed("component2", "setting unrecoverable release to failed for the next reconciliation").
			Return(nil)

		componentClientMock := newMockComponentClient(t)
		componentClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).
			Return(&v1.ComponentList{Items: []v1.Component{*component1, *component2}}, nil)

		sut := NewStartupManager(helmClientMock, componentClientMock)

		// when
		err := sut.setInstallingComponentsToFailed(testCtx)

		// then
		require.NoError(t, err)
	})
}
