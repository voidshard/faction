package api

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/search"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
)

type backgroundWorker interface {
	Kill()
}

type Server struct {
	structs.UnimplementedAPIServer

	cfg *Config

	db db.Database
	qu *Queue
	sb search.Search

	srv *grpc.Server
	log log.Logger

	tickManager *tickManager
	workers     []backgroundWorker
	apiSub      queue.Subscription
}

func NewServer(cfg *Config, db db.Database, qu queue.Queue, sb search.Search) (*Server, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.setDefaults()

	apiQueue, err := newQueue(qu)
	if err != nil {
		return nil, err
	}

	srv := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	me := &Server{
		cfg: cfg,
		db:  db,
		qu:  apiQueue,
		sb:  sb,
		srv: srv,
		log: log.Sublogger("api"),
	}

	structs.RegisterAPIServer(srv, me)
	return me, nil
}

func (s *Server) Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.log.Error().Int("port", port).Err(err).Msg("Failed to listen")
		return err
	}
	s.log.Info().Msg("Serving...")
	err = s.startWorkers()
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.log.Info().Msg("...Exiting")
	s.stopWorkers()
	s.srv.Stop()
}

func (s *Server) DeferChange(ctx context.Context, in *structs.DeferChangeRequest) (*structs.DeferChangeResponse, error) {
	pan := log.NewSpan(ctx, "api.DeferChange")
	defer pan.End()

	// validate the request
	err := validDeferChange(in)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.DeferChangeResponse{Error: toError(err)}, nil
	}

	pan.SetAttributes(map[string]interface{}{
		"world": in.Data.World, "id": in.Data.Id, "key": in.Data.Key, "to-tick": *in.ToTick, "by-tick": *in.ByTick,
	})

	// we've been given the specific tick to defer to (easy case)
	if *in.ToTick > 0 {
		err := s.qu.DeferChange(ctx, in.Data, *in.ToTick)
		pan.Err(err)
		return &structs.DeferChangeResponse{Error: toError(err)}, nil
	}

	// read world tick out of tick manager: since it keeps (mostly) up-to-date with the DB objects
	// nb. it actually doesn't matter if we defer 1 or 2 behind the current tick because the tick manager monitors
	// the last few world-tick queues to catch late comers.
	tick, err := s.tickManager.Tick(in.Data.World)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get world tick")
		pan.Err(err)
		return &structs.DeferChangeResponse{Error: toError(err)}, nil
	}

	err = s.qu.DeferChange(ctx, in.Data, tick+*in.ByTick)
	pan.Err(err)
	return &structs.DeferChangeResponse{Error: toError(err)}, nil
}

func (s *Server) Worlds(ctx context.Context, in *structs.GetWorldsRequest) (*structs.GetWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.Worlds", map[string]interface{}{"IdCount": len(in.Ids)})
	defer pan.End()

	err := validIDs(kindWorld, in)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.GetWorldsResponse{Error: toError(err)}, nil
	}

	data, err := s.db.Worlds(ctx, in.Ids)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get worlds")
		pan.Err(err)
	}
	return &structs.GetWorldsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetWorld(ctx context.Context, in *structs.SetWorldRequest) (*structs.SetWorldResponse, error) {
	obj := &structs.SetWorldResponse{}
	s.genericSet(ctx, kindWorld, in.Data.Id, []structs.Object{in.Data}, in, obj)
	return obj, nil
}

func (s *Server) DeleteWorld(ctx context.Context, in *structs.DeleteWorldRequest) (*structs.DeleteWorldResponse, error) {
	obj := &structs.DeleteWorldResponse{}
	s.genericDelete(ctx, kindWorld, in.Id, in.Id, in, obj)
	return obj, nil
}

func (s *Server) ListWorlds(ctx context.Context, in *structs.ListWorldsRequest) (*structs.ListWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListWorlds", map[string]interface{}{"Limit": in.Limit, "Offset": in.Offset})
	defer pan.End()
	data, err := s.db.ListWorlds(ctx, in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)))
	if err != nil {
		pan.Err(err)
	}
	return &structs.ListWorldsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) genericGet(ctx context.Context, kind, world string, ids hasIDs, data interface{}) error {
	pan := log.NewSpan(ctx, fmt.Sprintf("api.%s", kind), map[string]interface{}{"IdCount": len(ids.GetIds()), "World": world})
	defer pan.End()

	err := validIDs(kind, ids)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return err
	}

	return s.db.Get(ctx, world, kind, ids.GetIds(), data)
}

