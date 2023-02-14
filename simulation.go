package faction

import ()

type Simulation interface {
	// advance the simulation by one 'tick'
	Tick() error

	//
	AddArea(a *Area) error
	//Populate(id string, demo *Demographics) error

	Events() <-chan *Event
}

// New Simulation, the main doo-da
func New(settings *Settings) (Simulation, error) {
	// apply default settings
	return nil, nil
}
