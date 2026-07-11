package gateway

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	Secretkey              string
	ServerURLs             []string
	RagURLs                []string
	RateLimitRequests      int
	RateLimitWindow        time.Duration
	RateLimitCleanupWindow time.Duration
}

func LoadConfig() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Can not file .env file")
	}

	return &Config{
		Port:                   loadEnv("PORT"),
		Secretkey:              loadEnv("JWT_SECRET"),
		ServerURLs:             loadURLList("SERVER_URL", "http://localhost:8081"),
		RagURLs:                loadURLList("RAG_URL", "http://localhost:8000"),
		RateLimitRequests:      loadIntOrDefault("RATE_LIMIT_REQUESTS", 120),
		RateLimitWindow:        time.Duration(loadIntOrDefault("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
		RateLimitCleanupWindow: time.Duration(loadIntOrDefault("RATE_LIMIT_CLEANUP_SECONDS", 300)) * time.Second,
	}
}

func LoadSecretKey() []byte {
	return []byte(loadEnv("JWT_SECRET"))
}

func loadEnv(str string) string {
	value := os.Getenv(str)
	return value
}

func loadEnvOrDefault(key, fallback string) string {
	value := loadEnv(key)
	if value == "" {
		return fallback
	}
	return value
}

func loadURLList(key, fallback string) []string {
	raw := loadEnvOrDefault(key, fallback)
	parts := strings.Split(raw, ",")
	urls := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		urls = append(urls, trimmed)
	}

	if len(urls) == 0 {
		return []string{fallback}
	}

	return urls
}

func loadIntOrDefault(key string, fallback int) int {
	raw := strings.TrimSpace(loadEnv(key))
	if raw == "" {
		return fallback
	}

	var value int
	if _, err := fmt.Sscanf(raw, "%d", &value); err != nil || value <= 0 {
		return fallback
	}

	return value
}
