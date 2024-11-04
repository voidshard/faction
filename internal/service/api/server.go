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
	"github.com/voidshard/faction/pkg/util/uuid"
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

	key, _ := structs.Metakey_name[int32(in.Data.Key)] // validated above
	pan.SetAttributes(map[string]interface{}{"world": in.Data.World, "id": in.Data.Id, "key": key, "to-tick": *in.ToTick, "by-tick": *in.ByTick})

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
	pan := log.NewSpan(ctx, "api.Worlds", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()

	err := validIDs(in)
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
	pan := log.NewSpan(ctx, "api.SetWorld", map[string]interface{}{"id": in.Data.Id})
	defer pan.End()
	obj := &structs.SetWorldResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		pan.Err(fmt.Errorf(obj.Error.Message))
	}
	return obj, nil
}

func (s *Server) DeleteWorld(ctx context.Context, in *structs.DeleteWorldRequest) (*structs.DeleteWorldResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteWorld", map[string]interface{}{"id": in.Id})
	defer pan.End()

	if in == nil || !uuid.IsValidUUID(in.Id) {
		err := fmt.Errorf("%w %v", ErrInvalid, in)
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.DeleteWorldResponse{Error: toError(err)}, nil
	}

	obj := &structs.DeleteWorldResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		err := fmt.Errorf(obj.Error.Message)
		s.log.Error().Err(err).Msg("Failed to delete world")
		pan.Err(err)
	}
	return obj, nil
}

func (s *Server) ListWorlds(ctx context.Context, in *structs.ListWorldsRequest) (*structs.ListWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListWorlds", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	data, err := s.db.ListWorlds(ctx, in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)))
	if err != nil {
		pan.Err(err)
	}
	return &structs.ListWorldsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) Factions(ctx context.Context, in *structs.GetFactionsRequest) (*structs.GetFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.Factions", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()

	err := validIDs(in)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.GetFactionsResponse{Error: toError(err)}, nil
	}

	data := []*structs.Faction{}
	err := s.db.Get(ctx, in.World, "Faction", in.Ids, &data)
	return &structs.GetFactionsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetFactions(ctx context.Context, in *structs.SetFactionsRequest) (*structs.SetFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.SetFactions", map[string]interface{}{"id-count": len(in.Data)})
	defer pan.End()
	obj := &structs.SetFactionsResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		pan.Err(fmt.Errorf(obj.Error.Message))
	}
	return obj, nil
}

func (s *Server) DeleteFaction(ctx context.Context, in *structs.DeleteFactionRequest) (*structs.DeleteFactionResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteFaction", map[string]interface{}{"id": in.Id})
	defer pan.End()

	if in == nil || !uuid.IsValidUUID(in.Id) {
		err := fmt.Errorf("%w %v", ErrInvalid, in)
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.DeleteFactionResponse{Error: toError(err)}, nil
	}

	obj := &structs.DeleteFactionResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	return obj, nil
}

func (s *Server) ListFactions(ctx context.Context, in *structs.ListFactionsRequest) (*structs.ListFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListFactions", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	data := []*structs.Faction{}
	err := s.db.List(ctx, in.World, "Faction", in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)), &data)
	if err != nil {
		pan.Err(err)
	}
	return &structs.ListFactionsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) Actors(ctx context.Context, in *structs.GetActorsRequest) (*structs.GetActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.Actors", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()

	err := validIDs(in)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.GetActorsResponse{Error: toError(err)}, nil
	}

	data := []*structs.Actor{}
	err := s.db.Get(ctx, in.World, "Actor", in.Ids, &data)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get actors")
		pan.Err(err)
	}
	return &structs.GetActorsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetActors(ctx context.Context, in *structs.SetActorsRequest) (*structs.SetActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.SetActor", map[string]interface{}{"id-count": len(in.Data)})
	defer pan.End()
	obj := &structs.SetActorsResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		err := fmt.Errorf(obj.Error.Message)
		s.log.Error().Err(err).Msg("Failed to set actors")
		pan.Err(err)
	}
	return obj, nil
}

