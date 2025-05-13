package tracer

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ananaslegend/short-link/internal/link/domain"
)

type BaseService interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
	InsertLink(ctx context.Context, dto domain.InsertLink) (domain.AliasedLink, error)
}

type OtelDecorator struct {
	tracer trace.Tracer

	BaseService
}

func NewOtelDecorator(provider *sdktrace.TracerProvider, baseSrv BaseService) *OtelDecorator {
	return &OtelDecorator{
		tracer:      provider.Tracer("internal.link.service.Link"),
		BaseService: baseSrv,
	}
}

func (o OtelDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	ctx, span := o.tracer.Start(ctx, "internal.link.service.Link.GetLinkByAlias")
	defer span.End()

	link, err := o.BaseService.GetLinkByAlias(ctx, alias)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return "", err
	}

	return link, nil
}

func (o OtelDecorator) InsertLink(
	ctx context.Context,
	dto domain.InsertLink,
) (domain.AliasedLink, error) {
	ctx, span := o.tracer.Start(ctx, "internal.link.service.Link.InsertLink")
	defer span.End()

	insertedLink, err := o.BaseService.InsertLink(ctx, dto)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return domain.AliasedLink{}, err
	}

	return insertedLink, nil
}
