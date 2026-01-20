package statistic

import (
	"go.uber.org/fx"

	linkService "github.com/ananaslegend/short-link/internal/link/service"
	natsrepo "github.com/ananaslegend/short-link/internal/statistic/repository/nats"
	"github.com/ananaslegend/short-link/internal/statistic/service"
	serviceTracer "github.com/ananaslegend/short-link/internal/statistic/service/tracer"
)

func Module() fx.Option {
	return fx.Module(
		"internal.statistic",

		fx.Provide(
			fx.Annotate(natsrepo.NewRepository, fx.As(new(service.RedirectHandler))),
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
