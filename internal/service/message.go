package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/pubsub"
	"github.com/ahernandez9/rockets/internal/repository"
)

//go:generate go run go.uber.org/mock/mockgen -source=message.go -destination=mocks/mock_message_service.go -package=mocks

type MessageService interface {
	Start()
	Stop()
	PublishMessage(msg *models.RocketMessage) error
}

// messageService handles async message processing via pub/sub
type messageService struct {
	pubsub pubsub.Interface
	repo   repository.RocketRepository
	ctx    context.Context
	cancel context.CancelFunc
}

// NewMessageService creates a new message service
func NewMessageService(ps pubsub.Interface, r repository.RocketRepository) MessageService {
	ctx, cancel := context.WithCancel(context.Background())

	return &messageService{
		pubsub: ps,
		repo:   r,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins processing messages
func (s *messageService) Start() {
	log.Println("MessageService: Started message processor")

	if err := s.pubsub.Subscribe(s.ctx, s.handleMessage); err != nil {
		log.Printf("MessageService: Subscriber stopped: %v", err)
	}

	log.Println("MessageService: Message processor stopped")
}

// Stop gracefully stops the message service
func (s *messageService) Stop() {
	log.Println("MessageService: Stopping")
	s.cancel()
	s.pubsub.Close()
}

// PublishMessage publishes a message for async processing
func (s *messageService) PublishMessage(msg *models.RocketMessage) error {
	return s.pubsub.Publish(s.ctx, msg)
}

// handleMessage processes a single message (callback from subscriber)
// In a production scenario, would implement retry logic with exponential backoff for consistency
func (s *messageService) handleMessage(ctx context.Context, msg *models.RocketMessage) error {
	channelID := msg.Metadata.Channel

	existingRocket, _ := s.repo.FindByID(ctx, channelID)

	// Check for duplicates/out-of-order
	if existingRocket != nil && msg.Metadata.MessageNumber <= existingRocket.LastMessageNumber {
		log.Printf("MessageService: Ignoring old/duplicate message: channel=%s, msgNum=%d, lastProcessed=%d",
			channelID, msg.Metadata.MessageNumber, existingRocket.LastMessageNumber)
		return nil
	}

	switch msg.Metadata.MessageType {
	case "RocketLaunched":
		return s.handleRocketLaunched(ctx, channelID, msg)
	case "RocketSpeedIncreased", "RocketSpeedDecreased":
		return s.handleRocketSpeedChanged(ctx, channelID, msg)
	case "RocketExploded":
		return s.handleRocketExploded(ctx, channelID, msg)
	case "RocketMissionChanged":
		return s.handleRocketMissionChanged(ctx, channelID, msg)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Metadata.MessageType)
	}
}

// parseMessage helper function to parse and type assert message payload
func parseMessage[T any](msg *models.RocketMessage) (T, error) {
	var result T

	// msg.Message is interface{}, need to marshal/unmarshal
	// Type assertion won't work because Gin unmarshals to map[string]interface{}
	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return result, fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := json.Unmarshal(msgBytes, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return result, nil
}

// updateRocketMetadata is a helper to update common metadata fields
func updateRocketMetadata(rocket *models.Rocket, msg *models.RocketMessage) {
	rocket.LastMessageNumber = msg.Metadata.MessageNumber
	rocket.LastUpdated = msg.Metadata.MessageTime
}

func (s *messageService) handleRocketLaunched(ctx context.Context, channelID string, msg *models.RocketMessage) error {
	launchMsg, err := parseMessage[models.RocketLaunchedMessage](msg)
	if err != nil {
		return err
	}

	rocket := &models.Rocket{
		ID:      channelID,
		Type:    launchMsg.Type,
		Speed:   launchMsg.LaunchSpeed,
		Mission: launchMsg.Mission,
		Status:  models.StatusActive,
	}
	updateRocketMetadata(rocket, msg)

	log.Printf("MessageService: Rocket launched: %s (type=%s, speed=%d, mission=%s)",
		channelID, rocket.Type, rocket.Speed, rocket.Mission)

	return s.repo.Save(ctx, rocket)
}

func (s *messageService) handleRocketSpeedChanged(ctx context.Context, channelID string, msg *models.RocketMessage) error {
	rocket, err := s.repo.FindByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	speedMsg, err := parseMessage[models.RocketSpeedChangedMessage](msg)
	if err != nil {
		return err
	}

	// Apply speed change based on message type
	if msg.Metadata.MessageType == "RocketSpeedIncreased" {
		rocket.Speed += speedMsg.By
	} else {
		rocket.Speed -= speedMsg.By
	}
	updateRocketMetadata(rocket, msg)

	log.Printf("MessageService: Speed changed: %s (type=%s, by=%d, new speed=%d)",
		channelID, msg.Metadata.MessageType, speedMsg.By, rocket.Speed)

	return s.repo.Save(ctx, rocket)
}

func (s *messageService) handleRocketExploded(ctx context.Context, channelID string, msg *models.RocketMessage) error {
	rocket, err := s.repo.FindByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	explodedMsg, err := parseMessage[models.RocketExplodedMessage](msg)
	if err != nil {
		return err
	}

	rocket.Status = models.StatusExploded
	rocket.ExplosionReason = explodedMsg.Reason
	rocket.Speed = 0
	updateRocketMetadata(rocket, msg)

	log.Printf("MessageService: Rocket exploded: %s (reason=%s)", channelID, explodedMsg.Reason)

	return s.repo.Save(ctx, rocket)
}

func (s *messageService) handleRocketMissionChanged(ctx context.Context, channelID string, msg *models.RocketMessage) error {
	rocket, err := s.repo.FindByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("rocket not found: %s", channelID)
	}

	missionMsg, err := parseMessage[models.RocketMissionChangedMessage](msg)
	if err != nil {
		return err
	}

	rocket.Mission = missionMsg.NewMission
	updateRocketMetadata(rocket, msg)

	log.Printf("MessageService: Mission changed: %s (new mission=%s)", channelID, missionMsg.NewMission)

	return s.repo.Save(ctx, rocket)
}
