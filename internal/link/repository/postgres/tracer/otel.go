package tracer

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ananaslegend/short-link/internal/link/domain"
)

type Base interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
	InsertAliasedLink(ctx context.Context, dto domain.InsertAliasedLink) (domain.AliasedLink, error)
}
type OtelDecorator struct {
	tracer trace.Tracer

	base Base
}

func NewOtelDecorator(base Base, provider *sdktrace.TracerProvider) *OtelDecorator {
	return &OtelDecorator{
		base:   base,
		tracer: provider.Tracer("internal.link.repository.postgres.LinkRepository"),
	}
}

func (o OtelDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	ctx, span := o.tracer.Start(
		ctx,
		"internal.link.repository.postgres.LinkRepository.GetLinkByAlias",
	)
	defer span.End()

	link, err := o.base.GetLinkByAlias(ctx, alias)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return "", err
	}

	return link, nil
}

func (o OtelDecorator) InsertAliasedLink(
	ctx context.Context,
	dto domain.InsertAliasedLink,
) (domain.AliasedLink, error) {
	ctx, span := o.tracer.Start(
		ctx,
		"internal.link.repository.postgres.LinkRepository.InsertAliasedLink",
	)
	defer span.End()

	res, err := o.base.InsertAliasedLink(ctx, dto)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return domain.AliasedLink{}, err
	}

	return res, nil
}
