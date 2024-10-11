package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/uuid"
	"github.com/voidshard/faction/pkg/structs"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitRetryHeader = "x-retries"
	rabbitRetryMax    = 5
)

type rabbitMessage struct {
	msg     amqp.Delivery
	channel *amqp.Channel

	// parsed context
	context context.Context

	// remaining data
	data []byte
}

func NewRabbitMessage(m amqp.Delivery) (*rabbitMessage, error) {
	parentCtx, msgdata, err := extractTraceData(m.Body)
	if err != nil {
		return nil, err
	}
	return &rabbitMessage{context: parentCtx, data: msgdata, msg: m}, nil
}

func (m *rabbitMessage) setReplyChannel(channel *amqp.Channel) {
	m.channel = channel
}

func (m *rabbitMessage) Reply(ctx context.Context, data []byte) error {
	if m.channel == nil {
		return fmt.Errorf("no channel to reply on, only supported on API Request message replies")
	}
	msgdata, err := injectTraceData(ctx, data)
	if err != nil {
		return err
	}
	log.Debug().Str("MessageId", m.msg.MessageId).Int("bytes", len(msgdata)).Msg("Injected telemetry, replying to message")
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

func (m *rabbitMessage) Reject() error {
	retries, ok := m.msg.Headers[rabbitRetryHeader]
	if !ok {
		retries = 0
	}
	attempt := retries.(int)
	if attempt > rabbitRetryMax {
		log.Debug().Int("attempt", attempt).Str("MessageId", m.msg.MessageId).Msg("Message has exceeded retry limit, rejecting")
		return m.msg.Reject(false)
	}
	attempt += 1
	log.Debug().Int("attempt", attempt).Str("MessageId", m.msg.MessageId).Msg("Message rejected, requeueing")
	m.msg.Headers["x-retries"] = attempt
	return m.msg.Reject(true)
}
