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

## Todo

- setup list of metrics to scrape
- add visualization with grafana
- play with labels


## Success rate

```promql
(sum(api_requests_total{code=~"20+"} > 0) or vector(0)) / (sum(api_requests_total) or vector(0)) * 100
```

```
histogram_quantile(0.95, sum(rate(request_duration_seconds_bucket[5m])) by(le))
```
