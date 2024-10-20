package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
)

type tickManager struct {
	kill chan bool
	log  log.Logger

	db           db.Database
	qu           queue.Queue
	worldChanges queue.Subscription

	// cache of world id -> current tick, strictly increasing
	cache     map[string]int64
	cacheLock sync.Mutex

	// worldid,tick -> subscription (subscriptions to deferred changes for a given world/tick)
	subs     map[string]queue.Subscription
	subsLock sync.Mutex
}

func newTickManager(name string, db db.Database, qu queue.Queue) (*tickManager, error) {
	// ie. subscribe to all changes on all world objects
	sub, err := qu.SubscribeChange(&structs.Change{Key: structs.Metakey_KeyWorld}, "")
	if err != nil {
		return nil, err
	}
	return &tickManager{
		kill:         make(chan bool),
		log:          log.Sublogger(name),
		db:           db,
		qu:           qu,
		worldChanges: sub,
		cache:        make(map[string]int64),
		cacheLock:    sync.Mutex{},
		subs:         make(map[string]queue.Subscription),
		subsLock:     sync.Mutex{},
	}, nil
}

func (tc *tickManager) handleWorldChange(msg queue.Message) {
	pan := log.NewSpan(msg.Context(), "api.tickManager.handleWorldChange", map[string]interface{}{"mid": msg.Id()})
	defer pan.End()

	ch, err := msg.Change()
	if err != nil {
		tc.log.Error().Err(err).Msg("Failed to get change from message")
		pan.Err(err)
		return
	}

	worlds, err := tc.db.Worlds(msg.Context(), []string{ch.World})
	if err != nil {
		tc.log.Error().Str("World", ch.World).Err(err).Msg("Failed to get world")
		pan.Err(err)
		msg.Reject() // requeue
		return
	} else if worlds == nil || len(worlds) == 0 {
		tc.log.Error().Str("World", ch.World).Msg("World not found")
		pan.Err(fmt.Errorf("world not found %s", ch.World))
		msg.Ack() // no point retrying, world doesn't exist
		return
	}
	pan.SetAttributes(map[string]interface{}{"world": ch.World, "tick": worlds[0].Tick})

	tc.cacheLock.Lock()
	v, _ := tc.cache[ch.World]
	if worlds[0].Tick > v { // we only ever increase
		tc.cache[ch.World] = worlds[0].Tick
		tc.cacheLock.Unlock()
		err := tc.maybeAlterSubscriptions(ch.World, worlds[0].Tick)
		if err != nil {
			tc.log.Error().Str("World", ch.World).Int64("Tick", worlds[0].Tick).Err(err).Msg("Failed to alter subscriptions")
			pan.Err(err)
			msg.Reject() // requeue
			return
		}
	} else {
		tc.cacheLock.Unlock()
	}

	tc.log.Info().Str("MessageId", msg.Id()).Str("World", ch.World).Int64("Tick", worlds[0].Tick).Msg("Updated tick cache")
	msg.Ack()
}

func (tc *tickManager) watchSubscription(sub queue.Subscription) {
	// These are all deferred changes for some world,tick. Since we've got the message from
	// one of the queues we're subscribed to, and we only subscribe to world,tick to world,tick-3
	// we *know* that this message is either meant for the current tick (or is late and meant for a previous tick).
	//
	// So all we have to do is accept any incoming message and forward it into the world's change stream.
	// Easy peasy.
	for msg := range sub.Channel() {
		pan := log.NewSpan(msg.Context(), "api.tickManager.watchSubscription", map[string]interface{}{"mid": msg.Id()})

		ch := &structs.Change{} // nb. deferred changes store the change in the message body
		err := ch.UnmarshalJSON(msg.Data())
		if err != nil {
			tc.log.Error().Str("MessageId", msg.Id()).Err(err).Msg("Failed to unmarshal change")
			pan.Err(err)
			msg.Ack() // ack; we can't process this message so there's no point retrying it
			pan.End()
			continue
		}

		err = tc.qu.PublishChange(msg.Context(), ch)
		if err != nil {
			tc.log.Error().Str("MessageId", msg.Id()).Err(err).Msg("Failed to publish change")
			pan.Err(err)
			msg.Reject() // requeue; the message is valid, the queue might be unavailable
			pan.End()
			continue
		}

		msg.Ack()
		pan.End()
	}
}

