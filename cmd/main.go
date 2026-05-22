package main

import (
	"log/slog"

	"Effective-Mobile/internal/config"
	"Effective-Mobile/internal/logger"
	"Effective-Mobile/internal/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5"
)

func main() {
	// Logger
	logger.NewLogger()

	// Load configs
	cfg := config.LoadConfigs()

	// Migrations
	err := repository.RunMigrations(cfg.DSN)
	if err != nil {
		slog.Error("Ошибка с миграцией", slog.String("error", err.Error()))
		return
	}

	// Server initialization
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.Use(logger.MiddleWareLogger())
	server.Use(gin.Recovery())
}
