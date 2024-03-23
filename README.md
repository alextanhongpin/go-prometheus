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

## RED

For every service, monitor request:
- Rate - the number of requests per second
- Errors - the number of those requests that are failing
- Duration - the amount of time those requests are taking

```promql
# Rate
rate(api_requests_total[1m])

# Errors
rate(api_requests_total{code=~"2.."}[1m])

# Duration
histogram_quantile(0.95, sum(rate(request_duration_seconds_bucket[1m])) by (le))
```

## Metrics


We use `hey` to generate load on the server.

```bash
brew install hey
# -n number of requests
# -z duration of the test
# -q queries per second
hey -z 1m -n 100000 -q 100 -H "x-release-header: canary" http://localhost:8000/
hey -z 1m -n 1000 -q 25 -H "x-release-header: stable" http://localhost:8000/
```

### Number of requests per minute

Monitor the number of requests per minute to detect spikes in traffic

```promql
rate(api_requests_total[1m])

# By release
sum by(release) (rate(api_requests_total[5m]))
sum(rate(api_requests_total[5m])) by(release)
```

### Success rate

Success rate of requests in percentage in the last 1 minute.

```promql
(sum(rate(api_requests_total{code=~"20+"}[1m]) > 0) or vector(0)) / (sum(rate(api_requests_total[1m])) or vector(0)) * 100

# By release
sum(rate(api_requests_total{code=~"20+"}[1m])) by (release) / sum(rate(api_requests_total[1m])) by (release) * 100
```


### 95th percentile of request duration

```bash
histogram_quantile(0.95, sum(rate(request_duration_seconds_bucket[5m])) by (path))
```
