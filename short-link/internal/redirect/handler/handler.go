package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/ananaslegend/short-link/internal/redirect/service"
	"github.com/ananaslegend/short-link/pkg/clog"
)

var (
	ErrEmptyAlias = errors.New("empty alias")
)

type GetLinkService interface {
	GetLink(c context.Context, alias string) (string, error)
}

type Handler struct {
	linkService GetLinkService
}

func New(srv GetLinkService) *Handler {
	return &Handler{
		linkService: srv,
	}
}

// RedirectHandler godoc
//
//	@Summary		Redirect to the original link
//	@Description	Redirect to the original link by alias
//	@Tags			redirect
//	@Accept			json
//	@Produce		json
//	@Param			alias	path		string	true	"Alias"
//	@Success		302		{string}	string	"Redirect to the original link"
//	@Failure		400		{object}	any
//	@Failure		404		{object}	any
//	@Router			/{alias} [get]
func (h Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias, err := h.fetchAlias(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := clog.WithString(r.Context(), "alias", alias)

	link, err := h.linkService.GetLink(ctx, alias)
	if err != nil {
		h.renderError(ctx, w, err)
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}

func (h Handler) fetchAlias(ctx context.Context, r *http.Request) (string, error) {
	alias := r.PathValue("alias")

	if err := validateAlias(alias); err != nil {
		clog.Ctx(ctx).With("alias", alias).Error("error getting alias", clog.ErrorMsg(err))
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

func (h Handler) renderError(ctx context.Context, w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrAliasNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		clog.Ctx(ctx).Error("failed to get link", clog.ErrorMsg(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
