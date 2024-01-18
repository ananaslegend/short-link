package cache

import (
	"context"
	"fmt"
	"github.com/ananaslegend/short-link/internal/redirect/repository"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

type SelectRepository interface {
	SelectLink(ctx context.Context, alias string) (string, error)
}

type CachedRepository struct {
	repo  SelectRepository
	cache Cache
}

func NewCachedRepository(repo SelectRepository, cache Cache) *CachedRepository {
	return &CachedRepository{
		repo:  repo,
		cache: cache,
	}
}

func (cr CachedRepository) SelectLink(ctx context.Context, alias string) (string, error) {
	const op = "internal.redirect.repository.cache.CachedRepository.SelectLink"

	link, err := cr.cache.Get(ctx, alias)
	if err == nil {
		return link, nil
	}

	link, err = cr.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = cr.cache.Set(ctx, alias, link); err != nil {
		return link, fmt.Errorf("%s: %w, link: %s ; err: %w", op, repository.ErrCantSetToCache, link, err)
	}

	return link, nil
}
