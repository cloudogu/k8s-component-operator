# Developing with a proxy for external queries to dogu container or helm registry

Component can mount a proxy stored in the `ces-proxy` secret.
To test this behaviour be sure to configure the proxy in the setup values of the setup's helm chart,
or create the secret manually by:

`kubectl create secret generic ces-proxy --from-literal=url=http://test:test@192.168.56.1:3128 -n ecosystem`

## Setup local proxy in docker

### Run container with host network mode to reach the development registry in the cluster if required.

- `docker run --net=host -d --name squid -e TZ=UTC -p 3128:3128 ubuntu/squid:5.2-22.04_beta`

### Configure auth and connection to dev registry

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

### Check access logs

- `docker logs -f squid`
