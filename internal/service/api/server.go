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
	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/pkg/structs"
)

type Server struct {
	structs.UnimplementedAPIServer

	svc *Service
	srv *grpc.Server

	log log.Logger
}

func NewServer(svc *Service) *Server {
	srv := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	me := &Server{svc: svc, srv: srv, log: log.Sublogger("api")}
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
	err = s.svc.start()
	if err != nil {
		return err
	}
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.log.Info().Msg("...Exiting")
	s.svc.stop()
	s.srv.Stop()
}

func (s *Server) Worlds(ctx context.Context, in *structs.GetWorldsRequest) (*structs.GetWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.Worlds", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()
	r := s.svc.Worlds(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) SetWorld(ctx context.Context, in *structs.SetWorldRequest) (*structs.SetWorldResponse, error) {
	pan := log.NewSpan(ctx, "api.SetWorld", map[string]interface{}{"id": in.Data.Id})
	defer pan.End()
	r := s.svc.SetWorld(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) DeleteWorld(ctx context.Context, in *structs.DeleteWorldRequest) (*structs.DeleteWorldResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteWorld", map[string]interface{}{"id": in.Id})
	defer pan.End()
	r := s.svc.DeleteWorld(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) ListWorlds(ctx context.Context, in *structs.ListWorldsRequest) (*structs.ListWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListWorlds", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	r := s.svc.ListWorlds(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) Factions(ctx context.Context, in *structs.GetFactionsRequest) (*structs.GetFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.Factions", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()
	r := s.svc.Factions(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) SetFactions(ctx context.Context, in *structs.SetFactionsRequest) (*structs.SetFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.SetFactions", map[string]interface{}{"id-count": len(in.Data)})
	defer pan.End()
	r := s.svc.SetFactions(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) DeleteFaction(ctx context.Context, in *structs.DeleteFactionRequest) (*structs.DeleteFactionResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteFaction", map[string]interface{}{"id": in.Id})
	defer pan.End()
	r := s.svc.DeleteFaction(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) ListFactions(ctx context.Context, in *structs.ListFactionsRequest) (*structs.ListFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListFactions", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	r := s.svc.ListFactions(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) Actors(ctx context.Context, in *structs.GetActorsRequest) (*structs.GetActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.Actors", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()
	r := s.svc.Actors(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) SetActors(ctx context.Context, in *structs.SetActorsRequest) (*structs.SetActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.SetActor", map[string]interface{}{"id-count": len(in.Data)})
	defer pan.End()
	r := s.svc.SetActors(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) DeleteActor(ctx context.Context, in *structs.DeleteActorRequest) (*structs.DeleteActorResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteActor", map[string]interface{}{"id": in.Id})
	defer pan.End()
	r := s.svc.DeleteActor(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

func (s *Server) ListActors(ctx context.Context, in *structs.ListActorsRequest) (*structs.ListActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListActors", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	r := s.svc.ListActors(ctx, in)
	if r.Error != nil {
		pan.Err(errors.New(r.Error.Message))
	}
	return r, nil
}

/*
func (s *Server) Jobs(ctx context.Context, in *structs.GetJobsRequest) (*structs.GetJobsResponse, error) {
	data, err := s.svc.Jobs(ctx, in.IDs...)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get jobs")
	}
	s.log.Debug().Int("count", len(data)).Msg("fetched jobs")
	return &structs.GetJobsResponse{Data: data, Error: toError(err)}, nil
}

func (s *Server) SetJob(ctx context.Context, in *structs.SetJobRequest) (*structs.SetJobResponse, error) {
	etag, err := s.svc.SetJob(ctx, in.Data)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to set job")
	}
	s.log.Debug().Str("old-etag", in.Data.Meta.Etag).Str("new-etag", etag).Msg("set job")
	return &structs.SetJobResponse{Etag: etag, Error: toError(err)}, nil
}
*/

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
	sub, err := s.svc.OnChange(in)
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

func toError(err error, code ...structs.ErrorCode) *structs.Error {
	if err == nil {
		return nil
	}
	if len(code) > 0 {
		return &structs.Error{Code: code[0], Message: err.Error()}
	}
	if errors.Is(err, db.ErrNotFound) {
		return &structs.Error{Code: structs.ErrorCode_NOT_FOUND, Message: err.Error()}
	}
	if errors.Is(err, db.ErrEtagMismatch) {
		return &structs.Error{Code: structs.ErrorCode_CONFLICT, Message: err.Error()}
	}
	if errors.Is(err, db.ErrInvalid) {
		return &structs.Error{Code: structs.ErrorCode_INVALID_OBJECT, Message: err.Error()}
	}
	return &structs.Error{Code: structs.ErrorCode_INTERNAL, Message: err.Error()}
}
