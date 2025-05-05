package alias_generator

import (
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/alias_generator/service"
	linkService "github.com/ananaslegend/short-link/internal/link/service"
)

func Module() fx.Option {
	return fx.Module(
		"short-link.internal.alias_generator",

		fx.Provide(
			fx.Annotate(service.NewUUIDGenerated, fx.As(new(linkService.AliasGenerator))),
		),
	)
}
