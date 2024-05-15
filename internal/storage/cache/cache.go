package cache

import (
	"context"
	"github.com/allegro/bigcache"
	"github.com/ananaslegend/short-link/internal/config"
)

type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

// Cache store data in memory. It is used to reduce the number of requests to the database.
type Cache struct {
	Cacher
}

// NewCache creates a new Cache. It returns an ErrNotImplementedCacheType if the cache type is not implemented.
// The cache type is specified in the configuration file.
func NewCache(cfg config.Cache) (*Cache, error) {
	cache := &Cache{}

	var err error

	switch cfg.CacheType {
	case config.BigCache:
		if cache.Cacher, err = NewBigCache(bigcache.DefaultConfig(cfg.TTL)); err != nil {
			return nil, err
		}
	default:
		return nil, ErrNotImplementedCacheType
	}

	return cache, nil
}
