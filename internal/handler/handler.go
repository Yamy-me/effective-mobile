package handler

import (
	"log/slog"
	"strconv"

	"Effective-Mobile/internal/dto"
	"Effective-Mobile/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.ServiceSubs
}

func NewHandler(service *service.ServiceSubs) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateSubRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("ошибка валидации JSON", slog.String("error", err.Error()))
		c.JSON(400, gin.H{"error": "неверные данные запроса: " + err.Error()})
		return
	}

	id, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": "не удалось создать подписку"})
		return
	}

	c.JSON(201, gin.H{"id": id})
}

func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID должен быть числом"})
		return
	}

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	res := dto.SubResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   dto.MonthYear{Time: sub.StartDate},
	}

	if sub.EndDate != nil {
		res.EndDate = &dto.MonthYear{Time: *sub.EndDate}
	}

	c.JSON(200, res)
}

func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID должен быть числом"})
		return
	}

	var req dto.UpdateSubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("ошибка валидации JSON при обновлении", slog.String("error", err.Error()))
		c.JSON(400, gin.H{"error": "неверные данные запроса: " + err.Error()})
		return
	}

	err = h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(500, gin.H{"error": "не удалось обновить подписку: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "успешно обновлено"})
}

func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID должен быть числом"})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "не удалось удалить подписку: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "успешно удалено"})
}

func (h *Handler) List(c *gin.Context) {
	var req dto.ListFilterRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Error("ошибка валидации параметров фильтра списка", slog.String("error", err.Error()))
		c.JSON(400, gin.H{"error": "неверные параметры фильтрации: " + err.Error()})
		return
	}

	subs, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": "не удалось получить список подписок: " + err.Error()})
		return
	}

	res := make([]dto.SubResponse, 0, len(subs))
	for _, sub := range subs {
		item := dto.SubResponse{
			ID:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID,
			StartDate:   dto.MonthYear{Time: sub.StartDate},
		}
		if sub.EndDate != nil {
			item.EndDate = &dto.MonthYear{Time: *sub.EndDate}
		}
		res = append(res, item)
	}

	c.JSON(200, res)
}

func (h *Handler) GetTotalCost(c *gin.Context) {
	var req dto.TotalCostRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Error("ошибка валидации параметров подсчета стоимости", slog.String("error", err.Error()))
		c.JSON(400, gin.H{"error": "неверные параметры запроса: " + err.Error()})
		return
	}

	totalCost, err := h.service.GetTotalCost(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": "не удалось подсчитать стоимость: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"total_cost": totalCost})
}
