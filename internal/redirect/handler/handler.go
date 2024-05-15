package handler

import (
	"context"
	"errors"
	"github.com/ananaslegend/go-logs/v2"
	"github.com/ananaslegend/short-link/internal/redirect/service"
	"log/slog"
	"net/http"
)

var (
	ErrEmptyAlias = errors.New("empty alias")
)

type GetLinkService interface {
	GetLink(c context.Context, alias string) (string, error)
}

type Handler struct {
	linkService GetLinkService
	logger      *slog.Logger
}

func New(srv GetLinkService, log *slog.Logger) *Handler {
	return &Handler{
		linkService: srv,
		logger:      log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := logs.WithMetric(r.Context(), "handler", "redirect")

	alias, err := h.fetchAlias(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx = logs.WithMetric(ctx, "alias", alias)

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
		h.logger.DebugContext(ctx, "error getting alias", logs.ErrorMsg(err))
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
		h.logger.ErrorContext(logs.ErrorCtx(ctx, err), "failed to get link", logs.ErrorMsg(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
