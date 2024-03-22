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
