package redirect

import (
	"context"
	"fmt"
	"log"
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
	const op = "storage.cached.SelectLink"

	link, err := cr.cache.Get(ctx, alias)
	if err == nil {
		log.Printf("cache hit: %s", alias)

		return link, nil
	}

	log.Printf("cache miss: %s", alias)

	link, err = cr.repo.SelectLink(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = cr.cache.Set(ctx, alias, link); err != nil {
		return link, fmt.Errorf("%s: %w, link: %s ; err: %w", op, ErrCantSetToCache, link, err)
	}

	log.Printf("cache set: %s", alias)

	return link, nil
}
