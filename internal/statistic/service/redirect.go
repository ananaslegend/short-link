package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/ananaslegend/short-link/internal/statistic/domain"
)

type RedirectHandler interface {
	AddRedirectEvent(ctx context.Context, redirectEvent domain.RedirectEventStatistic) error
}

func (s Statistic) AddRedirectEvent(
	ctx context.Context,
	redirectEvent domain.RedirectEventStatistic,
) error {
	const op = "short-link.internal.statistic.service.Statistic.AddRedirectEvent"

	if err := s.redirectHandler.AddRedirectEvent(ctx, redirectEvent); err != nil {
		zerolog.Ctx(ctx).Error().Str("op", op).Err(err).Msg("error adding redirect event")

		return fmt.Errorf("%v: %w", op, err)
	}

	return nil
}
