# k8s-component-operator Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added 
- [#71] Add comments for RBAC permissions

### Removed
- [#71] Remove unused ClusterRole for reading metrics
- [#71] Remove unused leader-election along with its RBAC permissions

## [v1.2.1] - 2024-11-04
### Fixed
- [#69] Image path for rbac proxy registry in patch templates.

### Changed
- Update samples

## [v1.2.0] - 2024-10-28
### Changed
- [#66] Make imagePullSecrets configurable via helm values and use `ces-container-registries` as default.

## [v1.1.1] - 2024-10-07
### Changed
- Upgrade go to v1.23
- Upgrade golang-ci to v1.61.0
- Update Dependencies

## [v1.1.0] - 2024-09-19
### Changed
- Relicense to AGPL-3.0-only

## [v1.0.1] - 2024-06-10
### Changed
- Throw error if the current deployment and the component target namespaces do not match
- Upgrade to go 1.22
- Upgrade to latest Makefiles, ces-build-lib and linter version

## [v1.0.0] - 2024-03-21
### Added
- [#54] Helm-Credentials can be stored in base64 encoding and clarified escaping rules if not
  (see [here](docs/development/developing_the_operator_en.md) or [here](.env.template))
- [#56] installed version in component status, so that health checks can check for the actual installed version
- [#58] Set health status of components on startup and shutdown.
  ([see here](docs/operations/component_health_en.md))

## [v0.8.0] - 2024-01-30
### Changed
- [#51] Requeue components on errors during install, update and delete.
### Added
- [#48] Add labels with component name and version to all resources of a component
- [#48] Track health on component CR

## [v0.7.0] - 2023-12-08
### Added
- [#44] Patch-templates for mirroring this operator and its images
- [#46] Added documentation for upgrading kubebuilder components.
### Changed
- [#42] Replace monolithic K8s resource YAML into Helm templates
- Update Makefiles to 9.0.1

## [v0.6.0] - 2023-11-16
### Added
- [#40] Components default values.yaml can be overwritten (with the field valuesYamlOverwrite)
- [#38] Add [documentation](docs/operations/creating_components_en.md) for creating components and component-patch-templates

### Changed
- [#36] Allow insecure TLS certificates with configuration options

## [v0.5.1] - 2023-10-11
### Changed
- [#33] Replace mittwald-go-helm-client with a reduced implementation fitted to our needs

## [v0.5.0] - 2023-10-06
### Added
- [#29] Add permissions on all namespaces to install components in namespaces like monitoring or longhorn-system. Add a new property deployNamespace. It is used to determine where the helm deployment should go. If it is empty the namespace from the component-operator is used.

## [v0.4.0] - 2023-10-05
### Added
- [#30] Add CRD-Release to Jenkinsfile

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
