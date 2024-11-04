package queue

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/voidshard/faction/pkg/util/log"
)

type rabbitSubscription struct {
	kill    chan bool    // signal to stop the subscription
	out     chan Message // channel a user will read messages from
	created time.Time
	log     log.Logger
}

func newRabbitSubscription(name string, attrs ...map[string]interface{}) *rabbitSubscription {
	return &rabbitSubscription{
		out:     make(chan Message),
		created: time.Now(),
		log:     log.Sublogger(name, attrs...),
	}
}

func (rs *rabbitSubscription) consumeDeliveries(ch *rabbitChannel, deliveries <-chan amqp.Delivery) {
	rs.kill = make(chan bool)
	for {
		select {
		case <-rs.kill:
			rs.log.Debug().Str("Channel", ch.id).Msg("Rabbit subscription consuming deliveries killed")
			return
		case d, ok := <-deliveries:
			if !ok {
				continue
				//rs.log.Debug().Str("Channel", ch.id).Msg("Rabbit subscription consuming deliveries: delivery channel closed")
				//rs.Close()
				//return
			}
			msg, err := newRabbitMessage(d)
			if err != nil {
				rs.log.Warn().Err(err).Str("Channel", ch.id).Msg("Rabbit subscription failed to create message from delivery")
				continue
			}
			msg.setReplyChannel(ch)
			rs.out <- msg
		}
	}
}

func (rs *rabbitSubscription) Channel() <-chan Message {
	return rs.out
}

func (rs *rabbitSubscription) Close() error {
	rs.log.Debug().Msg("Rabbit subscription closing")
	if rs.kill != nil {
		rs.kill <- true
		close(rs.kill)
	}
	close(rs.out)
	return nil
}
