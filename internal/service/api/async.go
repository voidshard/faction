package api

import (
	"context"
	"fmt"

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
	GetError() *structs.Error
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
		return s.asyncSetWorld(ctx, msg, data)
	case kindDeleteWorld:
		return s.asyncDeleteWorld(ctx, msg, data)
	case kindSetFaction:
		return s.asyncSetFactions(ctx, msg, data)
	case kindDeleteFaction:
		return s.asyncDeleteFaction(ctx, msg, data)
	case kindSetActor:
		return s.asyncSetActors(ctx, msg, data)
	case kindDeleteActor:
		return s.asyncDeleteActor(ctx, msg, data)
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
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Err(encodeErr).Msg("Failed to marshal response")
		return encodeErr
	}

	s.log.Debug().Str("MessageId", m.Id()).Int("bytes", len(outdata)).Msg("SendAPIReply data encoded")
	sendErr := m.Reply(c, outdata)
	if sendErr != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Err(sendErr).Msg("Failed to send reply")
		return sendErr
	} else if err == nil {
		s.log.Debug().Str("MessageId", m.Id()).Err(m.Ack()).Msg("Ack")
	}

	return err
}
func (s *Server) asyncGenericSet(c context.Context, m queue.Message, world string, in []structs.Object, out marshalableReply) error {
	if len(in) == 0 {
		return s.sendApiReply(c, m, out, fmt.Errorf("no objects to set"))
	}

	kind := in[0].Kind()
	log := s.log.With().Int("Count", len(in)).Str("MessageId", m.Id()).Str("kind", kind).Str("world", world).Logger()

	// update db
	_, err := s.db.Set(c, world, m.Id(), in)
	if err != nil {
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
		err = s.qu.PublishChange(c, &structs.Change{World: world, Key: kind, Id: v.GetId()})
		if err != nil {
			log.Error().Err(m.Reject()).Err(err).Str("Id", v.GetId()).Msg("Failed to publish change, rejecting")
			return err
		}
	}

	return s.sendApiReply(c, m, out, nil)
}

func (s *Server) asyncGenericDelete(c context.Context, m queue.Message, world, kind, id string, out marshalableReply) error {
	log := s.log.With().Str("MessageId", m.Id()).Str("kind", kind).Str("world", world).Str("Id", id).Logger()

	// update db
	err := s.db.Delete(c, world, kind, id)
	if err != nil {
		log.Error().Err(m.Ack()).Err(err).Msg("Failed to delete faction, acking")
		return err
	}

	// update search
	err = s.sb.Delete(c, world, kind, id)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to delete faction from search, search will drift")
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: world, Key: kind, Id: id})
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to publish change, rejecting")
		return err
	}

	return s.sendApiReply(c, m, out, nil)
}

func (s *Server) asyncSetWorld(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	out := &structs.SetWorldResponse{Etag: m.Id()}
	in := &structs.SetWorldRequest{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Err(err).Msg("Failed to unmarshal request")
		return s.sendApiReply(c, m, out, err)
	}

	log := s.log.With().Str("MessageId", m.Id()).Str("kind", kindWorld).Str("world", in.Data.Id).Logger()

	// update db
	_, err = s.db.SetWorld(c, m.Id(), in.Data)
	if err != nil {
		log.Error().Err(m.Ack()).Msg("Failed to set world, acking")
		return s.sendApiReply(c, m, out, err)
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.Data.Id, Key: kindWorld, Id: in.Data.Id})
	if err != nil {
		// reject and requeue
		log.Error().Err(err).Err(m.Reject()).Msg("Failed to publish change, rejecting")
		return err
	}

	return s.sendApiReply(c, m, out, err)
}

func (s *Server) asyncDeleteWorld(c context.Context, m queue.Message, data []byte) error {
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
		log.Error().Err(err).Err(m.Ack()).Msg("Failed to delete world, acking")
		return err
	}

	// publish changes
	err = s.qu.PublishChange(c, &structs.Change{World: in.Id, Key: kindWorld, Id: in.Id})
	if err != nil {
		log.Error().Err(m.Reject()).Err(err).Msg("Failed to publish change, rejecting")
		return err
	}

	return s.sendApiReply(c, m, out, nil)
}

func (s *Server) asyncSetFactions(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.SetFactionsRequest{}
	out := &structs.SetFactionsResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		s.log.Error().Err(m.Ack()).Str("MessageId", m.Id()).Str("kind", kindFaction).Err(err).Msg("Failed to unmarshal request, acking")
		return s.sendApiReply(c, m, out, err)
	}

	// tidy data
	//for _, v := range in.Data {
	//	tidyRelations(v, factionMaxRelations)
	//	tidyMemories(v, factionMaxMemories)
	//}

	// perform the set
	objs := make([]structs.Object, len(in.Data))
	for i, v := range in.Data {
		objs[i] = v
	}
	return s.asyncGenericSet(c, m, in.World, objs, out)
}

func (s *Server) asyncDeleteFaction(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.DeleteFactionRequest{}
	out := &structs.DeleteFactionResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		// the message is invalid, we cannot process it so ack it, no requeue
		s.log.Error().Str("MessageId", m.Id()).Str("kind", kindFaction).Err(err).Msg("Failed to unmarshal request")
		return s.sendApiReply(c, m, out, err)
	}

	// perform the delete
	return s.asyncGenericDelete(c, m, in.World, kindFaction, in.Id, out)
}

func (s *Server) asyncSetActors(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.SetActorsRequest{}
	out := &structs.SetActorsResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Err(m.Ack()).Str("kind", kindActor).Err(err).Msg("Failed to unmarshal request, acking")
		return s.sendApiReply(c, m, out, err)
	}

	// tidy data
	//for _, v := range in.Data {
	//	tidyRelations(v, actorMaxRelations)
	//	tidyMemories(v, actorMaxMemories)
	//}

	// perform the set
	objs := make([]structs.Object, len(in.Data))
	for i, v := range in.Data {
		objs[i] = v
	}
	return s.asyncGenericSet(c, m, in.World, objs, out)
}

func (s *Server) asyncDeleteActor(c context.Context, m queue.Message, data []byte) error {
	// decode the request
	in := &structs.DeleteActorRequest{}
	out := &structs.DeleteActorResponse{}
	err := in.UnmarshalJSON(data)
	if err != nil {
		s.log.Error().Str("MessageId", m.Id()).Str("kind", kindActor).Err(err).Msg("Failed to unmarshal request")
		return s.sendApiReply(c, m, out, err)
	}

	return s.asyncGenericDelete(c, m, in.World, kindActor, in.Id, out)
}
