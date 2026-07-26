package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment          string
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	CORSOrigins          []string
	GeminiAPIKey         string
	GeminiChatModel      string
	GeminiEmbeddingModel string
	EmbeddingDimensions  int
	IPHashSecret         string
}

func Load() (Config, error) {
	cfg := Config{
		Environment:          envOr("APP_ENV", "development"),
		Port:                 envOr("PORT", "8080"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		AccessTokenTTL:       durationOr("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL:      durationOr("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		CORSOrigins:          csvOr("CORS_ORIGINS", []string{"http://localhost:5173"}),
		GeminiAPIKey:         os.Getenv("GEMINI_API_KEY"),
		GeminiChatModel:      envOr("GEMINI_CHAT_MODEL", "gemini-2.0-flash"),
		GeminiEmbeddingModel: envOr("GEMINI_EMBEDDING_MODEL", "text-embedding-004"),
		EmbeddingDimensions:  intOr("EMBEDDING_DIMENSIONS", 768),
		IPHashSecret:         os.Getenv("IP_HASH_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL no está configurada")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET debe tener al menos 32 caracteres")
	}
	if cfg.IPHashSecret == "" {
		cfg.IPHashSecret = cfg.JWTSecret
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func csvOr(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}

func durationOr(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func intOr(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
