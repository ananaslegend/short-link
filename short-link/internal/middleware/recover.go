package middleware

import (
	"github.com/ananaslegend/short-link/pkg/clog"
	"log/slog"
	"net/http"
)

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					clog.Ctx(r.Context()).Error("app in panic!", slog.Any("request:", r))
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}()
			next.ServeHTTP(w, r)
		},
	)
}
