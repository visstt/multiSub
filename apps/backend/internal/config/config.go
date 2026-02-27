package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        int
	DatabaseURL string
	JWTSecret   string

	// JWT сроки жизни (в минутах)
	AccessTokenTTL  int
	RefreshTokenTTL int // в днях

	// Bcrypt
	BcryptCost int

	// Rate limiting
	RateLimitRequests int
	RateLimitWindow   int // в секундах

	// CORS
	CORSOrigins string
}

func Load() (*Config, error) {
	// Загружаем .env файл (игнорируем ошибку — в production .env может не быть)
	_ = godotenv.Load()

	cfg := &Config{
		Env:               getEnv("APP_ENV", "development"),
		Port:              getEnvInt("PORT", 8080),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		AccessTokenTTL:    getEnvInt("ACCESS_TOKEN_TTL", 15),
		RefreshTokenTTL:   getEnvInt("REFRESH_TOKEN_TTL", 7),
		BcryptCost:        getEnvInt("BCRYPT_COST", 12),
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 60),
		RateLimitWindow:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		CORSOrigins:       getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("переменная DATABASE_URL обязательна")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("переменная JWT_SECRET обязательна")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
