up:
	@docker-compose up -d


down:
	@docker-compose down


install:
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promauto
	go get github.com/prometheus/client_golang/prometheus/promhttp

run:
	@go run main.go

fire:
	@go run cmd/fire.go


# Requires `brew install hey`
canary:
	hey -z 1m -n 100000 -q 100 -H "x-release-header: canary" http://localhost:8000/

stable:
	hey -z 1m -n 1000 -q 25 -H "x-release-header: stable" http://localhost:8000/


GRAFANA_URL=http://grafana:3000
GRAFANA_USER=admin
GRAFANA_TOKEN=admin
export

render:
	GRAFANA_URL=http://grafana:3000 docker run -v "$(shell pwd):/src/" grafana/grizzly:main-f431d43 apply src/dashboards/data-source.yaml
	@docker run -v "$(shell pwd):/src/" grafana/grizzly:main-f431d43 export dashboards/jsonnet/main.jsonnet dashboards/json/
	@cp dashboards/json/Dashboard/* dashboards/json
	@rm -rf dashboards/json/Dashboard
