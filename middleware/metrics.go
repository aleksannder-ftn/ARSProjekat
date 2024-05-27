package middleware

import (
	"ars_projekat/model"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func count(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := &model.ResponseWriter{w, http.StatusOK}
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
