package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ananaslegend/short-link/pkg/cslog"
)

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					cslog.Logger(r.Context()).Error("app in panic!", slog.Any("request:", r))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		},
	)
}
