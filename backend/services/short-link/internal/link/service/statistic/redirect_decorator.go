package statistic

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

type BaseService interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}

type RedirectStatisticProvider interface {
	AddRedirectEvent(ctx context.Context, redirectEvent domain.RedirectEventStatistic) error
}
type RedirectDecorator struct {
	base BaseService

	redirectStatisticProvider RedirectStatisticProvider
}

func NewRedirectDecorator(
	base BaseService,
	redirectProvider RedirectStatisticProvider,
) *RedirectDecorator {
	return &RedirectDecorator{
		base:                      base,
		redirectStatisticProvider: redirectProvider,
	}
}

func (d RedirectDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "internal.link.service.statistic.RedirectDecorator.GetLinkByAlias"

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

	ctx = context.WithoutCancel(ctx)

	go func() {
		if err = d.redirectStatisticProvider.AddRedirectEvent(ctx, event); err != nil {
			zerolog.Ctx(ctx).
				Error().
				Err(err).
				Str("op", op).
				Any("event", event).
				Msg("failed to add redirect event")
		}
	}()

	return link, nil
}
