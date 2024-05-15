package middleware

import (
	"context"
	"errors"
	"github.com/ananaslegend/go-logs/v2"
	"github.com/google/uuid"
	"net/http"
)

type requestID struct{}

var (
	ErrNoRequestID = errors.New("no request id in context")
)

// WithRequestID - middleware to provide unique uuid ver 2 and push it context.Context of http.Request.
func WithRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()
		ctx := context.WithValue(r.Context(), requestID{}, reqID)

		ctx = logs.WithMetric(ctx, "request_id", reqID)

		next(w, r.WithContext(ctx))
	}
}

func GetRequestID(ctx context.Context) (uuid.UUID, error) {
	reqID, ok := ctx.Value(requestID{}).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrNoRequestID
	}

	return reqID, nil
}
