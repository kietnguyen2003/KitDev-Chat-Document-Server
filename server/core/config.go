package core

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Secretkey string
	DBURL     string
	RagURL    string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
}

func LoadConfig() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Can not file .env file")
	}

	return &Config{
		Port:           loadEnv("PORT"),
		Secretkey:      loadEnv("JWT_SECRET"),
		DBURL:          loadEnv("DB_URL"),
		RagURL:         loadEnv("RAG_URL"),
		MinioEndpoint:  loadEnv("MINIO_ENDPOINT"),
		MinioAccessKey: loadEnv("MINIO_ACCESS_KEY"),
		MinioSecretKey: loadEnv("MINIO_SECRET_KEY"),
	}
}

func loadEnv(str string) string {
	value := os.Getenv(str)
	return value
}
