package server

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AESKeyB64   string
}

func LoadConfigFromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		Port:        port,
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/blueberry?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", "supersecret"),
		AESKeyB64:   getenv("AES_KEY_BASE64", ""),
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
