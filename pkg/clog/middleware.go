package clog

import (
	"context"
	"log/slog"
	"net/http"
)

func WithCtxLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = context.WithValue(ctx, Key{}, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
