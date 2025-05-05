package echowraper

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func SetupMiddleware(router *echo.Echo, logger zerolog.Logger) {
	router.Use(middleware.Recover())

	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().WithContext(logger.WithContext(c.Request().Context()))

			c.SetRequest(req)

			return next(c)
		}
	})

	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogProtocol: true,
		LogStatus:   true,
		LogMethod:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			zerolog.Ctx(c.Request().Context()).Info().
				Str("http_method", v.Method).
				Str("uri", v.URI).
				Str("protocol", v.Protocol).
				Int("status", v.Status).
				Send()

			return nil
		},
	}))

	router.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, s string) {
			c.SetRequest(
				c.Request().WithContext(
					zerolog.Ctx(c.Request().Context()).
						With().
						Str("req_id", s).
						Logger().
						WithContext(c.Request().Context()),
				),
			)
		},
	}))
}
