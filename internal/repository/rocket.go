package repository

import (
	"context"

	"github.com/ahernandez9/rockets/internal/models"
)

// RocketRepository defines the interface for rocket storage
type RocketRepository interface {
	Save(ctx context.Context, rocket *models.Rocket) error
	FindByID(ctx context.Context, id string) (*models.Rocket, error)
	FindAll(ctx context.Context) []*models.Rocket
	GetCount(ctx context.Context) int
}
