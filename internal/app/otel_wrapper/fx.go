package otelwrapper

import (
	"context"
	"fmt"
	"github.com/ananaslegend/short-link/internal/app/config"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"internal.app.otel",

		fx.Provide(func(ctx context.Context, cfg config.Config) *sdkresource.Resource {
			resource, err := NewResource(ctx, cfg)
			if err != nil {
				panic(fmt.Sprintf("failed to initialize otel resource: %v", err))
			}

			return resource
		}),

		fx.Provide(func(ctx context.Context, cfg config.Config) sdktrace.SpanExporter {
			spanExporter, err := NewSpanExporter(ctx, cfg)
			if err != nil {
				panic(fmt.Sprintf("failed to initialize otel stdout trace exporter: %v", err))
			}

			return spanExporter
		}),

		fx.Provide(NewTraceProvider),

		fx.Invoke(func(e *echo.Echo, traceProvider *sdktrace.TracerProvider) {
			e.Use(otelecho.Middleware("http-request", otelecho.WithTracerProvider(traceProvider)))
		}),
	)
}
