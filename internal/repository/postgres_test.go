package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	models "Effective-Mobile/internal/model"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Вспомогательная функция для подготовки чистой базы перед каждым тестом
func setupTestDB(t *testing.T) *sql.DB {
	dsn := "postgres://Atay:admin123@localhost:5433/something?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой БД: %v", err)
	}

	// Очищаем таблицу перед тестом
	_, err = db.Exec("TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("Не удалось очистить таблицу: %v", err)
	}

	return db
}

// ==========================================
// ТЕСТ: GetSubsByID (Получение по ID)
// ==========================================
func TestPostgresRepository_GetSubsByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db) // Учтено имя конструктора
	ctx := context.Background()

	testUUID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	expectedSub := models.Subscription{
		ServiceName: "Netflix",
		Price:       900,
		UserID:      testUUID,
		StartDate:   time.Now().Truncate(time.Second),
	}

	var insertedID int
	insertQuery := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date) 
		VALUES ($1, $2, $3, $4) RETURNING id`

	err := db.QueryRowContext(ctx, insertQuery, expectedSub.ServiceName, expectedSub.Price, expectedSub.UserID, expectedSub.StartDate).Scan(&insertedID)
	if err != nil {
		t.Fatalf("Не удалось вставить тестовую подписку: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		sub, err := repo.GetSubsByID(ctx, insertedID)
		if err != nil {
			t.Fatalf("Ожидался успех, но получена ошибка: %v", err)
		}
		if sub == nil {
			t.Fatal("Ожидалась структура подписки, но получен nil")
		}
		if sub.ServiceName != expectedSub.ServiceName {
			t.Errorf("Ожидалось имя %s, получено %s", expectedSub.ServiceName, sub.ServiceName)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		sub, err := repo.GetSubsByID(ctx, 999999)
		if err == nil {
			t.Error("Ожидалась ошибка 'не найдено', но получено err = nil")
		}
		if sub != nil {
			t.Errorf("Ожидался nil вместо подписки, получено: %v", sub)
		}
	})
}

// ==========================================
// ТЕСТ: UpdateSubs (Обновление)
// ==========================================
func TestPostgresRepository_UpdateSubs(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db)
	ctx := context.Background()

	testUUID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	originalSub := models.Subscription{
		ServiceName: "Old Service",
		Price:       100,
		UserID:      testUUID,
		StartDate:   time.Now().Truncate(time.Second),
	}

	var id int
	err := db.QueryRow(`INSERT INTO subscriptions (service_name, price, user_id, start_date) VALUES ($1, $2, $3, $4) RETURNING id`,
		originalSub.ServiceName, originalSub.Price, originalSub.UserID, originalSub.StartDate).Scan(&id)
	if err != nil {
		t.Fatalf("Ошибка вставки: %v", err)
	}

	t.Run("Success", func(t *testing.T) {
		updatedSub := originalSub
		updatedSub.ServiceName = "New Service"
		updatedSub.Price = 500

		err := repo.UpdateSubs(ctx, id, updatedSub)
		if err != nil {
			t.Fatalf("Ожидался успех, получена ошибка: %v", err)
		}

		var newName string
		var newPrice int
		err = db.QueryRow(`SELECT service_name, price FROM subscriptions WHERE id = $1`, id).Scan(&newName, &newPrice)
		if err != nil || newName != "New Service" || newPrice != 500 {
			t.Errorf("Данные в базе не обновились. Получено: %s, %d", newName, newPrice)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		err := repo.UpdateSubs(ctx, 99999, originalSub)
		if err != sql.ErrNoRows {
			t.Errorf("Ожидалась ошибка sql.ErrNoRows, получена: %v", err)
		}
	})
}

// ==========================================
// ТЕСТ: ListSubs (Получение списка с динамическими фильтрами)
// ==========================================
func TestPostgresRepository_ListSubs(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db)
	ctx := context.Background()

	user1 := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	user2 := uuid.MustParse("11111111-2bf1-4721-ae6f-7636e79a0cba")

	queries := []struct {
		name   string
		price  int
		userID uuid.UUID
	}{
		{"Yandex Plus", 300, user1},
		{"Netflix", 900, user1},
		{"Yandex Plus", 300, user2},
	}

	for _, q := range queries {
		_, err := db.Exec(`INSERT INTO subscriptions (service_name, price, user_id, start_date) VALUES ($1, $2, $3, $4)`,
			q.name, q.price, q.userID, time.Now())
		if err != nil {
			t.Fatalf("Ошибка вставки: %v", err)
		}
	}

	t.Run("Filter by UserID only", func(t *testing.T) {
		filter := models.ListFilter{UserID: &user1}

		subs, err := repo.ListSubs(ctx, filter)
		if err != nil {
			t.Fatalf("Ошибка: %v", err)
		}
		if len(subs) != 2 {
			t.Errorf("Ожидалось 2 подписки, получено %d", len(subs))
		}
	})

	t.Run("Filter by UserID and ServiceName", func(t *testing.T) {
		sName := "Netflix"
		filter := models.ListFilter{UserID: &user1, ServiceName: &sName}

		subs, err := repo.ListSubs(ctx, filter)
		if err != nil {
			t.Fatalf("Ошибка: %v", err)
		}
		if len(subs) != 1 || subs[0].ServiceName != "Netflix" {
			t.Errorf("Ожидалась 1 подписка Netflix, получено %d", len(subs))
		}
	})
}

// ==========================================
// ТЕСТ: TotalCost (Подсчет суммы за период)
// ==========================================
func TestPostgresRepository_TotalCost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepo(db)
	ctx := context.Background()

	testUUID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	layout := "2006-01-02"

	start1, _ := time.Parse(layout, "2025-01-01")
	end1, _ := time.Parse(layout, "2025-01-10")

	start2, _ := time.Parse(layout, "2025-01-15")
	end2, _ := time.Parse(layout, "2025-01-30")

	_, _ = db.Exec(`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)`,
		"Sub 1", 100, testUUID, start1, end1)
	_, _ = db.Exec(`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)`,
		"Sub 2", 200, testUUID, start2, end2)

	t.Run("Period covers both subscriptions", func(t *testing.T) {
		from, _ := time.Parse(layout, "2025-01-05")
		to, _ := time.Parse(layout, "2025-01-20")

		filter := models.TotalCostFilter{
			FromDate: &from,
			ToDate:   &to,
		}

		total, err := repo.TotalCost(ctx, testUUID, filter)
		if err != nil {
			t.Fatalf("Ошибка: %v", err)
		}

		if total != 300 {
			t.Errorf("Ожидалась сумма 300, получено %d", total)
		}
	})

	t.Run("Period covers only first subscription", func(t *testing.T) {
		from, _ := time.Parse(layout, "2025-01-01")
		to, _ := time.Parse(layout, "2025-01-12")

		filter := models.TotalCostFilter{
			FromDate: &from,
			ToDate:   &to,
		}

		total, err := repo.TotalCost(ctx, testUUID, filter)
		if err != nil {
			t.Fatalf("Ошибка: %v", err)
		}

		if total != 100 {
			t.Errorf("Ожидалась сумма 100, получено %d", total)
		}
	})
}
