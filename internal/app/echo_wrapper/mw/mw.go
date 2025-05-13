package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func SetupMiddleware(router *echo.Echo, logger zerolog.Logger, traceProvider *sdktrace.TracerProvider, metricProvider *sdkmetric.MeterProvider) {
	router.Use(middleware.Recover())

	router.Use(ZerologContextMiddleware(logger))

	router.Use(middleware.RequestLoggerWithConfig(SetupEchoRequestLoggerConfig()))

	router.Use(middleware.RequestIDWithConfig(SetupEchoRequestIDConfig()))

	router.Use(TracerMiddleware(traceProvider))
	router.Use(MetricsMiddleware(metricProvider))
}
