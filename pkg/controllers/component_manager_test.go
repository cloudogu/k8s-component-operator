package controllers

import (
	"context"
	v1 "github.com/cloudogu/k8s-component-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-component-operator/pkg/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewComponentManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// when
		sut := NewComponentManager(&config.OperatorConfig{}, nil, nil)

		// then
		require.NotNil(t, sut)
	})
}

func Test_componentManager_Install(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		component := getComponent("ecosystem", "dogu-op", "0.1.0")
		installManagerMock := NewMockInstallManager(t)
		installManagerMock.EXPECT().Install(context.TODO(), component).Return(nil)

		sut := &componentManager{
			installManager: installManagerMock,
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
		component := getComponent("ecosystem", "dogu-op", "0.1.0")
		upgradeManagerMock := NewMockUpgradeManager(t)
		upgradeManagerMock.EXPECT().Upgrade(context.TODO(), component).Return(nil)

		sut := &componentManager{
			upgradeManager: upgradeManagerMock,
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
		component := getComponent("ecosystem", "dogu-op", "0.1.0")
		deleteManagerMock := NewMockDeleteManager(t)
		deleteManagerMock.EXPECT().Delete(context.TODO(), component).Return(nil)

		sut := &componentManager{
			deleteManager: deleteManagerMock,
		}
		// when
		err := sut.Delete(context.TODO(), component)

		// then
		require.Nil(t, err)
	})
}

func getComponent(namespace string, name string, version string) *v1.Component {
	return &v1.Component{Spec: v1.ComponentSpec{
		Namespace: namespace,
		Name:      name,
		Version:   version,
	}}
}
