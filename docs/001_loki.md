# Loki

## Installation

Install the docker plugin that will read logs from Docker containers and send them to Loki.
```bash
$ docker plugin install grafana/loki-docker-driver:2.9.4 --alias loki --grant-all-permissions
```


Verify that the plugin is installed:

```bash
$ docker plugin ls
ID             NAME          DESCRIPTION           ENABLED
2448feba7973   loki:latest   Loki Logging Driver   true
```


To update the plugin:

```bash
docker plugin disable loki --force
docker plugin upgrade loki grafana/loki-docker-driver:2.9.4 --grant-all-permissions
docker plugin enable loki
systemctl restart docker
```


To uninstall the plugin:

```bash
docker plugin disable loki --force
docker plugin rm loki
```
