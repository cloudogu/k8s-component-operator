package health

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStartupHandler(t *testing.T) {
	// given
	clientSetMock := newMockEcosystemClientSet(t)
	appsV1Mock := newMockAppsV1Client(t)
	clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
	componentV1Mock := newMockComponentV1Alpha1Client(t)
	componentMock := newMockComponentClient(t)
	componentV1Mock.EXPECT().Components(testNamespace).Return(componentMock)
	clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1Mock)

	// when
	actual := NewStartupHandler(testNamespace, clientSetMock)

	// then
	assert.NotEmpty(t, actual)
	assert.NotEmpty(t, actual.manager)
}

func TestStartupHandler_Start(t *testing.T) {
	tests := []struct {
		name      string
		managerFn func(t *testing.T) ComponentManager
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "should fail",
			managerFn: func(t *testing.T) ComponentManager {
				managerMock := NewMockComponentManager(t)
				managerMock.EXPECT().UpdateComponentHealthAll(testCtx).Return(assert.AnError)
				return managerMock
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError, i) &&
					assert.ErrorContains(t, err, "failed to update health of components on startup", i)
			},
		},
		{
			name: "should succeed",
			managerFn: func(t *testing.T) ComponentManager {
				managerMock := NewMockComponentManager(t)
				managerMock.EXPECT().UpdateComponentHealthAll(testCtx).Return(nil)
				return managerMock
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StartupHandler{
				manager: tt.managerFn(t),
			}
			tt.wantErr(t, s.Start(testCtx))
		})
	}
}
