package mw

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"time"
)

func MetricsMiddleware(metricProvider *sdkmetric.MeterProvider) echo.MiddlewareFunc {
	meter := metricProvider.Meter("http-server")

	requestCounter, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		panic(err)
	}

	requestDuration, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
	)
	if err != nil {
		panic(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			method := c.Request().Method
			route := c.Request().URL.Path

			err = next(c)

			status := c.Response().Status
			duration := time.Since(start).Seconds()

			attrs := []attribute.KeyValue{
				attribute.String("method", method),
				attribute.String("route", route),
				attribute.Int("status", status),
			}

			requestCounter.Add(c.Request().Context(), 1, metric.WithAttributes(attrs...))
			requestDuration.Record(c.Request().Context(), duration, metric.WithAttributes(attrs...))

			return err
		}
	}
}
