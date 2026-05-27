package handler

import (
	"errors"
	"net/http"
	"strconv"

	"goapi/backend/internal/app"
)

type createItemRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (api *API) createItem(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	item, err := api.services.Items.Create(r.Context(), userIDFromContext(r.Context()), req.Title, req.Body)
	if errors.Is(err, app.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid_input")
		return
	}
	if err != nil {
		api.logger.Error("create item", "error", err, "request_id", requestIDFromContext(r.Context()))
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (api *API) listItems(w http.ResponseWriter, r *http.Request) {
	limit := int32(50)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_limit")
			return
		}
		limit = int32(parsed)
	}

	items, err := api.services.Items.List(r.Context(), userIDFromContext(r.Context()), limit)
	if errors.Is(err, app.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid_input")
		return
	}
	if err != nil {
		api.logger.Error("list items", "error", err, "request_id", requestIDFromContext(r.Context()))
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
