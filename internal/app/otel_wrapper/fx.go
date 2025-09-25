package otelwrapper

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ananaslegend/short-link/internal/app/config"
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
				panic(fmt.Sprintf("failed to initialize otel trace exporter: %v", err))
			}

			return spanExporter
		}),

		fx.Provide(NewTraceProvider),

		fx.Provide(func(ctx context.Context, cfg config.Config) sdkmetric.Exporter {
			metricExporter, err := NewMetricExporter(ctx, cfg)
			if err != nil {
				panic(fmt.Sprintf("failed to initialize otel metric exporter: %v", err))
			}

			return metricExporter
		}),

		fx.Provide(NewMetricProvider),

		fx.Provide(NewPropagator),
	)
}
