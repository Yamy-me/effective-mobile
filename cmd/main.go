package main

import (
	"Effective-Mobile/pkg/repository"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../.env")

	logger := NewLoger()
	slog.SetDefault(logger)

	DSN := os.Getenv("DSN")

	log.Println(DSN)

	err := repository.RunMigrations(DSN)
	if err != nil {
		slog.Error("Ошибка с миграцией", slog.String("error", err.Error()))
	}
}

func NewLoger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return log
}
