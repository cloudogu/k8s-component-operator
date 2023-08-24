package helm

import (
	"errors"
	"fmt"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"

	"github.com/Masterminds/semver/v3"
)

type dependencyChecker interface {
	// CheckSatisfied validates that all dependencies are installed in the required version.
	CheckSatisfied(dependencies []*chart.Dependency, deployedReleases []*release.Release) error
}

type installedDependencyChecker struct{}

// CheckSatisfied validates that all dependencies are installed in the required version.
func (d *installedDependencyChecker) CheckSatisfied(dependencies []*chart.Dependency, deployedReleases []*release.Release) error {
	var errs []error
	for _, dependency := range dependencies {
		isInstalled := false
		for _, deployedRelease := range deployedReleases {
			if dependency.Name == deployedRelease.Chart.Name() {
				isInstalled = true
				err := checkVersion(dependency, deployedRelease.Chart)
				errs = append(errs, err)

				break
			}
		}

		if !isInstalled {
			errs = append(errs, fmt.Errorf("dependency %s with version %s is not installed", dependency.Name, dependency.Version))
		}
	}

	return errors.Join(errs...)
}

func checkVersion(dependency *chart.Dependency, deployedChart *chart.Chart) error {
	constraint, err := semver.NewConstraint(dependency.Version)
	if err != nil {
		return fmt.Errorf("failed to parse constraint for dependency %s with version requirement %s: %w", dependency.Name, dependency.Version, err)
	}

	version, err := semver.NewVersion(deployedChart.Metadata.Version)
	if err != nil {
		return fmt.Errorf("failed to parse version of installed component %s with version %s: %w", deployedChart.Metadata.Name, deployedChart.Metadata.Version, err)
	}

	isSatisfied := constraint.Check(version)
	if !isSatisfied {
		return fmt.Errorf("installed dependency %s with version %s does not satisfy version requirement %s", deployedChart.Metadata.Name, deployedChart.Metadata.Version, dependency.Version)
	}

	return nil
}
