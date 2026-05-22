package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	models "Effective-Mobile/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (r *PostgresRepository) CreateSubs(ctx context.Context, subs *models.Subscription) (int, error) {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int

	err := r.db.QueryRowContext(ctx, query, subs.ServiceName, subs.Price, subs.UserID, subs.StartDate, subs.EndDate).Scan(&id)
	if err != nil {
		slog.Error("Database error", slog.String("error", err.Error()))
		return -1, err
	}

	slog.Info("Подписка успешна создана", slog.Any("Подписка", subs))
	return id, nil
}

func (r *PostgresRepository) GetSubsByID(ctx context.Context, id int) (*models.Subscription, error) {
	var Subscription models.Subscription
	query := `SELECT * FROM subscriptions WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&Subscription.ID,
		&Subscription.ServiceName,
		&Subscription.Price,
		&Subscription.UserID,
		&Subscription.StartDate,
		&Subscription.EndDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info("Подписка с ID %d не найдена", slog.Int("id", id))
			return nil, fmt.Errorf("Подписка с ID %d не найдена", id)
		}

		slog.Error("ошибка при получении подписки из БД",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	return &Subscription, nil
}
