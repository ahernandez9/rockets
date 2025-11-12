package handler

import (
	"net/http"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetRocket godoc
// @Summary Get rocket by ID
// @Description Retrieves the current state of a specific rocket
// @Tags rockets
// @Produce json
// @Param id path string true "Rocket ID (UUID)"
// @Success 200 {object} models.Rocket
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /rockets/{id} [get]
func GetRocket(rs service.RocketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid rocket ID",
				Message: "The rocket ID must be a valid UUID (e.g., 193270a9-c9cf-404a-8f83-838e71d9ae67)",
			})
			return
		}

		rocket, err := rs.GetRocket(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Rocket not found",
				Message: "No rocket exists with the provided ID. It may not have been launched yet.",
			})
			return
		}

		c.JSON(http.StatusOK, rocket)
	}
}

// ListRockets godoc
// @Summary List all rockets
// @Description Retrieves a list of all rockets in the system with optional sorting
// @Tags rockets
// @Produce json
// @Param sort query string false "Sort by field (type, speed, mission, status)" default(id)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Router /rockets [get]
func ListRockets(rs service.RocketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sortBy := c.DefaultQuery("sort", "id")

		validSortFields := map[string]bool{
			"id":      true,
			"type":    true,
			"speed":   true,
			"mission": true,
			"status":  true,
		}

		if !validSortFields[sortBy] {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid sort parameter",
				Message: "Sort parameter must be one of: id, type, speed, mission, status",
			})
			return
		}

		var rockets []*models.Rocket
		var err error
		if rockets, err = rs.ListRockets(c.Request.Context(), sortBy); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Failed to retrieve rockets",
				Message: "An error occurred while fetching the list of rockets. Please try again later.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"count":   len(rockets),
			"rockets": rockets,
			"sortBy":  sortBy,
		})
	}
}
