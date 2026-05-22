package repository

import (
	"context"

	models "Effective-Mobile/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSubs(ctx context.Context, subs *models.Subscription) (int, error)
	GetSubsByID(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubs(ctx context.Context, id int, subs models.Subscription) error
	DeleteSubs(ctx context.Context, id int) error
	ListSubs(ctx context.Context, filter models.ListFilter) ([]models.Subscription, error)
	TotalCost(ctx context.Context, userID uuid.UUID, filter models.TotalCostFilter) (int, error)
}
