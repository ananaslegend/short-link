package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func SetupEchoRequestLoggerConfig() middleware.RequestLoggerConfig {
	return middleware.RequestLoggerConfig{
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
	}
}
