package statistic

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/ananaslegend/short-link/internal/link/service"
	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

type RedirectStatisticProvider interface {
	AddRedirectEvent(ctx context.Context, redirectEvent domain.RedirectEventStatistic) error
}
type RedirectDecorator struct {
	base *service.Link

	redirectStatisticProvider RedirectStatisticProvider
}

func NewRedirectDecorator(
	base *service.Link,
	redirectProvider RedirectStatisticProvider,
) *RedirectDecorator {
	return &RedirectDecorator{base: base, redirectStatisticProvider: redirectProvider}
}

func (d RedirectDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "short-link.internal.link.service.statistic.RedirectDecorator.GetLinkByAlias"

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
