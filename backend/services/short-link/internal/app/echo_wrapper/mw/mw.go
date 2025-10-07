package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/propagation"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func SetupMiddleware(
	router *echo.Echo,
	logger zerolog.Logger,
	traceProvider *sdktrace.TracerProvider,
	metricProvider *sdkmetric.MeterProvider,
	tracePropagator propagation.TextMapPropagator,
) {
	router.Use(ZerologContextMiddleware(logger))

	router.Use(TracePropagationMiddleware(tracePropagator))
	router.Use(TracerMiddleware(traceProvider))

	router.Use(LogTraceIDFromContextMiddleware())
	router.Use(middleware.RequestLoggerWithConfig(SetupEchoRequestLoggerConfig()))

	router.Use(middleware.RecoverWithConfig(
		middleware.RecoverConfig{
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				zerolog.Ctx(c.Request().Context()).
					Err(err).
					Bytes("stack", stack).
					Msg("panic recovered")

				return nil
			},
		}),
	)

	router.Use(REDMetricsMiddleware(metricProvider))
}
