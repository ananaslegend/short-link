package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ananaslegend/short-link/internal/link/service"
)

var ErrEmptyAlias = errors.New("empty alias")

type LinkGetter interface {
	GetLinkByAlias(c context.Context, alias string) (string, error)
}

// RedirectHandler godoc
//
//	@Summary		LinkHandler to the original link
//	@Description	LinkHandler to the original link by alias
//	@Tags			redirect
//	@Accept			json
//	@Produce		json
//	@Param			alias	path		string	true	"Alias"
//	@Success		302		{string}	string	"LinkHandler to the original link"
//	@Failure		400		{object}	any
//	@Failure		404		{object}	any
//	@Router			/{alias} [get]
func (h LinkHandler) RedirectHandler(c echo.Context) error {
	alias, err := h.fetchAlias(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	link, err := h.linkGetter.GetLinkByAlias(c.Request().Context(), alias)
	if err != nil {
		return h.handleError(err)
	}

	return c.Redirect(http.StatusFound, link)
}

func (h LinkHandler) fetchAlias(c echo.Context) (string, error) {
	alias := c.Param("alias")

	if err := validateAlias(alias); err != nil {
		return "", err
	}

	return alias, nil
}

func validateAlias(alias string) error {
	if len(alias) == 0 {
		return ErrEmptyAlias
	}

	return nil
}

func (h LinkHandler) handleError(err error) error {
	switch {
	case errors.Is(err, service.ErrAliasNotFound):
		return echo.NewHTTPError(http.StatusNotFound, err)

	default:
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
}
