package controllers

import (
	"context"
	"testing"

	v1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewComponentManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// when
		sut := NewComponentManager(nil, nil, nil, nil, defaultHelmClientTimeoutMins, nil)

		// then
		require.NotNil(t, sut)
	})
}

func Test_componentManager_Install(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		installManagerMock := newMockInstallManager(t)
		installManagerMock.EXPECT().Install(context.TODO(), component).Return(nil)
		eventRecorderMock := newMockEventRecorder(t)
		eventRecorderMock.EXPECT().Event(component, "Normal", "Installation", "Starting installation...")

		sut := &DefaultComponentManager{
			installManager: installManagerMock,
			recorder:       eventRecorderMock,
		}

		// when
		err := sut.Install(context.TODO(), component)

		// then
		require.Nil(t, err)
	})
}

func Test_componentManager_Upgrade(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		upgradeManagerMock := newMockUpgradeManager(t)
		upgradeManagerMock.EXPECT().Upgrade(context.TODO(), component).Return(nil)
		eventRecorderMock := newMockEventRecorder(t)
		eventRecorderMock.EXPECT().Event(component, "Normal", "Upgrade", "Starting upgrade...")

		sut := &DefaultComponentManager{
			upgradeManager: upgradeManagerMock,
			recorder:       eventRecorderMock,
		}
		// when
		err := sut.Upgrade(context.TODO(), component)

		// then
		require.Nil(t, err)
	})
}

func Test_componentManager_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "", "dogu-op", "0.1.0")
		deleteManagerMock := newMockDeleteManager(t)
		deleteManagerMock.EXPECT().Delete(context.TODO(), component).Return(nil)
		eventRecorderMock := newMockEventRecorder(t)
		eventRecorderMock.EXPECT().Event(component, "Normal", "Deinstallation", "Starting deinstallation...")

		sut := &DefaultComponentManager{
			deleteManager: deleteManagerMock,
			recorder:      eventRecorderMock,
		}
		// when
		err := sut.Delete(context.TODO(), component)

		// then
		require.Nil(t, err)
	})
}

func getComponent(namespace string, helmNamespace string, deployNamespace string, name string, version string) *v1.Component {
	return &v1.Component{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: v1.ComponentSpec{
			Namespace:       helmNamespace,
			DeployNamespace: deployNamespace,
			Name:            name,
			Version:         version,
		}}
}
