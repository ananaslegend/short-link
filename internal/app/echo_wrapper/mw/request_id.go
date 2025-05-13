package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func SetupEchoRequestIDConfig() middleware.RequestIDConfig {
	return middleware.RequestIDConfig{
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
	}
}
