package pgx_wrapper

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/app/config"
)

func Module() fx.Option {
	return fx.Module(
		"short-link.internal.app.postgres",

		fx.Provide(func(lc fx.Lifecycle, cfg config.Config, logger zerolog.Logger) *Wrapper {
			wrapper := &Wrapper{}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					pgpool, err := pgxpool.New(ctx, cfg.DbConn)
					if err != nil {
						logger.Error().Err(err).Msg("failed to connect to db")

						return err
					}

					if err = pgpool.Ping(ctx); err != nil {
						logger.Error().Err(err).Msg("failed to ping db")

						return err
					}

					wrapper.Pool = pgpool

					return nil
				},
				OnStop: func(ctx context.Context) error {
					if wrapper.Pool != nil {
						wrapper.Pool.Close()
					}

					return nil
				},
			})

			return wrapper
		}),
	)
}
