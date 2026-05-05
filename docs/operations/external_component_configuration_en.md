# External Component Configuration

In some cases, it is necessary to configure components within the cluster that are dependent on other components.
ConfigMaps are an ideal choice for managing external configuration in such scenarios.
Configuration via ConfigMaps is currently possible using `.spec.valuesConfigRef`.
However, this requires the ConfigMap to be specified in the spec and would need to be configured with every deployment of the component.
This involves manual effort. Furthermore, this configuration option is intended more for individual customer configuration.

For example, ports need to be dynamically opened in `k8s-ces-gateway`. This cannot be done dynamically via Traefik itself.
The ports are statically defined in `values.yaml`. The configuration can only be updated by restarting the service.

The `k8s-component-operator` reconciles ConfigMaps with the label `k8s.cloudogu.com/component.config` (e.g. `k8s.cloudogu.com/component.config: k8s-ces-gateway`).
These ConfigMaps contain a `values.yaml` file for the component under the `values` key.
The component with the name matching the value of the label is reconciled.
In this process, the values from the ConfigMap are merged with the component’s other configuration values.
The values from the ConfigMap have a low priority, meaning they can, for example, be overwritten
by the `MappedValues`.