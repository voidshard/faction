package queue

import (
	"context"
	"time"
)

type Queue interface {
	// Request creates a new request on the given queue and waits for a reply
	Request(ctx context.Context, queue string, data []byte) (Subscription, error)

	// Enqueue adds a message to the queue, but does not wait nor expect a reply.
	Enqueue(ctx context.Context, queue string, data []byte) error

	// Dequeue returns a subscription to the given queue.
	Dequeue(queue string) (Subscription, error)

	// Publish sends a message to the given topic, these messages are fan-out to all subscribers.
	// The key accepts empty strings to match all keys at the given key index
	Publish(ctx context.Context, topic string, key []string, data []byte) error

	// Subscribe returns a subscription to the given topic matching the given key.
	// The key accepts empty strings to match all keys at the given key index
	Subscribe(queue, topic string, key []string, durable bool) (Subscription, error)

	// DeleteQueue deletes the given queue.
	DeleteQueue(queue string) error

	// Close all connections / channels ready for shutdown.
	Close() error
}

type closer interface {
	Close() error
}

type Subscription interface {
	Channel() <-chan Message
	Close() error
}

type Message interface {
	// Id returns the message id.
	Id() string

	// CorrelationId returns the message correlation id.
	CorrelationId() string

	// Reply sends a reply to the message.
	// Only supported on messages delivered via Dequeue
	Reply(context.Context, []byte) error

	// Subject returns the subject of the message.
	Subject() string

	// Data returns the data of the message
	Data() []byte

	// Ack acknowledges the message & removes it from the queue.
	// Must be called on successful processing of the message.
	Ack() error

	// Reject rejects the message, will be requeued
	Reject() error

	// Timestamp returns the time the message was sent.
	Timestamp() time.Time

	// Context returns the context of the message.
	Context() context.Context
}
