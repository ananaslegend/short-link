package handlers

import (
	"context"
	"errors"
	"github.com/ananaslegend/short-link/pkg/errs"
	"github.com/ananaslegend/short-link/pkg/logs"
	"log/slog"
	"net/http"
	"strings"
)

type LinkGetter interface {
	GetLink(c context.Context, alias string) (string, error)
}

func Redirect(w http.ResponseWriter, r *http.Request, log *slog.Logger, lg LinkGetter) {
	const op = "api.handlers.Redirect"
	log.With(slog.String("op", op))

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alias := segments[1]

	link, err := lg.GetLink(r.Context(), alias)
	if err != nil {
		if errors.Is(err, errs.ErrAliasNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("failed to get link", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, link, http.StatusFound)
}
