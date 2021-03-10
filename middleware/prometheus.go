package middleware

import (

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)
const (
	METRICS_URL string = "/metrics"
)

func MetricsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == METRICS_URL {
			promhttp.Handler().ServeHTTP(w, req)
			return
		}
		next.ServeHTTP(w, req)
	})
}
