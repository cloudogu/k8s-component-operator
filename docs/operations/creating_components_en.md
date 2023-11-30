# Creation of K8s components

K8s CES components provide required services for the Cloudogu EcoSystem (CES).

## Creating a new component
The following steps describe the creation of a general K8s CES component that can be operated in the Multinode CES:

- Create a new repo
  - Schema `k8s-NAME`
- Import or create basic files
  - README.md
  - Jenkinsfile
  - LICENSE
  - CHANGELOG.md
  - [Makefiles](https://github.com/cloudogu/makefiles)
  - .gitignore
- Determine the K8s resources of the component:
  - As K8s controller: Include the `k8s-controller.mk` Makefile to generate the K8s resources
- Create Helm chart `Chart.yaml` in `k8s/helm/` with `make helm-init-chart`
- If necessary, enter [Component Dependencies](#component-dependencies) in the `Chart.yaml`.
- Create a [Component Patch Template](#component-patch-template)

The following make targets can then be used:
- `helm-generate`: Assembles the finished Helm chart in the target folder from the resources under k8s/helm and the generated K8s resources
- `helm-apply`: Applies the chart in the local DEV cluster
- `component-apply`: Applying the chart in the local DEV cluster as an installation/upgrade via the component operator
- `helm-package`: Builds and packs the Helm chart as `.tgz` to release it into a Helm repository


### Create component for third-party applications
Additional steps are required to create third-party Helm charts as K8s CES components, described here using the example of `promtail`:

- Search for the official chart of the product (e.g. promtail) and insert it into your own `Chart.yaml` as `dependency`.
  ```yaml
  name: k8s-promtail
  ...
  dependencies:
    - name: promtail
      version: 6.15.2
      repository: https://grafana.github.io/helm-charts
  ```
- Write a make target that creates the `k8s/helm/charts` folder from the `dependencies` entry
  ```makefile
  .PHONY: ${K8S_HELM_RESSOURCES}/charts
  ${K8S_HELM_RESSOURCES}/charts: ${BINARY_HELM}
  @cd ${K8S_HELM_RESSOURCES} && ${BINARY_HELM} repo add grafana https://grafana.github.io/helm-charts && ${BINARY_HELM} dependency build
  ```
- Create "manifests" folder with dummy-yaml
  - Necessary because the Makefiles currently require K8s resources (yaml) to create the Helm chart
  - e.g. "promtail.yaml"
  ```yaml
  # This is a dummy file, required for the makefile's yaml file generation process.
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: promtail-dummy
  data: {}
  ```
- Overwrite `K8S_PRE_GENERATE_TARGETS` target in the Makefile with your own target
  - E.g. `K8S_PRE_GENERATE_TARGETS=generate-release-resource`.
  - Move the dummy-yaml to `K8S_RESOURCE_TEMP_YAML` in this target
- Create your own target `helm-NAME-apply` (e.g. `helm-promtail-apply`) in the Makefile
  - This works similar to "k8s-apply" from `k8s.mk`, but without the "image-import" target
  ```makefile
  .PHONY: helm-promtail-apply
  helm-promtail-apply: ${BINARY_HELM} ${K8S_HELM_RESSOURCES}/charts helm-generate $(K8S_POST_GENERATE_TARGETS) ## Generates and installs the helm chart.
  @echo "Apply generated helm chart"
  @${BINARY_HELM} upgrade -i ${ARTIFACT_ID} ${K8S_HELM_TARGET} --namespace ${NAMESPACE} 
  ```


## Component dependencies
K8s CES components may depend on other K8s CES components.
So that the component operator can check these dependencies when installing or upgrading components, they must be specified in the Helm chart as `annotation`.
Several dependencies can be specified.

The component dependency's `annotation` key must always be specified in the form `k8s.cloudogu.com/ces-dependency/<dependecy-name>`.
The `<dependency-name>` is the name of the K8s CES component on which the dependency exists.

The component dependency's `annotation` value contains the version of the dependent component.
[Semantic versioning](https://semver.org/) is used here so that version ranges can also be specified.

### Example component dependency in the `Chart.yaml`
```yaml
annotations:
  # Dependency for the Component-CRD.
  "k8s.cloudogu.com/ces-dependency/k8s-component-operator-crd": "1.x.x-0"
```

## Component patch template
In order for a K8s CES component to be mirrored into air-gapped environments with a Cloudogu application, it must contain a 'component patch template'.
This must be stored in a file with the name `component-patch-tpl.yaml` in the root directory of a Helm chart.
The `component-patch-template` contains a list of all necessary container images and template instructions to rewrite image references in Helm Values files during mirroring.

The structure is as follows:
```yaml
apiVersion: v1

values:
  images:
    <imageKey1>: "<imageRef1>"
    <imageKey2>: "<imageRef2>"

patches:
  <filename>:
    foo:
      bar:
        registry: "{{ registryFrom .images.imageKey1 }}"
        repository: "{{ repositoryFrom .images.imageKey1 }}"
        tag: "{{ tagFrom .images.imageKey1 }}"
```

### apiVersion
The `apiVersion` specifies the version of the patch API used in the template.
Version `v1` is currently supported.
The associated template functions are described under [patches](#patches).

### values
`values` contains a map of arbitrary values that can be used for templating the files specified in the [patches](#patches).
The `values` must contain at least one map `images`, which contains all container images to be mirrored.
The key of an entry in the `images` map can be chosen arbitrarily.
The value of an entry in the `images` map corresponds to a container image reference (e.g. `registry.cloudogu.com/k8s/k8s-dogu-operator:0.35.1`).

> **Important:**
> - The key of an entry in the `images` map must not contain any hyphens "-" so that processing in [Go-Templates](https://pkg.go.dev/text/template) is possible.
> - The value of an entry in the `images` map should always be specified as a string in double quotes to avoid problems when parsing as YAML.

### patches
`patches` contain individual templates for any YAML files of the Helm chart (e.g. the `values.yaml`).
Each template is stored under the file name of the file to be patched.
A template can contain any YAML structure.  
The [Go template language](https://godoc.org/text/template) is used.
The [`values`-map](#values) is available as data in the templating.

> **Important:**
> Go template functions (e.g. "{{ .Foo }}") must be specified as a string in double quotes to prevent problems when parsing as YAML.

In addition, the following template functions are available for parsing container image references. The [keys](#values) for container images already listed under `.values.images` should be used, e.g. in the form `.images.yourContainerImageKey`:

- **registryFrom <string>**: returns the registry of a container image reference (e.g. `registry.cloudogu.com`)
- **repositoryFrom <string>**: returns the repository of a container image reference (e.g. `k8s/k8s-dogu-operator`)
- **tagFrom <string>**: returns the tag of a container image reference (e.g. `0.35.1`)

After a template has been rendered, it is merged into the "original" YAML file of the Helm chart.
This preserves values in the "original" YAML file that are _not_ contained in the template.
Existing values are overwritten by the rendered template.

#### Example `component-patch-tpl.yaml`

```yaml
apiVersion: v1

values:
  images:
    engine: "longhornio/longhorn-engine:v1.5.1"
    manager: "longhornio/longhorn-manager:v1.5.1"
    ui: "longhornio/longhorn-ui:v1.5.1"

patches:
  values.yaml:
    longhorn:
      image:
        longhorn:
          engine:
            repository: "{{ registryFrom .images.engine }}/{{ repositoryFrom .images.engine }}"
            tag: "{{ tagFrom .images.engine }}"
          manager:
            repository: "{{ registryFrom .images.manager }}/{{ repositoryFrom .images.manager }}"
            tag: "{{ tagFrom .images.manager }}"
          ui:
            repository: "{{ registryFrom .images.ui }}/{{ repositoryFrom .images.ui }}"
            tag: "{{ tagFrom .images.ui }}"
```
