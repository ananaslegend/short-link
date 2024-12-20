package statistic

import (
	"context"
	"net/http"
)

func WithStatisticRow(statManager StatManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			row := NewRow(statManager.FlushTime())

			ctx := context.WithValue(r.Context(), statRowCtxKey{}, row)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithSendingStatistic(statManager StatManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			if row, ok := GetFromCtx(r.Context()); ok {
				if row.IsEmpty() {
					return
				}

				statManager.AppendRow(row)
			}
		})
	}
}
