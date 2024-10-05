package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/voidshard/faction/internal/uuid"
	"github.com/voidshard/faction/pkg/structs"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMessage struct {
	msg     amqp.Delivery
	channel *amqp.Channel

	// parsed context
	context context.Context

	// remaining data
	data []byte
}

func (m *rabbitMessage) Reply(ctx context.Context, data []byte) error {
	if m.channel == nil {
		return fmt.Errorf("no channel to reply on, only supported on API Request message replies")
	}
	msgdata, err := injectTraceData(ctx, data)
	if err != nil {
		return err
	}
	return m.channel.PublishWithContext(ctx,
		"",            // exchange
		m.msg.ReplyTo, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			MessageId:     uuid.NewID().String(),
			Timestamp:     time.Now(),
			CorrelationId: m.msg.CorrelationId,
			ContentType:   "text/plain",
			Body:          msgdata,
		},
	)
}

func (m *rabbitMessage) Id() string {
	return m.msg.MessageId
}

func (m *rabbitMessage) Timestamp() time.Time {
	return m.msg.Timestamp
}

func (m *rabbitMessage) Change() (*structs.Change, error) {
	return fromRabbitSubject(m.msg.RoutingKey)
}

func (m *rabbitMessage) Subject() string {
	return m.msg.RoutingKey
}

func (m *rabbitMessage) Context() context.Context {
	return m.context
}

func (m *rabbitMessage) Data() []byte {
	return m.data
}

func (m *rabbitMessage) Ack() error {
	return m.msg.Ack(false)
}

func (m *rabbitMessage) Reject(requeue bool) error {
	return m.msg.Reject(requeue)
}
