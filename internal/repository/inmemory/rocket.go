package inmemory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ahernandez9/rockets/internal/models"
)

// RocketRepository implements Repository with in-memory storage
type RocketRepository struct {
	rockets map[string]*models.Rocket
	mu      sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *RocketRepository {
	return &RocketRepository{
		rockets: make(map[string]*models.Rocket),
	}
}

// Save stores or updates a rocket
func (r *RocketRepository) Save(ctx context.Context, rocket *models.Rocket) error {
	if rocket == nil {
		return fmt.Errorf("cannot save nil rocket")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.rockets[rocket.ID] = rocket
	return nil
}

// FindByID retrieves a rocket by ID
func (r *RocketRepository) FindByID(ctx context.Context, id string) (*models.Rocket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rocket, exists := r.rockets[id]
	if !exists {
		return nil, fmt.Errorf("rocket not found: %s", id)
	}

	// Return a copy to prevent external modifications
	rocketCopy := *rocket
	return &rocketCopy, nil
}

// FindAll retrieves all rockets
func (r *RocketRepository) FindAll(ctx context.Context) []*models.Rocket {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rockets := make([]*models.Rocket, 0, len(r.rockets))
	for _, rocket := range r.rockets {
		rocketCopy := *rocket
		rockets = append(rockets, &rocketCopy)
	}

	// Default sort by ID
	sort.Slice(rockets, func(i, j int) bool {
		return rockets[i].ID < rockets[j].ID
	})

	return rockets
}

// GetCount returns the total number of rockets
func (r *RocketRepository) GetCount(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.rockets)
}
