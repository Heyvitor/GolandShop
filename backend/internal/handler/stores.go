package handler

import (
	"errors"
	"net/http"
	"strconv"

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
	slug := storeSlugFromRequest(r)
	if slug == "" {
		writeError(w, http.StatusBadRequest, "missing_slug")
		return
	}

	store, err := api.services.Stores.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "store_not_found")
		return
	}

	writeJSON(w, http.StatusOK, store)
}

func (api *API) getStoreCatalog(w http.ResponseWriter, r *http.Request) {
	slug := storeSlugFromRequest(r)
	if slug == "" {
		writeError(w, http.StatusBadRequest, "missing_slug")
		return
	}

	limit := int32(50)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_limit")
			return
		}
		limit = int32(parsed)
	}

	store, err := api.services.Stores.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, http.StatusNotFound, "store_not_found")
		return
	}

	items, err := api.services.Items.ListByStore(r.Context(), store.ID, limit)
	if errors.Is(err, app.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid_input")
		return
	}
	if err != nil {
		api.logger.Error("store catalog", "error", err, "request_id", requestIDFromContext(r.Context()))
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"store": store,
		"items": items,
	})
}

func storeSlugFromRequest(r *http.Request) string {
	if slug := r.PathValue("slug"); slug != "" {
		return slug
	}

	return r.URL.Query().Get("slug")
}

func (api *API) getMyStore(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r.Context())
	store, err := api.services.Stores.GetByOwner(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "store_not_found")
		return
	}

	writeJSON(w, http.StatusOK, store)
}
