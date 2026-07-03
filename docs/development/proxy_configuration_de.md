# Entwicklung mit einem Proxy für externe Abfragen an Dogu-Container oder Helm-Registry

Die Komponente kann einen Proxy einbinden, der im Secret `ces-proxy` gespeichert ist.
Um dieses Verhalten zu testen, stelle sicher, dass der Proxy in den Setup-Values des Helm-Charts des Setups konfiguriert ist,
oder erstelle das Secret manuell mit:

`kubectl create secret generic ces-proxy --from-literal=url=http://test:test@192.168.56.1:3128 -n ecosystem`

## Lokalen Proxy in Docker einrichten

### Container mit Host-Netzwerkmodus starten, um bei Bedarf die Entwicklungs-Registry im Cluster zu erreichen.

- `docker run --net=host -d --name squid -e TZ=UTC -p 3128:3128 ubuntu/squid:5.2-22.04_beta`

### Authentifizierung und Verbindung zur Entwicklungs-Registry konfigurieren

- `docker exec -it squid /bin/bash`
- `apt update && apt-get install -y apache2-utils`
- `htpasswd -c -d /etc/squid/passwords test`
- `chmod o+r /etc/squid/passwords`
- `apt-get install vim`
- `vi /etc/squid/conf.d/auth.conf`
    - ```
      acl allcomputers src all
      auth_param basic program /usr/lib/squid/basic_ncsa_auth /etc/squid/passwords
      auth_param basic realm proxy
      acl authenticated proxy_auth REQUIRED
      http_access allow authenticated allcomputers
      ```
- `echo "localhost k3ces.local" >> /etc/hosts`
- `squid -k reconfigure`

### Zugriffslogs prüfen

- `docker logs -f squid`

## Lokalen Proxy in k8s einrichten

#### 1. Squid-Proxy-Server installieren

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: squid-auth-config
data:
  auth.conf: |
    acl allcomputers src all
    http_access allow allcomputers
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: squid-hosts-config
data:
  hosts: |
    127.0.0.1   localhost
    ::1         localhost ip6-localhost ip6-loopback
    localhost ecosystem.svc.cluster.local

---
apiVersion: v1
kind: Pod
metadata:
  name: squid-proxy
  labels:
    app: squid
spec:
  containers:
    - name: squid
      image: ubuntu/squid:latest
      ports:
        - containerPort: 3128
          name: proxy
          protocol: TCP
      volumeMounts:
        - name: auth-config
          mountPath: /etc/squid/conf.d/auth.conf
          subPath: auth.conf
        - name: hosts-config
          mountPath: /etc/hosts
          subPath: hosts
      resources:
        requests:
          memory: "256Mi"
          cpu: "250m"
        limits:
          memory: "512Mi"
          cpu: "500m"
  volumes:
    - name: auth-config
      configMap:
        name: squid-auth-config
    - name: hosts-config
      configMap:
        name: squid-hosts-config

---
apiVersion: v1
kind: Service
metadata:
  name: squid-proxy-service
spec:
  selector:
    app: squid
  ports:
    - protocol: TCP
      port: 3128
      targetPort: 3128
  type: ClusterIP
```

#### 2. Proxy-Einstellungen zur global-config hinzufügen

```
    proxy:
      enabled: "true"
      server: "squid-proxy-service.ecosystem.svc.cluster.local"
      port: "3128"
```
