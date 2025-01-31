package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/ananaslegend/short-link/pkg/cslog"
)

type SelectRepository interface {
	SelectLink(ctx context.Context, alias string) (string, error)
}

type Redis struct {
	repo    SelectRepository
	client  *redis.Client
	timeout time.Duration
}

func NewRedisCache(repo SelectRepository, client *redis.Client) *Redis {
	return &Redis{
		repo:   repo,
		client: client,
	}
}

func (r Redis) SelectLink(ctx context.Context, alias string) (string, error) {
	link, err := r.client.Get(ctx, alias).Result()
	if err == nil {
		cslog.Logger(ctx).Info("got link from cache", slog.String("link", link))
		return link, nil
	}

	link, err = r.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", err
	}

	if err = r.client.Set(ctx, alias, link, r.timeout).Err(); err != nil {
		cslog.Logger(ctx).Error("cant set link to cache", cslog.Error(err))
	}

	return link, nil
}
