package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/ananaslegend/short-link/internal/link/repository/postgres"
)

type LinkGetter interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}

func (s Link) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "internal.link.service.Link.GetLinkByAlias"

	link, err := s.linkGetter.GetLinkByAlias(ctx, alias)
	if err != nil {
		if errors.Is(err, postgres.ErrAliasNotFound) {
			return "", ErrAliasNotFound
		}

		zerolog.Ctx(ctx).
			Error().
			Str("op", op).
			Err(err).
			Str("alias", alias).
			Msg("failed to get link by alias")

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}
