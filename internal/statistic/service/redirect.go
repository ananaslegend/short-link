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
	const op = "internal.statistic.service.Statistic.AddRedirectEvent"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	if err := s.redirectHandler.AddRedirectEvent(ctx, redirectEvent); err != nil {
		zerolog.Ctx(ctx).Error().Str("op", op).Err(err).Msg("error adding redirect event")

		return fmt.Errorf("%v: %w", op, err)
	}

	return nil
}
