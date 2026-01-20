package link

import (
	"time"

	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/link/handler/http"
	"github.com/ananaslegend/short-link/internal/link/repository/postgres"
	postgresTracer "github.com/ananaslegend/short-link/internal/link/repository/postgres/tracer"
	"github.com/ananaslegend/short-link/internal/link/repository/redis"
	redisTracer "github.com/ananaslegend/short-link/internal/link/repository/redis/tracer"
	"github.com/ananaslegend/short-link/internal/link/service"
	serviceTracer "github.com/ananaslegend/short-link/internal/link/service/tracer"
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
			fx.Annotate(postgres.NewLinkRepository, fx.As(new(postgresTracer.Base))),
			fx.Annotate(postgresTracer.NewOtelDecorator, fx.As(new(redis.BaseRepository))),

			fx.Annotate(
				redis.NewLinkRepositoryDecorator,
				fx.As(new(redisTracer.Base)), fx.ParamTags("", "", `name:"get_link_cache_ttl"`),
			),
			fx.Annotate(redisTracer.NewOtelDecorator, fx.As(new(service.LinkGetter))),

			fx.Annotate(postgresTracer.NewOtelDecorator, fx.As(new(service.AliasedLinkInserter))),
		),

		fx.Provide(
			fx.Annotate(service.New, fx.As(new(serviceTracer.BaseService))),
			fx.Annotate(serviceTracer.NewOtelDecorator, fx.As(new(http.LinkGetter))),
		),

		fx.Provide(
			fx.Annotate(serviceTracer.NewOtelDecorator, fx.As(new(http.LinkInserter))),
		),

		fx.Provide(http.NewHandler),
	)
}
