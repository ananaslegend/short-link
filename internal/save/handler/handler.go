package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ananaslegend/go-logs/v2"
	"github.com/ananaslegend/short-link/internal/save/service"
	"io"
	"log/slog"
	"net/http"
)

type LinkSetterService interface {
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
	service LinkSetterService
	log     *slog.Logger
}

func New(srv LinkSetterService, log *slog.Logger) *Handler {
	return &Handler{
		service: srv,
		log:     log,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "internal.save.handler.ServeHTTP"
	var (
		logger = h.log.With(slog.String("op", op))
		req    Request
		resp   Response
	)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("cant read body", logs.ErrorMsg(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(b, &req); err != nil { // TODO Add Validation
		logger.Info(fmt.Sprintf("cant unmarshal request body. body: %s", b), logs.ErrorMsg(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addedAlias, err := h.service.AddLink(r.Context(), req.Link, req.Alias)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAliasAlreadyExists):
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			resp.Error = "alias already exists"
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				logger.Error("cant encode json", logs.ErrorMsg(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		default:
			logger.Error("failed to add link", logs.ErrorMsg(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	resp.Alias = addedAlias
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("cant encode json", logs.ErrorMsg(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
