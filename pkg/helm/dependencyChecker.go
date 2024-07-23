package helm

import (
	"errors"
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"

	"github.com/Masterminds/semver/v3"
)

const cesDependencyAnnotationIdentifier = "k8s.cloudogu.com/ces-dependency/"

type Dependency struct {
	Name    string
	Version string
}

func getComponentDependencies(chart *chart.Chart) []Dependency {
	var dependencies []Dependency
	for key, value := range chart.Metadata.Annotations {
		componentName, found := strings.CutPrefix(key, cesDependencyAnnotationIdentifier)
		if found && componentName != "" {
			dependencies = append(dependencies, Dependency{Name: componentName, Version: value})
		}
	}

	return dependencies
}

type dependencyChecker interface {
	// CheckSatisfied validates that all dependencies are installed in the required version.
	CheckSatisfied(dependencies []Dependency, deployedReleases []*release.Release) error
}

type installedDependencyChecker struct{}

// CheckSatisfied validates that all dependencies are installed in the required version.
func (d *installedDependencyChecker) CheckSatisfied(dependencies []Dependency, deployedReleases []*release.Release) error {
	var errs []error
	for _, dependency := range dependencies {
		isInstalled := false
		for _, deployedRelease := range deployedReleases {
			if dependency.Name == deployedRelease.Chart.Name() {
				isInstalled = true
				err := checkVersion(dependency, deployedRelease.Chart)
				errs = append(errs, err)

				break
			} else if strings.Contains(dependency.Name, "crd") {
				isInstalled = true
				break
			}
		}

		if !isInstalled {
			errs = append(errs, fmt.Errorf("dependency %s with version %s is not installed", dependency.Name, dependency.Version))
		}
	}

	return errors.Join(errs...)
}

func checkVersion(dependency Dependency, deployedChart *chart.Chart) error {
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
