package search

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/martinlindhe/base36"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/pkg/structs"
)

type Opensearch struct {
	l   log.Logger
	cfg *OpensearchConfig

	api *opensearchapi.Client

	bulklock sync.Mutex
	bulk     map[string]*opensearchBulk
}

type OpensearchConfig struct {
	Address  string
	Username string
	Password string

	FlushInterval time.Duration
	WriteRoutines int
}

func NewOpensearch(cfg *OpensearchConfig) (*Opensearch, error) {
	if cfg == nil {
		cfg = &OpensearchConfig{}
	}
	if cfg.Address == "" {
		cfg.Address = "https://localhost:9200"
	}
	if cfg.Username == "" {
		cfg.Username = "admin"
	}
	if cfg.Password == "" {
		cfg.Password = "admin" // pretty sure this wont work, opensearch has password min requirements
	}
	if cfg.FlushInterval == 0 {
		cfg.FlushInterval = time.Second * 5
	}
	if cfg.WriteRoutines == 0 {
		cfg.WriteRoutines = 2
	}
	me := &Opensearch{
		cfg:      cfg,
		l:        log.Sublogger("opensearch", map[string]string{"username": cfg.Username, "address": cfg.Address}),
		bulklock: sync.Mutex{},
		bulk:     map[string]*opensearchBulk{},
	}
	go me.ping()
	me.connect()
	return me, nil
}

// world_index returns a name for a world, index tuple.
// ie. actors_world1, actors_world2 etc.
// This forcibly divides data from different worlds and simplifies data management.
func world_index(world, name string) string {
	return fmt.Sprintf("%s_%s", name, strings.ToLower(base36.EncodeBytes([]byte(world))))
}

func (s *Opensearch) IndexActors(ctx context.Context, world string, actors []*structs.Actor, flush bool) error {
	objects := make([]structs.Object, len(actors))
	for i, a := range actors {
		objects[i] = a
	}
	return s.index(ctx, world_index(world, "actors"), objects, flush)
}

func (s *Opensearch) DeleteActor(ctx context.Context, world string, id string) error {
	return s.delete(ctx, world_index(world, "actors"), id)
}

func (s *Opensearch) IndexFactions(ctx context.Context, world string, factions []*structs.Faction, flush bool) error {
	objects := make([]structs.Object, len(factions))
	for i, f := range factions {
		objects[i] = f
	}
	return s.index(ctx, world_index(world, "factions"), objects, flush)
}

func (s *Opensearch) DeleteFaction(ctx context.Context, world string, id string) error {
	return s.delete(ctx, world_index(world, "factions"), id)
}

func (s *Opensearch) ping() {
	for {
		time.Sleep(time.Second * 60)
		if s.api == nil {
			continue
		}
		_, err := s.api.Ping(context.Background(), &opensearchapi.PingReq{})
		if err != nil {
			s.l.Error().Err(err).Msg("failed to ping opensearch")
			s.connect() // since we failed, reconnect
		}
	}
}

func (s *Opensearch) connect() {
	var (
		api *opensearchapi.Client
		err error
	)
	i := 0
	for {
		i += 1
		if i > 1 {
			time.Sleep(time.Second * time.Duration(i) * time.Duration(i))
			s.l.Debug().Int("attempt", i).Msg("retrying connection to opensearch")
		}
		s.l.Info().Msg("connecting to opensearch")
		cfg := opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing only. Use certificate for validation.
			},
			Addresses: strings.Split(s.cfg.Address, ","),
			Username:  s.cfg.Username,
			Password:  s.cfg.Password,
		}

		api, err = opensearchapi.NewClient(opensearchapi.Config{Client: cfg})
		if err != nil {
			s.l.Error().Err(err).Msg("failed to create opensearch api client")
			continue
		}

		_, err = api.Ping(context.Background(), &opensearchapi.PingReq{})
		if err != nil {
			s.l.Error().Err(err).Msg("failed to ping opensearch")
			continue
		}

		s.api = api
		return
	}
}

func (s *Opensearch) delete(ctx context.Context, index string, id string) error {
	pan := log.NewSpan(ctx, "opensearch.delete", map[string]interface{}{
		"index":   index,
		"address": s.cfg.Address,
		"id":      id,
	})
	defer pan.End()

	_, err := s.api.Document.Delete(ctx, opensearchapi.DocumentDeleteReq{Index: index, DocumentID: id})
	if err != nil {
		s.l.Error().Err(err).Str("id", id).Str("index", index).Msg("failed to delete document")
		return pan.Err(err)
	}

	return nil
}

func (s *Opensearch) getBulk(index string) (*opensearchBulk, error) {
	s.bulklock.Lock()
	defer s.bulklock.Unlock()

	if b, ok := s.bulk[index]; ok {
		return b, nil
	}

	b, err := newOpensearchBulk(index, s.api, s.cfg)
	if err != nil {
		return nil, err
	}

	s.bulk[index] = b
	return b, nil
}

func (s *Opensearch) index(ctx context.Context, index string, objects []structs.Object, flush bool) error {
	s.l.Debug().Str("index", index).Int("count", len(objects)).Bool("flush", flush).Msg("indexing objects")
	defer s.l.Debug().Str("index", index).Int("count", len(objects)).Bool("flush", flush).Msg("indexed objects")

	pan := log.NewSpan(ctx, "opensearch.index", map[string]interface{}{"index": index, "address": s.cfg.Address, "flush": flush})
	defer pan.End()

	if objects == nil || len(objects) == 0 {
		return nil
	}

	bulk, err := s.getBulk(index)
	if err != nil {
		return err
	}

	if !flush {
		// if we're not waiting for flushing then just throw the objects over the fence
		for _, obj := range objects {
			bulk.Index(ctx, obj, nil)
		}
		return nil
	}

	// if we've been asked to flush, we need to wait for the flush to complete and check if
	// any of the docs we wrote returned a failure
	wg := sync.WaitGroup{}
	errchan := make(chan error)
	errsdone := make(chan bool)

	// roll up any errors we get from the bulk indexer
	var flusherr error
	go func() {
		for err := range errchan {
			if err == nil {
				continue
			}
			if flusherr == nil {
				flusherr = err
			} else {
				flusherr = fmt.Errorf("%s; %w", err, flusherr)
			}
		}
		errsdone <- true
	}()

	// pass the objects for indexing
	for _, obj := range objects {
		wg.Add(1)
		bulk.Index(ctx, obj, func(err error) {
			errchan <- err
			wg.Done()
		})
	}

	// wait for all the things to happen
	wg.Wait()
	close(errchan)
	<-errsdone

	return flusherr
}
