package handler

import (
	"errors"
	"net/http"

	"goapi/backend/internal/app"
)

type createStoreRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (api *API) createStore(w http.ResponseWriter, r *http.Request) {
	var req createStoreRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	userID := userIDFromContext(r.Context())
	store, err := api.services.Stores.Create(r.Context(), userID, req.Name, req.Slug)
	if errors.Is(err, app.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid_input")
		return
	}
	if err != nil {
		api.logger.Error("create store", "error", err, "request_id", requestIDFromContext(r.Context()))
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, store)
}

func (api *API) getStore(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		// Tentar buscar do contexto se for dono de loja logado? 
		// Por enquanto via query
	}

	store, err := api.services.Stores.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "store_not_found")
		return
	}

	writeJSON(w, http.StatusOK, store)
}