func (s *Server) DeleteActor(ctx context.Context, in *structs.DeleteActorRequest) (*structs.DeleteActorResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteActor", map[string]interface{}{"id": in.Id})
	defer pan.End()

	if in == nil || !uuid.IsValidUUID(in.Id) {
		err := fmt.Errorf("%w %v", ErrInvalid, in)
		s.log.Warn().Err(err).Msg("Invalid request")
		pan.Err(err)
		return &structs.DeleteActorResponse{Error: toError(err)}, nil
	}

	obj := &structs.DeleteActorResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		err := fmt.Errorf(obj.Error.Message)
		s.log.Error().Err(err).Msg("Failed to delete actor")
		pan.Err(err)
	}
	return obj, nil
}

func (s *Server) ListActors(ctx context.Context, in *structs.ListActorsRequest) (*structs.ListActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListActors", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	data := []*structs.Actor{}
	err := s.db.List(ctx, in.World, "Actor", in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)), &data)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to list actors")
		pan.Err(err)
	}
	return &structs.ListActorsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) OnChange(in *structs.OnChangeRequest, stream structs.API_OnChangeServer) error {
	// validate the request
	err := validOnChange(in)
	if err != nil {
		s.log.Warn().Err(err).Msg("Invalid request")
		strean.Send(&structs.OnChangeResponse{Error: toError(err)})
		return err
	}

	key, _ := structs.Metakey_name[int32(in.Data.Key)] // validated above

	// setup log & trace
	peerAddr := ""
	peer, ok := peer.FromContext(stream.Context())
	if ok {
		peerAddr = peer.Addr.String()
	}

	attrs := map[string]interface{}{
		"world": in.Data.World,
		"area":  in.Data.Area,
		"key":   key,
		"id":    in.Data.Id,
		"queue": in.Queue,
		"peer":  peerAddr,
	}
	l := log.Sublogger("api.OnChange", attrs)

	// start the change listener
	l.Debug().Msg("Starting change listener")
	defer l.Debug().Msg("Stopped change listener")
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
			pan := log.NewSpan(msg.Context(), "api.OnChange", attrs)

			change := &structs.Change{}
			err := change.UnmarshalJSON(msg.Data())
			if err != nil {
				l.Error().Err(err).Msg("Failed to unmarshal change")
				pan.Err(err)
				pan.End()
				continue
			}

			ackId := ""
			if in.Queue != "" {
				ackId = s.qu.DeferAck(msg)
			}

			err = stream.Send(&structs.OnChangeResponse{Data: change, Ack: ackId, Error: toError(err)})
			if err != nil {
				l.Error().Err(err).Msg("Failed to send change")
			}

			pan.Err(err)
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
			if err != nil {
				// tell the caller something went wrong, this doesn't stop the stream
				stream.Send(&structs.AckResponse{Error: toError(err)})
			}
		}
	}
}

// genericAsyncRequestResponse is a helper function to send a request to the queue and wait for a response.
// We do this for all write requests sent to the API Server.
//
// This allows us to retry, reject, ack or nack requests as internally as required.
// Ie. we can ensure that we write to the database and emit to the queue, if either fails we can reject and
// retry the request (the db write will be idempotent if the object Etag has not changed in the meantime).
//
// This causes a worker routine running on an API Server (even us) to call asyncAPIRequest defined in
// async.go -- from here the relevant functions defined in async.go are called.
func (s *Server) genericAsyncRequestResponse(c context.Context, in marshalable, out marshalableReply) {
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

	s.workers = make([]backgroundWorker, s.cfg.WorkersAPI+1)

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
	for i := 1; i < s.cfg.WorkersAPI; i++ {
		log.Debug().Int("Worker", i).Msg("Starting worker")

		qsub, err := s.qu.DequeueApiReq()
		if err != nil {
			defer s.stopWorkers()
			log.Error().Err(err).Msg("Failed to subscribe to queue")
			return err
		}

		wrk := newApiWorker(fmt.Sprintf("api-worker-%d", i), s, qsub)
		s.workers[i] = wrk

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
}
