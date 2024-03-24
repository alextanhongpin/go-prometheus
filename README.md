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
sum(rate(request_duration_seconds_count[1m])) by (release)

# Errors
sum(rate(request_duration_seconds_count{status=~"2.."}[1m])) by (release)

# Duration
histogram_quantile(0.95, sum(rate(request_duration_seconds_bucket[1m])) by (le, release))
```

Reference: https://grafana.com/files/grafanacon_eu_2018/Tom_Wilkie_GrafanaCon_EU_2018.pdf

## Custom Metrics

- <entity> created
- gmv recorded
- transition state (paid, refund, success, error)
- views

## Labels

Labels helps us group metrics so that they can be observed independently. For example, when monitoring HTTP error rate, we want to know which endpoint has the highest error rate. Adding a label `method` and `path` allows us to do so.

Without the labels, it will show the cumulative error rate for all endpoints.

When you have multiple microservice, it is important to share the same naming conventions, but use proper labels to differentiate the metrics.

An alternative is to just namespace the metrics.


Some useful labels includes
- app - the name of the app
- release - whether it is `stable` release or `canary` release. Allows us to understand the impact of the release during deployments and rollback when there is unexpected errors
- status - http status code, or just `failed`/`success` if we want to normalize the status
- path - the path template without query string and params, e.g. `/users/{id}`
- method - http method like GET/POST


## Simulate


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

## IAC

IAC allows us to create the dashboards easily. There are two ways to do it:
- using Terraform
- using grizzly


### Importing existing dashboard using terraform

First, add the new resource inside the `provider.tf`:

```tf
resource "grafana_dashboard" "another_dashboard" {
  folder = grafana_folder.my_folder.id
  config_json = file("${path.module}/another.json")
}
```

Then, go to the dashboard you wish to import:

```bash
http://localhost:3000/d/ee06ace5-cccd-4fc7-a761-303063f9f345/another-dashboard?orgId=1
```

The id `ee06ace5-cccd-4fc7-a761-303063f9f345` will be used when importing.

Run the import:

```bash
terraform import grafana_dashboard.another_dashboard  ee06ace5-cccd-4fc7-a761-303063f9f345
```

The `terraform.tfstate` should have the `config_json`. Copy paste it into the `another.json`.

Then run `terraform plan` and `terraform apply`.
