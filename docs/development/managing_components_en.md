# Manage components with the component operator

Here we describe how to install and delete k8s-CES components in the cluster using the component operator.

## Preparations

### Configure helm repository
- Create the file `.env` from the template `.env.template`
    - The variables HELM_REPO_ENDPOINT (e.g. https://registry.domain.test), HELM_REPO_USERNAME and HELM_REPO_PASSWORD are important.
    - In addition, NAMESPACE should be set correctly.
- Store credentials in the cluster: `make helm-repo-config`.

### Install component operator
- Build operator and install in cluster: `make k8s-helm-apply`

### Prepare component for test
- Open repository of component, e.g. k8s-etcd
- Create helm chart: `make k8s-helm-package-release`.
    - Generates a package according to the scheme COMPONENTNAME-VERSION.tgz
- Log in to the helm registry: e.g. `helm registry login registry.domain.test`.
- Push helm chart to registry: e.g. `helm push target/make/k8s/helm/k8s-etcd-3.5.9-1.tgz oci://registry.domain.test/testing/`
    - `testing` here is the namespace of the component in the helm registry and can be modified if necessary

## Installing the component
- Write custom resource (CR) for component. Example:

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-etcd
spec:
  name: k8s-etcd
  namespace: testing
  version: 3.5.9-1
```

- `namespace` here is the namespace of the component in the helm registry (see above)
- More examples can be found at config/samples

- Apply the CR to the cluster: e.g. `kubectl apply -f etcd.yaml`.
- The component operator will now start installing the component

## Uninstalling the component

- Delete the component CR from the cluster: e.g. `kubectl delete -f etcd.yaml`
- The component operator will now start uninstalling the component