func (s *Server) genericSet(ctx context.Context, kind, world string, objects interface{}, req marshalable, resp marshalableReply) {
	in, ok := objects.([]interface{})
	if !ok {
		in = []interface{}{objects}
	}

	pan := log.NewSpan(ctx, fmt.Sprintf("api.Set%s", kind), map[string]interface{}{"IdCount": len(in), "World": world})
	defer pan.End()

	if kind != kindWorld {
		_, err := s.tickManager.Tick(world)
		if errors.Is(err, ErrNotFound) {
			resp.SetError(toError(err))
			return
		} else if err != nil {
			s.log.Error().Str("World", world).Err(err).Msg("Failed to get world tick")
			pan.Err(err)
			resp.SetError(toError(err))
			return
		}
	}

	for _, raw := range in {
		o, ok := raw.(structs.Object)
		if !ok {
			s.log.Warn().Msg("Failed to cast given interface{} as structs.Object")
			continue
		}
		err := validObject(o)
		if err != nil {
			s.log.Warn().Err(err).Msg("Invalid request")
			pan.Err(err)
			resp.SetError(toError(err))
			return
		}
	}

	s.genericRequestResponse(ctx, req, resp)
	if resp.GetError() != nil {
		err := fmt.Errorf(resp.GetError().Message)
		s.log.Error().Err(err).Msg("Failed to set objects")
		pan.Err(err)
		return
	}

	return
}

func (s *Server) genericDelete(ctx context.Context, kind, world, id string, req marshalable, resp marshalableReply) {
	pan := log.NewSpan(ctx, fmt.Sprintf("api.Delete%s", kind), map[string]interface{}{"Id": id, "World": world})
	defer pan.End()

	err := validID(kind, id)
	if err != nil {
		err := fmt.Errorf("%w %v", ErrInvalid, id)
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		resp.SetError(toError(err))
		return
	}

	s.genericRequestResponse(ctx, req, resp)
	if resp.GetError() != nil {
		err := fmt.Errorf(resp.GetError().Message)
		s.log.Error().Err(err).Msg("Failed to delete object")
		pan.Err(err)
		return
	}

	return
}

func (s *Server) genericList(ctx context.Context, kind, world string, limit, offset *uint64, labels map[string]string, data interface{}) error {
	pan := log.NewSpan(ctx, fmt.Sprintf("api.List%s", kind), map[string]interface{}{"Limit": limit, "Offset": offset, "World": world})
	defer pan.End()
	err := s.db.List(ctx, world, kind, labels, int64(clamp(limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(offset, 0)), data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list objects")
		pan.Err(err)
	}
	return err
}

func (s *Server) Factions(ctx context.Context, in *structs.GetFactionsRequest) (*structs.GetFactionsResponse, error) {
	data := []*structs.Faction{}
	err := s.genericGet(ctx, kindFaction, in.World, in, &data)
	return &structs.GetFactionsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetFactions(ctx context.Context, in *structs.SetFactionsRequest) (*structs.SetFactionsResponse, error) {
	obj := &structs.SetFactionsResponse{}
	s.genericSet(ctx, kindFaction, in.World, in.Data, in, obj)
	return obj, nil
}

func (s *Server) DeleteFaction(ctx context.Context, in *structs.DeleteFactionRequest) (*structs.DeleteFactionResponse, error) {
	out := &structs.DeleteFactionResponse{}
	s.genericDelete(ctx, kindFaction, in.World, in.Id, in, out)
	return out, nil
}

func (s *Server) ListFactions(ctx context.Context, in *structs.ListFactionsRequest) (*structs.ListFactionsResponse, error) {
	data := []*structs.Faction{}
	err := s.genericList(ctx, kindFaction, in.World, in.Limit, in.Offset, in.Labels, &data)
	return &structs.ListFactionsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) Actors(ctx context.Context, in *structs.GetActorsRequest) (*structs.GetActorsResponse, error) {
	data := []*structs.Actor{}
	err := s.genericGet(ctx, kindActor, in.World, in, &data)
	return &structs.GetActorsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetActors(ctx context.Context, in *structs.SetActorsRequest) (*structs.SetActorsResponse, error) {
	data := &structs.SetActorsResponse{}
	s.genericSet(ctx, kindActor, in.World, in.Data, in, data)
	return data, nil
}

func (s *Server) DeleteActor(ctx context.Context, in *structs.DeleteActorRequest) (*structs.DeleteActorResponse, error) {
	data := &structs.DeleteActorResponse{}
	s.genericDelete(ctx, kindActor, in.World, in.Id, in, data)
	return data, nil
}

func (s *Server) ListActors(ctx context.Context, in *structs.ListActorsRequest) (*structs.ListActorsResponse, error) {
	data := []*structs.Actor{}
	err := s.genericList(ctx, kindActor, in.World, in.Limit, in.Offset, in.Labels, &data)
	return &structs.ListActorsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) Races(ctx context.Context, in *structs.GetRacesRequest) (*structs.GetRacesResponse, error) {
	data := []*structs.Race{}
	err := s.genericGet(ctx, kindRace, in.World, in, &data)
	return &structs.GetRacesResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetRaces(ctx context.Context, in *structs.SetRacesRequest) (*structs.SetRacesResponse, error) {
	data := &structs.SetRacesResponse{}
	s.genericSet(ctx, kindRace, in.World, in.Data, in, data)
	return data, nil
}

func (s *Server) OnChange(in *structs.OnChangeRequest, stream structs.API_OnChangeServer) error {
	// validate the request
	err := validOnChange(in)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		stream.Send(&structs.OnChangeResponse{Error: toError(err)})
		return err
	}

	// setup log & trace
	peerAddr := ""
	peer, ok := peer.FromContext(stream.Context())
	if ok {
		peerAddr = peer.Addr.String()
	}

	attrs := map[string]interface{}{
		"world": in.Data.World,
		"key":   in.Data.Key,
		"type":  in.Data.Type,
		"id":    in.Data.Id,
		"queue": in.Queue,
		"peer":  peerAddr,
	}
	l := log.Sublogger("api.OnChange", attrs)

	// start the change listener
	l.Debug().Msg("Starting change listener")

	sub, err := s.qu.SubscribeChange(in.Data, in.Queue)
	if err != nil {
		l.Error().Err(err).Msg("Failed to start change listener")
		return err
	}

	defer func() {
		l.Debug().Msg("Closing change listener")
		sub.Close()
	}()

	// while the stream is open, send changes to client
	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			s.log.Debug().Err(err).Msg("OnChange client disconnected")
			return err
		case msg, ok := <-sub.Channel():
			if !ok {
				return nil
			}
			l.Debug().Str("MessageId", msg.Id()).Msg("Received change")
			pan := log.NewSpan(msg.Context(), "api.OnChange", attrs, map[string]interface{}{"MessageId": msg.Id()})

			change := &structs.Change{}
			err := change.UnmarshalJSON(msg.Data())
			if err != nil {
				if in.Queue != "" {
					msg.Ack() // message is borked
				}
				l.Error().Err(err).Msg("Failed to unmarshal change")
				pan.Err(err)
				pan.End()
				continue
			}

			ackId := ""
			if in.Queue != "" {
				ackId = s.qu.DeferAck(msg, change)
			}

			err = stream.Send(&structs.OnChangeResponse{Data: change, Ack: ackId, Error: toError(err)})
			if err != nil {
				l.Error().Err(err).Msg("Failed to send change")
				pan.Err(err)
			}
			pan.End()
		}
	}
}

// AckStream is a stream that receives acks from the client.
// We pass these to the Queue to deal with. Messages will be dealt with locally if
// possible or published into a topic so whomever sent them can deal with them.
func (s *Server) AckStream(stream grpc.ClientStreamingServer[structs.AckRequest, structs.AckResponse]) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			s.log.Debug().Err(err).Msg("AckStream client disconnected")
			return err
		}
		for _, ackId := range msg.Ack {
			if ackId == "" {
				continue
			}
			err = s.qu.Ack(ackId)
			s.log.Debug().Err(err).Str("AckId", ackId).Msg("Received ack from client")
		}
	}
}

