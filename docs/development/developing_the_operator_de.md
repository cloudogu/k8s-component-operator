# Komponenten-Operator und Komponenten entwickeln

Dieses Dokument beschreibt sowohl, wie man den Komponenten-Operator entwickelt als auch Komponenten-spezifische Eigenheiten.

## Vorbereitungen

### Helm-Repository konfigurieren

- Die Datei `.env` aus dem Template `.env.template` erstellen
   - Wichtig sind die Variablen 
      - `HELM_REPO_ENDPOINT` (bspw. https://registry.cloudogu.com)
      - `HELM_REPO_USERNAME`
      - `HELM_REPO_PASSWORD`
      - `NAMESPACE`
- Credentials im Cluster ablegen: `make helm-repo-config`

### Den Komponenten-Operator lokal debuggen

1. Befolgen Sie die Installationsanweisungen von k8s-ecosystem
2. Öffnen Sie die Datei `.env.template` und folgen Sie den Anweisungen um eine
   Umgebungsvariablendatei mit persönlichen Informationen anzulegen
3. Löschen Sie eventuelle Komponenten-Operator-Deployments im Cluster, um Parallelisierungsfehler auszuschließen
   - `kubectl -n ecosystem delete deployment k8s-component-operator`

Nun haben Sie zwei Möglichkeiten den Operator lokal zu starten:
1. mit `make run`
2. mit IntelliJ (um den Code zu Debuggen)
   - Generieren Sie die Configmap und die lokale config.json mit `make helm-repo-config-local`
   - Legen Sie eine neue Debug-Konfiguration (z. B. in IntelliJ) an
      - diese Umgebungsvariablen werden benötigt (kann mit `make print-debug-info` generiert werden):
      - STAGE=production;NAMESPACE=ecosystem;KUBECONFIG=/path/to/kubeconfig/.kube/k3ces.local
   - Breakpoints setzen und ggf. ein Komponenten-CR auf den Cluster anwenden

### Komponenten-Operator installieren

- Operator bauen und im Cluster installieren: `make k8s-helm-apply`

### Komponente für Test vorbereiten

Am Beispiel von `k8s-dogu-operator`

1. Repository der Komponente öffnen
2. Ggf. im Verzeichnis `k8s/helm` ein `Chart.yaml` mit `make k8s-helm-init-chart` anlegen
3. Helm-Package erstellen: `make k8s-helm-package-release`
   - generiert ein Paket nach dem Schema KOMPONENTENNAME-VERSION.tgz
4. ggf. alle nötigen Nicht-Test-Komponenten [installieren](../operations/managing_components_de.md#komponenten-installieren-oder-aktualisieren)
   `kubectl -n ecosystem apply -f yourComponentCR.yaml`
5. Test-Komponente pushen:
   - `make chart-import`
6. die ConfigMap `component-operator-helm-repository` auf die cluster-lokale Registry richten
   - `kubectl -n ecosystem patch configmap component-operator-helm-repository -p '{"data": {"endpoint": "oci://k3ces.local:30099","plainHttp": "true"}}'`
7. YAML der Test-Komponente überprüfen und [installieren](../operations/managing_components_de.md#komponenten-installieren-oder-aktualisieren)
   `kubectl -n ecosystem apply -f k8s-dogu-operator.yaml`

## Abhängigkeiten in Komponenten darstellen

Komponenten müssen nicht unbedingt für sich alleine stehen, sondern können auch andere Komponenten erfordern. Dies wird als Abhängigkeit im Helm-Chart definiert:

```yaml
apiVersion: v2
name: k8s-dogu-operator
...
dependencies:
  - name: k8s/k8s-dogu-operator
    version: 3.*.*
    condition: false
```

Abhängigkeitsversionen sollten so gestaltet werden, dass sie nicht auf eine einzige Version fixiert werden, sondern unterschiedliche Versionsbereiche abdecken. Dies ermöglicht den Betrieb von Komponenten, selbst wenn Komponentenversionen mit kleineren Änderungen oder Fehlerbehebungen ausgebracht wurden.

Die Bibliothek [Masterminds/semver](https://github.com/Masterminds/semver#checking-version-constraints) beschreibt genauer, welche Versionseinschränkungen möglich sind.

Da wir die Abhängigkeitsdeklaration im Helm-Chart nur nutzen, um Abhängigkeiten für den Komponenten-Operator darzustellen, muss das Feld `.dependencies.[].condition` zwingend auf `false` gesetzt werden. Würde dieses Feld `true` sein, würde Helm die Abhängigkeit automatisch installieren und der Komponenten-Operator würde in seiner eigenen Tätigkeit gestört werden.
