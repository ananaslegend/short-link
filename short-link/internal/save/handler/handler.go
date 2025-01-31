package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ananaslegend/short-link/internal/save/service"
	"github.com/ananaslegend/short-link/pkg/cslog"
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
}

func New(srv LinkSetterService) *Handler {
	return &Handler{
		service: srv,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		cslog.Logger(r.Context()).Error("cant read body", cslog.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req Request
	if err = json.Unmarshal(b, &req); err != nil { // TODO Add Validation
		cslog.Logger(r.Context()).With("body", b).Error("cant unmarshal request body", cslog.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addedAlias, err := h.service.AddLink(r.Context(), req.Link, req.Alias)
	if err != nil {
		h.renderError(r.Context(), err, w)
		return
	}

	h.renderResponse(r.Context(), w, addedAlias)
}

func (h Handler) renderResponse(ctx context.Context, w http.ResponseWriter, addedAlias string) {
	var resp Response
	resp.Alias = addedAlias
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		cslog.Logger(ctx).Error("cant encode json", cslog.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h Handler) renderError(ctx context.Context, err error, w http.ResponseWriter) {
	var resp Response

	switch {
	case errors.Is(err, service.ErrAliasAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		resp.Error = "alias already exists"
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			cslog.Logger(ctx).Error("cant encode json", cslog.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}

	default:
		cslog.Logger(ctx).Error("failed to add link", cslog.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
