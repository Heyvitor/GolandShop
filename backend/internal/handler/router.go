package handler

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"goapi/backend/internal/app"
	"goapi/backend/internal/model"
	"goapi/backend/internal/security"
)

func NewRouter(services *app.Services, tokens *security.TokenService, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	api := &API{
		services: services,
		tokens:   tokens,
		logger:   logger,
	}

	// Rotas Públicas
	mux.HandleFunc("GET /healthz", api.health)
	mux.HandleFunc("POST /api/v1/auth/register", api.register)
	mux.HandleFunc("POST /api/v1/auth/login", api.login)
	mux.HandleFunc("GET /api/v1/stores/view", api.getStore)

	// Rotas Protegidas (Logado)
	authOnly := api.auth
	mux.Handle("POST /api/v1/auth/logout", authOnly(http.HandlerFunc(api.logout)))
	mux.Handle("GET /api/v1/auth/me", authOnly(http.HandlerFunc(api.me)))

	// Rotas de Admin
	adminOnly := api.requireRole(model.RoleAdmin)
	mux.Handle("GET /api/v1/admin/stats", authOnly(adminOnly(http.HandlerFunc(api.health)))) // Exemplo

	// Rotas de Usuário (Dono de Loja)
	userOnly := api.requireRole(model.RoleUser, model.RoleAdmin)
	mux.Handle("GET /api/v1/stores/mine", authOnly(userOnly(http.HandlerFunc(api.getMyStore))))
	mux.Handle("POST /api/v1/stores", authOnly(userOnly(http.HandlerFunc(api.createStore))))
	mux.Handle("POST /api/v1/items", authOnly(userOnly(http.HandlerFunc(api.createItem))))
	mux.Handle("GET /api/v1/items", authOnly(userOnly(http.HandlerFunc(api.listItems))))

	// Rotas de Cliente
	clientOnly := api.requireRole(model.RoleClient, model.RoleUser, model.RoleAdmin)
	mux.Handle("GET /api/v1/orders/my", authOnly(clientOnly(http.HandlerFunc(api.listItems)))) // Exemplo

	return gzipMiddleware(recoverPanic(requestID(logging(securityHeaders(mux), logger)), logger))
}

type API struct {
	services *app.Services
	tokens   *security.TokenService
	logger   *slog.Logger
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", requestIDFromContext(r.Context()),
		)
	})
}

func recoverPanic(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic", "error", err, "request_id", requestIDFromContext(r.Context()))
				writeError(w, http.StatusInternalServerError, "internal_error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, _ := strings.Cut(r.RemoteAddr, ":")
	return ip
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}
