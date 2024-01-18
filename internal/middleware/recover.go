package middleware

import (
	"log/slog"
	"net/http"
)

func WithRecover(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("app in panic!", slog.Any("request:", r))
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}()
			next.ServeHTTP(w, r)
		},
	)
}
