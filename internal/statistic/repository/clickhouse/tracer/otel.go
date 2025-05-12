package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

type BaseService interface {
	AddRedirectEvent(ctx context.Context, redirectEvent domain.RedirectEventStatistic) error
}

type OtelDecorator struct {
	tracer trace.Tracer

	BaseService
}

func NewOtelDecorator(provider *sdktrace.TracerProvider, baseSrv BaseService) *OtelDecorator {
	return &OtelDecorator{
		tracer: provider.Tracer(
			"internal.statistic.repository.clickhouse.StatisticRepository",
		),
		BaseService: baseSrv,
	}
}

func (o OtelDecorator) AddRedirectEvent(
	ctx context.Context,
	redirectEvent domain.RedirectEventStatistic,
) error {
	ctx, span := o.tracer.Start(
		ctx,
		"internal.statistic.repository.clickhouse.StatisticRepository.AddRedirectEvent",
	)
	defer span.End()

	err := o.BaseService.AddRedirectEvent(ctx, redirectEvent)
	if err != nil {
		span.RecordError(err)

		return err
	}

	return nil
}
