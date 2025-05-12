package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type BaseService interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}

type OtelDecorator struct {
	tracer trace.Tracer

	BaseService
}

func NewOtelDecorator(provider *sdktrace.TracerProvider, baseSrv BaseService) *OtelDecorator {
	return &OtelDecorator{
		tracer:      provider.Tracer("internal.link.service.statistic.RedirectDecorator"),
		BaseService: baseSrv,
	}
}

func (o OtelDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	ctx, span := o.tracer.Start(
		ctx,
		"internal.link.service.statistic.RedirectDecorator.GetLinkByAlias",
	)
	defer span.End()

	link, err := o.BaseService.GetLinkByAlias(ctx, alias)
	if err != nil {
		span.RecordError(err)

		return "", err
	}

	return link, nil
}
