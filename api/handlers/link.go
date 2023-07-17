package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ananaslegend/short-link/logs"
	"github.com/ananaslegend/short-link/usecases"
	"github.com/ananaslegend/short-link/utils/request"
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
			reqBody := request.AddLink{}
			if err = json.Unmarshal(b, &reqBody); err != nil {
				err = fmt.Errorf("%s:%w", op, err)
				log.Error("cant unmarshal request body", logs.Err(err))
				return
			}

			if err = usecases.AddLink(c, log, ls, reqBody.Link, reqBody.Alias); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK) // json
			return
		}
	}
}
