package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type RequestID struct{}

// WithRequestId - middleware to provide unique uuid ver 2 and push it context.Context of http.Request.
// Use constant middleware.RequestID to get it from context.
func WithRequestId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New().ID()
		ctx := context.WithValue(r.Context(), RequestID{}, reqId)

		next(w, r.WithContext(ctx))
	}
}
