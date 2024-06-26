package middleware

import (
	"context"
	"github.com/ananaslegend/short-link/pkg/clog"
	"github.com/google/uuid"
	"net/http"
)

type reqIDKey struct{}

var (
	RequestIDHeaderKey = "X-Request-ID"
	RequestIDLogKey    = "request_id"
)

// WithRequestID middleware try to get request ID from request with [RequestIDHeaderKey] or generates new one
// and sets it to response header with [RequestIDHeaderKey], log message with [RequestIDLogKey] and to context.
func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeaderKey)

		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := clog.WithString(r.Context(), RequestIDLogKey, reqID)

		ctx = context.WithValue(ctx, reqIDKey{}, reqID)

		SetRequestIDHeader(ctx, w.Header())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestID(ctx context.Context) string {
	reqID, _ := ctx.Value(reqIDKey{}).(string)

	if reqID == "" {
		clog.Ctx(ctx).Error("request_id not found")
	}

	return reqID
}

func SetRequestIDHeader(ctx context.Context, headers http.Header) {
	headers.Set(RequestIDHeaderKey, RequestID(ctx))
}
