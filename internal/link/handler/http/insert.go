package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ananaslegend/short-link/internal/link/domain"
)

type LinkInserter interface {
	InsertLink(ctx context.Context, dto domain.InsertLink) (domain.AliasedLink, error)
}

func (h Link) SaveLink(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := validateInsertRequest(c, &req); err != nil {
		return err
	}

	insertLink, err := h.linkInserter.InsertLink(c.Request().Context(), mapInsertRequest(req))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, mapSuccessResponse(insertLink))
}

func validateInsertRequest(c echo.Context, req *Request) error {
	if !strings.Contains(req.Link, "http") {
		req.Link = "https://" + req.Link
	}

	return c.Validate(req)
}

type Request struct {
	Link  string  `json:"link"  validate:"required"`
	Alias *string `json:"alias"`
}

func mapInsertRequest(r Request) domain.InsertLink {
	return domain.InsertLink{
		Link:  r.Link,
		Alias: r.Alias,
	}
}

type SuccessResponse struct {
	Link  string `json:"link"`
	Alias string `json:"alias"`
}

func mapSuccessResponse(r domain.AliasedLink) SuccessResponse {
	return SuccessResponse{
		Link:  r.Link,
		Alias: r.Alias,
	}
}
