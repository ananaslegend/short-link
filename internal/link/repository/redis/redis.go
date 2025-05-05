package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type BaseRepository interface {
	GetLinkByAlias(ctx context.Context, alias string) (string, error)
}

type RepositoryDecorator struct {
	base   BaseRepository
	client *redis.Client
	ttl    time.Duration
}

func NewRepositoryDecorator(
	repo BaseRepository,
	client *redis.Client,
	ttl time.Duration,
) *RepositoryDecorator {
	return &RepositoryDecorator{
		base:   repo,
		client: client,
		ttl:    ttl,
	}
}

func (r RepositoryDecorator) GetLinkByAlias(ctx context.Context, alias string) (string, error) {
	const op = "short-link.internal.link.repository.redis.redis.RepositoryDecorator.GetLinkByAlias"

	link, err := r.client.Get(ctx, alias).Result()
	if err == nil {
		zerolog.Ctx(ctx).Info().Str("alias", alias).Str("link", link).Msg("success")

		return link, nil
	}

	link, err = r.base.GetLinkByAlias(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("%v: %w", op, err)
	}

	if err = r.client.Set(ctx, alias, link, r.ttl).Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("alias", alias).Msg("failed to set link alias")
	}

	return link, nil
}
