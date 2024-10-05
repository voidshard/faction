package api

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

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
	return s.svc.Worlds(ctx, in), nil
}

func (s *Server) SetWorld(ctx context.Context, in *structs.SetWorldRequest) (*structs.SetWorldResponse, error) {
	pan := log.NewSpan(ctx, "api.SetWorld", map[string]interface{}{"id": in.Data.Id})
	defer pan.End()
	return s.svc.SetWorld(ctx, in), nil
}

func (s *Server) DeleteWorld(ctx context.Context, in *structs.DeleteWorldRequest) (*structs.DeleteWorldResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteWorld", map[string]interface{}{"id": in.Id})
	defer pan.End()
	return s.svc.DeleteWorld(ctx, in), nil
}

func (s *Server) ListWorlds(ctx context.Context, in *structs.ListWorldsRequest) (*structs.ListWorldsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListWorlds", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	return s.svc.ListWorlds(ctx, in), nil
}

func (s *Server) Factions(ctx context.Context, in *structs.GetFactionsRequest) (*structs.GetFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.Factions", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()
	return s.svc.Factions(ctx, in), nil
}

func (s *Server) SetFactions(ctx context.Context, in *structs.SetFactionsRequest) (*structs.SetFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.SetFactions", map[string]interface{}{"id-count": len(in.Data)})
	defer pan.End()
	return s.svc.SetFactions(ctx, in), nil
}

func (s *Server) DeleteFaction(ctx context.Context, in *structs.DeleteFactionRequest) (*structs.DeleteFactionResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteFaction", map[string]interface{}{"id": in.Id})
	defer pan.End()
	return s.svc.DeleteFaction(ctx, in), nil
}

func (s *Server) ListFactions(ctx context.Context, in *structs.ListFactionsRequest) (*structs.ListFactionsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListFactions", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	return s.svc.ListFactions(ctx, in), nil
}

func (s *Server) Actors(ctx context.Context, in *structs.GetActorsRequest) (*structs.GetActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.Actors", map[string]interface{}{"id-count": len(in.Ids)})
	defer pan.End()
	return s.svc.Actors(ctx, in), nil
}

func (s *Server) SetActor(ctx context.Context, in *structs.SetActorsRequest) (*structs.SetActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.SetActor", map[string]interface{}{"id-count": len(in.Data)})
	defer pan.End()
	return s.svc.SetActors(ctx, in), nil
}

func (s *Server) DeleteActor(ctx context.Context, in *structs.DeleteActorRequest) (*structs.DeleteActorResponse, error) {
	pan := log.NewSpan(ctx, "api.DeleteActor", map[string]interface{}{"id": in.Id})
	defer pan.End()
	return s.svc.DeleteActor(ctx, in), nil
}

func (s *Server) ListActors(ctx context.Context, in *structs.ListActorsRequest) (*structs.ListActorsResponse, error) {
	pan := log.NewSpan(ctx, "api.ListActors", map[string]interface{}{"limit": in.Limit, "offset": in.Offset})
	defer pan.End()
	return s.svc.ListActors(ctx, in), nil
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

func (s *Server) OnChange(in *structs.OnChangeRequest, stream structs.API_OnChangeServer) error {
	sub, err := s.svc.OnChange(in.Data)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to start change listener")
		return err
	}

	defer func() {
		s.log.Debug().Str("world", in.Data.World).Str("key", fmt.Sprintf("%d", in.Data.Key)).Str("id", in.Data.ID).Msg("Closing change listener")
		sub.Close()
	}()

	for {
		select {
		case <-stream.Context().Done():
			err := stream.Context().Err()
			s.log.Debug().Err(err).Msg("OnChange client disconnected")
			return err
		case change, ok := <-sub.Channel():
			s.log.Debug().Str("world", change.World).Str("key", fmt.Sprintf("%d", change.Key)).Str("id", change.ID).Msg("Sending change")
			if !ok {
				return nil
			}
			err := stream.Send(&structs.OnChangeResponse{Data: change})
			if err != nil {
				return err
			}
		}
	}
}
*/

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
	return &structs.Error{Code: structs.ErrorCode_INTERNAL, Message: err.Error()}
}
