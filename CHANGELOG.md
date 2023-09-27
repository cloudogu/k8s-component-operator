# k8s-component-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- [#29] Add permissions to the `longhorn-system` namespace and for `VolumeSnapshotClasses` resources.
  - With these permissions the component operator can install backup provider components required by the backup-operator.

## [v0.3.0] - 2023-09-15
### Changed
- [#25] Use component-dependencies from the annotations of a HelmChart instead of the Helm-dependencies
- [#27] updated go dependencies
- [#27] updated kube-rbac-proxy

### Fixed
- [#27] deprecation warning for argument `logtostderr` in kube-rbac-proxy

### Removed
- [#27] deprecated argument `logtostderr` from kube-rbac-proxy

## [v0.2.0] - 2023-09-07
### Added
- [#23] Validate that `metadata.Name` equals `spec.Name`.

## [v0.1.2] - 2023-09-01
### Fixed
- [#21] Fixes dependency-check for components with the version-format "x.x.x-x"
  - "x.x.x-x"-versions are not treated as "pre-release"-versions and are ordered accordingly

## [v0.1.1] - 2023-08-25
### Fixed
- [#19] Fixes operator configuration by splitting helm registry endpoint and schema into separate configmap fields
- Fixes K8s resource conflicts while updating components that may have changed concurrently

## [v0.1.0] - 2023-08-24
### Added
- [#15] Check if component dependencies are installed and if their version is appropriate
  - you can find more information about components in the [operations docs](docs/operations/managing_components_en.md)
- this release adds the ability to requeue CR requests

## [v0.0.3] - 2023-08-21
### Changed
- [#17] Make helmClient more generic to be usable by other components (e.g. "k8s-ces-setup")

## [v0.0.2] - 2023-07-14
### Added
- [#12] Add upgrade of components and self-upgrade of component-operator
- Add documentation for component operator usage in a development environment

### Fixed
- Operator finishes uninstallation steps even if component has been uninstalled already
- [#12] Fix the log-format for the logger used in the helm-client

## [v0.0.1] - 2023-07-07
### Changed
- [#8] Stabilise the installation process with atomic helm operations and a timeout for the underlying k8s client.

### Added
- [#4] Add Helm chart release process to project
- [#3] Initialize a first version for the `k8s-component-operator`. In contrast to the prior PoC status the operator pulls charts from an OCI registry.
