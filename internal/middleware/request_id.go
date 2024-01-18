package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type requestID struct{}

// WithRequestID - middleware to provide unique uuid ver 2 and push it context.Context of http.Request.
// Use constant GetRequestID to get it from context.
func WithRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New().ID()
		ctx := context.WithValue(r.Context(), requestID{}, reqId)

		next(w, r.WithContext(ctx))
	}
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(requestID{}).(string)
}
