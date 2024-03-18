# Setting up Prometheus


There are several examples on how to instrument a Golang application with prometheus.

We will use the example from the official prometheus documentation:

https://prometheus.io/docs/guides/go-application/



## Installation

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```
