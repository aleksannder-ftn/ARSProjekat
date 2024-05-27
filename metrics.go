package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ttp_total_request",
			Help: "Total number of HTTP requests in last 24h",
		},
		[]string{},
	)

	httpSuccessfulRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_successful_requests",
			Help: "Number of successful HTTP requests in last 24h (2xx, 3xx).",
		},
		[]string{},
	)

	httpUnsuccessfulRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_unsuccessful_requests",
			Help: "Number of unsuccessful HTTP requests in last 24h (4xx, 5xx).",
		},
		[]string{},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// this list will be collected
	metricsList = []prometheus.Collector{
		httpTotalRequests,
		httpSuccessfulRequests,
		httpUnsuccessfulRequests,
		httpRequestDuration,
	}

	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := &responseWriter{w, http.StatusOK}
		f(rw, r)

		duration := time.Since(start).Seconds()

		statusCode := rw.statusCode
		if statusCode >= 200 && statusCode < 400 {
			httpSuccessfulRequests.WithLabelValues().Inc()
		} else if statusCode >= 400 && statusCode < 600 {
			httpUnsuccessfulRequests.WithLabelValues().Inc()
		}

		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		httpTotalRequests.WithLabelValues().Inc()

	}
}
