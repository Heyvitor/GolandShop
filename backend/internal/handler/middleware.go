package handler

import "net/http"

func (api *API) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			writeError(w, http.StatusUnauthorized, "missing_token")
			return
		}

		claims, err := api.tokens.Parse(cookie.Value)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid_token")
			return
		}

		if api.services.Auth.IsTokenBlacklisted(r.Context(), claims.ID) {
			writeError(w, http.StatusUnauthorized, "token_revoked")
			return
		}

		next.ServeHTTP(w, r.WithContext(withUserID(r.Context(), claims.UserID)))
	})
}
