package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"pattern", "handler", "method"},
	)

	// responseSize has no labels, making it a zero-dimensional
	// ObserverVec.
	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)
)

func wrapHandler(pattern string, h http.HandlerFunc) http.Handler {
	// Get handler name from h
	handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()

	return promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{
			// Registers the URL pattern.
			"pattern": pattern,
			// Registers the handler name.
			"handler": handlerName,
		}),
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerResponseSize(responseSize, http.Handler(h)),
			),
		),
	)
}

func registerHandler(pattern string, h http.HandlerFunc) {
	http.Handle(pattern, wrapHandler(pattern, h))
}

func main() {
	reg := prometheus.NewRegistry()
	// Install the default prometheus collectors.
	reg.MustRegister(prometheus.NewGoCollector())
	// Install the custom metrics.
	reg.MustRegister(inFlightGauge, counter, duration, responseSize)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// Matches only the path '/'.
	registerHandler("GET /{$}", getHandler)
	registerHandler("POST /{$}", postHandler)

	log.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(100) > 90 {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(r.URL.Path)
	w.Write([]byte("hello world"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("created"))
}
