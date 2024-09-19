The files in this directory were copied from https://github.com/mittwald/go-helm-client/tree/15ee7e014f3c79d7b48b24fcd29a34a2990d4450.

Their original contents are [licensed MIT and copyrighted by Mittwald CM Service](./LICENSE) except stated otherwise.

Modifications are licensed AGPL-3.0-only under the [root license of this repository](../../../LICENSE).
Modifications include:
- support for plain http registries
- usage of the client's `action.Config` when getting charts
- simplification of `ChartSpec` by stripping unnecessary fields
- simplification of the installation and upgrade procedures by removing (for our use-case) unnecessary code
- removal of unnecessary functions from the client
- embedding the registry client into the helm client to act as a tag resolver
- Moving example code to `example_test.go`
- refactoring to improve testability
- adding tests
