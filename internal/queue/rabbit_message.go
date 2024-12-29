package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitRetryHeader = "x-retries"
	rabbitRetryMax    = 5
)

type rabbitMessage struct {
	msg     amqp.Delivery
	channel *rabbitChannel

	// parsed context
	context context.Context

	// remaining data
	data []byte
}

func newRabbitMessage(m amqp.Delivery) (*rabbitMessage, error) {
	parentCtx, msgdata, err := log.ExtractTraceData(m.Body)
	if err != nil {
		return nil, err
	}
	return &rabbitMessage{context: parentCtx, data: msgdata, msg: m}, nil
}

func (m *rabbitMessage) setReplyChannel(channel *rabbitChannel) {
	m.channel = channel
}

func (m *rabbitMessage) Reply(ctx context.Context, data []byte) error {
	if m.channel == nil {
		return fmt.Errorf("no channel to reply on, only supported on API Request message replies")
	}
	msgdata, err := log.InjectTraceData(ctx, data)
	if err != nil {
		return err
	}
	log.Debug().Str("ReplyTo", m.msg.ReplyTo).Str("MessageId", m.msg.MessageId).Int("bytes", len(msgdata)).Msg("Injected telemetry, replying to message")
	return m.channel.Channel().PublishWithContext(ctx,
		"",            // exchange
		m.msg.ReplyTo, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			MessageId:     uuid.NewID().String(),
			CorrelationId: m.msg.CorrelationId,
			Timestamp:     time.Now(),
			ContentType:   "text/plain",
			Body:          msgdata,
		},
	)
}

func (m *rabbitMessage) Id() string {
	return m.msg.MessageId
}

func (m *rabbitMessage) CorrelationId() string {
	return m.msg.CorrelationId
}

func (m *rabbitMessage) Timestamp() time.Time {
	return m.msg.Timestamp
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
	log.Debug().Str("MessageId", m.msg.MessageId).Msg("Message rejected, requeueing")
	return m.msg.Reject(true)
}
