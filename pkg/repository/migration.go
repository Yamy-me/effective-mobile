package repository

import (
	"Effective-Mobile/migrations"
	"database/sql"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	goose.SetBaseFS(migrations.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	slog.Info("Миграция успешно применена")
	return nil
}
