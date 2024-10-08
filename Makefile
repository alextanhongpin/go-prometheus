up:
	@docker compose up -d


down:
	@docker compose down


install:
	#https://grafana.com/docs/loki/latest/send-data/docker-driver/
	docker plugin install grafana/loki-docker-driver:2.9.2 --alias loki --grant-all-permissions
	docker plugin enable loki
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promauto
	go get github.com/prometheus/client_golang/prometheus/promhttp

run:
	@go run main.go

fire:
	@go run cmd/fire.go


# Requires `brew install hey`
canary:
	hey -z 1m -n 100000 -q 100 -H "x-release-header: canary" http://localhost:8080/

stable:
	hey -z 1m -n 1000 -q 25 -H "x-release-header: stable" http://localhost:8080/


GRAFANA_URL=http://grafana:3000
GRAFANA_USER=admin
GRAFANA_TOKEN=admin
export

grizzly:
	docker run -it -v "$(shell pwd):/src/" --entrypoint="/bin/sh" grafana/grizzly:main-f431d43

.PHONY: terraform
terraform:
	docker run -it -v "$(shell pwd)/terraform:/src" --entrypoint="/bin/sh" hashicorp/terraform:1.7.5

rebuild:
	docker compose build app
	docker compose up -d
