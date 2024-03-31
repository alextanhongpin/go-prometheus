package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = ":8080"

var (
	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status", "release"},
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

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func wrap(path string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inFlightGauge.Inc()
		defer inFlightGauge.Dec()

		start := time.Now()
		wr := newStatusCodeResponseWriter(w)
		size := computeApproximateRequestSize(r)
		responseSize.WithLabelValues().Observe(float64(size))

		defer func() {
			duration.WithLabelValues(
				r.Method,                         // method
				path,                             // path
				fmt.Sprintf("%d", wr.statusCode), // status
				r.Header.Get("x-release-header"), // release
			).Observe(time.Since(start).Seconds())
		}()

		next.ServeHTTP(wr, r)
	})
}

// Copied from prometheus source code.
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

func registerHandler(mux *http.ServeMux, path string, h http.HandlerFunc) {
	mux.Handle(path, wrap(path, h))
}

func main() {
	reg := prometheus.NewRegistry()
	// Install the default prometheus collectors.
	reg.MustRegister(prometheus.NewGoCollector())
	// Install the custom metrics.
	reg.MustRegister(inFlightGauge, duration, responseSize)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// Matches only the path '/'.
	registerHandler(mux, "GET /{$}", getHandler)
	registerHandler(mux, "POST /{$}", postHandler)

	logger.Info("Server is running on port " + port)
	graceful(port, mux)
}

func graceful(port string, h http.Handler) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    port,
		Handler: h,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	release := r.Header.Get("x-release-header")

	logger.Info("get handler", slog.String("release", release))
	threshold := 90
	if release == "canary" {
		threshold = 50
	}

	if rand.Intn(100) > threshold {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("hello world"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	logger.Info("post handler", slog.String("body", string(b)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("created"))
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newStatusCodeResponseWriter(w http.ResponseWriter) *statusCodeResponseWriter {
	return &statusCodeResponseWriter{w, http.StatusOK}
}

func (rw *statusCodeResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *statusCodeResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
