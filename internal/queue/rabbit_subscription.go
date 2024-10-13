package queue

import (
	"time"

	"github.com/voidshard/faction/pkg/util/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitSubscription struct {
	out  chan Message
	kill chan bool

	closed   bool
	closeOut bool

	created time.Time
}

func newRabbitTopicSubscription() *rabbitSubscription {
	return &rabbitSubscription{created: time.Now(), closeOut: true, out: make(chan Message)}
}

func newRabbitChannelSubscription(logname string, replyChan *amqp.Channel, in <-chan amqp.Delivery, onExit ...Closer) *rabbitSubscription {
	out := make(chan Message)
	kill := make(chan bool)
	go func() {
		log := log.Sublogger(logname)

		for _, v := range onExit {
			defer v.Close()
		}
		defer close(out)

		for {
			select {
			case <-kill:
				log.Debug().Msg("stopping subscription worker")
				return
			case m, ok := <-in:
				if !ok {
					continue
				}
				log.Debug().Str("subject", m.RoutingKey).Msg("subscription received message")

				rab, err := NewRabbitMessage(m)
				if err != nil {
					log.Error().Err(err).Msg("failed to create rabbit message")
					continue
				}
				rab.setReplyChannel(replyChan)
				out <- rab
			}
		}
	}()
	return &rabbitSubscription{created: time.Now(), out: out, kill: kill}
}

func (s *rabbitSubscription) Channel() <-chan Message {
	return s.out
}

func (s *rabbitSubscription) Close() {
	if !s.closed {
		log.Debug().Msg("closing subscription")
		if s.closeOut {
			// since the channel sub variant does this in it's goroutine
			close(s.out)
		}
		if s.kill != nil {
			s.kill <- true
			close(s.kill)
		}
		s.closed = true
	}
}
