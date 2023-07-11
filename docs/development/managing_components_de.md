# Komponenten mit dem Komponenten-Operator verwalten

Hier wird beschrieben, wie man mit dem Komponenten-Operator k8s-CES-Komponenten im Cluster installiert und löscht.

## Vorbereitungen

### Helm-Repository konfigurieren
- Die Datei `.env` aus dem Template `.env.template` erstellen
    - Wichtig sind die Variablen HELM_REPO_ENDPOINT (bspw. https://registry.domain.test), HELM_REPO_USERNAME und HELM_REPO_PASSWORD
- Credentials im Cluster ablegen: `make helm-repo-config`

### Komponenten-Operator installieren
- Operator bauen und im Cluster installieren: `make k8s-helm-apply`

### Komponente für Test vorbereiten
- Repository der Komponente öffnen, bspw. k8s-etcd
- Helm-Chart erstellen: `make k8s-helm-package-release`
    - Generiert ein Paket nach dem Schema KOMPONENTENNAME-VERSION.tgz
- An der Helm-Registry anmelden: bspw. `helm registry login registry.domain.test`
- Helm-Chart in Registry pushen: bspw. `helm push target/make/k8s/helm/k8s-etcd-3.5.9-1.tgz oci://registry.domain.test/testing/`
    - `testing` ist hier der Namespace der Komponente in der Helm-Registry und kann angepasst werden, falls nötig

## Komponente installieren
- Custom Ressource (CR) für Komponente schreiben. Beispiel:

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

- `namespace` ist hier der Namespace der Komponente in der Helm-Registry (s.o.)
- Weitere Beispiele finden sich unter config/samples


- Anwenden der CR auf den Cluster: bspw. `kubectl apply -f etcd.yaml --namespace ecosystem`
- Der Komponenten-Operator beginnt nun mit der Installation der Komponente

## Komponente deinstallieren

- Löschen der Komponenten-CR aus dem Cluster: bspw. `kubectl delete -f etcd.yaml --namespace ecosystem`
- Der Komponenten-Operator beginnt nun mit der Deinstallation der Komponente