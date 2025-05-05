package automaxprocs_wrapper

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/fx"

	zerologwrapper "github.com/ananaslegend/short-link/internal/app/zerolog_wrapper"
)

func Module() fx.Option {
	return fx.Module(
		"short-link.internal.app.automaxprocs",

		fx.Invoke(
			func(lc fx.Lifecycle, logger zerolog.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						if _, err := maxprocs.Set(maxprocs.Logger(zerologwrapper.ZerologMaxProcsLogger(logger))); err != nil {
							logger.Error().Err(err).Msg("failed to set GOMAXPROCS")
						}

						return nil
					},
				})
			},
		),
	)
}
