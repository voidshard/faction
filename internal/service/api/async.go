package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
)

// marshalable is an interface for types that can be marshaled and unmarshaled to/from bytes.
// All types we expect to transport over the wire should implement this interface.
type marshalable interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// marshalableReply is an interface for types that can be marshaled and unmarshaled to/from bytes, and can set an error.
// We include an optional error on all Response objects to tell a call what / if things go wrong.
type marshalableReply interface {
	marshalable
	SetError(*structs.Error)
}

func (s *Service) sendApiReply(c context.Context, m queue.Message, out marshalable, err error) error {
	s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("sendApiReply start")
	defer s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("sendApiReply done")

	outdata, encodeErr := out.MarshalJSON()
	if encodeErr != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Reject()).Err(encodeErr).Msg("Failed to marshal response")
		return encodeErr
	}

	s.log.Debug().Str("MessageId", m.Id()).Int("bytes", len(outdata)).Msg("SendAPIReply data encoded")
	sendErr := m.Reply(c, outdata)
	if sendErr != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Reject()).Err(sendErr).Msg("Failed to send reply")
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
		s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("setWorld response")
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.SetWorldRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Err(err).Msg("Failed to unmarshal request")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.Data.Id).Logger()

	// update db
	_, err = s.db.SetWorld(c, m.Id(), in.Data)
	if err != nil {
		if db.ErrorRetryable(err) {
			// reject and requeue
			log.Error().Err(m.Reject()).Msg("Failed to set world, rejecting")
			return err
		}
		// no retry, Etag mismatch means we cannot apply this update
		log.Error().Err(m.Ack()).Msg("Failed to set world, acking")
		return reply(err)
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.Data.Id, Key: structs.Metakey_KeyWorld, Id: in.Data.Id})
	if err != nil {
		// reject and requeue
		log.Error().Err(err).Err(m.Reject()).Msg("Failed to publish change, rejecting")
		return err
	}

	return reply(nil)
}

func (s *Service) deleteWorld(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyWorld.String()

	reply := func(err error) error {
		out := &structs.DeleteWorldResponse{Error: toError(err)}
		s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("deleteWorld response")
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.DeleteWorldRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(err).Err(m.Ack()).Msg("Failed to unmarshal request, acking")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.Id).Logger()

	// update db
	err = s.db.DeleteWorld(c, in.Id)
	if err != nil {
		log.Error().Err(err).Err(m.Reject()).Msg("Failed to delete world, rejecting")
		return err
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.Id, Key: structs.Metakey_KeyWorld, Id: in.Id})
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to publish change, rejecting")
		return err
	}

	return reply(err)
}

func (s *Service) setFactions(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyFaction.String()

	reply := func(err error) error {
		out := &structs.SetFactionsResponse{Etag: m.Id(), Error: toError(err)}
		s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("setFactions response")
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.SetFactionsRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		s.log.Error().Err(m.Ack()).Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request, acking")
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
		if db.ErrorRetryable(err) {
			// reject and requeue
			log.Error().Err(m.Reject()).Err(err).Msg("Failed to set factions, rejecting")
			return err
		}
		// no retry, Etag mismatch means we cannot apply this update
		log.Error().Err(m.Ack()).Err(err).Msg("Failed to set factions, acking")
		return reply(err)
	}

	// update search
	err = s.sb.IndexFactions(c, in.World, in.Data, s.cfg.FlushSearch)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to index factions, search will drift")
	}

	// publish changes
	for _, v := range in.Data {
		err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyFaction, Id: v.Id})
		if err != nil {
			log.Error().Err(m.Reject()).Err(err).Str("Id", v.Id).Msg("Failed to publish change, rejecting")
			return err
		}
	}

	return reply(nil)
}

