package channel

import (
	"context"
	"fmt"
	"log"

	"github.com/ahernandez9/rockets/internal/models"
	"github.com/ahernandez9/rockets/internal/pubsub"
)

// PubSub implements PubSub using Go channels
type PubSub struct {
	messageChan chan *models.RocketMessage
	closed      bool
}

// NewPubSub creates a new channel-based pub/sub
func NewPubSub(bufferSize int) *PubSub {
	return &PubSub{
		messageChan: make(chan *models.RocketMessage, bufferSize),
	}
}

// Publish sends a message to the channel
func (p *PubSub) Publish(ctx context.Context, msg *models.RocketMessage) error {
	select {
	case p.messageChan <- msg:
		log.Printf("Message published: channel=%s, type=%s, number=%d",
			msg.Metadata.Channel, msg.Metadata.MessageType, msg.Metadata.MessageNumber)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		log.Printf("Warning: message channel full, dropping message: channel=%s", msg.Metadata.Channel)
		return fmt.Errorf("channel full")
		// trade-off: we don't want to block HTTP handlers (bad UX) nor store overflow messages in memory (dangerous)
		// for a Production ready system, consider using a persistent message broker like RabbitMQ, or Redis Streams
	}
}

// Subscribe starts listening for messages and calls handler for each message
func (p *PubSub) Subscribe(ctx context.Context, handler pubsub.MessageHandler) error {
	for {
		select {
		case msg, ok := <-p.messageChan:
			if !ok {
				log.Println("PubSub: Channel closed")
				return nil
			}
			if err := handler(ctx, msg); err != nil {
				log.Printf("PubSub: Error handling message: %v", err)
			}
		case <-ctx.Done():
			log.Println("PubSub: Context canceled")
			return ctx.Err()
		}
	}
}

// Close closes the pub/sub channel
func (p *PubSub) Close() error {
	if !p.closed {
		p.closed = true
		close(p.messageChan)
	}
	return nil
}
