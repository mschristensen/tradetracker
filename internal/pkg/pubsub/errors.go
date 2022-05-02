package pubsub

import "errors"

// ErrTopicAlreadySubscribed indicates that a topic is already subscribed by another consumer.
var ErrTopicAlreadySubscribed error = errors.New("already subscribed")

// ErrTopicClosed indicates that a topic has been closed.
var ErrTopicClosed error = errors.New("topic closed")
