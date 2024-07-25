package config

import (
	"os"

	"github.com/joho/godotenv"
)

func GetSecret(key string) string {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv(key)
}
