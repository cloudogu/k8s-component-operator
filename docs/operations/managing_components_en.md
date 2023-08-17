# Using `k8s-component-operator`

The component operator `k8s-component-operator` is a component for the Kubernetes version of the Cloudogu EcoSystem. This operator enables you to install, update, or delete components in an easy fashion. These components in turn provide necessary services specific to the EcoSystem.

## Installation of the component operator

### Configure a Helm repository

```bash
$ helm repo add oci://registry.cloudogu.com/k8s ????
```

### Configure credentials

Gegebenenfalls für das Helm-Repository Zugangsdaten anlegen.

```bash
$ kubectl add secret ????
```

### Install the component operator

Entweder mittels Helm-Client installieren oder aktualisieren

```bash
$ helm install chart.tgz???
```

oder mittels `helm template` und einem darauf folgenden `kubectl`-Aufruf


```bash
$ helm template chart.tgz???
```

## Komponente installieren oder aktualisieren

Um Komponenten zu installieren oder zu aktualisieren, muss jeweils eine _Custom Resource_ (CR) für die gewünschte Komponente existieren. 

Beispiel einer Komponenten-Ressource (z. B. als `k8s-etcd.yaml`):

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-etcd
spec:
  name: k8s-etcd
  namespace: k8s
  version: 3.5.9-1
```

Diese CR kann dann auf den Cluster angewendet werden: `kubectl apply -f k8s-etcd.yaml`
- Der Komponenten-Operator beginnt nun mit der Installation der Komponente
- Abhängigkeiten zu anderen k8s-CES-Komponenten und deren Versionen müssen erfüllt sein (siehe hierzu [Abhängigkeiten zu anderen Komponenten](#Abhängigkeiten-zu-anderen-Komponenten))

Weitere Beispiele von Komponenten-Ressourcen befinden sich im [config/samples-Verzeichnis](../../config/samples)

### Felder und deren Bedeutung:

- `.metadata.name`: Der Komponentenname der Kubernetes-Resource. Dieser muss identisch mit `.spec.name` sein.
- `.spec.name`: Der Komponentenname wie er in der Komponenten-Registry lautet. Dieser muss identisch mit `.spec.name` sein.
- `.spec.namespace`: Der Namespace der Komponente in der Komponenten-Registry. 
  - Mittels unterschiedlicher Komponenten-Namespaces können unterschiedliche Versionen ausgebracht werden (z. B. zu Debugging-Zwecken). 
  - Es handelt sich hierbei _nicht_ um den Cluster-Namespace.
- `.spec.version`: Die Version der Komponente in der Komponenten-Registry.
 

## Komponente deinstallieren

- Löschen der Komponenten-CR aus dem Cluster: bspw. `kubectl delete -f etcd.yaml`
- Der Komponenten-Operator beginnt nun mit der Deinstallation der Komponente

```bash
$ kubectl api-resources --verbs=list -o name \
  | sort | xargs -t -n 1 \
  kubectl delete --ignore-not-found -l app.kubernetes.io/name: k8s-component-operator -n ecosystem
```

## Abhängigkeiten zu anderen Komponenten

K8s-CES-Komponenten können von anderen k8s-CES-Komponenten abhängen. Um sicherzustellen, dass eine Komponente voll funktionsfähig ist, wird während der Installation bzw. Aktualisierung geprüft, ob Komponenten vorhanden sind und eine korrekte Version aufweisen.

Sollte eine oder mehrere Komponenten fehlen oder nicht in der richtigen Version vorhanden sein, so müssen diese manuell [nachinstalliert](#Komponente-installieren-oder-aktualisieren) bzw. aktualisiert werden.

Die Versionen zu Abhängigkeiten werden während der Komponentenentwicklung im Helm-Chart hinterlegt. Abhängige Versionen können so gestaltet werden, dass sie nicht auf eine einzige Version fixiert werden, sondern unterschiedliche Versionsbereiche abdecken. Dies ermöglicht den Betrieb von Komponenten, selbst wenn Kompoentenversionen mit kleineren Änderungen oder Fehlerbehebungen ausgebracht wurden.    