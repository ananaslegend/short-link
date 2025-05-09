package link

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/link/handler/http"
	"github.com/ananaslegend/short-link/internal/link/repository/postgres"
	"github.com/ananaslegend/short-link/internal/link/repository/redis"
	"github.com/ananaslegend/short-link/internal/link/service"
	"github.com/ananaslegend/short-link/internal/link/service/statistic"
)

const (
	getLinkCacheTTL = 5 * time.Minute
)

func Module() fx.Option {
	return fx.Module(
		"internal.link",

		fx.Provide(
			fx.Annotate(
				func() time.Duration { return getLinkCacheTTL },
				fx.ResultTags(`name:"get_link_cache_ttl"`),
			),
			fx.Annotate(
				postgres.NewLinkRepository, fx.As(new(redis.BaseRepository)),
			),
			fx.Annotate(
				redis.NewLinkRepositoryDecorator,
				fx.As(new(service.LinkGetter)), fx.ParamTags("", "", `name:"get_link_cache_ttl"`),
			),

			fx.Annotate(postgres.NewLinkRepository, fx.As(new(service.AliasedLinkInserter))),
		),

		fx.Provide(
			service.New,
			fx.Annotate(statistic.NewRedirectDecorator, fx.As(new(http.LinkGetter))),

			fx.Annotate(service.New, fx.As(new(http.LinkInserter))),
		),

		fx.Provide(http.NewHandler),

		fx.Invoke(func(group *echo.Group, h *http.LinkHandler) {
			group.GET("/:alias", h.RedirectHandler)
			group.POST("", h.SaveLink)
		}),
	)
}
