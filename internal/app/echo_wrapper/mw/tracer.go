package mw

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TracerMiddleware(traceProvider *sdktrace.TracerProvider) echo.MiddlewareFunc {
	return otelecho.Middleware("http-request", otelecho.WithTracerProvider(traceProvider))
}
