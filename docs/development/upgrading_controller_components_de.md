# Upgrading Kubebuilder

In den meisten Fällen reicht es in der `go.mod` die Kubernetes-Abhängigkeiten zu aktualisieren.
Möchte man allerdings für den Controller Webhooks oder zum Beispiel einen Prometheus-Exporter konfigurieren,
ist es sinnvoll ein initiales Kubebuilder-Projekt zu erstellen und die Files anschließend von der Kustomize-Struktur im Ordner `config` in das Helm-Chart zu überführen.


## Projekterstellung

### Struktur

`kubebuilder init --domain k8s.cloudogu.com --repo github.com/cloudogu/k8s-component-operator`

> [Kubebuilder binaries](https://github.com/kubernetes-sigs/kubebuilder/releases)

### Api

`kubebuilder create api --group k8s.cloudogu.com --version v1beta1 --kind Component`

## YAML-Ressourcen erzeugen

`kustomize build config/default`

> [Kustomize binaries](https://github.com/kubernetes-sigs/kustomize/releases)

Das resultierende Artefakt kann dann in das Helm-Chart überführt werden.
Bei einer initialen Erstellung eines Operators kann [`helmify`](https://github.com/arttor/helmify/releases) verwendet werden.

`kustomize build config/default | helmify helm/component-operator`

## Zu beachten

1. Attribute, die über die `values.yaml` konfigurierbar sein sollten:
    - Images
    - Tags
    - PullPolicy
    - Ressourcenanforderungen
    - Umgebungsvariablen (LogLevel, Stage)
2. Die Templates sollten Default-Werte für Templating besitzen. Beispiel:

```
          env:
            - name: STAGE
              value: {{ quote .Values.manager.env.stage | default "production" }}
            - name: LOG_LEVEL
              value: {{ quote .Values.manager.env.logLevel | default "info"}}
```

3. Hilfsfunktionen sollten mindestens für den Namen, Labels und Selektorlabels existieren. Sie sollten für die entsprechenden Attribute verwendet werden Beispiel:

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