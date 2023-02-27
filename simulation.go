package faction

import ()

type Simulation interface {
	// Tick advances the simulation by one 'tick'
	Tick() error

	// AddArea makes a new area in the simulation. For our purposes it's enough to
	// know who lives there, what the area produces and what other area(s) it connects
	// to.
	AddArea(a *Area) error

	// Populate adds people to an area based on a general outline.
	Populate(areaID string, demo *Demographics) error

	//
	Events() <-chan *Event
}

// New Simulation, the main doo-da
func New(settings *Settings) (Simulation, error) {
	// apply default settings
	return nil, nil
}
