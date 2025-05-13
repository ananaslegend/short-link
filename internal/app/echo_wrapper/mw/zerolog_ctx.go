package mw

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func ZerologContextMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request().WithContext(logger.WithContext(c.Request().Context()))

			c.SetRequest(req)

			return next(c)
		}
	}
}
