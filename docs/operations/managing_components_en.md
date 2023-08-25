# Using `k8s-component-operator`

The component operator `k8s-component-operator` is a component for the Kubernetes version of the Cloudogu EcoSystem (K8s-CES). This operator allows to easily install, upgrade or delete components. These components in turn provide required services to the EcoSystem.

## Installing the component operator

### Configure helm repository

To initially install the component operator, a log-in to the Cloudogu Helm registry is required.

```bash
$ helm registry login -u myuser registry.cloudogu.com
Password: ************************
Login succeeded
```

For later K8s CES components, this helm repository log-in is unnecessary, since the component operator has its own configuration. See the next section [Configure credentials](#configure-credentials).

### Configure credentials

The component operator has its own configuration regarding endpoint and credentials. When the K8s-CES instance is able to access the internet, the endpoint and credentials are identical to those of the Dogu registry:
- Endpoint: `oci://registry.cloudogu.com`
- Credentials: The same user/password as those from the secret `k8s-dogu-operator-dogu-registry`

This configuration can be manually created for the cluster namespace `ecosystem` as follows:

```bash
$ kubectl -n ecosystem create configmap component-operator-helm-repository --from-literal=endpoint="${HELM_REPO_ENDPOINT}" --from-literal=schema=oci
$ kubectl -n ecosystem create secret generic component-operator-helm-registry \
  --from-literal=config.json='{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "$(shell printf "%s:%s" "${HELM_REPO_USERNAME}" "${HELM_REPO_PASSWORD}" | base64 -w0)"}}}'
```

### Install component operator

Normally the component operator is installed by `k8s-ces-setup`. This can be achieved in a manual way for the cluster namespace `ecosystem` and the helm registry namespace `k8s` as follows:

```bash
$ helm install -n ecosystem k8s-component-operator oci://${HELM_REPO_ENDPOINT}/k8s/k8s-component-operator --version ${DESIRED_VERSION}
```

### Uninstall component operator

```bash
$ helm uninstall -n ecosystem k8s-component-operator
```

## Install or upgrade components

To install or upgrade components, a _Custom Resource_ (CR) for each desired component must be applied to the cluster in the correct cluster namespace.

Example of a component resource (e.g. as `k8s-dogu-operator.yaml` and from the Helm registry namespace `k8s`):

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-dogu-operator
spec:
  name: k8s-dogu-operator
  namespace: k8s
  version: 0.35.0
```

CRs like this can then be applied to the cluster:

```bash
kubectl -n ecosystem apply -f k8s-dogu-operator.yaml
```

The component operator now starts installing the component. Dependencies to other k8s-CES components and their versions must be fulfilled (this is checked by the component operator). For more information on this topic can be found in the section [Dependencies to other components](#Dependencies-to-other-components).

Examples of component resources are located in the [config/samples directory](../../config/samples)

### Fields and their meaning:

A component CR consists of various fields. This section describes these:

- `.metadata.name`: The component name of the Kubernetes resource. This must be identical to `.spec.name`.
- `.spec.name`: The component name as it appears in the Helm registry. This must be identical to `.metadata.name`.
- `.spec.namespace`: The component namespace in the helm registry.
  - Using different component namespaces, different versions could be deployed (e.g. for debugging purposes).
  - This is _not_ the cluster namespace.
- `.spec.version`: The version of the component in the helm registry.

## Uninstall components

> [!WARNING]
> Deleting components that maintain a state may jeopardize the stability of the K8s-CES installation.
> This is especially (but not exclusively) true for the component `k8s-etcd`.

- Deleting a component CR from the cluster can be done in two ways:
  1. by deleting a component from an existing component CR file, e.g. `kubectl -n ecosystem delete -f k8s-dogu-operator.yaml`.
  2. by specifying `.metadata.name` of the components, e.g. `kubectl -n ecosystem delete component k8s-dogu-operator`.
- The component operator will now start uninstalling the component

## Dependencies to other components

K8s-CES components may depend on other k8s-CES components. To ensure that a component is fully functional, the component operator checks any dependency requirements during the installation/upgrade process to see if such component dependencies are present and that they have the correct version.

If one or more components are missing or do not have the correct version, a corresponding event will be written to the component resource. Such errors can be discovered by `kubectl describe`ing the component resource, like so:

```bash
$ kubectl -n ecosystem describe component k8s-dogu-operator
```

In that case, the components in question must be manually [installed or upgraded](#Install-or-upgrade-components).

The versions to dependencies are declared in the helm chart during the component development. These can usually not be changed at the time of installation.
