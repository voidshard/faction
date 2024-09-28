package queue

import (
	"time"

	"github.com/voidshard/faction/pkg/structs"
)

type Queue interface {
	// PublishApiReq publishes data to the API request stream.
	PublishApiReq(data []byte) (Subscription, error)

	// SubscribeApiReq subscribes to the API request stream.
	SubscribeApiReq() (Subscription, error)

	// PublishChange publishes a change to the change stream.
	PublishChange(ch *structs.Change) error

	// SubscribeChange subscribes to changes on the change stream.
	// queueName can be given to configure a durable queue. If not given
	// a temporary non-durable queue will be used.
	SubscribeChange(ch *structs.Change, queueName string) (Subscription, error)
}

type Closer interface {
	Close() error
}

type Subscription interface {
	Channel() <-chan Message
	Close()
}

type Message interface {
	Id() string

	// Reply sends a reply to the message.
	// This is only supported on messages from the API request stream.
	Reply([]byte) error

	// Change returns the change that triggered the message.
	// This is only set on messages from a change subscription.
	Change() (*structs.Change, error)

	// Subject returns the subject of the message.
	Subject() string

	// Data returns the data of the message
	Data() []byte

	// Ack acknowledges the message & removes it from the queue.
	// Must be called on successful processing of the message.
	Ack() error

	// Reject rejects the message & requeues it.
	Reject(bool) error

	// Timestamp returns the time the message was sent.
	Timestamp() time.Time
}
