package main

import (
	"database/sql"
	"log/slog"

	"Effective-Mobile/internal/config"
	"Effective-Mobile/internal/handler"
	"Effective-Mobile/internal/logger"
	"Effective-Mobile/internal/repository"
	"Effective-Mobile/internal/service"

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
	// Database Connection
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil{
		slog.Error("Подключение к ДБ отсуствует")
		return
	}

	// Repository
	repo := repository.NewPostgresRepo(db)

	// Service manager
	service := service.NewSubscriptionService(repo)

	// Handler 
	hnd := handler.NewHandler(service)

	// Server initialization
	gin.SetMode(gin.ReleaseMode)

	server := gin.New()
	server.Use(logger.MiddleWareLogger())
	server.Use(gin.Recovery())
	handler.InitRoutes(server, hnd)
}
