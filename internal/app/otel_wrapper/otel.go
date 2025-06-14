//nolint:ireturn
package otelwrapper

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"github.com/ananaslegend/short-link/internal/app/config"
)

func NewResource(ctx context.Context, cfg config.Config) (*sdkresource.Resource, error) {
	resource, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.DeploymentEnvironmentName(string(cfg.Environment)),
		),
	)
	if err != nil {
		return nil, err
	}

	return sdkresource.Merge(sdkresource.Default(), resource)
}

func NewSpanExporter(ctx context.Context, cfg config.Config) (sdktrace.SpanExporter, error) {
	if cfg.Otel.TraceGRCPAddr == "" {
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.Otel.TraceGRCPAddr),
	)
}

func NewTraceProvider(
	cfg config.Config,
	resource *sdkresource.Resource,
	spanExporter sdktrace.SpanExporter,
) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter,
			sdktrace.WithBatchTimeout(cfg.Otel.TraceFlushInterval),
		),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider
}

func NewMetricExporter(ctx context.Context, cfg config.Config) (sdkmetric.Exporter, error) {
	if cfg.Otel.MeterGRCPAddr == "" {
		return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	}

	return otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(cfg.Otel.MeterGRCPAddr),
	)
}

func NewMetricProvider(
	cfg config.Config,
	resource *sdkresource.Resource,
	exporter sdkmetric.Exporter,
) (*sdkmetric.MeterProvider, error) {
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(cfg.Otel.MetricFlushInterval),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}
