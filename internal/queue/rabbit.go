package queue

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	mqApiReq = "api-request"
	mqStream = "change-stream"
	mqDefer  = "defer-change"
)

type RabbitQueue struct {
	// https://pkg.go.dev/github.com/rabbitmq/amqp091-go
	// https://www.rabbitmq.com/tutorials/tutorial-two-go
	// https://www.rabbitmq.com/tutorials/tutorial-five-go
	cfg *RabbitConfig
	log log.Logger

	// change stream
	streamConn *amqp.Connection
	streamCh   *amqp.Channel

	// api request stream
	apiConn  *amqp.Connection
	apiCh    *amqp.Channel
	apiQueue *amqp.Queue

	replyQueue *amqp.Queue
	replyLock  sync.Mutex
	replyChans map[string]*rabbitSubscription
	replies    <-chan amqp.Delivery // apiCh replies

	// defer change stream
	defConn *amqp.Connection
	defCh   *amqp.Channel
}

type RabbitConfig struct {
	Username string
	Password string
	Host     string
	Port     int

	ReplyMaxTime  time.Duration
	ReplyRoutines int

	PrefetchAPIRequests  int // default 1
	PrefetchChangeStream int // default 20
	PrefetchDeferStream  int // default 500
}

func NewRabbitQueue(cfg *RabbitConfig) *RabbitQueue {
	if cfg.ReplyMaxTime == 0 {
		cfg.ReplyMaxTime = time.Second * 60
	}
	if cfg.ReplyRoutines < 1 {
		cfg.ReplyRoutines = 10
	}
	if cfg.PrefetchAPIRequests < 1 {
		cfg.PrefetchAPIRequests = 1
	}
	if cfg.PrefetchChangeStream < 1 {
		cfg.PrefetchChangeStream = 20
	}
	if cfg.PrefetchDeferStream < 1 {
		cfg.PrefetchDeferStream = 500
	}
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

	msgData, err := injectTraceData(context.Background(), errData)
	if err != nil {
		q.log.Warn().Err(err).Msg("failed to inject trace data, API error responses may be confused")
		msgData = errData
	}

	for {
		time.Sleep(q.cfg.ReplyMaxTime / 5)

		q.log.Debug().Msg("cleaning up old reply channels")
		q.replyLock.Lock()

		for id, sub := range q.replyChans {
			q.replyLock.Unlock()

			if time.Since(sub.created) > q.cfg.ReplyMaxTime {
				q.log.Debug().Str("id", id).Msg("cleaning up old reply channel")

				sub.out <- &rabbitMessage{context: context.Background(), data: errData, msg: amqp.Delivery{Body: msgData}}

				q.replyLock.Lock()
				delete(q.replyChans, id)
				sub.Close()
			}
		}

		q.replyLock.Unlock()
	}
}

