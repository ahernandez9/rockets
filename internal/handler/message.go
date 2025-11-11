package handler

import (
	"net/http"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/service"

	"github.com/gin-gonic/gin"
)

// PostMessage godoc
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
func PostMessage(ms *service.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg models.RocketMessage

		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid request body",
				Message: "The request body must be valid JSON matching the RocketMessage schema",
			})
			return
		}

		if err := validateMessageMetadata(msg.Metadata); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid message metadata",
				Message: err.Error(),
			})
			return
		}

		if err := validateMessageContent(&msg); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid message content",
				Message: err.Error(),
			})
			return
		}

		if err := ms.PublishMessage(&msg); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to publish message",
				Message: "The message could not be queued for processing. Please try again.",
			})
			return
		}

		// Return immediately while message is processed in background
		c.JSON(http.StatusAccepted, gin.H{
			"status":  "accepted",
			"message": "Message queued for processing",
		})
	}
}
