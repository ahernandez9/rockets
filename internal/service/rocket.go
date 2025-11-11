package service

import (
	"context"
	"sort"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/repository"
)

// RocketService handles rocket business logic and repository operations
type RocketService struct {
	repo repository.RocketRepository
}

// NewRocketService creates a new rocket service
func NewRocketService(repo repository.RocketRepository) *RocketService {
	return &RocketService{
		repo: repo,
	}
}

// GetRocket retrieves a rocket by ID
func (s *RocketService) GetRocket(ctx context.Context, id string) (*models.Rocket, error) {
	return s.repo.FindByID(ctx, id)
}

// ListRockets retrieves all rockets with optional sorting
func (s *RocketService) ListRockets(ctx context.Context, sortBy string) ([]*models.Rocket, error) {
	rockets := s.repo.FindAll(ctx)

	switch sortBy {
	case "type":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Type < rockets[j].Type
		})
	case "speed":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Speed < rockets[j].Speed
		})
	case "mission":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Mission < rockets[j].Mission
		})
	case "status":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Status < rockets[j].Status
		})
	}

	return rockets, nil
}

// UpdateRocket updates or creates a rocket
func (s *RocketService) UpdateRocket(ctx context.Context, rocket *models.Rocket) error {
	return s.repo.Save(ctx, rocket)
}

// GetCount returns the number of rockets
func (s *RocketService) GetCount(ctx context.Context) int {
	return s.repo.GetCount(ctx)
}
