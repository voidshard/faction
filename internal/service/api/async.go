package api

import (
	"context"
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
)

// marshalable is an interface for types that can be marshaled and unmarshaled to/from bytes.
// All types we expect to transport over the wire should implement this interface.
type marshalable interface {
	Kind() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// marshalableReply is an interface for types that can be marshaled and unmarshaled to/from bytes, and can set an error.
// We include an optional error on all Response objects to tell a call what / if things go wrong.
type marshalableReply interface {
	marshalable
	SetError(*structs.Error)
}

// asyncAPIRequest is a helper function to process a message from the queue.
// It logically does the inverse of genericAsyncRequestResponse, calling the appropriate function to handle the request.
func (s *Server) asyncAPIRequest(ctx context.Context, msg queue.Message) error {
	kind, data, err := decodeRequest(msg.Data())
	if err != nil {
		return err
	}
	switch kind {
	case kindSetWorld:
		return s.setWorld(ctx, msg, data)
	case kindSetFaction:
		return s.setFactions(ctx, msg, data)
	case kindSetActor:
		return s.setActors(ctx, msg, data)
	default:
		return fmt.Errorf("unknown request kind: %s", kind)
	}
}

func (s *Server) sendApiReply(c context.Context, m queue.Message, out marshalableReply, err error) error {
	s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("sendApiReply start")
	defer s.log.Debug().Str("MessageId", m.Id()).Err(err).Msg("sendApiReply done")

	out.SetError(toError(err))

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
func (s *Server) genericSet(c context.Context, m queue.Message, key structs.Metakey, world string, in []structs.Object, out marshalableReply) error {
	if len(in) == 0 {
		return s.sendApiReply(c, m, out, fmt.Errorf("no objects to set"))
	}

	kind := in[0].Kind()
	log := s.log.With().Int("Count", len(in)).Str("MessageId", m.Id()).Str("kind", kind).Str("world", world).Logger()

	// update db
	_, err := s.db.Set(c, world, m.Id(), in)
	if err != nil {
		if db.ErrorRetryable(err) {
			// reject and requeue
			log.Error().Err(m.Reject()).Err(err).Msg("Failed to set factions, rejecting")
			return err
		}
		// no retry, Etag mismatch means we cannot apply this update
		log.Error().Err(m.Ack()).Err(err).Msg("Failed to set factions, acking")
		return s.sendApiReply(c, m, out, err)
	}

	// update search
	err = s.sb.Index(c, world, in, s.cfg.FlushSearch)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to index factions, search will drift")
	}

	// publish changes
	for _, v := range in {
		err = s.qu.PublishChange(c, &structs.Change{World: world, Key: key, Id: v.GetId()})
		if err != nil {
			log.Error().Err(m.Reject()).Err(err).Str("Id", v.GetId()).Msg("Failed to publish change, rejecting")
			return err
		}
	}

	return s.sendApiReply(c, m, out, nil)
}

func (s *Server) genericDelete(c context.Context, m queue.Message, key structs.Metakey, world, kind, id string, out marshalableReply) error {
	log := s.log.With().Str("MessageId", m.Id()).Str("kind", kind).Str("world", world).Str("Id", id).Logger()

	// update db
	err := s.db.Delete(c, world, kind, id)
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to delete faction, rejecting")
		return err
	}

	// update search
	err = s.sb.Delete(c, world, kind, id)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to delete faction from search, search will drift")
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: world, Key: structs.Metakey_KeyFaction, Id: id})
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to publish change, rejecting")
		return err
	}

	return s.sendApiReply(c, m, out, nil)
}

func (s *Server) setWorld(c context.Context, m queue.Message, data []byte) error {
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

func (s *Server) deleteWorld(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.DeleteWorldRequest{}
	out := &structs.DeleteWorldResponse{}

	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(err).Err(m.Ack()).Msg("Failed to unmarshal request, acking")
		return s.sendApiReply(c, m, out, err)
	}
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

	return s.sendApiReply(c, m, out, nil)
}

func (s *Server) setFactions(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.SetFactionsRequest{}
	out := &structs.SetFactionsResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		k := structs.Metakey_KeyFaction.String()
		s.log.Error().Err(m.Ack()).Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request, acking")
		return s.sendApiReply(c, m, out, err)
	}

	// tidy data
	for _, v := range in.Data {
		tidyRelations(v, factionMaxRelations)
		tidyMemories(v, factionMaxMemories)
	}

	// perform the set
	objs := make([]structs.Object, len(in.Data))
	for i, v := range in.Data {
		objs[i] = v
	}
	return s.genericSet(c, m, structs.Metakey_KeyFaction, in.World, objs, out)
}

func (s *Server) deleteFaction(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.DeleteFactionRequest{}
	out := &structs.DeleteFactionResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		k := structs.Metakey_KeyFaction.String()
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		return s.sendApiReply(c, m, out, err)
	}

	// perform the delete
	return s.genericDelete(c, m, structs.Metakey_KeyFaction, in.World, kindFaction, in.Id, out)
}

func (s *Server) setActors(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.SetActorsRequest{}
	out := &structs.SetActorsResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		k := structs.Metakey_KeyActor.String()
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Str("key", k).Err(err).Msg("Failed to unmarshal request, acking")
		return s.sendApiReply(c, m, out, err)
	}

	// tidy data
	for _, v := range in.Data {
		tidyRelations(v, actorMaxRelations)
		tidyMemories(v, actorMaxMemories)
	}

	// perform the set
	objs := make([]structs.Object, len(in.Data))
	for i, v := range in.Data {
		objs[i] = v
	}
	return s.genericSet(c, m, structs.Metakey_KeyActor, in.World, objs, out)
}

func (s *Server) deleteActor(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.DeleteActorRequest{}
	out := &structs.DeleteActorResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		k := structs.Metakey_KeyActor.String()
		s.log.Error().Str("MessageId", m.Id()).Str("key", k).Err(err).Msg("Failed to unmarshal request")
		return s.sendApiReply(c, m, out, err)
	}

	return s.genericDelete(c, m, structs.Metakey_KeyActor, in.World, kindActor, in.Id, out)
}
