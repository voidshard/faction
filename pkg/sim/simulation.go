package faction

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Simulation is our raison d'être and provides a single interface for any action a
// user might want to perform.
type Simulation interface {
	// Tick advances the simulation by one 'tick'
	Tick() error

	// AddArea makes a new area in the simulation. For our purposes it's enough to
	// know who lives there, what the area produces and what other area(s) it connects
	// to.
	AddArea(a *structs.Area) error

	// Populate adds people to an area based on a general outline.
	// Nb. you can Populate the same area multiple times, if you wanted to overlay
	// multiple cultures (for example).
	Populate(areaID string, demo *structs.Demographics) error

	//
	Events() <-chan *structs.Event
}

// New Simulation, the main doo-da
func New(cfg *config.Simulation) (Simulation, error) {
	// apply default settings
	return nil, nil
}