func (tc *tickManager) maybeAlterSubscriptions(world string, tick int64) error {
	// since tick is always increasing, we should always sub onto the next tick
	// and remove the oldest tick (tick - 4).
	// ie. if we're at tick 10, we'll sub to 7,8,9,10 and remove 6
	//
	// In this way we maintain a rolling subscription to the current world tick, minus
	// a few ticks to allow for some late messages.

	tc.subsLock.Lock()
	defer tc.subsLock.Unlock()

	// tick to tick - 3 we'll maintain subs for
	// tick -4 and below we'll discard
	for i := tick - 3; i <= tick; i++ {
		if i < 0 {
			continue
		}
		key := fmt.Sprintf("%s,%d", world, i)

		if _, ok := tc.subs[key]; ok {
			continue // already subscribed, nothing to do
		}

		sub, err := tc.qu.SubscribeDeferredChanges(world, i)
		if err != nil {
			log.Warn().Str("World", world).Int("Tick", int(i)).Err(err).Msg("Failed to subscribe to deferred changes")
			return err
		}
		go tc.watchSubscription(sub)
		tc.subs[key] = sub
	}

	if tick-4 < 0 {
		return nil
	}

	key := fmt.Sprintf("%s,%d", world, tick-4)
	if sub, ok := tc.subs[key]; ok {
		sub.Close()
		delete(tc.subs, key)

		err := tc.qu.DeleteDeferredChangeQueue(world, tick-4)
		if err != nil {
			tc.log.Warn().Str("World", world).Int("Tick", int(tick-4)).Err(err).Msg("Failed to delete deferred change queue")
		}
	}

	return nil
}

func (tc *tickManager) Run() {
	// kick off a routine to listen for changes in worlds and update our cache
	go func() {
		defer tc.log.Debug().Msg("Tick manager worker stopped")
		for {
			select {
			case <-tc.kill:
				return
			case msg := <-tc.worldChanges.Channel():
				tc.log.Debug().Str("MessageId", msg.Id()).Msg("Tick manager recieved world change message")
				tc.handleWorldChange(msg)
			}
		}
	}()

	// read all worlds into cache (world id -> tick), retry on any errors forever
	tc.populateCache()
}

func (tc *tickManager) populateCache() {
	ctx := context.Background()
	pan := log.NewSpan(ctx, "api.tickManager.populateCache", nil)
	defer pan.End()

	var limit int64 = 1000

	var offset int64
	for {
		tc.log.Debug().Int("Offset", int(offset)).Int("Limit", int(limit)).Msg("Populating tick manager cache, listing worlds")
		worlds, err := tc.db.ListWorlds(ctx, nil, limit, offset)
		if err != nil {
			tc.log.Warn().Err(err).Msg("Failed to list worlds")
			pan.Err(err)
			time.Sleep(time.Second * 2)
			continue
		}

		tc.cacheLock.Lock()
		for _, w := range worlds {
			tc.cache[w.Id] = w.Tick
		}
		tc.cacheLock.Unlock()

		for _, w := range worlds {
			err = tc.maybeAlterSubscriptions(w.Id, w.Tick)
			if err != nil {
				tc.log.Warn().Str("World", w.Id).Int64("Tick", w.Tick).Err(err).Msg("Populating tick cache, failed to alter subscriptions")
				pan.Err(err)
				time.Sleep(time.Second * 2)
				continue
			}
		}

		if len(worlds) < int(limit) { // we've iterated all worlds
			tc.log.Info().Msg("Populated tick manager cache")
			return
		}

		offset += int64(len(worlds))
	}
}

func (tc *tickManager) Kill() {
	tc.log.Debug().Msg("Killing tick cache worker")
	defer tc.log.Debug().Msg("Tick cache worker killed")

	tc.kill <- true
	close(tc.kill)
}
