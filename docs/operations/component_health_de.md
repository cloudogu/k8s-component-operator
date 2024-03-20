# Komponenten-Health

Der Komponenten-Operator schreibt für im Multinode-EcoSystem installierte Komponenten
den aktuellen Health-Zustand in den Status der `Component`-Ressource.
Aktuell gibt es folgende mögliche Health-Zustände:
- Ausstehend (`""`) - das Feld ist bisher noch nicht befüllt.
- Verfügbar (`"available"`)
- Nicht verfügbar (`"unavailable"`)
- Unbekannt (`"unknown"`) - der Health-Zustand der Komponente ist unbekannt oder kann nicht garantiert werden,
  [siehe unten](#sonderfall-abschalten-des-komponenten-operators).

## Sonderfall: Abschalten des Komponenten-Operators

Wird der Komponenten-Operator regulär beendet (mittels `sigint` oder `sigterm`),
dann setzt er seine eigene Health auf `unavailable` und die aller anderen Komponenten auf `unknown`,
um irreführende Zustände zu vermeiden.