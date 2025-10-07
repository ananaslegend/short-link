package echowraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/app/config"
)

func NewEchoRouter(logger zerolog.Logger) *echo.Echo {
	echoRouter := echo.New()
	echoRouter.Logger = lecho.New(logger.With().Str("component", "echo_http_server").Logger())
	echoRouter.HideBanner = true

	return echoRouter
}

func NewEchoAPIGroup(router *echo.Echo) *echo.Group {
	return router.Group("/api")
}

func RunEchoServer(lc fx.Lifecycle, e *echo.Echo, cfg config.Config, logger zerolog.Logger) {
	lc.Append(fx.Hook{
		OnStart: StartEchoServer(e, cfg, logger),
		OnStop:  StopEchoServer(e, logger),
	})
}

func StartEchoServer(
	e *echo.Echo,
	cfg config.Config,
	logger zerolog.Logger,
) func(ctx context.Context) error {
	return func(_ context.Context) error {
		go func() {
			if err := e.Start(fmt.Sprintf(":%v", cfg.HttpServer.Port)); err != nil &&
				!errors.Is(err, http.ErrServerClosed) {
				logger.Error().Err(err).Msg("http server error")
			}
		}()

		return nil
	}
}

func StopEchoServer(e *echo.Echo, logger zerolog.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return e.Shutdown(ctx)
	}
}
