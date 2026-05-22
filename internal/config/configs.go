package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	DSN string
}

func LoadConfigs() *Configs {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error(".env Problem", slog.String("error", err.Error()))
		return nil
	}
	return &Configs{
		DSN: os.Getenv("DSN"),
	}
}
