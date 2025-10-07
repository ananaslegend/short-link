package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	traceIDKey = "trace_id"
)

func TracerMiddleware(traceProvider *sdktrace.TracerProvider) echo.MiddlewareFunc {
	return otelecho.Middleware("short-link", otelecho.WithTracerProvider(traceProvider))
}

func TracePropagationMiddleware(propagator propagation.TextMapPropagator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := propagator.Extract(
				c.Request().Context(),
				propagation.HeaderCarrier(c.Request().Header),
			)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func LogTraceIDFromContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := trace.SpanFromContext(c.Request().Context())

			if span != nil && span.SpanContext().IsValid() {
				traceID := span.SpanContext().TraceID().String()

				loggingTraceIdCtx := zerolog.Ctx(c.Request().Context()).
					With().
					Str(traceIDKey, traceID).
					Logger().
					WithContext(c.Request().Context())

				c.SetRequest(c.Request().WithContext(loggingTraceIdCtx))
				c.Set(traceIDKey, traceID)
			}

			return next(c)
		}
	}
}
