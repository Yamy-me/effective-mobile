package main

import (
	"log/slog"
	"os"
	"time"

	"Effective-Mobile/internal/config"
	"Effective-Mobile/internal/repository"
)

func main() {
	logger := NewLoger()
	slog.SetDefault(logger)
	cfg := config.RunConfigs()

	err := repository.RunMigrations(cfg.DSN)
	if err != nil {
		slog.Error("Ошибка с миграцией", slog.String("error", err.Error()))
	}
}

func NewLoger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String(a.Key, a.Value.Time().Format(time.DateTime))
			}
			return a
		},
	}))

	return log
}
