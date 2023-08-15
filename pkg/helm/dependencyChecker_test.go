package helm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"testing"
)

func Test_installedDependencyChecker_CheckSatisfied(t *testing.T) {
	type args struct {
		dependencies     []*chart.Dependency
		deployedReleases []*release.Release
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should succeed if dependencies and releases is nil",
			args: args{
				dependencies:     nil,
				deployedReleases: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should succeed if dependencies is nil",
			args: args{
				dependencies:     nil,
				deployedReleases: []*release.Release{createRelease("k8s-etcd", "3.0.0")},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should fail if no dependency is installed",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("not_installed", ">1.2.3")},
				deployedReleases: []*release.Release{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "dependency k8s-etcd with version ~3.0.0 is not installed\ndependency not_installed with version >1.2.3 is not installed", i)
			},
		},
		{
			name: "should fail if one dependency is not installed",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("not_installed", ">1.2.3")},
				deployedReleases: []*release.Release{createRelease("k8s-etcd", "3.0.0")},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "dependency not_installed with version >1.2.3 is not installed", i)
			},
		},
		{
			name: "should succeed if all dependencies are installed",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("k8s-dogu-operator", ">1.2.3")},
				deployedReleases: []*release.Release{createRelease("k8s-dogu-operator", "2.1.0"), createRelease("k8s-etcd", "3.0.0")},
			},
			wantErr: assert.NoError,
		},
		{
			name: "should fail to parse version requirement",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("k8s-dogu-operator", "invalid")},
				deployedReleases: []*release.Release{createRelease("k8s-dogu-operator", "2.1.0"), createRelease("k8s-etcd", "3.0.0")},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse constraint for dependency k8s-dogu-operator with version requirement invalid", i)
			},
		},
		{
			name: "should fail to parse version",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("k8s-dogu-operator", ">1.2.3")},
				deployedReleases: []*release.Release{createRelease("k8s-dogu-operator", "2.1.0"), createRelease("k8s-etcd", "invalid")},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse version of installed component k8s-etcd with version invalid", i)
			},
		},
		{
			name: "should fail if one version requirement is not satisfied",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("k8s-dogu-operator", ">1.2.3")},
				deployedReleases: []*release.Release{createRelease("k8s-dogu-operator", "2.1.0"), createRelease("k8s-etcd", "2.0.0")},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "installed dependency k8s-etcd with version 2.0.0 does not satisfy version requirement ~3.0.0", i)
			},
		},
		{
			name: "should fail if two version requirements are not satisfied",
			args: args{
				dependencies:     []*chart.Dependency{createDependency("k8s-etcd", "~3.0.0"), createDependency("k8s-dogu-operator", ">1.2.3")},
				deployedReleases: []*release.Release{createRelease("k8s-dogu-operator", "1.2.2"), createRelease("k8s-etcd", "2.0.0")},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "installed dependency k8s-etcd with version 2.0.0 does not satisfy version requirement ~3.0.0\ninstalled dependency k8s-dogu-operator with version 1.2.2 does not satisfy version requirement >1.2.3", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &installedDependencyChecker{}
			tt.wantErr(t, d.CheckSatisfied(tt.args.dependencies, tt.args.deployedReleases), fmt.Sprintf("CheckSatisfied(%v, %v)", tt.args.dependencies, tt.args.deployedReleases))
		})
	}
}

func createDependency(name, version string) *chart.Dependency {
	return &chart.Dependency{
		Name:    name,
		Version: version,
	}
}

func createRelease(name, version string) *release.Release {
	return &release.Release{Chart: &chart.Chart{Metadata: &chart.Metadata{
		Name:    name,
		Version: version,
	}}}
}