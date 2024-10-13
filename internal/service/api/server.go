package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/search"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/log"
)

const (
	defaultRoutines = 5
	defaultMaxAge   = 60 * time.Minute // implies something is horribly wrong
	defaultLimit    = 100
	defaultMaxLimit = 1000
)

type Config struct {
	// Routines is the number of worker routines to use for processing API requests.
	Routines int

	// MaxMessageAge is the maximum age of a message before it is considered stale.
	MaxMessageAge time.Duration

	// FlushSearch waits for writes to searchbase before returning API writes (slow).
	FlushSearch bool
}

type Server struct {
	structs.UnimplementedAPIServer

	cfg *Config

	db db.Database
	qu queue.Queue
	sb search.Search

	srv *grpc.Server
	log log.Logger

	workers []*worker
}

func NewServer(cfg *Config, db db.Database, qu queue.Queue, sb search.Search) *Server {
	srv := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	me := &Server{
		cfg: cfg,
		db:  db,
		qu:  qu,
		sb:  sb,
		srv: srv,
		log: log.Sublogger("api"),
	}

	//me.asyncFuncs[(&structs.SetWorldRequest{}).Kind()] = me.setWorld
	//me.asyncFuncs[(&structs.DeleteWorldRequest{}).Kind()] = me.deleteWorld
	structs.RegisterAPIServer(srv, me)
	return me
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

func (s *Server) Worlds(ctx context.Context, in *structs.GetWorldsRequest) (*structs.GetWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.Worlds", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()
	data, err := s.db.Worlds(ctx, in.Ids)
	if err != nil {
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
	obj := &structs.DeleteWorldResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		pan.Err(fmt.Errorf(obj.Error.Message))
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
	data := []*structs.Actor{}
	err := s.db.Get(ctx, in.World, "Actor", in.Ids, &data)
	if err != nil {
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
		pan.Err(fmt.Errorf(obj.Error.Message))
	}
	return obj, nil
}

func (s *Server) DeleteActor(ctx context.Context, in *structs.DeleteActorRequest) (*structs.DeleteActorResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteActor", map[string]interface{}{"id": in.Id})
	defer pan.End()
	obj := &structs.DeleteActorResponse{}
	s.genericAsyncRequestResponse(ctx, in, obj)
	if obj.Error != nil {
		pan.Err(fmt.Errorf(obj.Error.Message))
	}
	return obj, nil
}

func (s *Server) ListActors(ctx context.Context, in *structs.ListActorsRequest) (*structs.ListActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListActors", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	data := []*structs.Actor{}
	err := s.db.List(ctx, in.World, "Actor", in.Labels, int64(clamp(in.Limit, 1, defaultMaxLimit, defaultLimit)), int64(defaultValue(in.Offset, 0)), &data)
	if err != nil {
		pan.Err(err)
	}
	return &structs.ListActorsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) OnChange(in *structs.OnChangeRequest, stream structs.API_OnChangeServer) error {
	// setup a sub logger
	key, ok := structs.Metakey_name[int32(in.Data.Key)]
	if !ok {
		return fmt.Errorf("invalid key %d", in.Data.Key)
	}

	peerAddr := ""
	peer, ok := peer.FromContext(stream.Context())
	if ok {
		peerAddr = peer.Addr.String()
	}

	l := log.Sublogger("api.OnChange", map[string]string{
		"world": in.Data.World,
		"area":  in.Data.Area,
		"key":   key,
		"id":    in.Data.Id,
		"queue": in.Queue,
		"peer":  peerAddr,
	})
	traceAttrs := map[string]interface{}{
		"world": in.Data.World,
		"area":  in.Data.Area,
		"key":   key,
		"id":    in.Data.Id,
		"queue": in.Queue,
		"peer":  peerAddr,
	}

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
			pan := log.NewSpan(msg.Context(), "api.OnChange", traceAttrs)
			change, err := msg.Change()
			l.Debug().Err(err).Msg("Received change")
			err = stream.Send(&structs.OnChangeResponse{Data: change, Error: toError(err)})
			if err != nil {
				l.Error().Err(err).Msg("Failed to send change")
			}
			pan.Err(err)
			pan.End()
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

func (s *Server) startWorkers() error {
	if s.workers != nil {
		s.stopWorkers()
	}

	s.workers = make([]*worker, s.cfg.Routines)

	s.log.Debug().Int("Routines", s.cfg.Routines).Msg("Starting workers")
	for i := 0; i < s.cfg.Routines; i++ {
		log.Debug().Int("Worker", i).Msg("Starting worker")

		qsub, err := s.qu.SubscribeApiReq()
		if err != nil {
			defer s.stopWorkers()
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

func (s *Server) stopWorkers() {
	if s.workers == nil {
		return
	}
	for _, w := range s.workers {
		w.Kill()
	}
	s.workers = nil
}
