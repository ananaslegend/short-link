package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Base interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}
type OtelDecorator struct {
	tracer trace.Tracer

	base Base
}

func NewOtelDecorator(base Base, provider *sdktrace.TracerProvider) *OtelDecorator {
	return &OtelDecorator{
		base:   base,
		tracer: provider.Tracer("internal.link.repository.redis.LinkRepositoryDecorator"),
	}
}

func (o OtelDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	ctx, span := o.tracer.Start(
		ctx,
		"internal.link.repository.redis.LinkRepositoryDecorator.GetLinkByAlias",
	)
	defer span.End()

	link, err := o.base.GetLinkByAlias(ctx, alias)
	if err != nil {
		span.RecordError(err)

		return "", err
	}

	return link, nil
}
