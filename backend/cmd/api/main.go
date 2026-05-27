package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goapi/backend/internal/app"
	"goapi/backend/internal/config"
	"goapi/backend/internal/database"
	"goapi/backend/internal/handler"
	"goapi/backend/internal/repository"
	"goapi/backend/internal/security"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(ctx, cfg.Database)
	if err != nil {
		logger.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	passwords := security.NewPasswordHasher(cfg.Security.BcryptCost)
	tokens := security.NewTokenService(cfg.Security.JWTSecret, cfg.Security.JWTTTL)

	services := app.NewServices(userRepo, itemRepo, passwords, tokens, rdb)
	router := handler.NewRouter(services, tokens, logger)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		logger.Info("api listening", "addr", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("api stopped")
}
ogger.Info("api stopped")
}


	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("api stopped")
}
