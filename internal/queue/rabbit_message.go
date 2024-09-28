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
}

func (m *rabbitMessage) Reply(data []byte) error {
	if m.channel == nil {
		return fmt.Errorf("no channel to reply on, only supported on API Request message replies")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
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
			Body:          data,
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

func (m *rabbitMessage) Data() []byte {
	return m.msg.Body
}

func (m *rabbitMessage) Ack() error {
	return m.msg.Ack(false)
}

func (m *rabbitMessage) Reject(requeue bool) error {
	return m.msg.Reject(requeue)
}
