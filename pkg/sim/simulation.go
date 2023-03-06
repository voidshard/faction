package sim

import (
	"fmt"

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

func (s *simulationImpl) AddGovernment(in *structs.Government) error {
	if in.ID == "" {
		in.ID = structs.NewID()
	}
	if !structs.IsValidID(in.ID) {
		return fmt.Errorf("invalid government id: %s", in.ID)
	}
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetGovernments(in)
	})
}

func (s *simulationImpl) AddArea(in *structs.Area, resources ...string) error {
	if in.ID == "" {
		in.ID = structs.NewID()
	}
	if !structs.IsValidID(in.ID) {
		return fmt.Errorf("invalid area id: %s", in.ID)
	}
	if in.GoverningFactionID != "" && !structs.IsValidID(in.GoverningFactionID) {
		return fmt.Errorf("invalid government id: %s", in.GoverningFactionID)
	}
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetAreas(in)
	})
}

func (s *simulationImpl) AddRoutes(in []*structs.Route) error {
	for _, r := range in {
		if !structs.IsValidID(r.SourceAreaID) {
			return fmt.Errorf("invalid source area id: %s", r.SourceAreaID)
		}
		if !structs.IsValidID(r.TargetAreaID) {
			return fmt.Errorf("invalid target area id: %s", r.TargetAreaID)
		}
		if r.SourceAreaID == r.TargetAreaID {
			return fmt.Errorf("source and target area ids are the same: %s", r.SourceAreaID)
		}
		if r.TravelTime < 0 { // we do not allow time travel, but instantaneous teleportation is fine
			return fmt.Errorf("invalid travel time, expected >= 0: %d", r.TravelTime)
		}
	}
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
