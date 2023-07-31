package mw

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const RequestId = "requestId"

// WithRequestId - middleware to provide unique uuid ver 2 and push it context.Context of http.Request.
// Use constant mw.RequestId to get it from context.
func WithRequestId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New().ID()
		ctx := context.WithValue(r.Context(), RequestId, reqId)

		next(w, r.WithContext(ctx))
	}
}
