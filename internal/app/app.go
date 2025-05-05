package app

import (
	"go.uber.org/fx"

	aliasgenerator "github.com/ananaslegend/short-link/internal/alias_generator"
	automaxprocswrapper "github.com/ananaslegend/short-link/internal/app/automaxprocs_wrapper"
	clickhousewrapper "github.com/ananaslegend/short-link/internal/app/clickhouse_wrapper"
	"github.com/ananaslegend/short-link/internal/app/config"
	echowrapper "github.com/ananaslegend/short-link/internal/app/echo_wrapper"
	pgxwrapper "github.com/ananaslegend/short-link/internal/app/pgx_wrapper"
	rediswrapper "github.com/ananaslegend/short-link/internal/app/redis_wrapper"
	zerologwrapper "github.com/ananaslegend/short-link/internal/app/zerolog_wrapper"
	"github.com/ananaslegend/short-link/internal/link"
	"github.com/ananaslegend/short-link/internal/statistic"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(config.MustLoadConfig),

		fx.Provide(zerologwrapper.SetupZerolog),
		fx.WithLogger(zerologwrapper.WithZerologFx),

		automaxprocswrapper.Module(),

		pgxwrapper.Module(),
		rediswrapper.Module(),
		clickhousewrapper.Module(),

		echowrapper.Module(),

		statistic.Module(),

		aliasgenerator.Module(),
		link.Module(),
	)
}
