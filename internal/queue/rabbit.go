package queue

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/uuid"
	"github.com/voidshard/faction/pkg/structs"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mqApiReq = "api-request"
	mqStream = "change-stream"

	mqReplyRoutines = 10
)

type RabbitQueue struct {
	// https://pkg.go.dev/github.com/rabbitmq/amqp091-go
	// https://www.rabbitmq.com/tutorials/tutorial-two-go
	// https://www.rabbitmq.com/tutorials/tutorial-five-go
	cfg *RabbitConfig
	log log.Logger

	streamConn *amqp.Connection
	streamCh   *amqp.Channel

	apiConn  *amqp.Connection
	apiCh    *amqp.Channel
	apiQueue *amqp.Queue

	replyQueue   *amqp.Queue
	replyLock    sync.Mutex
	replyChans   map[string]*rabbitSubscription
	replies      <-chan amqp.Delivery
	replyMaxTime time.Duration
}

type RabbitConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

func NewRabbitQueue(cfg *RabbitConfig) *RabbitQueue {
	return &RabbitQueue{
		cfg: cfg,
		log: log.Sublogger("rabbit-queue", map[string]string{
			"host": cfg.Host,
			"port": fmt.Sprintf("%d", cfg.Port),
		}),
		replyChans: map[string]*rabbitSubscription{},
		replyLock:  sync.Mutex{},
	}
}

func (q *RabbitQueue) replyCleanRoutine() {
	// just so we can't have lots of queues sitting around in memory for messages that are
	// never going to arrive
	for {
		time.Sleep(60 * time.Second)

		q.log.Debug().Msg("cleaning up old reply channels")
		q.replyLock.Lock()

		for id, sub := range q.replyChans {
			q.replyLock.Unlock()

			if time.Since(sub.created) > q.replyMaxTime {
				q.log.Debug().Str("id", id).Msg("cleaning up old reply channel")
				sub.out <- &rabbitMessage{msg: amqp.Delivery{Body: []byte("timeout waiting for reply")}}

				q.replyLock.Lock()
				delete(q.replyChans, id)
				sub.Close()
			}
		}

		q.replyLock.Unlock()
	}
}

func (q *RabbitQueue) replyRoutine() {
	for msg := range q.replies {
		q.replyLock.Lock()
		ch, ok := q.replyChans[msg.CorrelationId]
		if !ok || ok && ch.closed {
			q.replyLock.Unlock()
			continue
		}
		q.replyLock.Unlock()

		// unlock to send incase no one is listening / delays on the other end
		ch.out <- &rabbitMessage{msg: msg}

		q.replyLock.Lock()
		delete(q.replyChans, msg.CorrelationId)
		ch.Close()
		q.replyLock.Unlock()
	}
}

func (q *RabbitQueue) Close() {
	for _, ch := range []*amqp.Channel{q.streamCh, q.apiCh} {
		if ch != nil {
			defer ch.Close()
		}
	}
	for _, conn := range []*amqp.Connection{q.streamConn, q.apiConn} {
		if conn != nil {
			defer conn.Close()
		}
	}
}

func (q *RabbitQueue) channel(prefetch int) (*amqp.Connection, *amqp.Channel, error) {
	q.log.Info().Str("username", q.cfg.Username).Str("host", q.cfg.Host).Int("port", q.cfg.Port).Int("prefetch", prefetch).Msg("connecting to rabbitmq")

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", q.cfg.Username, q.cfg.Password, q.cfg.Host, q.cfg.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, ch.Qos(prefetch, 0, false)
}

// connectApi connects to the RabbitMQ API queue
func (q *RabbitQueue) connectApi() error {
	if q.apiConn != nil {
		return nil
	}

	// setup api request queue: where we send requests to
	apiConn, apiChan, err := q.channel(1) // prefetch 1 message
	if err != nil {
		return err
	}
	q.apiConn = apiConn
	q.apiCh = apiChan
	apiQueue, err := apiChan.QueueDeclare(
		mqApiReq, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}
	q.apiQueue = &apiQueue

	// setup api response queue: where we listen for responses
	replyQueue, err := apiChan.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	q.replyQueue = &replyQueue

	replies, err := apiChan.Consume(
		replyQueue.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return err
	}
	q.replies = replies

	// kick off routines to reply to api responses
	for i := 0; i < mqReplyRoutines; i++ {
		go q.replyRoutine()
	}
	go q.replyCleanRoutine()

	return err
}

// connectStream connects to the RabbitMQ change stream - we use this to send messages to
// the change stream (only).
func (q *RabbitQueue) connectStream() error {
	if q.streamConn != nil {
		return nil
	}
	conn, ch, err := q.newChangeStream()
	if err != nil {
		return err
	}

	q.streamConn = conn
	q.streamCh = ch
	return nil
}

