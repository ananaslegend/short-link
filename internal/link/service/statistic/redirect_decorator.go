package statistic

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/ananaslegend/short-link/internal/link/service"
	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

type RedirectStatisticProvider interface {
	AddRedirectEvent(ctx context.Context, redirectEvent domain.RedirectEventStatistic) error
}
type RedirectDecorator struct {
	base *service.Link

	redirectStatisticProvider RedirectStatisticProvider

	tracer trace.Tracer
}

func NewRedirectDecorator(
	base *service.Link,
	redirectProvider RedirectStatisticProvider,
	traceProvider *sdktrace.TracerProvider,
) *RedirectDecorator {
	return &RedirectDecorator{
		base:                      base,
		redirectStatisticProvider: redirectProvider,
		tracer: traceProvider.Tracer(
			"internal.link.service.statistic.RedirectStatisticProvider",
		),
	}
}

func (d RedirectDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "internal.link.service.statistic.RedirectDecorator.GetLinkByAlias"

	ctx, span := d.tracer.Start(ctx, op)
	defer span.End()

	link, err := d.base.GetLinkByAlias(ctx, alias)
	if err != nil {
		zerolog.Ctx(ctx).
			Error().
			Err(err).
			Str("op", op).
			Str("alias", alias).
			Msg("failed to get link by alias")

		return "", err
	}

	event := domain.RedirectEventStatistic{
		Alias: alias,
		Link:  link,
	}

	if err = d.redirectStatisticProvider.AddRedirectEvent(ctx, event); err != nil {
		zerolog.Ctx(ctx).
			Error().
			Err(err).
			Str("op", op).
			Any("event", event).
			Msg("failed to add redirect event")

		return link, nil
	}

	return link, nil
}
