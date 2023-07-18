package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/core"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/storage"
	"github.com/ananaslegend/short-link/usecases"
	"github.com/ananaslegend/short-link/utils/request"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type LinkSaver interface {
	AddLink(link, alias string) error
}

func Link(c context.Context, log *slog.Logger, ls LinkSaver) func(http.ResponseWriter, *http.Request) {
	const op = "api.handlers.Link"

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			b, err := io.ReadAll(r.Body)
			if err != nil {
				err = fmt.Errorf("%s:%w", op, err)
				log.Error("cant read body", logs.Err(err))
				return
			}
			model := request.AddLink{}
			if err = json.Unmarshal(b, &model); err != nil {
				err = fmt.Errorf("%s:%w", op, err)
				log.Error("cant unmarshal request body", logs.Err(err))
				return
			}

			var autoAlias bool
			if len(model.Alias) == 0 {
				model.Alias = core.MakeShorter(uuid.New().ID())
				autoAlias = true
			}

			if err = usecases.AddLink(c, ls, model.Link, model.Alias); err != nil {
				if errors.Is(err, storage.ErrAliasExists) {
					if !autoAlias {
						w.WriteHeader(http.StatusConflict)

						// JSON
						return
					}

					log.Error(fmt.Sprintf("auto generated alias: %s already exists", model.Alias), logs.Err(err))
				}

				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			h := http.Response{ //TODO remove
				Status:           "",
				StatusCode:       0,
				Proto:            "",
				ProtoMajor:       0,
				ProtoMinor:       0,
				Header:           nil,
				Body:             nil,
				ContentLength:    0,
				TransferEncoding: nil,
				Close:            false,
				Uncompressed:     false,
				Trailer:          nil,
				Request:          nil,
				TLS:              nil,
			}
			h.Write(w)
			w.WriteHeader(http.StatusCreated) // json
			return
		}
	}
}
