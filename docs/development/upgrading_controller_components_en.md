# Upgrading Kubebuilder

In most cases, it is sufficient to update the Kubernetes dependencies in `go.mod`.
However, if you want to configure webhooks or, for example, a Prometheus exporter for the controller,
it makes sense to create an initial Kubebuilder project and then transfer the files from the Kustomize structure in the `config` folder to the Helm chart.


## Project creation

### Structure

`kubebuilder init --domain k8s.cloudogu.com --repo github.com/cloudogu/k8s-component-operator`

> [Kubebuilder binaries](https://github.com/kubernetes-sigs/kubebuilder/releases)

### Api

`kubebuilder create api --group k8s.cloudogu.com --version v1beta1 --kind Component`

## Create YAML resources

`kustomize build config/default`

> [Kustomize binaries](https://github.com/kubernetes-sigs/kustomize/releases)

The resulting artifact can then be transferred to the Helm chart.
When initially creating an operator, [`helmify`](https://github.com/arttor/helmify/releases) can be used.

`kustomize build config/default | helmify helm/component-operator`

## Please note

1. Attributes that should be configurable via the `values.yaml`:
   - Images
   - Tags
   - PullPolicy
   - Resource requirements
   - Environment variables (LogLevel, Stage)

2. The templates should have default values for templating. Example:

```
          env:
            - name: STAGE
              value: {{ quote .Values.manager.env.stage | default "production" }}
            - name: LOG_LEVEL
              value: {{ quote .Values.manager.env.logLevel | default "info"}}
```

3. Auxiliary functions should exist at least for the name, labels and selector labels. They should be used for the corresponding attributes. Example:

```
{{/* Chart basics
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec) starting from
Kubernetes 1.4+.
*/}}
{{- define "k8s-component-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/* All-in-one labels */}}
{{- define "k8s-component-operator.labels" -}}
app: ces
{{ include "k8s-component-operator.selectorLabels" . }}
helm.sh/chart: {{- printf " %s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Selector labels */}}
{{- define "k8s-component-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "k8s-component-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
```