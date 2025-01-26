package queue

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

type RabbitConfig struct {
	Username string
	Password string
	Host     string
	Port     int

	ReplyMaxTime  time.Duration
	ReplyRoutines int

	PrefetchDequeue   int // default 1
	PrefetchSubscribe int // default 20
}

type RabbitQueue struct {
	cfg *RabbitConfig
	log log.Logger

	conn *rabbitChannel // connection for sending messages

	replyChan  *rabbitChannel // connection for receiving replies
	replyQueue string
	replyLock  sync.Mutex
	replyChans map[string]*rabbitSubscription
}

func NewRabbitQueue(cfg *RabbitConfig) (*RabbitQueue, error) {
	if cfg.ReplyMaxTime == 0 {
		cfg.ReplyMaxTime = time.Minute
	}
	if cfg.ReplyRoutines < 1 {
		cfg.ReplyRoutines = 1
	}

	conn, err := newRabbitChannel(cfg, "send", 1)
	if err != nil {
		return nil, err
	}

	return &RabbitQueue{
		cfg:  cfg,
		conn: conn,
		log: log.Sublogger("rabbit-queue", map[string]interface{}{
			"host": cfg.Host,
			"port": cfg.Port,
		}),
		replyQueue: fmt.Sprintf("internal.reply.%s", uuid.NewID().String()),
		replyLock:  sync.Mutex{},
		replyChans: make(map[string]*rabbitSubscription),
	}, nil
}

func (q *RabbitQueue) Close() error {
	q.replyLock.Lock()
	defer q.replyLock.Unlock()

	for _, sub := range q.replyChans {
		sub.Close()
	}

	if q.replyChan != nil {
		q.replyChan.Close()
	}
	q.conn.Close()
	return nil
}

// DeleteQueue deletes a queue.
func (q *RabbitQueue) DeleteQueue(queue string) error {
	_, err := q.conn.Channel().QueueDelete(queue, false, false, false)
	return err
}

// Enqueue sends a message to a queue.
func (q *RabbitQueue) Enqueue(ctx context.Context, queue string, data []byte) error {
	return q.enqueue(ctx, uuid.NewID().String(), queue, "", data)
}

// Request sends a message to a queue and waits for a reply.
func (q *RabbitQueue) Request(ctx context.Context, queue string, data []byte) (Subscription, error) {
	cid := uuid.NewID().String()

	// since we're going to be waiting for a reply we need to ensure we're ready to receive them
	err := q.ensureReadyForReplies()
	if err != nil {
		return nil, err
	}

	// store reply channel for our reply thread
	sub := newRabbitSubscription("request", map[string]interface{}{"CorrelationId": cid, "Queue": queue})
	q.replyLock.Lock()
	q.replyChans[cid] = sub
	q.replyLock.Unlock()

	return sub, q.enqueue(ctx, cid, queue, q.replyQueue, data)
}

// enqueue sends a message to a queue.
func (q *RabbitQueue) enqueue(ctx context.Context, cid, queue, replyTo string, data []byte) error {
	ch := q.conn.Channel()

	// name, durable, autoDelete, exclusive, noWait, args
	qu, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// prepare reply channel & message ID
	mid := uuid.NewID().String()

	// log & trace
	defer q.log.Debug().Str("MessageId", mid).Str("CorrelationId", cid).Str("Queue", queue).Int("Data", len(data)).Msg("enqueued message")
	pan := log.NewSpan(ctx, "rabbit.Enqueue", map[string]interface{}{"MessageId": mid, "CorrelationId": cid, "Queue": queue, "ReplyTo": replyTo, "bytes": len(data)})
	defer pan.End()

	mdata, err := log.InjectTraceData(ctx, data)
	if err != nil {
		pan.Err(err)
		return err
	}

	log.Debug().Str("ReplyTo", q.replyQueue).Str("MessageId", mid).Str("CorrelationId", cid).Int("bytes", len(mdata)).Msg("Injected telemetry, enqueuing message")
	err = ch.PublishWithContext(ctx,
		"",      // exchange
		qu.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			Timestamp:     time.Now(),
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			Body:          mdata,
			MessageId:     mid,
			CorrelationId: cid,
			ReplyTo:       q.replyQueue,
		},
	)
	pan.Err(err)

	return err
}

