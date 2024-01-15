package metrics

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Listen(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, mux)
}
