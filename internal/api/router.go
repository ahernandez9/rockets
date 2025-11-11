package api

import (
	"github.com/ahernandez9/rockets/internal/handler"
	"github.com/ahernandez9/rockets/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin router with explicit dependency injection
func SetupRouter(messageService *service.MessageService, rocketService *service.RocketService) *gin.Engine {
	router := gin.Default()

	router.GET("/health", handler.Healthcheck())

	// Message ingestion uses MessageService
	router.POST("/messages", handler.PostMessage(messageService))

	// Queries use RocketService
	router.GET("/rockets", handler.ListRockets(rocketService))
	router.GET("/rockets/:id", handler.GetRocket(rocketService))

	return router
}
