# Verwendung von `k8s-component-operator`

Der Komponenten-Operator `k8s-component-operator` ist eine Komponente für die Kubernetes-Version des Cloudogu EcoSystems. Dieser Operator ermöglicht es, Komponenten auf einfache Weise zu installieren, zu aktualisieren oder zu löschen. Diese Komponenten stellen ihrerseits erforderliche Dienste für das EcoSystem bereit.

## Installation des Komponenten-Operators

### Helm-Repository konfigurieren

```bash
$ helm repo add oci://registry.cloudogu.com/k8s ????
```

### Zugangsdaten konfigurieren

Gegebenenfalls für das Helm-Repository Zugangsdaten anlegen.

siehe make helm-repo-config

```bash
$ kubectl add secret ????
```

### Komponenten-Operator installieren

setup hinweis

bei fehlern

Entweder mittels Helm-Client installieren oder aktualisieren

- An der Helm-Registry anmelden: bspw. `helm registry login registry.domain.test`
    - oder helm registry login -u myuser localhost:5000****
- Helm-Chart in Registry pushen: bspw. `helm push target/make/k8s/helm/k8s-etcd-3.5.9-1.tgz oci://registry.domain.test/testing/`

```bash
$ helm install k8s-component-operator oci://registry.domain.test/testing/k8s-component-operator --version 0.2.0
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

Die Versionen zu Abhängigkeiten werden während der Komponentenentwicklung im Helm-Chart hinterlegt. Abhängige Versionen können so gestaltet werden, dass sie nicht auf eine einzige Version fixiert werden, sondern unterschiedliche Versionsbereiche abdecken. Dies ermöglicht den Betrieb von Komponenten, selbst wenn Komponentenversionen mit kleineren Änderungen oder Fehlerbehebungen ausgebracht wurden.    