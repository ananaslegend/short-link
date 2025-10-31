package app

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/fx"

	"github.com/ananaslegend/statistic/internal/app/api"
	automaxprocswrapper "github.com/ananaslegend/statistic/internal/app/automaxprocs_wrapper"
	clickhousewrapper "github.com/ananaslegend/statistic/internal/app/clickhouse_wrapper"
	"github.com/ananaslegend/statistic/internal/app/config"
	echowrapper "github.com/ananaslegend/statistic/internal/app/echo_wrapper"
	otelwrapper "github.com/ananaslegend/statistic/internal/app/otel_wrapper"
	zerologwrapper "github.com/ananaslegend/statistic/internal/app/zerolog_wrapper"
	"github.com/ananaslegend/statistic/internal/statistic"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(func(lc fx.Lifecycle) context.Context {
			ctx, cancel := signal.NotifyContext(
				context.Background(),
				syscall.SIGINT,
				syscall.SIGTERM,
			)

			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					cancel()

					return nil
				},
			})

			return ctx
		}),

		fx.Provide(config.MustLoadConfig),

		zerologwrapper.Module(),

		automaxprocswrapper.Module(),

		otelwrapper.Module(),

		clickhousewrapper.Module(),

		echowrapper.Module(),

		statistic.Module(),

		api.Module(),
	)
}
