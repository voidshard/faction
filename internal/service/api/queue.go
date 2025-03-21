package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jellydator/ttlcache/v3"

	"github.com/voidshard/faction/internal/queue"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

const (
	topicEvents    = "internal.apiserver-events"
	topicEventsAck = "internal.apiserver-events-ack"
)

// Queue presents queue functionality the way the APISever expects to use it.
//
// Ie. it handles publishing and subscribing to "events" rather than "topics" or "keys"
type Queue struct {
	qu  queue.Queue
	id  string
	log log.Logger

	ackCache *ttlcache.Cache[string, queue.Message]
	ackSub   queue.Subscription
}

func newQueue(qu queue.Queue) (*Queue, error) {
	// generate a unique id for this server
	id := fmt.Sprintf("%s", uuid.New())
	l := log.Sublogger("apiserver.Queue", map[string]interface{}{"id": id})

	// prep a cache for acks that this server is waiting on
	cache := ttlcache.New[string, queue.Message](
		ttlcache.WithTTL[string, queue.Message](5 * time.Minute),
	)

	// subscribe to an ack topic for this server
	sub, err := qu.Subscribe(id, topicEventsAck, []string{id, ""}, false)
	if err != nil {
		return nil, err
	}

	// setup & kick off routines
	me := &Queue{qu: qu, id: id, log: l, ackCache: cache, ackSub: sub}

	go func() {
		// These are messages from other APIServers who wish to Ack messages we were
		// delivered from Rabbit.
		for msg := range sub.Channel() { // no need to ack (queue to sub is "")
			// pull out and ack the message if we're waiting for it
			ackId := string(msg.Data())
			log.Debug().Str("AckId", ackId).Msg("Received ack from topic")
			me.Ack(ackId)
		}
	}()

	cache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, queue.Message]) {
		if reason == ttlcache.EvictionReasonDeleted {
			return // we deleted it, no need to ack / nack
		}
		// otherwise reject so it is redelivered
		item.Value().Reject()
		log.Warn().Str("AckId", item.Key()).Str("reason", fmt.Sprintf("%s", reason)).Msg("AckId evicted from cache")
	})
	go cache.Start() // cache cleaning

	return me, nil
}

func (q *Queue) Close() {
	q.qu.DeleteQueue(q.id)
	q.ackSub.Close()
	q.qu.Close()
}

func (q *Queue) Ack(id string) error {
	bits := strings.SplitN(id, ".", 2)
	if len(bits) != 2 {
		return fmt.Errorf("invalid ack id %s", id)
	}
	if uuid.IsValidUUID(bits[0]) && uuid.IsValidUUID(bits[1]) {
		if bits[0] == q.id { // this host sent the original message
			item, ok := q.ackCache.GetAndDelete(bits[1])
			if ok {
				return item.Value().Ack()
			}
			return nil // ack already processed
		} else { // another host sent the message: publish
			q.log.Debug().Str("AckId", id).Msg("Non-local ack, forwarding to topic")
			return q.qu.Publish(context.Background(), topicEventsAck, bits, []byte(id))
		}
	}
	return fmt.Errorf("invalid ack id %s", id)
}

// NewAckId generates a new ack id for a given message.
//
// We keep tabs on this Id & msg so that the caller can send us an 'ack' for it later.
func (q *Queue) NewAckId(msg queue.Message, event *v1.Event) string {
	ackId := fmt.Sprintf("%s.%s", q.id, msg.Id())
	q.ackCache.Set(msg.Id(), msg, ttlcache.DefaultTTL)
	return ackId
}

// PublishEvent publishes a event to the event stream.
func (q *Queue) PublishEvent(ctx context.Context, ch *v1.Event) error {
	key := toQueueKey(ch)
	data, err := json.Marshal(ch)
	if err != nil {
		return err
	}
	return q.qu.Publish(ctx, topicEvents, key, data)
}

// SubscribeEvent subscribes to events on the event stream.
// queueName can be given to configure a durable queue. If not given
// a temporary non-durable queue will be used.
func (q *Queue) SubscribeEvent(ch *v1.Event, queueName string, durable bool) (queue.Subscription, error) {
	key := toQueueKey(ch)
	queueName = fmt.Sprintf("subscribe-event.apiserver.%s", queueName)
	return q.qu.Subscribe(queueName, topicEvents, key, durable)
}

// DeferEvent defers a event to be processed at given tick.
func (q *Queue) DeferEvent(ctx context.Context, ch *v1.Event, tick uint64) error {
	data, err := json.Marshal(ch)
	if err != nil {
		return err
	}
	qname, err := deferredQueueName(ch.World, tick)
	if err != nil {
		return err
	}
	return q.qu.Enqueue(ctx, qname, data)
}

// SubscribeDeferredEvents subscribes to events that have been deferred to a given tick.
func (q *Queue) SubscribeDeferredEvents(world string, tick uint64) (queue.Subscription, error) {
	qname, err := deferredQueueName(world, tick)
	if err != nil {
		return nil, err
	}
	return q.qu.Dequeue(qname)
}

// DeleteDeferredEventQueue deletes the queue for deferred events at given tick.
// This should be called when the queue is no longer needed to tidy old queues we're
// never going to use again.
func (q *Queue) DeleteDeferredEventQueue(world string, tick uint64) error {
	qname, err := deferredQueueName(world, tick)
	if err != nil {
		return err
	}
	return q.qu.DeleteQueue(qname)
}

func deferredQueueName(world string, tick uint64) (string, error) {
	return fmt.Sprintf("internal.defer.%s.%d", world, tick), nil
}

// toQueueKey converts a Event into a queue Key ([]string)
func toQueueKey(ch *v1.Event) []string {
	return []string{ch.World, ch.Kind, ch.Controller, ch.Id}
}
