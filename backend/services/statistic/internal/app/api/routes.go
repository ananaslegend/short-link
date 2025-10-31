package api

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"internal.api.routes",
		// TODO: Add statistic handlers here later
	)
}
