package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `json:"id" db:"id"`
	ServiceName string     `json:"service_name" db:"service_name"`
	Price       int        `json:"price" db:"price"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
}

type ListFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	FromDate    *time.Time
	ToDate      *time.Time
}

type TotalCostFilter struct {
	ServiceName *string
	FromDate    *time.Time
	ToDate      *time.Time
}
