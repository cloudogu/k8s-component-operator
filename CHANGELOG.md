# k8s-component-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- [#4] Add Helm chart release process to project
### Changed
- [#8] Stabilise the installation process with atomic helm operations and a timeout for the underlying k8s client.
### Added
- [#3] Initialize a first version for the `k8s-component-operator`. In contrast to the prior poc status the operator pulls charts from an oci registry.
