package middleware

import (
	"net/http"

	"github.com/ananaslegend/short-link/pkg/clog"
)

type wrappedRespWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedRespWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func WithLoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := clog.Ctx(r.Context())

		wrappedWriter := &wrappedRespWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(wrappedWriter, r)

		logger.
			With("method", r.Method).
			With("path", r.URL.String()).
			//With("body", r.Body).
			With("status", wrappedWriter.statusCode).
			Info("request")
	})
}
