package sim

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// simulationImpl implements Simulation
type simulationImpl struct {
	cfg *config.Simulation

	dbconn *db.FactionDB
}

// New Simulation, the main doo-da
func New(cfg *config.Simulation) (Simulation, error) {
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

	// we always start at tick 1
	err = dbconn.InTransaction(func(tx db.ReaderWriter) error {
		tick, err := tx.Tick()
		if err != nil {
			return err
		}
		if tick <= 0 {
			return tx.SetTick(1)
		}
		return nil
	})

	return &simulationImpl{
		cfg:    cfg,
		dbconn: dbconn,
	}, nil
}

func (s *simulationImpl) SetAreas(in ...*structs.Area) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetAreas(in...)
	})
}

func (s *simulationImpl) SetLandRights(in ...*structs.LandRight) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetLandRights(in...)
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

func (s *simulationImpl) SetGoverningFaction(factionID string, areas ...string) error {
	return s.dbconn.ChangeGoverningFaction(factionID, areas)
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
