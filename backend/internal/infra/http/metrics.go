package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServeMetrics(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, mux)
}
