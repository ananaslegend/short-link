package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/ananaslegend/short-link/internal/link/domain"
	"github.com/ananaslegend/short-link/internal/link/repository/postgres"
	statisticDomain "github.com/ananaslegend/short-link/internal/statistic/domain"
)

type LinkGetter interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}

type RedirectStatisticProvider interface {
	AddRedirectEvent(
		ctx context.Context,
		redirectEvent statisticDomain.RedirectEventStatistic,
	) error
}

func (s Link) GetLinkByAlias(ctx context.Context, dto domain.GetLinkDTO) (string, error) {
	const op = "internal.link.service.Link.GetLinkByAlias"

	link, err := s.linkGetter.GetLinkByAlias(ctx, dto.Alias)
	if err != nil {
		if errors.Is(err, postgres.ErrAliasNotFound) {
			return "", ErrAliasNotFound
		}

		zerolog.Ctx(ctx).
			Error().
			Str("op", op).
			Err(err).
			Str("alias", dto.Alias).
			Msg("failed to get link by alias")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.handleStatistic(ctx, dto, link)

	return link, nil
}

func (s Link) handleStatistic(ctx context.Context, dto domain.GetLinkDTO, link string) {
	if !dto.SendStatistic {
		return
	}

	const op = "internal.link.service.Link.sendRedirectStatistic"

	ctx = context.WithoutCancel(ctx)

	go func() {
		defer recoverRedirectStatisticSender(ctx, op)

		event := statisticDomain.RedirectEventStatistic{
			Alias: dto.Alias,
			Link:  link,
		}

		if err := s.redirectStatisticProvider.AddRedirectEvent(ctx, event); err != nil {
			zerolog.Ctx(ctx).
				Error().
				Err(err).
				Str("op", op).
				Any("event", event).
				Msg("failed to add redirect event")
		}
	}()
}

func recoverRedirectStatisticSender(ctx context.Context, op string) {
	if r := recover(); r != nil {
		zerolog.Ctx(ctx).
			Error().
			Str("op", op).
			Any("recover", r).
			Msg("panic recovered in redirect statistic sender")
	}
}
