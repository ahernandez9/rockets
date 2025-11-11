package handler

import (
	"net/http"

	"github.com/ahernandez9/rockets/internal/models"

	"github.com/gin-gonic/gin"
)

// Healthcheck godoc
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func Healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, models.HealthResponse{
			Status:  "ok",
			Service: "rockets",
		})
	}
}
