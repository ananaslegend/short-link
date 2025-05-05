package clickhouse_wrapper

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/app/config"
)

func Module() fx.Option {
	return fx.Module(
		"short-link.internal.app.clickhouse",
		fx.Provide(func(lc fx.Lifecycle, cfg config.Config, logger zerolog.Logger) driver.Conn {
			conn, err := clickhouse.Open(&clickhouse.Options{
				Addr: []string{cfg.ClickHouse.Host},
				Auth: clickhouse.Auth{
					Username: cfg.ClickHouse.User,
					Password: cfg.ClickHouse.Pass,
					Database: cfg.ClickHouse.Db,
				},
				Compression: &clickhouse.Compression{
					Method: clickhouse.CompressionLZ4,
				},
				Debug: cfg.Environment == config.Local,
			})
			if err != nil {
				logger.Fatal().Err(err).Msg("failed to connect to ClickHouse")
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return conn.Ping(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return conn.Close()
				},
			})

			return conn
		}),
	)
}
