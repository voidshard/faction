package sim

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/premade"
	"github.com/voidshard/faction/pkg/structs"
)

// simulationImpl implements Simulation
type simulationImpl struct {
	cfg  *config.Simulation
	eco  Economy
	tech Technology

	dbconn *db.FactionDB
}

// New Simulation, the main doo-da
func New(cfg *config.Simulation, opts ...simOption) (Simulation, error) {
	// apply default settings
	if cfg == nil {
		cfg = &config.Simulation{}
	}
	if cfg.Database == nil {
		cfg.Database = config.DefaultDatabase()
	}

	// setup sim using settings
	dbconn, err := db.New(cfg.Database)
	if err != nil {
		return nil, err
	}

	me := &simulationImpl{
		cfg:    cfg,
		dbconn: dbconn,
		// default economy / tech
		eco:  premade.NewFantasyEconomy(),
		tech: premade.NewFantasyTechnology(),
	}

	// apply given options
	for _, opt := range opts {
		err = opt(me)
		if err != nil {
			return nil, err
		}
	}

	return me, nil
}

func (s *simulationImpl) Factions(ids ...string) ([]*structs.Faction, error) {
	if len(ids) == 0 {
		return nil, nil // we don't want to return everything
	}
	f := make([]*db.FactionFilter, len(ids))
	for i, id := range ids {
		f[i].ID = id
	}
	out, _, err := s.dbconn.Factions("", f...)
	return out, err
}

func (s *simulationImpl) FactionSummaries(in ...string) ([]*structs.FactionSummary, error) {
	return s.dbconn.FactionSummary(db.FactionSummaryRelations, in...)
}

func (s *simulationImpl) SetAreas(in ...*structs.Area) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetAreas(in...)
	})
}

func (s *simulationImpl) SetPlots(in ...*structs.Plot) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetPlots(in...)
	})
}

func (s *simulationImpl) SetGovernments(in ...*structs.Government) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetGovernments(in...)
	})
}

func (s *simulationImpl) SetFactions(in ...*structs.Faction) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetFactions(in...)
	})
}

func (s *simulationImpl) SetAreaGovernment(governmentID string, areas ...string) error {
	return s.dbconn.SetAreaGovernment(governmentID, areas)
}

func (s *simulationImpl) SetRoutes(in ...*structs.Route) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetRoutes(in...)
	})
}

func (s *simulationImpl) Tick() (tick int, err error) {
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

func (s *simulationImpl) Demographics(areas ...string) (*structs.Demographics, error) {
	return s.dbconn.Demographics(&db.DemographicQuery{Areas: areas})
}
