package handlers

import (
	"context"
	"golang.org/x/exp/slog"
	"net/http"
)

type LinkSaver interface {
	AddLink(link, alias string) error
}

func Link(c context.Context, log *slog.Logger, ls LinkSaver) func(http.ResponseWriter, *http.Request) {
	const op = "api.handlers.Link"

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
