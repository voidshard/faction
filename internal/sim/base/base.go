package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/economy"
	fantasy "github.com/voidshard/faction/pkg/premade/fantasy"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/technology"
)

// Base is our base implementation of Simulation
type Base struct {
	cfg  *config.Simulation
	eco  economy.Economy
	tech technology.Technology

	dbconn *db.FactionDB
}

// New Simulation, the main doo-da
func New(cfg *config.Simulation) (*Base, error) {
	dbconn, err := db.New(cfg.Database)
	return &Base{
		cfg:    cfg,
		dbconn: dbconn,
		// default tech / eco
		eco:  fantasy.NewEconomy(),
		tech: fantasy.NewTechnology(),
	}, err
}

func (s *Base) SetTechnology(tech technology.Technology) error {
	s.tech = tech
	return nil
}

func (s *Base) SetEconomy(eco economy.Economy) error {
	s.eco = eco
	return nil
}

func (s *Base) Factions(ids ...string) ([]*structs.Faction, error) {
	if len(ids) == 0 {
		return nil, nil // we don't want to return everything
	}
	f := db.Q(db.F(db.ID, db.In, ids))
	out, _, err := s.dbconn.Factions("", f)
	return out, err
}

func (s *Base) FactionSummaries(in ...string) ([]*structs.FactionSummary, error) {
	return s.dbconn.FactionSummary(db.FactionSummaryRelations, in...)
}

func (s *Base) SetAreas(in ...*structs.Area) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetAreas(in...)
	})
}

func (s *Base) SetPlots(in ...*structs.Plot) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetPlots(in...)
	})
}

func (s *Base) SetGovernments(in ...*structs.Government) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetGovernments(in...)
	})
}

func (s *Base) SetFactions(in ...*structs.Faction) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetFactions(in...)
	})
}

func (s *Base) SetAreaGovernment(governmentID string, areas ...string) error {
	return s.dbconn.SetAreaGovernment(governmentID, areas)
}

func (s *Base) SetRoutes(in ...*structs.Route) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetRoutes(in...)
	})
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
	return
}

func (s *Base) Demographics(areas ...string) (*structs.Demographics, error) {
	return s.dbconn.Demographics(&db.DemographicQuery{Areas: areas})
}