// newChangeStream creates a new change stream exchange & returns the connection & channel.
func (q *RabbitQueue) newChangeStream() (*amqp.Connection, *amqp.Channel, error) {
	conn, ch, err := q.channel(20) // prefetch 20 messages
	if err != nil {
		return nil, nil, err
	}
	return conn, ch, ch.ExchangeDeclare(
		mqStream, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}

func (q *RabbitQueue) PublishApiReq(data []byte) (Subscription, error) {
	err := q.connectApi() // default api queue
	if err != nil {
		return nil, err
	}

	// prepare reply channel & message ID
	mid := uuid.NewID().String()
	correlationID := uuid.NewID().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer q.log.Debug().Int("data-size", len(data)).Msg("published api request")

	// store reply channel for our reply thread
	q.replyLock.Lock()
	defer q.replyLock.Unlock()

	sub := newRabbitTopicSubscription()
	q.replyChans[correlationID] = sub

	return sub, q.apiCh.PublishWithContext(ctx,
		"",              // exchange
		q.apiQueue.Name, // routing key
		false,           // mandatory
		false,
		amqp.Publishing{
			Timestamp:     time.Now(),
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			Body:          data,
			MessageId:     mid,
			CorrelationId: correlationID,
			ReplyTo:       q.replyQueue.Name,
		},
	)
}

func (q *RabbitQueue) SubscribeApiReq() (Subscription, error) {
	err := q.connectApi() // default api queue
	if err != nil {
		return nil, err
	}

	messages, err := q.apiCh.Consume(
		q.apiQueue.Name, // queue
		"",              // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return nil, err
	}
	return newRabbitChannelSubscription(fmt.Sprintf("rabbit-sub-%s", q.apiQueue.Name), q.apiCh, messages), nil
}

func (q *RabbitQueue) PublishChange(ch *structs.Change) error {
	subject, err := makeRabbitSubject(ch)
	if err != nil {
		return err
	}

	err = q.connectStream() // default change stream (used for sending only)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return q.streamCh.PublishWithContext(ctx,
		mqStream, // exchange
		subject,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Timestamp:     time.Now(),
			MessageId:     uuid.NewID().String(),
			CorrelationId: uuid.NewID().String(),
			ContentType:   "text/plain",
			Body:          []byte{},
		},
	)
}

func (q *RabbitQueue) SubscribeChange(ch *structs.Change, queueName string) (Subscription, error) {
	subject, err := makeRabbitSubject(ch)
	if err != nil {
		return nil, err
	}

	// make a new connection & channel for this subscription, so we can sub to
	// a given routing key & close it when we're done.
	conn, channel, err := q.newChangeStream()
	if err != nil {
		return nil, err
	}

	durable := queueName != ""

	subQueue, err := channel.QueueDeclare(
		queueName, // name
		durable,   // durable
		!durable,  // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		subQueue.Name, // queue name
		subject,       // routing key
		mqStream,      // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	messages, err := channel.Consume(
		subQueue.Name, // queue
		queueName,     // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return nil, err
	}

	return newRabbitChannelSubscription(
		fmt.Sprintf("rabbit-sub-%s", queueName),
		nil,
		messages,
		conn,
		channel,
	), nil
}

// makeRabbitSubject converts a Change into a RabbitMQ subject.
func makeRabbitSubject(ch *structs.Change) (string, error) {
	// "World" cannot be wildcarded.
	// https://www.rabbitmq.com/tutorials/tutorial-five-go
	// * (star) can substitute for exactly one word.
	// # (hash) can substitute for zero or more words.
	if ch.World == "" {
		return "", fmt.Errorf("world must be specified")
	}
	key, ok := structs.Metakey_name[int32(ch.Key)]
	if !ok {
		return "", fmt.Errorf("invalid key %d", ch.Key)
	}
	area := ch.Area
	if area == "" {
		area = "*"
	}
	id := ch.Id
	if id == "" {
		id = "*"
	}
	return fmt.Sprintf("%s.%s.%s.%s", ch.World, area, key, id), nil
}

// fromRabbitSubject converts a RabbitMQ subject into a Change.
func fromRabbitSubject(subject string) (*structs.Change, error) {
	parts := strings.Split(subject, ".")
	key, ok := structs.Metakey_value[parts[2]]
	if !ok {
		return nil, fmt.Errorf("invalid key %s", parts[2])
	}
	return &structs.Change{
		World: parts[0],
		Area:  parts[1],
		Key:   structs.Metakey(key),
		Id:    parts[3],
	}, nil
}