// genericRequestResponse is a helper function to send a request to the queue and wait for a response.
// We do this for all write requests (Set & Delete) sent to the API Server.
//
// This allows us to retry, reject, ack or nack requests as internally as required.
// Ie. we can ensure that we write to the database and emit to the queue, if either fails we can reject and
// retry the request (the db write will be idempotent if the object Etag has not changed in the meantime).
//
// This causes a worker routine running on an API Server (even us) to call asyncAPIRequest defined in
// async.go -- from here the relevant functions defined in async.go are called.
func (s *Server) genericRequestResponse(c context.Context, in marshalable, out marshalableReply) {
	pan := log.NewSpan(c, "api.genericRequestResponse")
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
	rchan, err := s.qu.EnqueueApiReq(c, data)
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

func (s *Server) startWorkers() error {
	if s.workers != nil {
		s.stopWorkers()
	}
	if s.apiSub != nil {
		s.apiSub.Close()
	}

	s.workers = make([]backgroundWorker, s.cfg.WorkersAPI+1)

	qsub, err := s.qu.DequeueApiReq()
	if err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to queue")
		return err
	}
	s.apiSub = qsub

	s.log.Debug().Msg("Starting tick manager")
	tc, err := newTickManager("tick-manager", s.db, s.qu)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create tick cache")
		return err
	}
	go tc.Run()
	s.workers[0] = tc
	s.tickManager = tc

	s.log.Debug().Int("Routines", s.cfg.WorkersAPI).Msg("Starting api workers")
	for i := 0; i < s.cfg.WorkersAPI; i++ {
		log.Debug().Int("Worker", i).Msg("Starting worker")

		wrk := newApiWorker(fmt.Sprintf("api-worker-%d", i), s, qsub)
		s.workers[i+1] = wrk

		go func(w *apiWorker) {
			w.Run()
		}(wrk)
	}

	return nil
}

func (s *Server) stopWorkers() {
	if s.workers == nil {
		return
	}
	for _, w := range s.workers {
		if w != nil {
			w.Kill()
		}
	}
	s.workers = nil
	if s.apiSub != nil {
		s.apiSub.Close()
	}
	s.apiSub = nil
}
