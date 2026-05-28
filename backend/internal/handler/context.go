package handler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userRoleKey  contextKey = "user_role"
	userNameKey  contextKey = "user_name"
	userEmailKey contextKey = "user_email"
	requestIDKey contextKey = "request_id"
)

func withUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}

func withUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

func userRoleFromContext(ctx context.Context) string {
	role, _ := ctx.Value(userRoleKey).(string)
	return role
}

func withUserName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, userNameKey, name)
}

func userNameFromContext(ctx context.Context) string {
	name, _ := ctx.Value(userNameKey).(string)
	return name
}

func withUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailKey, email)
}

func userEmailFromContext(ctx context.Context) string {
	email, _ := ctx.Value(userEmailKey).(string)
	return email
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	})
}

func requestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}
