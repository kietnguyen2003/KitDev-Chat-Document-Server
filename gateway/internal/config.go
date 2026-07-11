package gateway

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Secretkey string
	ServerURL string
	RagURL    string
}

func LoadConfig() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Can not file .env file")
	}

	return &Config{
		Port:      loadEnv("PORT"),
		Secretkey: loadEnv("JWT_SECRET"),
		ServerURL: loadEnvOrDefault("SERVER_URL", "http://localhost:8081"),
		RagURL:    loadEnvOrDefault("RAG_URL", "http://localhost:8000"),
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
