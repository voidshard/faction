package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	defaultRoutines = 5
	defaultMaxAge   = 60 * time.Minute // implies something is horribly wrong
	defaultLimit    = 100
	defaultMaxLimit = 1000
)

type Config struct {
	Routines      int
	MaxMessageAge time.Duration
}

type Service struct {
	cfg *Config

	db db.Database
	qu queue.Queue

	log log.Logger

	workers []*worker
}

// marshalable is an interface for types that can be marshaled and unmarshaled to/from bytes.
// All types we expect to transport over the wire should implement this interface.
type marshalable interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// marshalableReply is an interface for types that can be marshaled and unmarshaled to/from bytes, and can set an error.
// We include an optional error on all Response objects to tell a call what / if things go wrong.
type marshalableReply interface {
	marshalable
	SetError(*structs.Error)
}

// NewService creates a new API service.
func NewService(cfg *Config, db db.Database, qu queue.Queue) *Service {
	// Read functions trigger us to go straight to the database.
	// Write functions are internally queued, though results are returned synchronously.
	if cfg == nil {
		cfg = &Config{Routines: defaultRoutines}
	}
	if cfg.Routines <= 0 {
		cfg.Routines = defaultRoutines
	}
	return &Service{cfg: cfg, db: db, qu: qu, log: log.Sublogger("api-service")}
}

func (s *Service) Worlds(c context.Context, in *structs.GetWorldsRequest) *structs.GetWorldsResponse {
	data, err := s.db.Worlds(c, in.Ids)
	return &structs.GetWorldsResponse{Data: data, Error: toError(err)}
}

func (s *Service) SetWorld(c context.Context, in *structs.SetWorldRequest) *structs.SetWorldResponse {
	obj := &structs.SetWorldResponse{}
	s.genericAsyncRequestResponse(c, in, obj)
	return obj
}

func (s *Service) ListWorlds(c context.Context, in *structs.ListWorldsRequest) *structs.ListWorldsResponse {
	data, err := s.db.ListWorlds(c, in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)))
	return &structs.ListWorldsResponse{Data: data, Error: toError(err)}
}

func (s *Service) DeleteWorld(c context.Context, in *structs.DeleteWorldRequest) *structs.DeleteWorldResponse {
	obj := &structs.DeleteWorldResponse{}
	s.genericAsyncRequestResponse(c, in, obj)
	return obj
}

func (s *Service) Factions(c context.Context, in *structs.GetFactionsRequest) *structs.GetFactionsResponse {
	data, err := s.db.Factions(c, in.World, in.Ids)
	return &structs.GetFactionsResponse{Data: data, Error: toError(err)}
}

func (s *Service) SetFactions(c context.Context, in *structs.SetFactionsRequest) *structs.SetFactionsResponse {
	obj := &structs.SetFactionsResponse{}
	s.genericAsyncRequestResponse(c, in, obj)
	return obj
}

func (s *Service) ListFactions(c context.Context, in *structs.ListFactionsRequest) *structs.ListFactionsResponse {
	data, err := s.db.ListFactions(c, in.World, in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)))
	return &structs.ListFactionsResponse{Data: data, Error: toError(err)}
}

func (s *Service) DeleteFaction(c context.Context, in *structs.DeleteFactionRequest) *structs.DeleteFactionResponse {
	obj := &structs.DeleteFactionResponse{}
	s.genericAsyncRequestResponse(c, in, obj)
	return obj
}

func (s *Service) Actors(c context.Context, in *structs.GetActorsRequest) *structs.GetActorsResponse {
	data, err := s.db.Actors(c, in.World, in.Ids)
	return &structs.GetActorsResponse{Data: data, Error: toError(err)}
}

func (s *Service) SetActors(c context.Context, in *structs.SetActorsRequest) *structs.SetActorsResponse {
	obj := &structs.SetActorsResponse{}
	s.genericAsyncRequestResponse(c, in, obj)
	return obj
}

func (s *Service) ListActors(c context.Context, in *structs.ListActorsRequest) *structs.ListActorsResponse {
	data, err := s.db.ListActors(c, in.World, in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)))
	return &structs.ListActorsResponse{Data: data, Error: toError(err)}
}

