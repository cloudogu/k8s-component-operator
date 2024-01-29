package health

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewController(t *testing.T) {
	appsV1Mock := newMockAppsV1Client(t)
	componentMock := newMockComponentClient(t)
	componentV1Alpha1Mock := newMockComponentV1Alpha1Client(t)
	componentV1Alpha1Mock.EXPECT().Components(testNamespace).Return(componentMock)
	clientSetMock := newMockEcosystemClientSet(t)
	clientSetMock.EXPECT().AppsV1().Return(appsV1Mock)
	clientSetMock.EXPECT().ComponentV1Alpha1().Return(componentV1Alpha1Mock)

	assert.NotEmpty(t, NewController(testNamespace, clientSetMock))
}

func Test_metaController_SetupWithManager(t *testing.T) {
	anError2 := fmt.Errorf("another error for testing: %w", assert.AnError)

	tests := []struct {
		name          string
		controllersFn func(t *testing.T) []RegistrableController
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "should fail with single error",
			controllersFn: func(t *testing.T) []RegistrableController {
				controller1 := NewMockRegistrableController(t)
				controller1.EXPECT().SetupWithManager(mock.Anything).Return(assert.AnError)
				controller2 := NewMockRegistrableController(t)
				controller2.EXPECT().SetupWithManager(mock.Anything).Return(nil)
				return []RegistrableController{controller1, controller2}
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError) &&
					assert.ErrorContains(t, err, "failed to setup controller \"*health.MockRegistrableController\"")
			},
		},
		{
			name: "should fail with multiple errors",
			controllersFn: func(t *testing.T) []RegistrableController {
				controller1 := NewMockRegistrableController(t)
				controller1.EXPECT().SetupWithManager(mock.Anything).Return(assert.AnError)
				controller2 := NewMockRegistrableController(t)
				controller2.EXPECT().SetupWithManager(mock.Anything).Return(anError2)
				return []RegistrableController{controller1, controller2}
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, assert.AnError) &&
					assert.ErrorIs(t, err, anError2) &&
					assert.ErrorContains(t, err, "failed to setup controller \"*health.MockRegistrableController\"")
			},
		},
		{
			name: "should succeed",
			controllersFn: func(t *testing.T) []RegistrableController {
				controller1 := NewMockRegistrableController(t)
				controller1.EXPECT().SetupWithManager(mock.Anything).Return(nil)
				controller2 := NewMockRegistrableController(t)
				controller2.EXPECT().SetupWithManager(mock.Anything).Return(nil)
				return []RegistrableController{controller1, controller2}
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := metaController{
				controllers: tt.controllersFn(t),
			}
			tt.wantErr(t, m.SetupWithManager(nil))
		})
	}
}
