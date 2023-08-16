# Komponenten mit dem Komponenten-Operator verwalten

Hier wird beschrieben, wie man mit dem Komponenten-Operator k8s-CES-Komponenten im Cluster installiert, aktualisiert und löscht.

## Installation des Komponenten-Operators

### Helm-Repository konfigurieren

```bash
$ helm repo add oci://registry.cloudogu.com/k8s ????
```

### Zugangsdaten konfigurieren

Gegebenenfalls für das Helm-Repository Zugangsdaten anlegen.

```bash
$ kubectl add secret ????
```

### Komponenten-Operator installieren

Entweder mittels Helm-Client installieren oder aktualisieren

```bash
$ helm install chart.tgz???
```

oder mittels `helm template` und einem darauf folgenden `kubectl`-Aufruf


```bash
$ helm template chart.tgz???
```

## Komponente installieren oder aktualisieren
- Custom Ressource (CR) für Komponente schreiben. Beispiel:

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

- `namespace` ist hier der Namespace der Komponente in der Helm-Registry (s. o.)
- Weitere Beispiele finden sich unter config/samples

- Anwenden der CR auf den Cluster: bspw. `kubectl apply -f etcd.yaml`
- Der Komponenten-Operator beginnt nun mit der Installation der Komponente
- Abhängigkeiten zu anderen k8s-CES-Komponenten und deren Versionen müssen erfüllt sein (siehe hierzu [Abhängigkeiten zu anderen Komponenten](#Abhängigkeiten-zu-anderen-Komponenten))

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
