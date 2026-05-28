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

		ctx := withUserID(r.Context(), claims.UserID)
		ctx = withUserRole(ctx, claims.Role)
		ctx = withUserName(ctx, claims.Name)
		ctx = withUserEmail(ctx, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) requireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := userRoleFromContext(r.Context())
			
			allowed := false
			for _, role := range roles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				writeError(w, http.StatusForbidden, "insufficient_permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
