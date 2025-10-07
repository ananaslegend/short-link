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
		"internal.app.postgres",

		fx.Provide(
			func(ctx context.Context, lc fx.Lifecycle, cfg config.Config, logger zerolog.Logger) *pgxpool.Pool {
				pgpool, err := pgxpool.New(ctx, cfg.DbConn)
				if err != nil {
					logger.Fatal().Err(err).Msg("failed to connect to db")
				}

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return pgpool.Ping(ctx)
					},
					OnStop: func(ctx context.Context) error {
						pgpool.Close()

						return nil
					},
				})

				return pgpool
			},
		),
	)
}
