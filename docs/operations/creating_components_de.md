# Erstellung von K8s-CES-Komponenten

K8s-CES-Komponenten stellen erforderliche Dienste für das Cloudogu EcoSystem bereit.
Dabei handelt es sich um Helm-Charts, die in einer OCI-Registry verfügbar sind.

## Eine neue Komponente erstellen
Die folgenden Schritte beschreiben die Erstellung einer allgemeinen K8s-CES-Komponente, die im Cloudogu EcoSystem betrieben werden kann:

- neues Github-Repo anlegen
   - Schema `k8s-NAME`
   - Readme kann mit angelegt werden, alles andere machen wir händisch danach
- "develop"-Branch anlegen und zum Default machen (in den Github-Einstellungen des Repos)
- Vom develop einen neuen Branch für den initialen Sourcecode abzweigen
- Grundsätzliche Dateien importieren bzw. erstellen
   - README.md
   - Jenkinsfile
   - LICENSE
   - CHANGELOG.md
   - Makefiles
   - .gitignore
- Die K8s-Ressourcen der Komponente bestimmen:
  - K8s-Controller: Einbindung des `k8s-controller.mk` Makefiles zur Generierung der K8s-Ressourcen
  - Das Make-Target `k8s-create-temporary-resource` erstellen, das für die Erstellung der K8s-Ressourcen zuständig ist
- Im `k8s`-Ordner Helm-Chart `helm/Chart.yaml` erzeugen
   - mit "make helm-init-chart" erzeugen
- Ggf. [Component Dependencies](#component-dependencies) in der `Chart.yaml` eintragen
- Ein [Component Patch Template](#component-patch-template) erstellen

Anschließend können die folgenden Make-Targets eingesetzt werden:
   - `helm-generate`: Baut im target-Ordner das fertige Helm-Chart aus den Ressourcen unter k8s/helm und den generierten K8s-Ressourcen zusammen
   - `helm-apply`: Anwenden des Charts im lokalen DEV-Cluster
   - `component-apply`: Anwenden des Charts im lokalen DEV-Cluster als Installation/Upgrade über den Komponenten-Operator
   - `helm-package-release`: Baut und packt das Helm-Chart als `.tgz` um es in ein Helm-Repository zu releasen


### Komponente für Fremd-Anwendungen erstellen
Um Fremd-Anwendungen als K8s-CES-Komponente zu erstellen sind einige extra Schritte nötig.
Hier am Beispiel für `promtail` beschrieben:

- "manifests"-Ordner mit dummy-yaml anlegen
  - bspw. "promtail.yaml"
  ```yaml
  # This is a dummy file, required for the makefile's yaml file generation process for the dogu-controller.
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: promtail-dummy
  data: {}
  ```
- Offizielles Chart des Produkts (bspw. promtail) suchen und in `Chart.yaml` einfügen
  ```yaml
  dependencies:
    - name: promtail
      version: 6.15.2
      repository: https://grafana.github.io/helm-charts
  ```
- Make-Target schreiben, was den `k8s/helm/charts`-Ordner aus dem `dependencies`-Eintrag erzeugt
  ```makefile
  .PHONY: ${K8S_HELM_RESSOURCES}/charts
  ${K8S_HELM_RESSOURCES}/charts: ${BINARY_HELM}
  @cd ${K8S_HELM_RESSOURCES} && ${BINARY_HELM} repo add grafana https://grafana.github.io/helm-charts && ${BINARY_HELM} dependency build
  ```
- `K8S_PRE_GENERATE_TARGETS`-Target in Makefile mit eigenem Target überschreiben
  - Bspw. `K8S_PRE_GENERATE_TARGETS=generate-release-resource`
  - In diesem Target die dummy-yaml nach `K8S_RESOURCE_TEMP_YAML` verschieben
- Eigenes Target `helm-NAME-apply` (bspw. `helm-promtail-apply`) im Makefile erstellen
  - Funktioniert analog zu "k8s-apply" aus k8s.mk, aber ohne das "image-import"-Target
  ```makefile
  .PHONY: helm-promtail-apply
  helm-promtail-apply: ${BINARY_HELM} ${K8S_HELM_RESSOURCES}/charts helm-generate $(K8S_POST_GENERATE_TARGETS) ## Generates and installs the helm chart.
  @echo "Apply generated helm chart"
  @${BINARY_HELM} upgrade -i ${ARTIFACT_ID} ${K8S_HELM_TARGET} --namespace ${NAMESPACE} 
  ```


### Component-Dependencies

> TODO

### Component-Patch-Template
Damit eine K8s-CES-Komponente von `ces-mirror` gespiegelt werden kann, muss sie ein `Component-Patch-Template`enthalten.
Diese muss in einer Datei mit dem Namen `component-patch-tpl.yaml` im Root-Verzeichnis eines Helm-Charts abgelegt werden.
Das `Component-Patch-Template` enthält eine Liste aller nötigen Container-Images und Template-Anweisungen, um Image-Referenzen in Helm-Values-Dateien während der Spiegelung umzuschreiben.

Der Aufbau sieht wie folgt aus:
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

#### apiVersion
Die `apiVersion` gibt die im Template verwendete Version der Patch-API an.
Derzeit wird von `ces-mirror` die Version `v1` unterstützt. 
Die zugehörigen Template-Funktionen werden unter [patches](#patches) beschrieben. 

#### values
`values` enthält eine Map von beliebigen Werten, die für das Templating der in den [patches](#patches) angegebenen Dateien verwendet werden können.
Die `values` müssen mindestens eine Map `images` enthalten, die alle zu spiegelnden Container-Images enthält.
Der Key eines Eintrags in der `images`-Map kann beliebig gewählt werden.
Der Value eines Eintrags in der `images`-Map entspricht einer Container-Image-Referenz (z.B. `registry.cloudogu.com/k8s/k8s-dogu-operator:0.35.1`).

#### patches
`patches` enthalten einzelne Templates für beliebige YAML-Dateien des Helm-Charts (z.B. die `values.yaml`).
Jedes Template ist unter dem Dateinamen der zu patchenden Datei abgelegt.
Ein Template kann eine beliebige YAML-Struktur enthalten.  
Es wird die [Go template language](https://godoc.org/text/template) verwendet. 
Die [`values`-map](#values) ist als Daten im Templating verfügbar.

Zusätzlich stehe folgende Template-Funktionen zum Parsen von Container-Image-Referenzen bereit:
  - **registryFrom <string>**: liefert die Registry einer Container-Image-Referenz (z.B. `registry.cloudogu.com`)
  - **repositoryFrom <string>**: liefert das Repository einer Container-Image-Referenz (z.B. `k8s/k8s-dogu-operator`)
  - **tagFrom <string>**: liefert den Tag einer Container-Image-Referenz (z.B. `0.35.1`)

Nachdem ein Template gerendert wurde, wird es in die "originale" YAML-Datei des Helm-Charts gemerged. 
So bleiben Werte in der "originalen" YAML-Datei erhalten, die _nicht_ im Template enthalten sind.
Bereit vorhandene Werte werden vom gerenderten Template überschrieben.

##### Beispiel `component-patch-tpl.yaml`

```yaml
apiVersion: v1

values:
  images:
    engine: longhornio/longhorn-engine:v1.5.1
    manager: longhornio/longhorn-manager:v1.5.1
    ui: longhornio/longhorn-ui:v1.5.1

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
