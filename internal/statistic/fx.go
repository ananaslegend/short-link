package statistic

import (
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/link/service/statistic"
	"github.com/ananaslegend/short-link/internal/statistic/repository/clickhouse"
	"github.com/ananaslegend/short-link/internal/statistic/service"
)

func Module() fx.Option {
	return fx.Module(
		"short-link.internal.statistic",

		fx.Provide(
			fx.Annotate(clickhouse.NewRepository, fx.As(new(service.RedirectHandler))),
		),

		fx.Provide(
			fx.Annotate(service.NewStatistic, fx.As(new(statistic.RedirectStatisticProvider))),
		),
	)
}
