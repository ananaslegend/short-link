package statistic

import (
	"go.uber.org/fx"

	"github.com/ananaslegend/statistic/internal/statistic/repository/clickhouse"
	repoTracer "github.com/ananaslegend/statistic/internal/statistic/repository/clickhouse/tracer"
	"github.com/ananaslegend/statistic/internal/statistic/service"
	serviceTracer "github.com/ananaslegend/statistic/internal/statistic/service/tracer"
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
			serviceTracer.NewOtelDecorator,
		),
	)
}
