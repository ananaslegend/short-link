package save

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/errs"
	"github.com/ananaslegend/short-link/logs"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type LinkAdder interface {
	AddLink(c context.Context, link, alias string) (string, error)
}

type Request struct {
	Link  string `json:"link"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Alias string `json:"alias,omitempty"`
	Error string `json:"error,omitempty"`
}

func Handle(w http.ResponseWriter, r *http.Request, log *slog.Logger, service LinkAdder) {
	const op = "api.handlers.save.link.Handle"
	log.With(slog.String("op", op))
	var (
		req  Request
		resp Response
	)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("%s:%w", op, err)
		log.Error("cant read body", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(b, &req); err != nil { // TODO Add Validation
		err = fmt.Errorf("%s:%w", op, err)
		log.Error("cant unmarshal request body", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addedAlias, err := service.AddLink(r.Context(), req.Link, req.Alias)
	if err != nil {
		if errors.Is(err, errs.ErrAutoAliasAlreadyExists) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if errors.Is(err, errs.ErrAliasExists) {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			resp.Error = "alias already exists"
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				log.Error("cant encode json", logs.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	resp.Alias = addedAlias
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("cant encode json", logs.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	return

}
