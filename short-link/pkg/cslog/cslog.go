package cslog

import (
	"context"
	"log/slog"
	"net/http"
)

type ContextSlog struct {
	slog.Handler
}

func New(handler slog.Handler) ContextSlog {
	return ContextSlog{Handler: handler}
}

func (cs ContextSlog) Handle(ctx context.Context, rec slog.Record) error {
	attrs := Attrs(ctx)

	for _, attr := range attrs {
		rec.Add(attr)
	}

	return cs.Handler.Handle(ctx, rec)
}

type loggerContextKey struct{}
type loggerAttrsKey struct{}

func Logger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(&loggerContextKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

func SetContextLogger(ctx context.Context, logger *slog.Logger) context.Context {
	ctx = context.WithValue(ctx, &loggerContextKey{}, logger)

	return context.WithValue(ctx, &loggerAttrsKey{}, []slog.Attr{})
}

func With(ctx context.Context, args ...slog.Attr) context.Context {
	if attrs, ok := ctx.Value(&loggerAttrsKey{}).([]slog.Attr); ok {
		attrs = append(attrs, args...)

		return context.WithValue(ctx, &loggerContextKey{}, attrs)
	}

	return ctx
}

func Attrs(ctx context.Context) []slog.Attr {
	if attrs, ok := ctx.Value(&loggerAttrsKey{}).([]slog.Attr); ok {
		return attrs
	}

	return nil
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func Middleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = SetContextLogger(ctx, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