// replyRoutine listens for replies to API requests & forwards them to the correct channel.
//
// When we send an API request (PublishApiReq) we store the subscription channel in a map (correlation ID -> channel).
// We listen for all replies to API requests here & forward them to the correct channel.
func (q *RabbitQueue) replyRoutine() {
	// https://www.rabbitmq.com/tutorials/tutorial-six-go
	//
	// Each client that calls PublishApiReq will create a new queue & consumer for itself (ie. reply to worker-14)
	// and subscribe to the generic work queue (ie. mqApiReq). This way a client can receive replies to its own
	// work requests. However we still need to work out exactly which caller on the client side to forward the message on to.
	//
	// That is, if you call PublishApiReq(do-foo id: x) you want a reply to that request not the answer to
	// PublishApiReq(do-bar id: y) that this client may have sent (in another routine or something).
	//
	// So;
	// [specific function call in client A] -> [generic work queue] -> [reply queue for client A] -> [specific function call in client A]
	//
	// Here we're getting replies to API requests on *this* clients reply queue, and working out which
	// specific subscription to pass the message to, based on the correlation ID.
	//
	// Note that when we publish and API Request message we set a correlation ID, since this is returned to us
	// in the reply, we know who to send the reply to.
	for msg := range q.replies {
		q.replyLock.Lock()
		ch, ok := q.replyChans[msg.CorrelationId]
		if !ok || ok && ch.closed {
			q.replyLock.Unlock()
			continue
		}
		// unlock to send incase no one is listening / delays on the other end
		q.replyLock.Unlock()

		rab, err := NewRabbitMessage(msg)
		if err != nil {
			log.Error().Str("MessageId", msg.MessageId).Err(err).Msg("failed to create rabbit message")
			continue
		}
		q.log.Debug().Err(err).Str("MessageId", msg.MessageId).Msg("received reply, forwarding message to caller")
		ch.out <- rab

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
	i := 0
	for {
		i += 1
		if i > 1 {
			time.Sleep(time.Duration(2*i*i) * time.Second)
			q.log.Debug().Int("attempt", i).Msg("retrying connection to rabbitmq")
		}

		url := fmt.Sprintf("amqp://%s:%s@%s:%d/", q.cfg.Username, q.cfg.Password, q.cfg.Host, q.cfg.Port)
		conn, err := amqp.Dial(url)
		if err != nil {
			q.log.Error().Err(err).Msg("failed to connect to rabbitmq")
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			q.log.Error().Err(err).Msg("failed to create channel")
			conn.Close()
			continue
		}

		err = ch.Qos(prefetch, 0, false)
		if err != nil {
			q.log.Error().Err(err).Msg("failed to set prefetch")
			ch.Close()
			conn.Close()
			continue
		}

		return conn, ch, nil
	}
}

func (q *RabbitQueue) connectDefer() error {
	if q.defConn != nil {
		return nil
	}

	defConn, defCh, err := q.channel(q.cfg.PrefetchAPIRequests)
	if err != nil {
		return err
	}
	q.defConn = defConn
	q.defCh = defCh

	return nil
}

// connectApi connects to the RabbitMQ API queue
func (q *RabbitQueue) connectApi() error {
	if q.apiConn != nil {
		return nil
	}

	// setup api request queue: where we send requests to
	apiConn, apiChan, err := q.channel(q.cfg.PrefetchAPIRequests)
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
	for i := 0; i < q.cfg.ReplyRoutines; i++ {
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
	conn, ch, err := q.channel(q.cfg.PrefetchChangeStream)
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

func (q *RabbitQueue) PublishApiReq(ctx context.Context, data []byte) (Subscription, error) {
	err := q.connectApi() // default api queue
	if err != nil {
		return nil, err
	}

	// prepare reply channel & message ID
	mid := uuid.NewID().String()
	cid := uuid.NewID().String()

	// log & trace
	defer q.log.Debug().Str("MessageId", mid).Str("CorrelationId", cid).Int("data-size", len(data)).Msg("published api request")

	pan := log.NewSpan(ctx, "rabbit.PublishApiReq", map[string]interface{}{"mid": mid, "id": cid})
	defer pan.End()

	msgdata, err := injectTraceData(ctx, data)
	if err != nil {
		return nil, err
	}

	// store reply channel for our reply thread
	q.replyLock.Lock()
	defer q.replyLock.Unlock()

	sub := newRabbitTopicSubscription()
	q.replyChans[cid] = sub

	return sub, q.apiCh.PublishWithContext(ctx,
		"",              // exchange
		q.apiQueue.Name, // routing key
		false,           // mandatory
		false,
		amqp.Publishing{
			Timestamp:     time.Now(),
			DeliveryMode:  amqp.Persistent,
			ContentType:   "text/plain",
			Body:          msgdata,
			MessageId:     mid,
			CorrelationId: cid,
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

func (q *RabbitQueue) DeferChange(ctx context.Context, ch *structs.Change, tick int64) error {
	// Each world + tick has it's own queue that is durable, this way we can write changes into a
	// queue and leave reading them (ie. subscribing to the queue) until desired (ie. the desired tick
	// in the given world).
	//
	// Subscribers to 'deferred changes' are internal processes in the API server that watch each world
	// and it's current tick & subscribe ticks up to the current tick. For any message on any of the queues
	// subscribed to, they call "PublishChange" to send the change to the change stream exactly as if it
	// "happend now".
	//
	// Ie. someone says "queue change X in world W for time n + 2" - we write that change to the queue
	// "defer.W.n+2" and when we get to time n + 2 we read the queue and publish the change to the change stream.

	defQueue, err := makeDeferredQueueName(ch.World, tick)
	if err != nil {
		return err
	}

	err = q.connectDefer()
	if err != nil {
		return err
	}

	queue, err := q.defCh.QueueDeclare(
		defQueue, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	mid := uuid.NewID().String()
	cid := uuid.NewID().String()

	pan := log.NewSpan(ctx, "rabbit.DeferChange", map[string]interface{}{"mid": mid, "id": cid})
	defer pan.End()

	// marshal & inject trace data
	data, err := ch.MarshalJSON()
	if err != nil {
		return err
	}
	msgdata, err := injectTraceData(ctx, data)
	if err != nil {
		return err
	}

	q.log.Debug().Str("Queue", defQueue).Str("MessageId", mid).Str("CorrelationId", cid).Msg("deferring change")
	return q.defCh.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Timestamp:     time.Now(),
			MessageId:     mid,
			CorrelationId: cid,
			ContentType:   "text/plain",
			Body:          msgdata,
		})
}

func (q *RabbitQueue) DeleteDeferredChangeQueue(world string, tick int64) error {
	defQueue, err := makeDeferredQueueName(world, tick)
	if err != nil {
		return err
	}

	err = q.connectDefer()
	if err != nil {
		return err
	}

	// https://github.com/rabbitmq/amqp091-go/blob/main/channel.go#L1030
	q.log.Info().Str("Queue", defQueue).Msg("deleting deferred change queue")
	_, err = q.defCh.QueueDelete(defQueue, false, false, true)
	return err
}

func (q *RabbitQueue) SubscribeDeferredChanges(world string, tick int64) (Subscription, error) {
	defQueue, err := makeDeferredQueueName(world, tick)
	if err != nil {
		return nil, err
	}

	err = q.connectDefer()
	if err != nil {
		return nil, err
	}

	queue, err := q.defCh.QueueDeclare(
		defQueue, // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	messages, err := q.defCh.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}
	return newRabbitChannelSubscription(fmt.Sprintf("rabbit-defer-%s", defQueue), q.defCh, messages), nil
}

// publishChange publishes a change to the change stream.
func (q *RabbitQueue) PublishChange(ctx context.Context, ch *structs.Change) error {
	subject, err := makeRabbitSubject(ch)
	if err != nil {
		return err
	}

	err = q.connectStream() // default change stream (used for sending only)
	if err != nil {
		return err
	}

	mid := uuid.NewID().String()
	cid := uuid.NewID().String()

	pan := log.NewSpan(ctx, "rabbit.PublishChange", map[string]interface{}{"mid": mid, "id": cid})
	defer pan.End()

	msgdata, err := injectTraceData(ctx, []byte{})
	if err != nil {
		return err
	}

	q.log.Debug().Str("subject", subject).Str("MessageId", mid).Str("CorrelationId", cid).Msg("publishing change")
	return q.streamCh.PublishWithContext(ctx,
		mqStream, // exchange
		subject,  // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Timestamp:     time.Now(),
			MessageId:     mid,
			CorrelationId: cid,
			ContentType:   "text/plain",
			Body:          msgdata,
		},
	)
}

// SubscribeChange subscribes to changes on the change stream.
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

	q.log.Debug().Str("subject", subject).Str("queue", queueName).Msg("subscribing to change")
	return newRabbitChannelSubscription(
		fmt.Sprintf("rabbit.%s.%s", subject, queueName),
		nil,
		messages,
		conn,
		channel,
	), nil
}

// injectTraceData injects trace data, prepending it to the given data.
func injectTraceData(ctx context.Context, data []byte) ([]byte, error) {
	td, err := log.MarshalTrace(ctx)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(td)
	buf.Write([]byte("|"))
	buf.Write(data)

	return buf.Bytes(), nil
}

// extractTraceData extracts trace data from the given data
func extractTraceData(data []byte) (context.Context, []byte, error) {
	parts := bytes.SplitN(data, []byte("|"), 2)
	ctx, err := log.UnmarshalTrace(parts[0])
	if err != nil {
		return nil, nil, err
	}
	if len(parts) == 2 {
		return ctx, parts[1], nil
	}
	return ctx, nil, fmt.Errorf("no data found")
}

func makeDeferredQueueName(world string, tick int64) (string, error) {
	if world == "" {
		return "", fmt.Errorf("world must be specified")
	}
	return fmt.Sprintf("defer.%s.%d", world, tick), nil
}

// makeRabbitSubject converts a Change into a RabbitMQ subject.
func makeRabbitSubject(ch *structs.Change) (string, error) {
	// "World" cannot be wildcarded.
	// https://www.rabbitmq.com/tutorials/tutorial-five-go
	// * (star) can substitute for exactly one word.
	// # (hash) can substitute for zero or more words.
	world := ch.World
	if world == "" {
		world = "*"
	}
	key, ok := structs.Metakey_name[int32(ch.Key)]
	if !ok {
		return "", fmt.Errorf("invalid key %d", ch.Key)
	}
	if ch.Key == structs.Metakey_KeyNone {
		key = "*"
	}
	area := ch.Area
	if area == "" {
		area = "*"
	}
	id := ch.Id
	if id == "" {
		id = "*"
	}
	return fmt.Sprintf("%s.%s.%s.%s", world, area, key, id), nil
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
