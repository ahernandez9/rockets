package handler

import (
	"encoding/json"
	"fmt"

	"github.com/ahernandez9/rockets/internal/models"

	"github.com/google/uuid"
)

// validateMessageMetadata validates the metadata fields
func validateMessageMetadata(metadata models.MessageMetadata) error {
	if _, err := uuid.Parse(metadata.Channel); err != nil {
		return fmt.Errorf("channel must be a valid UUID, got: %s", metadata.Channel)
	}

	if metadata.MessageNumber <= 0 {
		return fmt.Errorf("messageNumber must be positive, got: %d", metadata.MessageNumber)
	}

	if metadata.MessageTime.IsZero() {
		return fmt.Errorf("messageTime is required and cannot be zero")
	}

	if metadata.MessageType == "" {
		return fmt.Errorf("messageType is required and cannot be empty")
	}

	validTypes := map[string]bool{
		"RocketLaunched":       true,
		"RocketSpeedIncreased": true,
		"RocketSpeedDecreased": true,
		"RocketExploded":       true,
		"RocketMissionChanged": true,
	}

	if !validTypes[metadata.MessageType] {
		return fmt.Errorf("messageType must be one of: RocketLaunched, RocketSpeedIncreased, "+
			"RocketSpeedDecreased, RocketExploded, RocketMissionChanged, got: %s", metadata.MessageType)
	}

	return nil
}

// validateMessageContent validates the message content based on type
func validateMessageContent(msg *models.RocketMessage) error {
	if msg.Message == nil {
		return fmt.Errorf("message content is required")
	}

	msgBytes, err := json.Marshal(msg.Message)
	if err != nil {
		return fmt.Errorf("invalid message format: %w", err)
	}

	switch msg.Metadata.MessageType {
	case "RocketLaunched":
		var launchMsg models.RocketLaunchedMessage
		if err := json.Unmarshal(msgBytes, &launchMsg); err != nil {
			return fmt.Errorf("invalid RocketLaunched message: %w", err)
		}
		if launchMsg.Type == "" {
			return fmt.Errorf("RocketLaunched message: 'type' field is required")
		}
		if launchMsg.LaunchSpeed < 0 {
			return fmt.Errorf("RocketLaunched message: 'launchSpeed' must be non-negative")
		}
		if launchMsg.Mission == "" {
			return fmt.Errorf("RocketLaunched message: 'mission' field is required")
		}

	case "RocketSpeedIncreased":
		var speedMsg models.RocketSpeedChangedMessage
		if err := json.Unmarshal(msgBytes, &speedMsg); err != nil {
			return fmt.Errorf("invalid RocketSpeedIncreased message: %w", err)
		}
		if speedMsg.By <= 0 {
			return fmt.Errorf("RocketSpeedIncreased message: 'by' must be positive")
		}

	case "RocketSpeedDecreased":
		var speedMsg models.RocketSpeedChangedMessage
		if err := json.Unmarshal(msgBytes, &speedMsg); err != nil {
			return fmt.Errorf("invalid RocketSpeedDecreased message: %w", err)
		}
		if speedMsg.By <= 0 {
			return fmt.Errorf("RocketSpeedDecreased message: 'by' must be positive (will be subtracted)")
		}

	case "RocketExploded":
		var explodedMsg models.RocketExplodedMessage
		if err := json.Unmarshal(msgBytes, &explodedMsg); err != nil {
			return fmt.Errorf("invalid RocketExploded message: %w", err)
		}
		if explodedMsg.Reason == "" {
			return fmt.Errorf("RocketExploded message: 'reason' field is required")
		}

	case "RocketMissionChanged":
		var missionMsg models.RocketMissionChangedMessage
		if err := json.Unmarshal(msgBytes, &missionMsg); err != nil {
			return fmt.Errorf("invalid RocketMissionChanged message: %w", err)
		}
		if missionMsg.NewMission == "" {
			return fmt.Errorf("RocketMissionChanged message: 'newMission' field is required")
		}
	}

	return nil
}
