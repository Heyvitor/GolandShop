package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv   string
	HTTPAddr string
	Database DatabaseConfig
	Security SecurityConfig
}

type DatabaseConfig struct {
	URL      string
	MaxConns int32
	MinConns int32
}

type SecurityConfig struct {
	JWTSecret  string
	JWTTTL     time.Duration
	BcryptCost int
}

func Load() (Config, error) {
	jwtTTL, err := time.ParseDuration(getenv("JWT_TTL", "24h"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:   getenv("APP_ENV", "development"),
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		Database: DatabaseConfig{
			URL:      getenv("DATABASE_URL", ""),
			MaxConns: int32(getenvInt("DB_MAX_CONNS", 40)),
			MinConns: int32(getenvInt("DB_MIN_CONNS", 5)),
		},
		RedisURL: getenv("REDIS_URL", "redis://localhost:6379/0"),
		Security: SecurityConfig{
			JWTSecret:  getenv("JWT_SECRET", ""),
			JWTTTL:     jwtTTL,
			BcryptCost: getenvInt("BCRYPT_COST", 12),
		},
	}

	if cfg.Database.URL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if len(cfg.Security.JWTSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must have at least 32 characters")
	}
	if cfg.Security.BcryptCost < 10 || cfg.Security.BcryptCost > 15 {
		return Config{}, errors.New("BCRYPT_COST must be between 10 and 15")
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