// Dequeue returns a subscription to the given queue.
func (q *RabbitQueue) Dequeue(queue string) (Subscription, error) {
	rchan, err := newRabbitChannel(q.cfg, queue, q.cfg.PrefetchDequeue)
	if err != nil {
		return nil, err
	}

	notify := make(chan *amqp.Channel)
	sub := newRabbitSubscription("dequeue", map[string]interface{}{"Queue": queue})

	go func() {
		for {
			select {
			case <-sub.kill:
				rchan.Close()
				sub.Close()
				return
			case ch := <-rchan.ChannelStream(notify):
				q.log.Debug().Msg("Dequeue channel change detected")

				// name, durable, autoDelete, exclusive, noWait, args
				qu, err := ch.QueueDeclare(queue, true, false, false, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to declare queue")
					continue
				}

				// name, consumer, auto-ack, exclusive, no-local, no-wait, args
				msgs, err := ch.Consume(qu.Name, "", false, false, false, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to consume queue")
					continue
				}

				q.log.Info().Str("queue", qu.Name).Msg("subscribed to queue")
				go sub.consumeDeliveries(rchan, msgs)
			}
		}
	}()

	notify <- rchan.Channel() // pass in the current channel to kick off
	return sub, nil
}

// Publish sends a message to a topic.
func (q *RabbitQueue) Publish(ctx context.Context, topic string, key []string, data []byte) error {
	ch := q.conn.Channel()

	err := ch.ExchangeDeclare(
		topic,   // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	mid := uuid.NewID().String()
	cid := uuid.NewID().String()
	rkey := toRabbitKey(key)

	pan := log.NewSpan(ctx, "rabbit.Publish", map[string]interface{}{
		"MessageId": mid, "CorrelationId": cid, "Topic": topic, "Key": rkey, "bytes": len(data),
	})
	defer pan.End()

	mdata, err := log.InjectTraceData(ctx, data)
	if err != nil {
		pan.Err(err)
		return err
	}

	err = ch.PublishWithContext(ctx,
		topic, // exchange
		rkey,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Timestamp:     time.Now(),
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			Body:          mdata,
			MessageId:     mid,
			CorrelationId: cid,
		},
	)
	pan.Err(err)

	return err
}

// Subscribe returns a subscription to the given topic.
func (q *RabbitQueue) Subscribe(queue, topic string, key []string, durable bool) (Subscription, error) {
	rchan, err := newRabbitChannel(q.cfg, fmt.Sprintf("%s:%s", queue, topic), q.cfg.PrefetchSubscribe)
	if err != nil {
		return nil, err
	}

	rkey := toRabbitKey(key)

	notify := make(chan *amqp.Channel)
	sub := newRabbitSubscription("subscribe", map[string]interface{}{"Queue": queue, "Topic": topic, "Key": rkey})
	go func() {
		for {
			select {
			case <-sub.kill:
				rchan.Close()
				sub.Close()
				return
			case ch := <-rchan.ChannelStream(notify):
				q.log.Debug().Msg("Subscribe channel change detected")

				// name, type, durable, auto-deleted, internal, no-wait, arguments
				err = ch.ExchangeDeclare(topic, "topic", true, false, false, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to declare exchange")
					continue
				}

				// name, durable, autoDelete, exclusive, noWait, args
				qu, err := ch.QueueDeclare(queue, durable, !durable, false, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to declare queue")
					continue
				}

				// name, routingKey, exchange, noWait, args
				err = ch.QueueBind(qu.Name, rkey, topic, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to bind queue")
					continue
				}

				// name, consumer, auto-ack, exclusive, no-local, no-wait, args
				msgs, err := ch.Consume(qu.Name, "", !durable, false, false, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to consume queue")
					continue
				}

				q.log.Info().Str("queue", qu.Name).Str("topic", topic).Str("key", rkey).Msg("bound to topic")

				// nb. we shouldn't leak routines unless the rabbit lib fails to close the chan it's
				// made for us when the underlying connection is closed
				go sub.consumeDeliveries(rchan, msgs)
			}
		}
	}()

	notify <- rchan.Channel() // pass in the current channel to kick off
	return sub, nil
}

