package remove

import (
	"context"
	"golang.org/x/exp/slog"
	"net/http"
)

type LinkByAliasRemover interface {
	RemoveLinkByAlias(alias string) error
}

func Handle(c context.Context, w http.ResponseWriter, r *http.Request, log *slog.Logger, remover LinkByAliasRemover) {
	const op = "api.handlers.remove.Handle"

	alias := r.URL.Query().Get("alias")

	if err := remover.RemoveLinkByAlias(alias); err != nil {

	}

}
