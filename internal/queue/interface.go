package queue

import (
	"context"
	"time"

	"github.com/voidshard/faction/pkg/structs"
)

type Queue interface {
	// PublishApiReq publishes data to the API request stream.
	PublishApiReq(c context.Context, data []byte) (Subscription, error)

	// SubscribeApiReq subscribes to the API request stream.
	SubscribeApiReq() (Subscription, error)

	// PublishChange publishes a change to the change stream.
	PublishChange(c context.Context, ch *structs.Change) error

	// SubscribeChange subscribes to changes on the change stream.
	// queueName can be given to configure a durable queue. If not given
	// a temporary non-durable queue will be used.
	SubscribeChange(ch *structs.Change, queueName string) (Subscription, error)

	// DeferChange defers a change to be processed at given tick.
	DeferChange(c context.Context, ch *structs.Change, tick int64) error

	// SubscribeDeferredChanges subscribes to changes that have been deferred to a given tick.
	SubscribeDeferredChanges(world string, tick int64) (Subscription, error)

	// DeleteDeferredChangeQueue deletes the queue for deferred changes at given tick.
	// This should be called when the queue is no longer needed to tidy old queues we're
	// never going to use again.
	DeleteDeferredChangeQueue(world string, tick int64) error
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
	Reply(context.Context, []byte) error

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

	// Reject rejects the message & should requeue unless the message has exceeded it's retry limit.
	Reject() error

	// Timestamp returns the time the message was sent.
	Timestamp() time.Time

	// Context returns the context of the message.
	Context() context.Context
}
