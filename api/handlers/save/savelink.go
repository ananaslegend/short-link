package save

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/storage"
	"github.com/ananaslegend/short-link/usecases"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type LinkSaver interface {
	AddLink(link, alias string) error
}

type Request struct {
	Link  string `json:"link"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Alias string `json:"alias,omitempty"`
	Error string `json:"error,omitempty"`
}

func Handle(c context.Context, log *slog.Logger, ls LinkSaver) func(http.ResponseWriter, *http.Request) {
	const op = "api.handlers.save.link.Handle"

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			err = fmt.Errorf("%s:%w", op, err)
			log.Error("cant read body", logs.Err(err))
			return
		}
		model := Request{}
		if err = json.Unmarshal(b, &model); err != nil {
			err = fmt.Errorf("%s:%w", op, err)
			log.Error("cant unmarshal request body", logs.Err(err))
			return
		}

		addedAlias, err := usecases.AddLink(c, log, ls, model.Link, model.Alias)
		if err != nil {
			if errors.Is(err, usecases.ErrAutoAliasAlreadyExists) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if errors.Is(err, storage.ErrAliasExists) {
				w.WriteHeader(http.StatusConflict)
				w.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(w).Encode(map[string]string{
					"alias": addedAlias,
				})
				if err != nil {
					log.Error("cant encode json", logs.Err(err))
				}
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated) // json
		return
	}
}
