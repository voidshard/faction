package base

import (
	"errors"
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/log"
	demographics "github.com/voidshard/faction/internal/random/demographics"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/economy"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/queue"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/technology"
)

// Base is our base implementation of Simulation
type Base struct {
	cfg  *config.Simulation
	eco  economy.Economy
	tech technology.Technology

	dbconn *db.FactionDB
	queue  queue.Queue

	dice *demographics.Dice
}

// New Simulation, the main doo-da
func New(cfg *config.Simulation) (*Base, error) {
	dbconn, err := db.New(cfg.Database)
	if err != nil {
		return nil, err
	}
	que, err := queue.New(cfg.Queue)
	if err != nil {
		return nil, err
	}

	me := &Base{
		cfg:    cfg,
		dbconn: dbconn,
		// default tech / eco
		eco:   fantasy.NewEconomy(),
		tech:  fantasy.NewTechnology(),
		queue: que,
		// dice for sim configs
		dice: demographics.New(cfg),
	}

	if cfg.Queue.Driver == config.QueueInMemory {
		// for local mode, we wont actually be calling "start processing events" because
		// we don't expect there to be some other client(s) listening & waiting to process.
		me.registerTasksWithQueue()
	}

	return me, nil
}

func (s *Base) StartProcessingEvents() error {
	err := s.registerTasksWithQueue()
	log.Debug().Err(err).Msg()("registering events with queue")
	if err != nil {
		return err
	}
	err = s.queue.Start()
	log.Debug().Err(err).Msg()("start processing events")
	return err
}

func (s *Base) StopProcessingEvents() error {
	err := s.queue.Stop()
	log.Debug().Err(err).Msg()("stop processing events")
	return err
}

func (s *Base) FireEvents() error {
	log.Debug().Msg()("fire events")
	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	q := db.Q(db.F(db.Tick, db.Equal, tick))
	var (
		events []*structs.Event
		token  string
	)

	count := 0
	for {
		events, token, err = s.dbconn.Events(token, q)
		if err != nil {
			log.Error().Err(err).Msg()("error reading events")
			return err
		}

		// TODO: Perhaps we could batch launch
		for _, e := range events {
			data, err := e.MarshalJson()
			if err != nil {
				log.Error().Err(err).Str("event", fmt.Sprintf("%v", e)).Msg()("error marshalling event")
				return err
			}
			count += 1

			_, err = s.queue.Enqueue(eventTask(e.Type), data)
			if errors.Is(err, queue.ErrNoHandler) {
				// no post processing is needed for this event type
				continue
			} else if err != nil {
				log.Error().Err(err).Str("event", fmt.Sprintf("%v", e)).Msg()("error enqueuing event")
				return err
			}
		}

		if token == "" {
			break
		}
	}
	return nil
}

func (s *Base) SetTechnology(tech technology.Technology) error {
	log.Debug().Msg()(fmt.Sprintf("set technology: %v", tech))
	s.tech = tech
	return nil
}

func (s *Base) SetEconomy(eco economy.Economy) error {
	log.Debug().Msg()(fmt.Sprintf("set economy: %v", eco))
	s.eco = eco
	return nil
}

func (s *Base) SetQueue(q queue.Queue) error {
	log.Debug().Msg()(fmt.Sprintf("set queue: %v", q))
	s.queue = q
	return nil
}

func (s *Base) Factions(ids ...string) ([]*structs.Faction, error) {
	if len(ids) == 0 {
		return nil, nil // we don't want to return everything
	}
	f := db.Q(db.F(db.ID, db.In, ids))
	out, _, err := s.dbconn.Factions("", f)
	log.Debug().Int("in", len(ids)).Int("out", len(out)).Err(err).Msg()(fmt.Sprintf("read factions by id"))
	return out, err
}

func (s *Base) FactionSummaries(in ...string) ([]*structs.FactionSummary, error) {
	sum, err := s.dbconn.FactionSummary(db.FactionSummaryRelations, in...)
	log.Debug().Int("in", len(in)).Int("out", len(sum)).Err(err).Msg()(fmt.Sprintf("read faction summaries by id"))
	return sum, err
}

func (s *Base) SetAreas(in ...*structs.Area) error {
	err := s.dbconn.SetAreas(in...)
	log.Debug().Int("in", len(in)).Err(err).Msg()(fmt.Sprintf("set areas"))
	return err
}

func (s *Base) SetPlots(in ...*structs.Plot) error {
	if s.eco != nil {
		// whenever we write a plot, set the valuation
		for _, p := range in {
			p.Value = simutil.PlotValuation(p, s.eco, 0)
		}
	}
	err := s.dbconn.SetPlots(in...)
	log.Debug().Int("in", len(in)).Err(err).Msg()(fmt.Sprintf("set plots"))
	return err
}

func (s *Base) SetGovernments(in ...*structs.Government) error {
	err := s.dbconn.SetGovernments(in...)
	log.Debug().Int("in", len(in)).Err(err).Msg()(fmt.Sprintf("set governments"))
	return err
}

func (s *Base) SetFactions(in ...*structs.Faction) error {
	err := s.dbconn.SetFactions(in...)
	log.Debug().Int("in", len(in)).Err(err).Msg()(fmt.Sprintf("set factions"))
	return err
}

func (s *Base) SetAreaGovernment(governmentID string, areas ...string) error {
	err := s.dbconn.SetAreaGovernment(governmentID, areas)
	log.Debug().Int("in", len(areas)).Err(err).Msg()(fmt.Sprintf("set area government"))
	return err
}

func (s *Base) Tick() (tick int, err error) {
	s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		tick, err = tx.Tick()
		if err != nil {
			return err
		}
		tick += 1
		return tx.SetTick(tick)
	})
	log.Debug().Int("tick", tick).Err(err).Msg()(fmt.Sprintf("set current tick"))
	return
}

func (s *Base) Demographics(areas ...string) (*structs.Demographics, error) {
	demo, err := s.dbconn.Demographics(&db.DemographicQuery{Areas: areas})
	log.Debug().Int("in", len(areas)).Err(err).Msg()(fmt.Sprintf("read demographics"))
	return demo, err
}