// replyCleanRoutine cleans up old reply channels.
func (q *RabbitQueue) replyCleanRoutine() {
	// just so we can't have lots of queues sitting around in memory for messages that are
	// never going to arrive

	// our standard "time out" message; we can marshal it only once
	apiErr := &structs.Error{Message: "timeout waiting for reply", Code: structs.ErrorCode_TIMEOUT}
	errData, err := apiErr.MarshalJSON()
	if err != nil {
		q.log.Warn().Err(err).Msg("failed to marshal static error data, API error responses may be confused")
		errData = []byte("timeout")
	}

	msgData, err := log.InjectTraceData(context.Background(), errData)
	if err != nil {
		q.log.Warn().Err(err).Msg("failed to inject trace data, API error responses may be confused")
		msgData = errData
	}

	for {
		time.Sleep(q.cfg.ReplyMaxTime / 5)

		q.log.Debug().Msg("cleaning up old reply channels")
		q.replyLock.Lock()

		for id, sub := range q.replyChans {
			if time.Since(sub.created) > q.cfg.ReplyMaxTime {
				q.log.Debug().Str("id", id).Msg("cleaning up old reply channel")
				sub.out <- &rabbitMessage{
					context: context.Background(),
					data:    errData,
					msg:     amqp.Delivery{Body: msgData},
				}
				delete(q.replyChans, id)
				sub.Close()
			}
		}

		q.replyLock.Unlock()
	}
}

func (q *RabbitQueue) ensureReadyForReplies() error {
	if q.replyChan != nil {
		return nil
	}

	err := q.setupReplyQueue()
	if err != nil {
		return err
	}
	go q.replyCleanRoutine()

	return nil
}

func (q *RabbitQueue) setupReplyQueue() error {
	rchan, err := newRabbitChannel(q.cfg, "reply", 1)
	if err != nil {
		return err
	}

	q.replyChan = rchan
	notify := make(chan *amqp.Channel)

	sub := newRabbitSubscription("reply", map[string]interface{}{"Queue": q.replyQueue})
	for i := 0; i < q.cfg.ReplyRoutines; i++ {
		go q.replyRoutine(sub)
	}

	go func() {
		for {
			select {
			case <-sub.kill:
				rchan.Close()
				sub.Close()
				return
			case ch := <-rchan.ChannelStream(notify):
				q.log.Debug().Msg("Reply channel change detected")

				// name, durable, autoDelete, exclusive, noWait, args
				qu, err := ch.QueueDeclare(q.replyQueue, false, true, true, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to declare queue")
					continue
				}

				// name, consumer, auto-ack, exclusive, no-local, no-wait, args
				msgs, err := ch.Consume(qu.Name, "", true, false, false, false, nil)
				if err != nil {
					q.log.Warn().Err(err).Msg("failed to consume queue")
					continue
				}

				q.log.Info().Str("queue", qu.Name).Msg("subscribed to reply queue")
				go sub.consumeDeliveries(rchan, msgs)
			}
		}
	}()

	notify <- rchan.Channel() // pass in the current channel to kick off
	return nil
}

// replyRoutine listens for replies to Enqueue messages & forwards them to the correct channel.
//
// Each client makes a queue to reply to it (ie. reply to Pod-14), and listens on that queue.
// That allows us to get back a message to the correct client. Messages have a CorrelationId
// which is how we work out which call (ie. Enqueue) the reply is for.
func (q *RabbitQueue) replyRoutine(sub *rabbitSubscription) {
	// https://www.rabbitmq.com/tutorials/tutorial-six-go
	for msg := range sub.Channel() {
		q.replyLock.Lock()
		ch, ok := q.replyChans[msg.CorrelationId()]
		if !ok {
			q.replyLock.Unlock()
			continue
		}
		// unlock to send incase no one is listening / delays on the other end
		q.replyLock.Unlock()

		ch.out <- msg

		q.replyLock.Lock()
		delete(q.replyChans, msg.CorrelationId())
		q.replyLock.Unlock()
		ch.Close()
	}
}

func toRabbitKey(key []string) string {
	for i, k := range key {
		if k == "" {
			key[i] = "*"
		}
	}
	return strings.Join(key, ".")
}
