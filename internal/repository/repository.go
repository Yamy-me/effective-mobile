package repository


import (
	models "Effective-Mobile/internal/model"
	"context"
)

type Repository interface {
	Create(context.Context, models.Subscription)
}
