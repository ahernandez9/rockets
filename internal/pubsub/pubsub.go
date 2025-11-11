package pubsub

import (
	"log"

	"github.com/ahernandez9/rockets/internal/models"
)

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(msg models.RocketMessage) error
	Close()
}

// Subscriber defines the interface for subscribing to messages
type Subscriber interface {
	Subscribe() <-chan models.RocketMessage
	Close()
}

// PubSub combines Publisher and Subscriber
type PubSub interface {
	Publisher
	Subscriber
}

// ChannelPubSub implements PubSub using Go channels
type ChannelPubSub struct {
	messageChan chan models.RocketMessage
	closed      bool
}

// NewChannelPubSub creates a new channel-based pub/sub
func NewChannelPubSub(bufferSize int) *ChannelPubSub {
	return &ChannelPubSub{
		messageChan: make(chan models.RocketMessage, bufferSize),
	}
}

// Publish sends a message to the channel
func (p *ChannelPubSub) Publish(msg models.RocketMessage) error {
	if p.closed {
		return nil // Silently ignore if closed
	}

	select {
	case p.messageChan <- msg:
		log.Printf("Message published: channel=%s, type=%s, number=%d",
			msg.Metadata.Channel, msg.Metadata.MessageType, msg.Metadata.MessageNumber)
		return nil
	default:
		log.Printf("Warning: message channel full, dropping message: channel=%s", msg.Metadata.Channel)
		return nil // Non-blocking, drop if full
	}
}

// Subscribe returns the receive-only channel for consuming messages
func (p *ChannelPubSub) Subscribe() <-chan models.RocketMessage {
	return p.messageChan
}

// Close closes the pub/sub channel
func (p *ChannelPubSub) Close() {
	if !p.closed {
		p.closed = true
		close(p.messageChan)
	}
}
