package repository

import (
	"context"
	"database/sql"
	"log/slog"

	models "Effective-Mobile/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (r *PostgresRepository) CreateSubs(ctx context.Context, subs *models.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, subs.ServiceName, subs.Price, subs.UserID, subs.StartDate, subs.EndDate)
	if err != nil {
		slog.Info("Database error", slog.String("error", err.Error()))
		return err
	}
	return nil
}