func (s *Service) DeleteActor(c context.Context, in *structs.DeleteActorRequest) *structs.DeleteActorResponse {
	obj := &structs.DeleteActorResponse{}
	s.genericAsyncRequestResponse(c, in, obj)
	return obj
}

func (s *Service) sendApiReply(c context.Context, m queue.Message, out marshalable, err error) error {
	outdata, encodeErr := out.Marshal()
	if encodeErr != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(encodeErr).Msg("Failed to marshal response")
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Reject(true)).Msg("Rejecting")
		return encodeErr
	}

	sendErr := m.Reply(c, outdata)
	if sendErr != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(sendErr).Msg("Failed to send reply")
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Reject(true)).Msg("Rejecting")
		return sendErr
	} else if err == nil {
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Ack()).Msg("Ack")
	}

	return err
}

func (s *Service) setWorld(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyWorld.String()

	reply := func(err error) error {
		out := &structs.SetWorldResponse{Etag: m.Id(), Error: toError(err)}
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.SetWorldRequest{}
	err := in.Unmarshal(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(err).Msg("Failed to unmarshal request")
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Ack()).Msg("Ack")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.Data.Id).Logger()

	// update db
	_, err = s.db.SetWorld(c, m.Id(), in.Data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set world")
		if errors.Is(err, db.ErrEtagMismatch) {
			// no retry, Etag mismatch means we cannot apply this update
			log.Debug().Err(m.Ack()).Msg("Acking")
		} else {
			// reject and requeue
			log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		}
		return reply(err)
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.Data.Id, Key: structs.Metakey_KeyWorld, Id: in.Data.Id})
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish change")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		return reply(err)
	}

	return reply(nil)
}

func (s *Service) deleteWorld(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyWorld.String()

	reply := func(err error) error {
		out := &structs.DeleteWorldResponse{Error: toError(err)}
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.DeleteWorldRequest{}
	err := in.Unmarshal(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(err).Msg("Failed to unmarshal request")
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Ack()).Msg("Ack")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.Id).Logger()

	// update db
	err = s.db.DeleteWorld(c, in.Id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete world")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		return reply(err)
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.Id, Key: structs.Metakey_KeyWorld, Id: in.Id})
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish change")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
	}

	return reply(err)
}

func (s *Service) setFactions(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyFaction.String()

	reply := func(err error) error {
		out := &structs.SetFactionsResponse{Etag: m.Id(), Error: toError(err)}
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.SetFactionsRequest{}
	err := in.Unmarshal(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		s.log.Debug().Str("MessageId", m.Id()).Str("key", k).Err(m.Ack()).Msg("Ack")
		return reply(err)
	}

	log := s.log.With().Int("Count", len(in.Data)).Str("MessageId", m.Id()).Str("key", k).Str("world", in.World).Logger()

	// tidy data
	for _, v := range in.Data {
		tidyRelations(v, factionMaxRelations)
		tidyMemories(v, factionMaxMemories)
	}

	// update db
	_, err = s.db.SetFactions(c, in.World, m.Id(), in.Data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set factions")
		if errors.Is(err, db.ErrEtagMismatch) {
			// no retry, Etag mismatch means we cannot apply this update
			log.Debug().Err(m.Ack()).Msg("Acking")
		} else {
			// reject and requeue
			log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		}
		return reply(err)
	}

	// publish changes
	for _, v := range in.Data {
		err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyFaction, Id: v.Id})
		if err != nil {
			log.Error().Err(err).Str("Id", v.Id).Msg("Failed to publish change")
			log.Debug().Err(m.Reject(true)).Str("Id", v.Id).Msg("Rejecting")
			return reply(err)
		}
	}

	return reply(nil)
}

func (s *Service) deleteFaction(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyFaction.String()

	reply := func(err error) error {
		out := &structs.DeleteFactionResponse{Error: toError(err)}
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.DeleteFactionRequest{}
	err := in.Unmarshal(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.World).Str("Id", in.Id).Logger()

	// update db
	err = s.db.DeleteFaction(c, in.World, in.Id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete faction")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		return reply(err)
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyFaction, Id: in.Id})
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish change")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
	}

	return reply(err)
}

