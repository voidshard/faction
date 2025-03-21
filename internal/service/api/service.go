package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/search"
	"github.com/voidshard/faction/pkg/kind"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"
	"github.com/voidshard/faction/pkg/util/uuid"
)

// Service handles the business logic of the API.
type Service struct {
	cfg *Config
	log log.Logger

	db db.Database
	qu *Queue
	sb search.Search

	publisher chan (*publishWork)

	tickManager *tickManager

	shutdownLock       sync.RWMutex
	shuttingDown       bool
	shutdownPublishers sync.WaitGroup
}

type publishWork struct {
	ctx    context.Context
	events []*v1.Event
}

func newService(cfg *Config, db db.Database, qu queue.Queue, sb search.Search) (*Service, error) {
	apiQueue, err := newQueue(qu)
	if err != nil {
		return nil, err
	}

	tm, err := newTickManager("tick-manager", db, apiQueue)
	if err != nil {
		return nil, err
	}
	go tm.Run()

	me := &Service{
		cfg:                cfg,
		log:                log.Sublogger("api.service", map[string]interface{}{}),
		db:                 db,
		qu:                 apiQueue,
		sb:                 sb,
		publisher:          make(chan (*publishWork)),
		tickManager:        tm,
		shutdownLock:       sync.RWMutex{},
		shutdownPublishers: sync.WaitGroup{},
	}

	for i := 0; i < cfg.PublishRoutines; i++ {
		me.shutdownPublishers.Add(1)
		go me.publishEvents()
	}

	return me, nil
}

func (s *Service) Shutdown() {
	if s.shuttingDown {
		return
	}
	s.shuttingDown = true

	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	close(s.publisher)
	s.shutdownPublishers.Wait()
	s.tickManager.Shutdown()
}

