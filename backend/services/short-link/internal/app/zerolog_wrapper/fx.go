package zerolog_wrapper

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.NopLogger,

		fx.Provide(SetupZerolog),

		fx.Invoke(func(lc fx.Lifecycle, logger zerolog.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info().Msg("app starting")

					return nil
				},

				OnStop: func(ctx context.Context) error {
					logger.Info().Msg("app stopping")

					return nil
				},
			})
		}),
	)
}
