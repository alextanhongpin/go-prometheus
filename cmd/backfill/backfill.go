package main

import (
	"math/rand"
	"os"
	"text/template"
	"time"
)

// # HELP http_requests_total The total number of HTTP requests.
// # TYPE http_requests_total counter http_requests_total{code="200",service="user"} 123 1609954636 http_requests_total{code="500",service="user"} 456 1609954730
// # EOF
type Metric struct {
	UnixTimestamp int64
	Value         int
	Status        int
}

type Metrics struct {
	Data []Metric
}

const t = `# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
{{- range .Data}}
http_requests_total{code="{{.Status}}",service="user"} {{.Value}} {{.UnixTimestamp}}
{{- end}}
# EOF`

// promtool tsdb create-blocks-from openmetrics <input file> [<output directory>]
// promtool tsdb create-blocks-from openmetrics out.txt /prometheus
func main() {
	statuses := []int{200, 500, 404, 503}

	start := time.Now().Unix()

	delta := int64((24 * time.Hour).Seconds())
	batch := int(delta / 5)

	data := make([]Metric, batch)
	var value int
	for i := 0; i < batch; i++ {
		unix := start - delta + int64(i*5)
		status := rand.Intn(len(statuses))
		value += rand.Intn(1_000)
		data[i] = Metric{
			UnixTimestamp: unix,
			Value:         value,
			Status:        statuses[status],
		}
	}

	tpl := template.Must(template.New("").Parse(t))
	tpl.Execute(os.Stdout, Metrics{
		Data: data,
	})
}
