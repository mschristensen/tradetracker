package pubsub

// PubSub supports publishing and subscribing to messages on a topic.
type PubSub interface {
	// Publish publishes a message on a topic.
	Publish(msg Message) error
	// Subscribe subscribes to messages on a topic.
	Subscribe(topic string, handler func(m *Message)) error
}

// Message describes the message topic and payload.
type Message struct {
	Topic string
	Value interface{}
}

// Channel is a channel for publishing messages.
type Channel struct {
	ch chan Message
}

// MemoryPubSub is a simple in-memory PubSub implementation.
// Note that this naïve implementation only supports one consumer per topic.
type MemoryPubSub struct {
	topics map[string]*Channel
}

// NewMemoryPubSub creates a new MemoryPubSub.
func NewMemoryPubSub() *MemoryPubSub {
	return &MemoryPubSub{
		topics: make(map[string]*Channel),
	}
}

func (s *MemoryPubSub) Subscribe(topic string, handler func(m *Message)) error {
	if _, ok := s.topics[topic]; ok {
		return ErrTopicAlreadySubscribed
	}
	s.topics[topic] = &Channel{
		ch: make(chan Message),
	}
	go func() {
		for {
			c := <-s.topics[topic].ch
			handler(&c)
		}
	}()
	return nil
}

func (s *MemoryPubSub) Publish(msg Message) error {
	if _, ok := s.topics[msg.Topic]; !ok {
		return ErrTopicClosed
	}
	s.topics[msg.Topic].ch <- msg
	return nil
}
