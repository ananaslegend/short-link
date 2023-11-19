package redirect

import (
	"context"
	"errors"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
	"net/http"
	"strings"
)

type GetLinkUseCase interface {
	GetLink(c context.Context, alias string) (string, error)
}

type Handler struct {
	useCase GetLinkUseCase
	log     *slog.Logger
}

func NewHandler(uc GetLinkUseCase, log *slog.Logger) *Handler {
	return &Handler{
		useCase: uc,
		log:     log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "api.handlers.ServeHTTP"
	log := h.log.With(slog.String("op", op))

	path := r.URL.Path
	pathSegments := strings.Split(path, "/")
	if len(pathSegments) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alias := pathSegments[1]

	link, err := h.useCase.GetLink(r.Context(), alias)
	if err != nil {
		switch {
		case errors.Is(err, ErrAliasNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		case errors.Is(err, ErrCantSetToCache):
			http.Redirect(w, r, link, http.StatusFound)
			return
		default:
			log.Error("failed to get link", logs.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, link, http.StatusFound)
}
