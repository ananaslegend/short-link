package redis_wrapper

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/app/config"
)

func Module() fx.Option {
	return fx.Module(
		"internal.app.redis",
		fx.Provide(
			func(lc fx.Lifecycle, cfg config.Config, logger zerolog.Logger) *redis.Client {
				client := redis.NewClient(&redis.Options{
					Addr:     cfg.Redis.Addr,
					Password: cfg.Redis.Password,
				})

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						res := client.Ping(ctx)
						if err := res.Err(); err != nil {
							logger.Error().Err(err).Msg("failed to ping redis")

							return fmt.Errorf("failed to ping redis: %w", err)
						}

						return nil
					},
					OnStop: func(ctx context.Context) error {
						if err := client.Close(); err != nil {
							logger.Error().Err(err).Msg("failed to close redis client")

							return fmt.Errorf("failed to close redis client: %w", err)
						}

						return nil
					},
				})

				return client
			},
		),
	)
}
