package controllers

import (
	"context"
	"fmt"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/go-logr/logr"
)

// handlePendingRelease sets the pending release as failed, waits for it to update
func handlePendingRelease(logger logr.Logger, component *k8sv1.Component, helmCtx context.Context, helmClient helmClient, timeout time.Duration) error {
	logger.Info(fmt.Sprintf("marking pending release for component %q as failed before reinstall", component.Spec.Name))
	err := helmClient.MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall")
	if err != nil {
		return &genericRequeueableError{"failed to mark release as failed", err}
	}
	waitCtx, cancel := context.WithTimeout(helmCtx, timeout)
	defer cancel()

	done := false
	for !done {
		select {
		case <-waitCtx.Done():
			return &genericRequeueableError{
				"timed out waiting for release status update after marking as failed",
				waitCtx.Err(),
			}
		case <-time.After(2 * time.Second):
			updatedRelease, getErr := helmClient.GetRelease(component.Spec.Name)
			if getErr != nil {
				return &genericRequeueableError{
					"failed to get release while waiting for status update",
					getErr,
				}
			}

			if !updatedRelease.Info.Status.IsPending() {
				logger.Info(fmt.Sprintf("release status for component %q updated to %q", component.Spec.Name, updatedRelease.Info.Status))
				done = true
			}
		}
	}
	return nil
}
