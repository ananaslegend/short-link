package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/ananaslegend/short-link/internal/link/handler/http"
)

func Module() fx.Option {
	return fx.Module(
		"internal.api.routes",

		fx.Invoke(func(group *echo.Group, h *http.LinkHandler) {
			group.GET("/:alias", h.RedirectHandler)
			group.POST("", h.SaveLink)
		}),
	)
}
