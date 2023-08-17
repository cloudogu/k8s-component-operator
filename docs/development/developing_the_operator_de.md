# Komponenten mit dem Komponenten-Operator verwalten

Hier wird beschrieben, wie man mit dem Komponenten-Operator k8s-CES-Komponenten im Cluster installiert und löscht.

## Vorbereitungen

### Helm-Repository konfigurieren
- Die Datei `.env` aus dem Template `.env.template` erstellen
    - Wichtig sind die Variablen HELM_REPO_ENDPOINT (bspw. https://registry.domain.test), HELM_REPO_USERNAME und HELM_REPO_PASSWORD
    - Außerdem sollte NAMESPACE korrekt gesetzt sein
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

## Komponenten verwalten

Siehe hierzu die Anmerkungen im [Operations-Dokument](../operations/managing_components_de.md)

## Abhängigkeiten in Komponenten darstellen

```yaml
apiVersion: v2
name: k8s-dogu-operator
...
dependencies:
- name: k8s/k8s-etcd
  version: 3.*.*
  condition: false
```

Versionsmöglichkeiten und evtl. best practices oder Empfehlungen hier beschreiben

## Den Komponenten-Operator mit anderen Komponenten lokal testen

irgendwelche Magie mit der cluster-lokalen Registry...