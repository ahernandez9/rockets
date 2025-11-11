package models

import "time"

// MessageMetadata contains metadata about the rocket message
type MessageMetadata struct {
	Channel       string    `json:"channel" example:"193270a9-c9cf-404a-8f83-838e71d9ae67"`
	MessageNumber int64     `json:"messageNumber" example:"1"`
	MessageTime   time.Time `json:"messageTime" example:"2022-02-02T19:39:05.86337+01:00"`
	MessageType   string    `json:"messageType" example:"RocketLaunched"`
}

// RocketMessage represents an incoming rocket message
type RocketMessage struct {
	Metadata MessageMetadata `json:"metadata"`
	Message  interface{}     `json:"message"`
}

// RocketLaunchedMessage represents a rocket launch event
type RocketLaunchedMessage struct {
	Type        string `json:"type" example:"Falcon-9"`
	LaunchSpeed int    `json:"launchSpeed" example:"500"`
	Mission     string `json:"mission" example:"ARTEMIS"`
}

// RocketSpeedChangedMessage represents a speed change event (increase or decrease)
type RocketSpeedChangedMessage struct {
	By int `json:"by" example:"3000"` // Always positive; message type determines if added or subtracted
}

// RocketExplodedMessage represents a rocket explosion event
type RocketExplodedMessage struct {
	Reason string `json:"reason" example:"PRESSURE_VESSEL_FAILURE"`
}

// RocketMissionChangedMessage represents a mission change event
type RocketMissionChangedMessage struct {
	NewMission string `json:"newMission" example:"SHUTTLE_MIR"`
}

// RocketStatus represents the possible states of a rocket
type RocketStatus string

const (
	StatusActive   RocketStatus = "ACTIVE"
	StatusExploded RocketStatus = "EXPLODED"
)

// Rocket represents the current state of a rocket
type Rocket struct {
	ID                string       `json:"id" example:"193270a9-c9cf-404a-8f83-838e71d9ae67"`
	Type              string       `json:"type" example:"Falcon-9"`
	Speed             int          `json:"speed" example:"3500"`
	Mission           string       `json:"mission" example:"ARTEMIS"`
	Status            RocketStatus `json:"status" example:"ACTIVE"`
	ExplosionReason   string       `json:"explosionReason,omitempty" example:"PRESSURE_VESSEL_FAILURE"`
	LastMessageNumber int64        `json:"lastMessageNumber" example:"42"`
	LastUpdated       time.Time    `json:"lastUpdated" example:"2022-02-02T19:39:05.86337+01:00"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid message format"`
	Message string `json:"message,omitempty" example:"The provided message could not be parsed"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"rockets"`
}
