// Package pubsub implents a dummy pub sub system to demonstrate how the application
// might integrate with a real system like Kafka.
package pubsub

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// Publisher supports publishing messages to a topic.
type Publisher interface {
	Publish(msg Message) error
}

// Subscriber supports subscribing to messages on a topic.
type Subscriber interface {
	Subscribe(ctx context.Context, topic Topic, handler Handler) error
}

// Closer supports closing a topic.
type Closer interface {
	Close(ctx context.Context, topic Topic) error
}

// PublisherSubscriber supports both publishing and subscribing to messages on a topic,
// as well as closing the topic.
type PublisherSubscriber interface {
	Publisher
	Subscriber
	Closer
}

// Handler is a function that handles a message.
type Handler func(m Message) error

// Message describes the message topic and payload.
type Message struct {
	Topic Topic
	Value interface{}
}

// MemoryPubSub is a simple in-memory PubSub implementation.
// Note that this naïve implementation only supports one consumer per topic.
type MemoryPubSub struct {
	topics map[Topic]chan Message
	mu     sync.Mutex
}

// NewMemoryPubSub creates a new MemoryPubSub.
func NewMemoryPubSub() *MemoryPubSub {
	return &MemoryPubSub{
		topics: make(map[Topic]chan Message),
	}
}

// Publish publishes a message on a topic.
// It returns an error if the topic is closed.
func (s *MemoryPubSub) Publish(msg Message) error {
	if _, ok := s.topics[msg.Topic]; !ok {
		return ErrTopicClosed
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics[msg.Topic] <- msg
	return nil
}

// Subscribe subscribes to messages on a topic.
// It blocks until the topic is closed or the context is cancelled.
func (s *MemoryPubSub) Subscribe(ctx context.Context, topic Topic, handler Handler) error {
	if _, ok := s.topics[topic]; ok {
		return ErrTopicAlreadySubscribed
	}
	s.topics[topic] = make(chan Message)
	for {
		select {
		case c, ok := <-s.topics[topic]:
			if !ok {
				return nil
			}
			if err := handler(c); err != nil {
				return errors.Wrap(err, "handler failed")
			}
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "context cancelled")
		}
	}
}

// Close closes the topic.
func (s *MemoryPubSub) Close(_ context.Context, topic Topic) error {
	if _, ok := s.topics[topic]; !ok {
		return ErrTopicClosed
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.topics[topic])
	return nil
}
