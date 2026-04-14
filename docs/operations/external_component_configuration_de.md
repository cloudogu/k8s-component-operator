# Externe Komponenten Konfiguration

In manchen Fällen ist es nötig, dass Konfigurationen von Komponenten innerhalb des Cluster an anderen Komponenten vorgenommen werden müssen.
Dafür sind ConfigMaps eine ideale Wahl zur Verwaltung der externen Konfiguration.
Für die Konfiguration über ConfigMaps gibt es zurzeit die Möglichkeit über `.spec.valuesConfigRef`. 
Dabei wird die ConfigMap jedoch in der Spec mit angegeben und müsste bei jedem Deployen der Komponente mit konfiguriert werden.
Dies geht mit manuellem Aufwand einher. Zudem ist diese Konfigurationsmöglichkeit eher für die individuelle Konfiguration des Kunden gedacht.

Bspw. soll es in `k8s-ces-gateway` Ports dynamisch freizugeben. Dies ist über traefik selbst nicht dynamisch möglich.
Die Ports werden in der `values.yaml` statisch festgelegt. Nur über einen Neustart ist die Konfiguration möglich.

Der `k8s-componnent-operator` reconciled ConfigMaps mit dem Label `k8s.cloudogu.com/component.config` (z.B. `k8s.cloudogu.com/component.config: k8s-ces-gateway`).
Diese ConfigMaps enthalten unter dem Key `values` eine `values.yaml` der Komponente.
Die Komponente mit dem Namen des Values des Labels wird gereconciled. 
Dabei werden die Werte aus der ConfigMap mit den anderen Konfigurationswerten der Komponente gemergt.
Die Werte aus der ConfigMap habe eine niedrige Priorität, sodass sie bspw. mit den `MappedValues` überschrieben 
werden können.