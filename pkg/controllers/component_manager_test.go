package controllers

import (
	"context"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestNewComponentManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// when
		sut := NewComponentManager(nil, nil, nil)

		// then
		require.NotNil(t, sut)
	})
}

func Test_componentManager_Install(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.1.0")
		installManagerMock := NewMockInstallManager(t)
		installManagerMock.EXPECT().Install(context.TODO(), component).Return(nil)
		eventRecorderMock := NewMockEventRecorder(t)
		eventRecorderMock.EXPECT().Event(component, "Normal", "Installation", "Starting installation...")

		sut := &componentManager{
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
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.1.0")
		upgradeManagerMock := NewMockUpgradeManager(t)
		upgradeManagerMock.EXPECT().Upgrade(context.TODO(), component).Return(nil)
		eventRecorderMock := NewMockEventRecorder(t)
		eventRecorderMock.EXPECT().Event(component, "Normal", "Upgrade", "Starting upgrade...")

		sut := &componentManager{
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
		component := getComponent("ecosystem", "k8s", "dogu-op", "0.1.0")
		deleteManagerMock := NewMockDeleteManager(t)
		deleteManagerMock.EXPECT().Delete(context.TODO(), component).Return(nil)
		eventRecorderMock := NewMockEventRecorder(t)
		eventRecorderMock.EXPECT().Event(component, "Normal", "Deinstallation", "Starting deinstallation...")

		sut := &componentManager{
			deleteManager: deleteManagerMock,
			recorder:      eventRecorderMock,
		}
		// when
		err := sut.Delete(context.TODO(), component)

		// then
		require.Nil(t, err)
	})
}

func getComponent(namespace string, helmNamespace string, name string, version string) *v1.Component {
	return &v1.Component{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: v1.ComponentSpec{
			Namespace: helmNamespace,
			Name:      name,
			Version:   version,
		}}
}