func (s *Service) publishEvents() {
	defer s.shutdownPublishers.Done()

	for work := range s.publisher {
		if work.events == nil || len(work.events) == 0 {
			return
		}
		l := log.Sublogger("api.publishEvents", map[string]interface{}{"world": work.events[0].World, "kind": work.events[0].Kind})
		for _, evt := range work.events {
			for {
				// if the queue is down we will keep retrying. Note that this
				// could result in our write operations blocking on pushing to
				// the publisher channel, this is intentional.
				err := s.qu.PublishEvent(work.ctx, evt)
				l.Debug().Err(err).Str("id", evt.Id).Str("controller", evt.Controller).Msg("publish event")
				if err == nil {
					break
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *Service) subscribeToEvents(req *api.StreamEvents) (<-chan *v1.Event, chan<- bool, error) {
	durable := true
	if req.Queue == "" {
		durable = false
		req.Queue = fmt.Sprintf("client.%s", uuid.New())
	}
	sub, err := s.qu.SubscribeEvent(
		&v1.Event{World: req.World, Kind: req.Kind, Id: req.Id, Controller: req.Controller},
		req.Queue,
		durable,
	)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan *v1.Event)
	kill := make(chan bool)

	go func() {
		attrs := map[string]interface{}{"world": req.World, "kind": req.Kind, "id": req.Id, "controller": req.Controller, "queue": req.Queue, "durable": durable}
		l := log.Sublogger("api.subscribeToEvents", attrs)

		for {
			select {
			case <-kill:
				sub.Close()
				return
			case msg, ok := <-sub.Channel():
				if !ok {
					return
				}
				l.Debug().Str("MessageId", msg.Id()).Msg("Received change")
				pan := log.NewSpan(msg.Context(), "service.subscribeToEvents", attrs, map[string]interface{}{"MessageId": msg.Id()})

				evt := &v1.Event{}
				err := json.Unmarshal(msg.Data(), evt)
				if err != nil {
					if durable {
						msg.Ack() // message is borked, throw away
					}
					l.Error().Err(err).Msg("Failed to unmarshal event")
					pan.Err(err)
					pan.End()
					continue
				}

				if durable {
					evt.AckId = s.qu.NewAckId(msg, evt)
				}

				events <- evt
				pan.End()
			}
		}
	}()

	return events, kill, nil
}

func (s *Service) ackEvent(ackId string) error {
	return s.qu.Ack(ackId)
}

func (s *Service) deferEvent(ctx context.Context, req *api.DeferEventRequest, rsp *api.DeferEventResponse) error {
	pan := log.NewSpan(ctx, "service.deferEvent")
	defer pan.End()

	// bonus validation that Validate package doesn't seem to support
	if req.ToTick <= 0 && req.ByTick <= 0 {
		return fmt.Errorf("%w either to_tick or by_tick must be set", ErrInvalid)
	} else if req.ToTick > 0 && req.ByTick > 0 {
		return fmt.Errorf("%w only one of to_tick or by_tick can be set", ErrInvalid)
	}

	// figure out the tick to set
	currentTick, err := s.tickManager.Tick(req.World)
	if err != nil {
		return err
	}
	if req.ToTick < currentTick-2 {
		// tick manager holds the last 4 but we regard any messages this old to be
		// erroneous, as opposed to simply late.
		return fmt.Errorf("%w cannot defer to a tick in the past", ErrPrecondition)
	}

	toTick := req.ToTick
	if toTick <= 0 {
		if err != nil {
			return err
		}
		toTick = currentTick + req.ByTick
	}
	rsp.ToTick = toTick
	pan.SetAttributes(map[string]interface{}{"world": req.World, "kind": req.Kind, "controller": req.Controller, "id": req.Id, "to_tick": toTick})

	// queue the event
	return s.qu.DeferEvent(ctx, &v1.Event{
		World:      req.World,
		Kind:       req.Kind,
		Controller: req.Controller,
		Id:         req.Id,
	}, toTick)
}

func (s *Service) searchKind(ctx context.Context, world string, req *api.SearchRequest, rsp *api.SearchResponse) error {
	pan := log.NewSpan(ctx, "service.searchKind", map[string]interface{}{"world": world, "kind": req.Kind})

	ids, err := s.sb.Find(ctx, world, &req.Query)
	if err != nil {
		pan.Err(err)
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	result := []map[string]interface{}{}
	err = s.db.Get(ctx, world, req.Kind, ids, &result)
	if err != nil {
		pan.Err(err)
		return err
	}

	pan.SetAttributes(map[string]interface{}{"data": len(result)})

	objects := []interface{}{}
	for _, obj := range result {
		objects = append(objects, obj)
	}
	rsp.Data = objects

	return nil
}

func (s *Service) getKind(ctx context.Context, k string, req *api.GetRequest, rsp *api.GetResponse) error {
	pan := log.NewSpan(ctx, "service.getKind", map[string]interface{}{"world": req.World, "kind": k})

	// we don't need a worldspace if it's a global kind
	isGlobal := kind.IsGlobal(k)
	if isGlobal {
		req.World = ""
	}

	// either get by ID(s) or list w/ labels + pagination
	result := []map[string]interface{}{}
	if len(req.Ids) > 0 {
		err := s.db.Get(ctx, req.World, k, req.Ids, &result)
		if err != nil {
			pan.Err(err)
			return err
		}
	} else {
		err := s.db.List(ctx, req.World, k, req.Labels, req.Limit, req.Offset, &result)
		if err != nil {
			pan.Err(err)
			return err
		}
	}
	pan.SetAttributes(map[string]interface{}{"data": len(result)})

	objects := []interface{}{}
	for _, obj := range result {
		objects = append(objects, obj)
	}
	rsp.Data = objects

	return nil
}

func (s *Service) setKind(ctx context.Context, k string, req *api.SetRequest, rsp *api.SetResponse) error {
	pan := log.NewSpan(ctx, "service.setKind", map[string]interface{}{"world": req.World, "kind": k})

	// we don't need a worldspace if it's a global kind
	if kind.IsGlobal(k) {
		req.World = ""
	}

	s.shutdownLock.RLock()
	defer s.shutdownLock.RUnlock()
	if s.shuttingDown {
		return ErrShuttingDown
	}

	events := []*v1.Event{}
	objects := []v1.Object{}
	for _, rawObj := range req.Data {
		// ensure each object is valid kind
		obj, err := kind.New(k, rawObj)
		if err != nil {
			pan.Err(err)
			return err
		}

		// validate object
		err = kind.Validate(k, obj)
		if err != nil {
			return fmt.Errorf("object invalid: %w %v", ErrInvalid, err)
		}
		objects = append(objects, obj)
		events = append(events, &v1.Event{World: req.World, Kind: k, Controller: obj.GetController(), Id: obj.GetId()})
	}

	// set some attributes for the span
	pan.SetAttributes(map[string]interface{}{"data": len(req.Data), "world": req.World})

	// write data to the database
	etag := uuid.New()
	_, err := s.db.Set(ctx, req.World, etag, objects)
	if err != nil {
		pan.Err(err)
		return err
	}
	if kind.IsSearchable(k) {
		err = s.sb.Index(ctx, req.World, objects, false)
		if err != nil {
			s.log.Warn().Str("world", req.World).Str("kind", k).Err(err).Msg("failed to index")
		}
	}

	s.publisher <- &publishWork{ctx: ctx, events: events}

	return nil
}

func (s *Service) deleteKind(ctx context.Context, k string, req *api.DeleteRequest, rsp *api.DeleteResponse) error {
	pan := log.NewSpan(ctx, "service.deleteKind", map[string]interface{}{"world": req.World, "kind": k})

	// we don't need a worldspace if it's a global kind
	if kind.IsGlobal(k) {
		req.World = ""
	}

	s.shutdownLock.RLock()
	defer s.shutdownLock.RUnlock()
	if s.shuttingDown {
		return ErrShuttingDown
	}

	// get data from the DB - we need the Controller fields
	result := []map[string]interface{}{}
	err := s.db.Get(ctx, req.World, k, req.Ids, &result)
	if err != nil {
		pan.Err(err)
		return err
	}

	events := []*v1.Event{}
	for _, item := range result {
		id, ok := item["_id"].(string)
		if !ok {
			s.log.Warn().Msg("missing _id field in object")
			continue
		}
		controller, ok := item["_controller"].(string)
		if !ok {
			s.log.Warn().Msg("missing _controller field in object")
			controller = "default"
		}

		// delete from db
		err = s.db.Delete(ctx, req.World, k, id)
		if err != nil {
			return err
		}

		// delete from search
		if kind.IsSearchable(k) {
			err = s.sb.Delete(ctx, req.World, k, id)
			if err != nil {
				s.log.Warn().Str("world", req.World).Str("kind", k).Str("id", id).Err(err).Msg("failed to delete from search")
			}
		}

		// queue event for publishing
		events = append(events, &v1.Event{World: req.World, Kind: k, Id: id, Controller: controller})
	}

	s.publisher <- &publishWork{ctx: ctx, events: events}

	return nil
}
