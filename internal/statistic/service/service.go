package service

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Statistic struct {
	redirectHandler RedirectHandler

	tracer trace.Tracer
}

func NewStatistic(redirectHandler RedirectHandler, traceProvider *sdktrace.TracerProvider) *Statistic {
	return &Statistic{redirectHandler: redirectHandler, tracer: traceProvider.Tracer("internal.statistic.service.Statistic")}
}
