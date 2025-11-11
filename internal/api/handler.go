package api

import (
	"net/http"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler contains HTTP handlers for the API
type Handler struct {
	rocketService *service.RocketService
}

// NewHandler creates a new Handler instance
func NewHandler(rocketService *service.RocketService) *Handler {
	return &Handler{
		rocketService: rocketService,
	}
}

// ReceiveMessage godoc
// @Summary Receive rocket telemetry message
// @Description Accepts rocket telemetry messages from the test program and publishes them asynchronously
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.RocketMessage true "Rocket message"
// @Success 202 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /messages [post]
func (h *Handler) ReceiveMessage(c *gin.Context) {
	var msg models.RocketMessage

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Publish the message asynchronously (non-blocking)
	if err := h.rocketService.PublishMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to publish message",
			Message: err.Error(),
		})
		return
	}

	// Return immediately while message is processed in background
	c.JSON(http.StatusAccepted, gin.H{
		"status": "accepted",
	})
}

// GetRocket godoc
// @Summary Get rocket by ID
// @Description Retrieves the current state of a specific rocket
// @Tags rockets
// @Produce json
// @Param id path string true "Rocket ID (channel)"
// @Success 200 {object} models.Rocket
// @Failure 404 {object} models.ErrorResponse
// @Router /rockets/{id} [get]
func (h *Handler) GetRocket(c *gin.Context) {
	id := c.Param("id")

	rocket, err := h.rocketService.GetRocket(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Rocket not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rocket)
}

// ListRockets godoc
// @Summary List all rockets
// @Description Retrieves a list of all rockets in the system with optional sorting
// @Tags rockets
// @Produce json
// @Param sort query string false "Sort by field (type, speed, mission, status)" default(id)
// @Success 200 {array} models.Rocket
// @Router /rockets [get]
func (h *Handler) ListRockets(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "id")

	rockets := h.rocketService.ListRockets(sortBy)

	c.JSON(http.StatusOK, gin.H{
		"count":   len(rockets),
		"rockets": rockets,
	})
}

// HealthCheck godoc
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:  "ok",
		Service: "rockets",
	})
}
