The files in this directory were copied from https://github.com/mittwald/go-helm-client/tree/15ee7e014f3c79d7b48b24fcd29a34a2990d4450.

They are [licensed MIT](./LICENSE) except stated otherwise.

Modifications include:
- support for plain http registries
- usage of the client's `action.Config` when getting charts
- simplification of `ChartSpec` by stripping unnecessary fields
- simplification of the installation and upgrade procedures by removing (for our use-case) unnecessary code
- removal of unnecessary functions from the client
- embedding the registry client into the helm client to act as a tag resolver