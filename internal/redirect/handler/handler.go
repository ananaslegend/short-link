package handler

import (
	"context"
	"errors"
	"github.com/ananaslegend/short-link/internal/redirect/service"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
	"net/http"
	"strings"
)

type GetLinkService interface {
	GetLink(c context.Context, alias string) (string, error)
}

type Handler struct {
	linkService GetLinkService
	log         *slog.Logger
}

func NewHandler(srv GetLinkService, log *slog.Logger) *Handler {
	return &Handler{
		linkService: srv,
		log:         log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "internal.redirect.handler.Handler.ServeHTTP"
	log := h.log.With(slog.String("op", op))

	alias, err := getAliasFromUrlPath(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := h.linkService.GetLink(r.Context(), alias)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAliasNotFound):
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Error("failed to get link", logs.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, link, http.StatusFound)
}

func getAliasFromUrlPath(urlPath string) (string, error) {
	pathSegments := strings.Split(urlPath, "/")
	if len(pathSegments) < 2 {
		return "", ErrInvalidPathRequest
	}

	alias := pathSegments[1]
	if alias == "" {
		return "", ErrInvalidPathRequest
	}

	return alias, nil
}
