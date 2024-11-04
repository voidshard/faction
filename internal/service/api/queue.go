package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jellydator/ttlcache/v3"

	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

const (
	queueAPIRequests = "apiserver-api-requests"
	topicChanges     = "apiserver-changes"
	topicChangesAck  = "apiserver-changes-ack"
)

// Queue presents queue functionality the way the APISever expects to use it.
//
// Ie. it handles publishing and subscribing to "changes" rather than "topics" or "keys"
type Queue struct {
	qu  queue.Queue
	id  string
	log log.Logger

	ackCache *ttlcache.Cache[string, queue.Message]
	ackSub   queue.Subscription
}

func newQueue(qu queue.Queue) (*Queue, error) {
	// generate a unique id for this server
	id := uuid.NewID().String()
	l := log.Sublogger("apiserver.Queue", map[string]interface{}{"id": id})

	// prep a cache for acks that this server is waiting on
	cache := ttlcache.New[string, queue.Message](
		ttlcache.WithTTL[string, queue.Message](10 * time.Minute),
	)

	// subscribe to an ack topic for this server
	sub, err := qu.Subscribe("", topicChangesAck, []string{id, ""})
	if err != nil {
		return nil, err
	}

	// setup & kick off routines
	me = &Queue{qu: qu, id: id, log: l, ackCache: cache, ackSub: sub}

	go func() {
		for msg := range sub.Channel() { // no need to ack (queue to sub is "")
			// pull out and ack the message if we're waiting for it
			ackId := string(msg.Data())
			log.Debug().Str("AckId", ackId).Msg("Received ack from topic")
			me.Ack(ackId)
		}
	}()

	go cache.Start() // cache cleaning

	return me, nil
}

func (q *Queue) Close() {
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
			return q.qu.Publish(context.Background(), topicChangesAck, bits, []byte(id))
		}
	}
	return fmt.Errorf("invalid ack id %s", id)
}

func (q *Queue) DeferAck(msg queue.Message) string {
	ackId := fmt.Sprintf("%s.%s", q.id, msg.Id())
	q.ackCache.Set(msg.Id(), msg, ttlcache.DefaultTTL)
	return ackId
}

// EnqueueApiReq publishes data to the API request stream.
func (q *Queue) EnqueueApiReq(ctx context.Context, data []byte) (queue.Subscription, error) {
	return q.qu.Request(ctx, queueAPIRequests, data)
}

// DequeueApiReq subscribes to the API request stream.
func (q *Queue) DequeueApiReq() (queue.Subscription, error) {
	return q.qu.Dequeue(queueAPIRequests)
}

// PublishChange publishes a change to the change stream.
func (q *Queue) PublishChange(ctx context.Context, ch *structs.Change) error {
	key, err := toQueueKey(ch)
	if err != nil {
		return err
	}
	data, err := ch.MarshalJSON()
	if err != nil {
		return err
	}
	return q.qu.Publish(ctx, topicChanges, key, data)
}

// SubscribeChange subscribes to changes on the change stream.
// queueName can be given to configure a durable queue. If not given
// a temporary non-durable queue will be used.
func (q *Queue) SubscribeChange(ch *structs.Change, queueName string) (queue.Subscription, error) {
	key, err := toQueueKey(ch)
	if err != nil {
		return nil, err
	}
	if queueName != "" { // ensure that queue name cannot clash with something internal
		queueName = fmt.Sprintf("subscribe-change.apiserver.%s", queueName)
	}
	return q.qu.Subscribe(queueName, topicChanges, key)
}

// DeferChange defers a change to be processed at given tick.
func (q *Queue) DeferChange(ctx context.Context, ch *structs.Change, tick int64) error {
	data, err := ch.MarshalJSON()
	if err != nil {
		return err
	}
	qname, err := deferredQueueName(ch.World, tick)
	if err != nil {
		return err
	}
	return q.qu.Enqueue(ctx, qname, data)
}

// SubscribeDeferredChanges subscribes to changes that have been deferred to a given tick.
func (q *Queue) SubscribeDeferredChanges(world string, tick int64) (queue.Subscription, error) {
	qname, err := deferredQueueName(world, tick)
	if err != nil {
		return nil, err
	}
	return q.qu.Dequeue(qname)
}

// DeleteDeferredChangeQueue deletes the queue for deferred changes at given tick.
// This should be called when the queue is no longer needed to tidy old queues we're
// never going to use again.
func (q *Queue) DeleteDeferredChangeQueue(world string, tick int64) error {
	qname, err := deferredQueueName(world, tick)
	if err != nil {
		return err
	}
	return q.qu.DeleteQueue(qname)
}

func deferredQueueName(world string, tick int64) (string, error) {
	if !uuid.IsValidUUID(world) {
		return "", fmt.Errorf("invalid world id %s", world)
	}
	return fmt.Sprintf("apiserver-defer-%s.%d", world, tick), nil
}

// toQueueKey converts a Change into a queue Key ([]string)
func toQueueKey(ch *structs.Change) ([]string, error) {
	key := make([]string, 4) // world, area, (metadata) key, id

	key[0] = ch.World
	key[1] = ch.Area
	key[3] = ch.Id

	metakey, ok := structs.Metakey_name[int32(ch.Key)]
	if !ok {
		return nil, fmt.Errorf("invalid key %d", ch.Key)
	}
	if ch.Key == structs.Metakey_KeyNone {
		key[2] = ""
	} else {
		key[2] = metakey
	}

	return key, nil
}
