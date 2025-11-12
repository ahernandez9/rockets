package service

import (
	"context"
	"sort"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/repository"
)

//go:generate go run go.uber.org/mock/mockgen -source=rocket.go -destination=mocks/mock_rocket_service.go -package=mocks

// RocketService defines the methods for rocket service (mostly to ease mocking in tests)
type RocketService interface {
	GetRocket(ctx context.Context, id string) (*models.Rocket, error)
	ListRockets(ctx context.Context, sortBy string) ([]*models.Rocket, error)
	UpdateRocket(ctx context.Context, rocket *models.Rocket) error
	GetCount(ctx context.Context) int
}

// RocketService handles rocket business logic and repository operations
type rocketService struct {
	repo repository.RocketRepository
}

// NewRocketService creates a new rocket service
func NewRocketService(repo repository.RocketRepository) RocketService {
	return &rocketService{
		repo: repo,
	}
}

// GetRocket retrieves a rocket by ID
func (s *rocketService) GetRocket(ctx context.Context, id string) (*models.Rocket, error) {
	return s.repo.FindByID(ctx, id)
}

// ListRockets retrieves all rockets with optional sorting
func (s *rocketService) ListRockets(ctx context.Context, sortBy string) ([]*models.Rocket, error) {
	rockets := s.repo.FindAll(ctx)

	switch sortBy {
	case "type":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Type < rockets[j].Type
		})
	case "speed":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Speed > rockets[j].Speed
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
func (s *rocketService) UpdateRocket(ctx context.Context, rocket *models.Rocket) error {
	return s.repo.Save(ctx, rocket)
}

// GetCount returns the number of rockets
func (s *rocketService) GetCount(ctx context.Context) int {
	return s.repo.GetCount(ctx)
}
