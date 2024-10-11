package api

import (
	"context"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/queue"
	"github.com/voidshard/faction/internal/search"
	"github.com/voidshard/faction/pkg/structs"
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

type Service struct {
	cfg *Config

	db db.Database
	qu queue.Queue
	sb search.Search

	log log.Logger

	workers []*worker
}

// NewService creates a new API service.
func NewService(cfg *Config, db db.Database, qu queue.Queue, sb search.Search) *Service {
	// Read functions trigger us to go straight to the database.
	// Write functions are internally queued, though results are returned synchronously.
	if cfg == nil {
		cfg = &Config{Routines: defaultRoutines}
	}
	if cfg.Routines <= 0 {
		cfg.Routines = defaultRoutines
	}
	return &Service{cfg: cfg, db: db, qu: qu, sb: sb, log: log.Sublogger("api-service")}
}

func (s *Service) OnChange(in *structs.OnChangeRequest) (queue.Subscription, error) {
	return s.qu.SubscribeChange(in.Data, in.Queue)
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
