package redirect

import (
	"context"
	"errors"
	"github.com/ananaslegend/short-link/errs"
	"github.com/ananaslegend/short-link/logs"
	"golang.org/x/exp/slog"
	"net/http"
	"strings"
)

type LinkGetter interface {
	GetLink(alias string) (string, error)
}

func Handle(c context.Context, w http.ResponseWriter, r *http.Request, log *slog.Logger, lg LinkGetter) {
	const op = "api.handlers.redirect.Handle"
	log.With(slog.String("op", op))

	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	alias := segments[1]

	link, err := lg.GetLink(alias)
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
