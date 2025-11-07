# Verwendung von `k8s-component-operator`

Der Komponenten-Operator `k8s-component-operator` ist eine Komponente für die Kubernetes-Version des Cloudogu EcoSystems (K8s-CES). Dieser Operator ermöglicht es, Komponenten auf einfache Weise zu installieren, zu aktualisieren oder zu löschen. Diese Komponenten stellen ihrerseits erforderliche Dienste für das EcoSystem bereit.

## Installation des Komponenten-Operators

### Helm-Repository konfigurieren

Um initial den Komponenten-Operator zu installieren, muss die Cloudogu Helm-Registry bekannt gemacht werden.

```bash
$ helm registry login -u myuser registry.cloudogu.com
Password: ************************
Login succeeded
```

Für spätere K8s-CES-Komponenten ist dieses Helm-Repository unnötig, da der Komponenten-Operator über seine eigene Konfiguration verfügt. Siehe hierzu den nächsten Abschnitt [Zugangsdaten konfigurieren](#Zugangsdaten-konfigurieren).

### Zugangsdaten konfigurieren

Der Komponenten-Operator verfügt über seine eigene Konfiguration hinsichtlich Endpunkt und Zugangsdaten. Wenn die K8s-CES-Instanz auf das Internet zugreifen kann, so sind Endpunkt und Zugangsdaten identisch mit denen der Dogu-Registry:
- Endpunkt: `oci://registry.cloudogu.com`
- Zugangsdaten: Der gleiche Benutzer/Passwort, wie die aus dem Secret `k8s-dogu-operator-dogu-registry`

Diese Konfiguration kann für den Cluster-Namespace `ecosystem` wie folgt manuell erzeugt werden:

```bash
$ kubectl -n ecosystem create configmap component-operator-helm-repository --from-literal=endpoint="${HELM_REPO_ENDPOINT}" --from-literal=schema=oci
$ kubectl -n ecosystem create secret generic component-operator-helm-registry \
  --from-literal=config.json='{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "$(shell printf "%s:%s" "${HELM_REPO_USERNAME}" "${HELM_REPO_PASSWORD}" | base64 -w0)"}}}'
```

### Komponenten-Operator installieren

Normalerweise wird der Komponenten-Operator vom `k8s-ces-setup` installiert. Manuell geschieht dies für den Cluster-Namespace `ecosystem` und den Helm-Registry-Namespace `k8s` wie folgt:

```bash
$ helm install -n ecosystem k8s-component-operator oci://${HELM_REPO_ENDPOINT}/k8s/k8s-component-operator --version ${DESIRED_VERSION}
```

### Komponenten-Operator deinstallieren

```bash
$ helm uninstall -n ecosystem k8s-component-operator
```

## Komponenten installieren oder aktualisieren

Um Komponenten zu installieren oder zu aktualisieren, muss jeweils eine _Custom Resource_ (CR) für die gewünschte Komponente auf den Cluster im korrekten Cluster-Namespace angewendet werden.

Beispiel einer Komponenten-Ressource (z. B. als `k8s-longhorn.yaml` und aus dem Helm-Registry-Namespace `k8s`):

```yaml
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-longhorn
spec:
  name: k8s-longhorn
  namespace: k8s
  version: 1.5.1-1
  deployNamespace: longhorn-system
  mappedValues:
    mainLogLevel: debug
  valuesYamlOverwrite: |
    longhorn:
      defaultSettings:
        backupTargetCredentialSecret: my-longhorn-backup-target
```

> [!IMPORTANT]
> `metadata.name` und `spec.name` müssen gleich sein.
> Ansonsten schlägt die Installation fehl.

Diese CR kann dann auf den Cluster angewendet werden:

```bash
kubectl -n ecosystem apply -f k8s-longhorn.yaml
```

Der Komponenten-Operator beginnt nun mit der Installation der Komponente. Abhängigkeiten zu anderen k8s-CES-Komponenten und deren Versionen müssen erfüllt sein (dies überprüft der Komponenten-Operator). Weitere Informationen zu diesem Thema befinden sich im Abschnitt [Abhängigkeiten zu anderen Komponenten](#Abhängigkeiten-zu-anderen-Komponenten).

Beispiele von Komponenten-Ressourcen befinden sich im [config/samples-Verzeichnis](../../config/samples)

### Felder und deren Bedeutung:

Ein Komponenten-CR besteht aus unterschiedlichen Feldern. Dieser Abschnitt erläutert diese:

- `.metadata.name`: Der Komponentenname der Kubernetes-Resource. Dieser muss identisch mit `.spec.name` sein.
- `.spec.name`: Der Komponentenname, wie er in der Helm-Registry lautet. Dieser muss identisch mit `.metadata.name` sein.
- `.spec.namespace`: Der Namespace der Komponente in der Helm-Registry. 
  - Mittels unterschiedlicher Komponenten-Namespaces können unterschiedliche Versionen ausgebracht werden (z. B. zu Debugging-Zwecken). 
  - Es handelt sich hierbei _nicht_ um den Cluster-Namespace.
- `.spec.version`: Die Version der Komponente in der Helm-Registry.
- `.spec.deployNamespace`: (optional) Der k8s-Namespace, in dem alle Ressourcen der Komponente deployed werden sollen. Wenn dieser leer ist, wird der Namespace des Komponenten-Operators verwendet.
- `.spec.mappedValues`: (optional) Helm-Werte zum Überschreiben von Konfigurationen aus der Helm-Datei values.yaml. Diese Werte werden durch die Konfiguration in component-values-metadata.yaml gemappt. 
- `.spec.valuesYamlOverwrite`: (optional) Helm-Werte zum Überschreiben von Konfigurationen aus der Helm-Datei values.yaml. Sollte aus Gründen der Lesbarkeit als [multiline-yaml](https://yaml-multiline.info/) geschrieben werden.
- `.spec.valuesConfigRef`: (optional) Referenz auf einen Key einer ConfigMap (`valuesConfigRef.name`, `valuesConfigRef.key`). zum Überschreiben von Konfigurationen aus der Helm-Datei values.yaml. Die Werte sollten als [multiline-yaml](https://yaml-multiline.info/) in den Key geschrieben sein.



> [!WARNING]
> `.spec.mappedValues`, `.spec.valuesYamlOverwrite` und `.spec.valuesConfigRef` sollten nicht gleichzeitig verwendet werden. Sind beide Werte konfiguriert, so bekommen die mappedValues den Vorzug. 
> Priorität haben `mappedValues` vor `valuesYamlOverwrite` und `valuesYamlOverwrite` vor `valuesConfigRef`

> [!WARNING]
> `.spec.mappedValues`, `.spec.valuesYamlOverwrite` und `.spec.valuesConfigRef` dürfen keine Listeneinträge überschreiben. Es ist durch die Struktur von Yaml nicht möglich einzelne Elemente innerhalb einer Liste zu setzen. 
>  Es kann immer nur die gesamte Liste überschrieben werden.

## Komponenten deinstallieren

> [!WARNING]
> Löschen von Komponenten, die einen Zustand besitzen, kann die Stabilität der K8s-CES-Installation betriebsverhindernd stören.
> Dies gilt insbesondere (aber nicht ausschließlich) für die Komponente `k8s-etcd`.

- Löschen der Komponenten-CR aus dem Cluster kann auf zwei Arten erfolgen:
   1. durch Anwendung einer existierenden Komponenten-CR-Datei, z. B. `kubectl -n ecosystem delete -f k8s-dogu-operator.yaml`
   2. durch Angabe von `.metadata.name` der Komponenten, z. B. `kubectl -n ecosystem delete component k8s-dogu-operator`
- Der Komponenten-Operator beginnt nun mit der Deinstallation der Komponente

## Abhängigkeiten zu anderen Komponenten

K8s-CES-Komponenten können von anderen k8s-CES-Komponenten abhängen. Um sicherzustellen, dass eine Komponente voll funktionsfähig ist, wird während der Installation bzw. Aktualisierung geprüft, ob Komponentenabhängigkeiten vorhanden sind und diese eine korrekte Version aufweisen.

Sollte eine oder mehrere Komponenten fehlen oder nicht in der richtigen Version vorhanden sein, so wird ein Event an der betroffenen Komponenten-Resource angefügt. Fehler wie diese können durch ein `kubectl describe` an der jeweiligen Komponten-Resource erkannt werden:

```bash
$ kubectl -n ecosystem describe component k8s-dogu-operator
```

In diesem Fall müssen die betroffenen Komponenten manuell [nachinstalliert oder aktualisiert](#Komponenten-installieren-oder-aktualisieren) werden.

Die Versionen zu Abhängigkeiten werden während der Komponentenentwicklung im Helm-Chart hinterlegt. Diese können i. d. R. nicht zum Installationszeitpunkt geändert werden.

## Konfigurationswerte mappen

Um zur Laufzeit Werte der values.yaml überschreiben zu können, kann das Feld `.spec.mappedValues` verwendet werden. 
Dies setzt jedoch vorraus, dass die entsprechende Komponente auch eine `component-values-metadata.yaml` 
Datei im Helm-Chart bereitstellt.

Eine Konfiguration innerhalb einer CR sieht beispiel haft wie folgt aus:

```yaml
spec:
  mappedValues:
    mainLogLevel: debug
```

Eine dazugehörige Mapping-Datei muss den gemappten Wert aufgreifen und entsprechend konfigurieren.

```yaml
apiVersion: v1
metavalues:
  mainLogLevel:
    name: Log-Level
    description: The central configuration value to set the log level for this component
    keys:
      - path: controllerManager.env.logLevel
        mapping:
          debug: trace
          info: info
          warn: warn
          error: error
      - path: manager.env.logLevel
        mapping:
          debug: debug
          panic: error
```

Hier wird nun der ursprüngliche Helm-Wert `.controllerManager.env.logLevel` durch den Wert aus dem CR für mainLogLevel ersetzt.
Dabei wird der Wert noch gegen ein Liste von Value-Mappings verglichen und entsprechend angepasst.

Der finale Eintrag für `.controllerManager.env.logLevel` in obrigen Beispiel, enthält somit den Wert `trace`.
Es kann für einen Mapping-Eintrag auch mehrere zumappenden Schlüssel geben. Dabei muss jeder Schlüssel sein eigenes Value-Mapping definieren.

### Besonderheiten
Da durch diesen Mechanismus sowohl durch `mappedValues` als auch durch `valuesYamlOverwrite` die selben Werte gesetzt werden können, kann es zu
Konflikten kommen.
In diesem Fall wird vom component-operator automatisch geprüft, ob es einen Konflikt gibt und eine entsprechende Fehlermeldung ausgegeben.
Dann hat der Wert, der in `mappedValues` eingetragen wurde Priorität gegenüber dem Wert aus `valuesYamlOverwrite`.
Es führt darüber hinaus aber nicht zu einem Fehlverhalten, der Konflikt wird lediglich im Log sichtbar.



