package repository

import (
	"github.com/ahernandez9/rockets/internal/models"
)

// RocketRepository defines the interface for rocket storage
type RocketRepository interface {
	Save(rocket *models.Rocket) error
	FindByID(id string) (*models.Rocket, error)
	FindAll() []*models.Rocket
	GetCount() int
}
