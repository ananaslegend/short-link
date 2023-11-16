package save

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ananaslegend/short-link/pkg/logs"
	"io"
	"log/slog"
	"net/http"
)

type LinkSetterUseCase interface {
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

type Handler struct {
	useCase LinkSetterUseCase
	log     *slog.Logger
}

func NewHandler(uc LinkSetterUseCase, log *slog.Logger) *Handler {
	return &Handler{
		useCase: uc,
		log:     log,
	}
}

func (h Handler) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "internal.save.handler.HandleHTTP"
	log := h.log.With(slog.String("op", op))
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
		log.Info(fmt.Sprintf("cant unmarshal request body. body: %s", b), logs.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addedAlias, err := h.useCase.AddLink(r.Context(), req.Link, req.Alias)
	if err != nil {
		if errors.Is(err, ErrAutoAliasAlreadyExists) {
			log.Error("auto generated alias already exists", logs.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if errors.Is(err, ErrAliasAlreadyExists) {
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

		w.WriteHeader(http.StatusInternalServerError)
		return
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
