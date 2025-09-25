package mw

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func REDMetricsMiddleware(metricProvider *sdkmetric.MeterProvider) echo.MiddlewareFunc {
	meter := metricProvider.Meter("http-server")

	requestCounter, err := meter.Int64Counter(
		"http.server.requests",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		panic(err)
	}

	errorCounter, err := meter.Int64Counter(
		"http.server.errors",
		metric.WithDescription("Total number of 5xx HTTP responses"),
	)
	if err != nil {
		panic(err)
	}

	requestDuration, err := meter.Float64Histogram(
		"http.server.duration",
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
				attribute.String("http.method", method),
				attribute.String("http.route", route),
				attribute.Int("http.status_code", status),
				attribute.Float64("http.duration", duration),
				attribute.String("trace_id", fmt.Sprintf("%s", c.Get(traceIDKey))),
			}

			requestCounter.Add(c.Request().Context(), 1, metric.WithAttributes(attrs...))

			if status >= http.StatusInternalServerError {
				errorCounter.Add(c.Request().Context(), 1, metric.WithAttributes(attrs...))
			}

			requestDuration.Record(c.Request().Context(), duration, metric.WithAttributes(attrs...))

			return err
		}
	}
}
