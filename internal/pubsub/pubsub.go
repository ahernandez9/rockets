package pubsub

import (
	"context"

	"github.com/ahernandez9/rockets/internal/models"
)

// MessageHandler processes received messages (callback function)
type MessageHandler func(ctx context.Context, msg *models.RocketMessage) error

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(ctx context.Context, msg *models.RocketMessage) error
	Close() error
}

// Subscriber defines the interface for subscribing to messages
type Subscriber interface {
	Subscribe(ctx context.Context, handler MessageHandler) error
	Close() error
}

// Interface combines Publisher and Subscriber
type Interface interface {
	Publisher
	Subscriber
}
