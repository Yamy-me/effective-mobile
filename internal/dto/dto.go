package dto

import "github.com/google/uuid"

type CreateSubRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   MonthYear `json:"start_date" binding:"required"`
}

type UpdateSubRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required,min=0"`
	UserID      uuid.UUID  `json:"user_id" binding:"required"`
	StartDate   MonthYear  `json:"start_date" binding:"required"`
	EndDate     *MonthYear `json:"end_date"`
}

type SubResponse struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date"`
}
