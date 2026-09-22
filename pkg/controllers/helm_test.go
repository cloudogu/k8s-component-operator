package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	helm "helm.sh/helm/v3/pkg/release"
)

func TestMarkUninstallingReleaseAsFailed(t *testing.T) {
	t.Run("A Helm release status is marked as 'failed'", func(t *testing.T) {
		logger, componentName, reason := fixtureForMarkUninstallingReleaseAsFailedTest()

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().
			MarkReleaseAsFailed(componentName, reason).
			Return(nil)
		helmRelease := &helm.Release{
			Info: &helm.Info{
				Status: helm.StatusFailed,
			},
		}
		helmClientMock.EXPECT().
			GetRelease(componentName).
			Return(helmRelease, nil)

		err := markUninstallingReleaseAsFailed(logger, componentName, context.Background(), helmClientMock, 10*time.Second)

		assert.NoError(t, err)
	})

	t.Run("Returns an RequeueableError if setting the release status fails", func(t *testing.T) {
		logger, componentName, reason := fixtureForMarkUninstallingReleaseAsFailedTest()

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().
			MarkReleaseAsFailed(componentName, reason).
			Return(assert.AnError)

		err := markUninstallingReleaseAsFailed(logger, componentName, context.Background(), helmClientMock, 10*time.Second)

		assert.Error(t, err)
		assert.IsType(t, err, &genericRequeueableError{})
	})

	t.Run("Returns an RequeueableError if getting the release status fails", func(t *testing.T) {
		logger, componentName, reason := fixtureForMarkUninstallingReleaseAsFailedTest()

		helmClientMock := newMockHelmClient(t)
		helmClientMock.EXPECT().
			MarkReleaseAsFailed(componentName, reason).
			Return(nil)
		helmClientMock.EXPECT().
			GetRelease(componentName).
			Return(nil, assert.AnError)

		err := markUninstallingReleaseAsFailed(logger, componentName, context.Background(), helmClientMock, 10*time.Second)

		assert.Error(t, err)
		assert.IsType(t, err, &genericRequeueableError{})
	})
}

func fixtureForMarkUninstallingReleaseAsFailedTest() (logr.Logger, string, string) {
	logger := logr.Discard()
	componentName := "dogu-operator"
	reason := "failing uninstalling release before reinstall"
	return logger, componentName, reason
}
