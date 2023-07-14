# k8s-component-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.0.2] - 2023-07-14
### Added
- Add documentation for component operator usage in dev environment
- [#12] Add upgrade of components and self-upgrade of component-operator

### Fixed
- Operator finishes uninstallation steps even if component has been uninstalled already
- [#12] Fix the log-format for the logger used in the helm-client

## [v0.0.1] - 2023-07-07
### Changed
- [#8] Stabilise the installation process with atomic helm operations and a timeout for the underlying k8s client.
### Added
- [#4] Add Helm chart release process to project
- [#3] Initialize a first version for the `k8s-component-operator`. In contrast to the prior poc status the operator pulls charts from an oci registry.
