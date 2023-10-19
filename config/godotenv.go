package config

import (
	"os"

	"github.com/joho/godotenv"
)

func RenderEnv(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	return os.Getenv(key)
}
