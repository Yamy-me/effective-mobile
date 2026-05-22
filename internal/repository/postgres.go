package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	models "Effective-Mobile/internal/model"

	"github.com/google/uuid"
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
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

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
			return nil, fmt.Errorf("подписка с ID %d не найдена", id)
		}

		slog.Error("ошибка при получении подписки из БД",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)

		return nil, err
	}

	slog.Info("Подписка успешна найдена")
	return &Subscription, nil
}

func (r *PostgresRepository) UpdateSubs(ctx context.Context, id int, subs models.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(ctx, query,
		subs.ServiceName,
		subs.Price,
		subs.UserID,
		subs.StartDate,
		subs.EndDate,
		id,
	)
	if err != nil {
		slog.Error("ошибка при выполнении SQL-запроса на обновление",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("ошибка при получении количества обновленных строк",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	if rowsAffected == 0 {
		slog.Info("попытка обновить несуществующую подписку", slog.Int("id", id))
		return fmt.Errorf("попытка обновить несуществующую подписку")
	}

	slog.Info("подписка успешно обновлена", slog.Int("id", id))
	return nil
}

func (r *PostgresRepository) DeleteSubs(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("ошибка удаления записи",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("ошибка при получении количества обновленных строк",
			slog.Int("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	if rowsAffected == 0 {
		slog.Info("попытка удалить несуществующую подписку", slog.Int("id", id))
		return sql.ErrNoRows
	}

	slog.Info("запись удалена", slog.Int("id", id))
	return nil
}

func (r *PostgresRepository) ListSubs(ctx context.Context, filter models.ListFilter) ([]models.Subscription, error) {
	var subs []models.Subscription

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE 1=1`

	var args []any
	argID := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argID)
		args = append(args, *filter.UserID)
		argID++
	}

	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argID)
		args = append(args, *filter.ServiceName)
		argID++
	}

	if filter.FromDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argID)
		args = append(args, *filter.FromDate)
		argID++
	}

	if filter.ToDate != nil {
		query += fmt.Sprintf(" AND start_date <= $%d", argID)
		args = append(args, *filter.ToDate)
		argID++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Error("ошибка при получении списка подписок", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *PostgresRepository) TotalCost(ctx context.Context, userID uuid.UUID, filter models.TotalCostFilter) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE user_id = $1`

	args := []any{userID}
	argID := 2

	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argID)
		args = append(args, *filter.ServiceName)
		argID++
	}

	if filter.ToDate != nil {
		query += fmt.Sprintf(" AND start_date <= $%d", argID)
		args = append(args, *filter.ToDate)
		argID++
	}

	if filter.FromDate != nil {
		query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", argID)
		args = append(args, *filter.FromDate)
		argID++
	}

	var totalCost int

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		slog.Error("ошибка при подсчете общей стоимости подписок",
			slog.String("user_id", userID.String()),
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	slog.Info("успешно подсчитана стоимость подписок",
		slog.String("user_id", userID.String()),
		slog.Int("total_cost", totalCost),
	)

	return totalCost, nil
}
