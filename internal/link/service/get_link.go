package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type LinkGetter interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}

func (s Link) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "internal.link.service.Link.GetLinkByAlias"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	link, err := s.linkGetter.GetLinkByAlias(ctx, alias)
	if err != nil {
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
