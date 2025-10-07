package statistic

import (
	"go.uber.org/fx"

	linkService "github.com/ananaslegend/short-link/internal/link/service/statistic"
	"github.com/ananaslegend/short-link/internal/statistic/repository/clickhouse"
	repoTracer "github.com/ananaslegend/short-link/internal/statistic/repository/clickhouse/tracer"
	"github.com/ananaslegend/short-link/internal/statistic/service"
	serviceTracer "github.com/ananaslegend/short-link/internal/statistic/service/tracer"
)

func Module() fx.Option {
	return fx.Module(
		"internal.statistic",

		fx.Provide(
			fx.Annotate(clickhouse.NewRepository, fx.As(new(repoTracer.BaseService))),
			fx.Annotate(repoTracer.NewOtelDecorator, fx.As(new(service.RedirectHandler))),
		),

		fx.Provide(
			fx.Annotate(service.NewStatistic, fx.As(new(serviceTracer.BaseService))),
			fx.Annotate(
				serviceTracer.NewOtelDecorator,
				fx.As(new(linkService.RedirectStatisticProvider)),
			),
		),
	)
}
