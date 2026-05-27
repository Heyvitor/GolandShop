package handler

import (
	"errors"
	"net/http"

	"goapi/backend/internal/app"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) register(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	if !api.services.Auth.AllowRequest(r.Context(), ip, 5, time.Minute) {
		writeError(w, http.StatusTooManyRequests, "too_many_requests")
		return
	}

	var req registerRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	result, err := api.services.Auth.Register(r.Context(), req.Name, req.Email, req.Password)
	if errors.Is(err, app.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid_input")
		return
	}
	if errors.Is(err, app.ErrEmailAlreadyExists) {
		writeError(w, http.StatusConflict, "email_already_exists")
		return
	}
	if err != nil {
		api.logger.Error("register", "error", err, "request_id", requestIDFromContext(r.Context()))
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	result, err := api.services.Auth.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, app.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}
	if err != nil {
		api.logger.Error("login", "error", err, "request_id", requestIDFromContext(r.Context()))
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    result.Token,
		Path:     "/",
		Expires:  result.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, map[string]any{"user": result.User})
}

func (api *API) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Value != "" {
		claims, err := api.tokens.Parse(cookie.Value)
		if err == nil {
			// Blacklist token in Redis
			_ = api.services.Auth.Logout(r.Context(), claims.ID, claims.ExpiresAt.Time)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, map[string]any{"message": "logged_out"})
}

