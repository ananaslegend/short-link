package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ananaslegend/short-link/internal/config"
)

func SetupServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadHeaderTimeout: config.DefaultReadHeaderRequestTimeout,
		ReadTimeout:       config.DefaultReadRequestTimeout,
		WriteTimeout:      config.DefaultWriteTimeout,
		IdleTimeout:       config.DefaultIdleTimeout,
	}

	return s
}
