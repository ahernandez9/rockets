package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/pubsub"
	"github.com/ahernandez9/rockets/internal/repository"
)

// RocketService handles business logic for rocket state management
type RocketService struct {
	repo   repository.RocketRepository
	pubsub pubsub.PubSub
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRocketService creates a new instance of RocketService
func NewRocketService(repo repository.RocketRepository, ps pubsub.PubSub) *RocketService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &RocketService{
		repo:   repo,
		pubsub: ps,
		ctx:    ctx,
		cancel: cancel,
	}

	// Start message processor in background goroutine
	go service.processMessages()

	log.Println("RocketService: Message processor started")

	return service
}

// Stop gracefully stops the service
func (s *RocketService) Stop() {
	log.Println("RocketService: Stopping service")
	s.cancel()
	s.pubsub.Close()
}

// PublishMessage publishes a message to the pub/sub channel (async)
func (s *RocketService) PublishMessage(msg models.RocketMessage) error {
	return s.pubsub.Publish(msg)
}

// processMessages runs in a goroutine and processes messages from the channel
func (s *RocketService) processMessages() {
	log.Println("RocketService: Starting to listen for messages")

	messageChan := s.pubsub.Subscribe()

	for {
		select {
		case <-s.ctx.Done():
			log.Println("RocketService: Message processor stopped")
			return
		case msg, ok := <-messageChan:
			if !ok {
				log.Println("RocketService: Message channel closed")
				return
			}
			if err := s.handleMessage(msg); err != nil {
				log.Printf("RocketService: Error processing message: %v", err)
			}
		}
	}
}

// handleMessage processes a single message and updates the repository
func (s *RocketService) handleMessage(msg models.RocketMessage) error {
	channelID := msg.Metadata.Channel

	// Get existing rocket (if any)
	existingRocket, _ := s.repo.FindByID(channelID)

	// Ignore messages with lower message numbers (out of order or duplicates)
	if existingRocket != nil && msg.Metadata.MessageNumber <= existingRocket.LastMessageNumber {
		log.Printf("RocketService: Ignoring old/duplicate message: channel=%s, msgNum=%d, lastProcessed=%d",
			channelID, msg.Metadata.MessageNumber, existingRocket.LastMessageNumber)
		return nil
	}

	// Process message based on type
	switch msg.Metadata.MessageType {
	case "RocketLaunched":
		return s.handleRocketLaunched(channelID, msg)
	case "RocketSpeedIncreased":
		return s.handleRocketSpeedIncreased(channelID, msg)
	case "RocketSpeedDecreased":
		return s.handleRocketSpeedDecreased(channelID, msg)
	case "RocketExploded":
		return s.handleRocketExploded(channelID, msg)
	case "RocketMissionChanged":
		return s.handleRocketMissionChanged(channelID, msg)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Metadata.MessageType)
	}
}

func (s *RocketService) handleRocketLaunched(channelID string, msg models.RocketMessage) error {
	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var launchMsg models.RocketLaunchedMessage
	if err := json.Unmarshal(msgBytes, &launchMsg); err != nil {
		return fmt.Errorf("failed to parse RocketLaunched message: %w", err)
	}

	rocket := &models.Rocket{
		ID:                channelID,
		Type:              launchMsg.Type,
		Speed:             launchMsg.LaunchSpeed,
		Mission:           launchMsg.Mission,
		Status:            models.StatusActive,
		LastMessageNumber: msg.Metadata.MessageNumber,
		LastUpdated:       msg.Metadata.MessageTime,
	}

	log.Printf("RocketService: Rocket launched: %s (type=%s, speed=%d, mission=%s)",
		channelID, rocket.Type, rocket.Speed, rocket.Mission)

	return s.repo.Save(rocket)
}

func (s *RocketService) handleRocketSpeedIncreased(channelID string, msg models.RocketMessage) error {
	rocket, err := s.repo.FindByID(channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var speedMsg models.RocketSpeedIncreasedMessage
	if err := json.Unmarshal(msgBytes, &speedMsg); err != nil {
		return fmt.Errorf("failed to parse RocketSpeedIncreased message: %w", err)
	}

	rocket.Speed += speedMsg.By
	rocket.LastMessageNumber = msg.Metadata.MessageNumber
	rocket.LastUpdated = msg.Metadata.MessageTime

	log.Printf("RocketService: Speed increased: %s (new speed=%d)", channelID, rocket.Speed)

	return s.repo.Save(rocket)
}

func (s *RocketService) handleRocketSpeedDecreased(channelID string, msg models.RocketMessage) error {
	rocket, err := s.repo.FindByID(channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var speedMsg models.RocketSpeedDecreasedMessage
	if err := json.Unmarshal(msgBytes, &speedMsg); err != nil {
		return fmt.Errorf("failed to parse RocketSpeedDecreased message: %w", err)
	}

	rocket.Speed -= speedMsg.By
	rocket.LastMessageNumber = msg.Metadata.MessageNumber
	rocket.LastUpdated = msg.Metadata.MessageTime

	log.Printf("RocketService: Speed decreased: %s (new speed=%d)", channelID, rocket.Speed)

	return s.repo.Save(rocket)
}

func (s *RocketService) handleRocketExploded(channelID string, msg models.RocketMessage) error {
	rocket, err := s.repo.FindByID(channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var explodedMsg models.RocketExplodedMessage
	if err := json.Unmarshal(msgBytes, &explodedMsg); err != nil {
		return fmt.Errorf("failed to parse RocketExploded message: %w", err)
	}

	rocket.Status = models.StatusExploded
	rocket.ExplosionReason = explodedMsg.Reason
	rocket.Speed = 0 // Exploded rockets have no speed
	rocket.LastMessageNumber = msg.Metadata.MessageNumber
	rocket.LastUpdated = msg.Metadata.MessageTime

	log.Printf("RocketService: Rocket exploded: %s (reason=%s)", channelID, explodedMsg.Reason)

	return s.repo.Save(rocket)
}

func (s *RocketService) handleRocketMissionChanged(channelID string, msg models.RocketMessage) error {
	rocket, err := s.repo.FindByID(channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var missionMsg models.RocketMissionChangedMessage
	if err := json.Unmarshal(msgBytes, &missionMsg); err != nil {
		return fmt.Errorf("failed to parse RocketMissionChanged message: %w", err)
	}

	rocket.Mission = missionMsg.NewMission
	rocket.LastMessageNumber = msg.Metadata.MessageNumber
	rocket.LastUpdated = msg.Metadata.MessageTime

	log.Printf("RocketService: Mission changed: %s (new mission=%s)", channelID, missionMsg.NewMission)

	return s.repo.Save(rocket)
}

// GetRocket retrieves a rocket by ID (for HTTP API)
func (s *RocketService) GetRocket(id string) (*models.Rocket, error) {
	return s.repo.FindByID(id)
}

// ListRockets retrieves all rockets with optional sorting (for HTTP API)
func (s *RocketService) ListRockets(sortBy string) []*models.Rocket {
	rockets := s.repo.FindAll()

	// Sort based on parameter
	switch sortBy {
	case "type":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Type < rockets[j].Type
		})
	case "speed":
		sort.Slice(rockets, func(i, j int) bool {
			return rockets[i].Speed > rockets[j].Speed // Descending order
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
	// Default sorting by ID is done in repository

	return rockets
}

// GetRocketCount returns the total number of rockets
func (s *RocketService) GetRocketCount() int {
	return s.repo.GetCount()
}
