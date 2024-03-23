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
