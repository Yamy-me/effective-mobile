package handler

import "github.com/gin-gonic/gin"

func InitRoutes(server *gin.Engine, handler *Handler) {
	subs := server.Group("/subscriptions")
	{
		subs.POST("/", handler.Create)
		subs.GET("/:id", handler.GetByID)
		subs.GET("/", handler.List)
		subs.PUT("/:id", handler.Update)
		subs.DELETE("/:id", handler.Delete)
		subs.GET("/total", handler.GetTotalCost)
	}
}
