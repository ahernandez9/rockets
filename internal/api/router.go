package api

import (
	"github.com/ahernandez9/rockets/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(rocketService *service.RocketService) *gin.Engine {
	router := gin.Default()

	// Create handler
	handler := NewHandler(rocketService)

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// Message ingestion endpoint
	router.POST("/messages", handler.ReceiveMessage)

	// Rocket endpoints
	router.GET("/rockets", handler.ListRockets)
	router.GET("/rockets/:id", handler.GetRocket)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
