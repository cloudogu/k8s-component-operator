package controllers

import (
	"context"
	"fmt"
	"time"

	k8sv1 "github.com/cloudogu/k8s-component-lib/api/v1"
	"github.com/go-logr/logr"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

// handlePendingRelease sets the pending release as failed, waits for it to update
func handlePendingRelease(logger logr.Logger, component *k8sv1.Component, ctx context.Context, helmClient helmClient, timeout time.Duration) error {
	logger.Info(fmt.Sprintf("marking pending release for component %q as failed before reinstall", component.Spec.Name))

	err := helmClient.MarkReleaseAsFailed(component.Spec.Name, "failing pending release before reinstall")
	if err != nil {
		return &genericRequeueableError{"failed to mark release as failed", err}
	}

	releaseStatus, err := waitForReleaseStatusUpdate(ctx, timeout, helmClient, component.Spec.Name, func(status helmRelease.Status) bool {
		return !status.IsPending()
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("release status for component %q updated to %q", component.Spec.Name, releaseStatus))

	return nil
}

// handleUninstallingRelease sets the uninstalling release as failed, waits for it to update
func handleUninstallingRelease(logger logr.Logger, component *k8sv1.Component, ctx context.Context, helmClient helmClient, timeout time.Duration) error {
	logger.Info(fmt.Sprintf("marking uninstalling release for component %q as failed before reinstall", component.Spec.Name))

	err := helmClient.MarkReleaseAsFailed(component.Spec.Name, "failing uninstalling release before reinstall")
	if err != nil {
		return &genericRequeueableError{"failed to mark release as failed", err}
	}

	releaseStatus, err := waitForReleaseStatusUpdate(ctx, timeout, helmClient, component.Spec.Name, func(status helmRelease.Status) bool {
		return status != helmRelease.StatusUninstalling
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("release status for component %q updated to %q", component.Spec.Name, releaseStatus))

	return nil
}

func waitForReleaseStatusUpdate(
	ctx context.Context,
	timeout time.Duration,
	helmClient helmClient,
	releaseName string,
	statusCheckFn func(status helmRelease.Status) bool,
) (helmRelease.Status, error) {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-waitCtx.Done():
			return "", &genericRequeueableError{"Timeout while getting the Helm release", waitCtx.Err()}
		case <-time.After(2 * time.Second):
			release, err := helmClient.GetRelease(releaseName)
			if err != nil {
				return "", &genericRequeueableError{"Error while getting the Helm release", err}
			}

			if statusCheckFn(release.Info.Status) {
				return release.Info.Status, nil
			}
		}
	}
}

func markReleaseAsFailed() error {
	return nil
}