func (s *Service) deleteFaction(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyFaction.String()

	reply := func(err error) error {
		out := &structs.DeleteFactionResponse{Error: toError(err)}
		s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("deleteFaction response")
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.DeleteFactionRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.World).Str("Id", in.Id).Logger()

	// update db
	err = s.db.DeleteFaction(c, in.World, in.Id)
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to delete faction, rejecting")
		return err
	}

	// update search
	err = s.sb.DeleteFaction(c, in.World, in.Id)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to delete faction from search, search will drift")
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyFaction, Id: in.Id})
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to publish change, rejecting")
		return err
	}

	return reply(err)
}

func (s *Service) setActors(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyActor.String()

	reply := func(err error) error {
		out := &structs.SetActorsResponse{Etag: m.Id(), Error: toError(err)}
		s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("setActors response")
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.SetActorsRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Str("key", k).Err(err).Msg("Failed to unmarshal request, acking")
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
		if db.ErrorRetryable(err) {
			log.Error().Err(m.Reject()).Err(err).Msg("Failed to set actors, rejecting")
			return err
		}
		log.Error().Err(m.Ack()).Err(err).Msg("Failed to set actors, acking")
		return reply(err)
	}

	// update search
	err = s.sb.IndexActors(c, in.World, in.Data, s.cfg.FlushSearch)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to index actors, search will drift")
	}

	// publish changes
	for _, v := range in.Data {
		err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyActor, Id: v.Id})
		if err != nil {
			log.Error().Err(m.Reject()).Err(err).Str("Id", v.Id).Msg("Failed to publish change, rejecting")
			return err
		}
	}

	return reply(nil)
}

func (s *Service) deleteActor(c context.Context, m queue.Message, data []byte) error {
	k := structs.Metakey_KeyActor.String()

	reply := func(err error) error {
		out := &structs.DeleteActorResponse{Error: toError(err)}
		s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("deleteActor response")
		return s.sendApiReply(c, m, out, err)
	}

	// decode the request
	in := &structs.DeleteActorRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		return reply(err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("key", k).Str("world", in.World).Str("Id", in.Id).Logger()

	// update db
	err = s.db.DeleteActor(c, in.World, in.Id)
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to delete actor, rejecting")
		return err
	}

	// update search
	err = s.sb.DeleteActor(c, in.World, in.Id)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to delete actor from search, search will drift")
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.World, Key: structs.Metakey_KeyActor, Id: in.Id})
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to publish change, rejecting")
		return err
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
		pan.Err(err)
		return
	}

	// publish the request
	s.log.Debug().Int("Bytes", len(data)).Msg("Publishing request")
	rchan, err := s.qu.PublishApiReq(c, data)
	if err != nil {
		out.SetError(toError(err))
		log.Error().Err(err).Msg("Failed to publish request")
		pan.Err(err)
		return
	}

	// wait for a reply
	reply := <-rchan.Channel()
	if reply == nil {
		err := errors.New("no data in reply")
		out.SetError(&structs.Error{Message: err.Error(), Code: 500})
		log.Error().Err(err).Msg("Nil reply received")
		pan.Err(err)
		return
	}
	if reply.Data == nil || len(reply.Data()) == 0 {
		log.Debug().Str("MessageId", reply.Id()).Msg("Received reply with no data")
		return // no data in reply
	}

	log.Debug().Str("MessageId", reply.Id()).Int("Bytes", len(reply.Data())).Msg("Received reply")

	err = out.UnmarshalJSON(reply.Data())
	if err == nil {
		return // seems legit
	}

	// before we complain that the message is invalid, try to unmarshal it as an error
	errStruct := &structs.Error{}
	suberr := errStruct.UnmarshalJSON(reply.Data())
	if suberr == nil { // ok something is telling us it broke
		out.SetError(errStruct)
		log.Error().Err(suberr).Msg("unexpected error message returned")
		pan.Err(suberr)
		return
	}

	// um, no idea what happened
	out.SetError(toError(err))
	log.Error().Err(err).Msg("Failed to unmarshal reply")
	pan.Err(err)
	return
}
