package service

import (
	"context"
	"fmt"
	"time"

	"Effective-Mobile/internal/dto"
	models "Effective-Mobile/internal/model"
	"Effective-Mobile/internal/repository"
)

type ServiceSubs struct {
	repo repository.Repository
}

func NewSubscriptionService(repo repository.Repository) *ServiceSubs {
	return &ServiceSubs{repo: repo}
}

func (s *ServiceSubs) Create(ctx context.Context, req dto.CreateSubRequest) (int, error) {
	model := models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate.Time,
	}

	return s.repo.CreateSubs(ctx, &model)
}

func (s *ServiceSubs) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	return s.repo.GetSubsByID(ctx, id)
}

func (s *ServiceSubs) Update(ctx context.Context, id int, req dto.UpdateSubRequest) error {
	sub := models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate.Time,
	}

	if req.EndDate != nil {
		sub.EndDate = &req.EndDate.Time
	}

	return s.repo.UpdateSubs(ctx, id, sub)
}

func (s *ServiceSubs) Delete(ctx context.Context, id int) error {
	return s.repo.DeleteSubs(ctx, id)
}

func (s *ServiceSubs) List(ctx context.Context, req dto.ListFilterRequest) ([]models.Subscription, error) {
	layout := "01-2006"
	filter := models.ListFilter{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
	}

	if req.FromDate != nil {
		t, err := time.Parse(layout, *req.FromDate)
		if err != nil {
			return nil, fmt.Errorf("неверный формат from_date: %w", err)
		}
		normalized := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
		filter.FromDate = &normalized
	}

	if req.ToDate != nil {
		t, err := time.Parse(layout, *req.ToDate)
		if err != nil {
			return nil, fmt.Errorf("неверный формат to_date: %w", err)
		}
		normalized := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
		filter.ToDate = &normalized
	}

	return s.repo.ListSubs(ctx, filter)
}

func (s *ServiceSubs) GetTotalCost(ctx context.Context, req dto.TotalCostRequest) (int, error) {
	layout := "01-2006"

	fromTime, err := time.Parse(layout, req.FromDate)
	if err != nil {
		return 0, fmt.Errorf("неверный формат from_date: %w", err)
	}

	toTime, err := time.Parse(layout, req.ToDate)
	if err != nil {
		return 0, fmt.Errorf("неверный формат to_date: %w", err)
	}

	// Нормализуем даты
	fromTime = time.Date(fromTime.Year(), fromTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	toTime = time.Date(toTime.Year(), toTime.Month(), 1, 0, 0, 0, 0, time.UTC)

	filter := models.TotalCostFilter{
		ServiceName: req.ServiceName,
		FromDate:    &fromTime,
		ToDate:      &toTime,
	}

	return s.repo.TotalCost(ctx, req.UserID, filter)
}
