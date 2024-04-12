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


<img width="1440" alt="image" src="https://github.com/alextanhongpin/go-prometheus/assets/6033638/7d14f782-2fa0-46dd-ae14-7a959a3a0891">


<img width="1440" alt="image" src="https://github.com/alextanhongpin/go-prometheus/assets/6033638/0ec819a8-b3f7-4ca4-979f-699963f27814">