func (s *Service) setActors(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyActor.String()

	reply := func(err error) error {
		out := &structs.SetActorsResponse{Etag: m.Id(), Error: toError(err)}
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.SetActorsRequest{}
	err := in.Unmarshal(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Ack()).Msg("Ack")
		return reply(err)
	}

	log := s.log.With().Int("Count", len(in.Data)).Str("MessageId", m.Id()).Str("key", k).Str("world", in.World).Logger()

	// tidy data
	for _, v := range in.Data {
		tidyRelations(v, actorMaxRelations)
		tidyMemories(v, actorMaxMemories)
	}

	// update db
	_, err = s.db.SetActors(c, in.World, m.Id(), in.Data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set actors")
		if errors.Is(err, db.ErrEtagMismatch) {
			log.Debug().Err(m.Ack()).Msg("Acking")
		} else {
			log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		}
		return reply(err)
	}

	// publish changes
	for _, v := range in.Data {
		err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyActor, Id: v.Id})
		if err != nil {
			log.Error().Err(err).Str("Id", v.Id).Msg("Failed to publish change")
			log.Debug().Err(m.Reject(true)).Str("Id", v.Id).Msg("Rejecting")
			return reply(err)
		}
	}

	return reply(nil)
}

func (s *Service) deleteActor(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyActor.String()

	reply := func(err error) error {
		out := &structs.DeleteActorResponse{Error: toError(err)}
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.DeleteActorRequest{}
	err := in.Unmarshal(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.World).Str("Id", in.Id).Logger()

	// update db
	err = s.db.DeleteActor(c, in.World, in.Id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete actor")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
		return reply(err)
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyActor, Id: in.Id})
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish change")
		log.Debug().Err(m.Reject(true)).Msg("Rejecting")
	}

	return reply(err)
}

func (s *Service) start() error {
	if s.workers != nil {
		s.stop()
	}

	s.workers = make([]*worker, s.cfg.Routines)

	s.log.Debug().Int("Routines", s.cfg.Routines).Msg("Starting workers")
	for i := 0; i < s.cfg.Routines; i++ {
		log.Debug().Int("Worker", i).Msg("Starting worker")

		qsub, err := s.qu.SubscribeApiReq()
		if err != nil {
			defer s.stop()
			log.Error().Err(err).Msg("Failed to subscribe to queue")
			return err
		}

		wrk := newWorker(fmt.Sprintf("api-worker-%d", i), s, qsub)
		s.workers[i] = wrk

		go func(w *worker) {
			w.Run()
		}(wrk)
	}

	return nil
}

func (s *Service) stop() {
	if s.workers == nil {
		return
	}
	for _, w := range s.workers {
		w.Kill()
	}
	s.workers = nil
}

func (s *Service) genericAsyncRequestResponse(c context.Context, in marshalable, out marshalableReply) {
	pan := log.NewSpan(c, "api.genericAsyncRequestResponse")
	defer pan.End()

	data, err := encodeRequest(in)
	if err != nil {
		out.SetError(toError(err))
		log.Error().Err(err).Msg("Failed to encode request")
		return
	}

	// publish the request
	rchan, err := s.qu.PublishApiReq(c, data)
	if err != nil {
		out.SetError(toError(err))
		log.Error().Err(err).Msg("Failed to publish request")
		return
	}

	reply := <-rchan.Channel()
	if reply == nil {
		out.SetError(&structs.Error{Message: "no data in reply", Code: 500})
		log.Error().Msg("Nil reply received")
		return
	}
	if reply.Data == nil || len(reply.Data()) == 0 {
		return // no data in reply
	}

	log.Debug().Str("MessageId", reply.Id()).Str("DATA", string(reply.Data())).Msg("Received reply")

	err = out.Unmarshal(reply.Data())
	if err != nil {
		out.SetError(toError(err))
		log.Error().Err(err).Msg("Failed to unmarshal reply")
		return
	}
}

/*
func (s *Service) Jobs(c context.Context, id ...string) ([]*structs.Job, error) {
	return s.db.Jobs(c, id)
}

func (s *Service) SetJob(c context.Context, in *structs.Job) (string, error) {
	// TODO: validate / check stuff

	etag, err := s.db.SetJob(c, in)
	if err != nil {
		return etag, err
	}

	return etag, err
}

func (s *Service) OnChange(pattern *structs.Change) (*OnChangeSubscription, error) {
	return newOnChangeSubscription(s.qu, pattern)
}
*/
