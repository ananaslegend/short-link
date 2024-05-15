package cache

import (
	"context"
	"github.com/allegro/bigcache"
)

type BigCache struct {
	BigCache *bigcache.BigCache
}

func NewBigCache(config bigcache.Config) (*BigCache, error) {
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	return &BigCache{BigCache: cache}, nil
}

func (bc *BigCache) Get(ctx context.Context, key string) (string, error) {
	val, err := bc.BigCache.Get(key)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (bc *BigCache) Set(ctx context.Context, key, value string) error {
	return bc.BigCache.Set(key, []byte(value))
}
