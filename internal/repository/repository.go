package repository

import (
	"context"

	models "Effective-Mobile/internal/model"
)

type Repository interface {
	CreateSubs(ctx context.Context, subs *models.Subscription) (int, error)
	GetSubsByID(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubs(ctx context.Context, id int, subs models.Subscription) error
	DeleteSubs(ctx context.Context, id int) error
}
