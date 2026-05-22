package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	DSN string
}

func RunConfigs() *Configs {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("Configs error", slog.String("error", err.Error()))
	}
	slog.Info("Loading....", slog.String("DSN", os.Getenv("DSN")))
	return &Configs{
		DSN: os.Getenv("DSN"),
	}
}
