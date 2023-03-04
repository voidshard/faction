package sim

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Simulation is our raison d'être and provides a single interface for any action a
// user might want to perform.
type Simulation interface {

	//
	AddGovernment(g *structs.Government) error

	// AddArea inserts an area with the given harvestable resources (commodity names).
	// Ie. the strings in `resources` should be items known to our Economy (see economy.go)
	AddArea(a *structs.Area, resources ...string) error

	// AddRoutes adds links between two areas
	AddRoutes([]*structs.Route) error

	// Populate adds people to area(s) based on a general 'Demographics' outline.
	//
	// - This is to seed initial populations; people are crafted from thin air.
	// - You can Populate the same area multiple times, this operation is strictly
	// additive.
	// - People are spread evenly over the given areas.
	// - At least one area ID is required.
	// - If you're aiming for unevenly distributed population centres (ie. large cities)
	// then probably you want to call this again for those (or dedicate a call for each).
	// - The people count here is the number of people we want *alive* at the end
	// we can create more people than this; as we sometimes create people dead (as part
	// of a family) in order to simulate some past family tragedy or something.
	//
	// TODO: consider a func to determine current demographics given an area id(s)
	Populate(people int, demo *structs.Demographics, areas ...string) error

	// Tick advances the simulation by one 'tick' and returns the current tick.
	// This kicks off a full simulation loop asyncrhonously.
	Tick() (int, error)

	//
	//	Events() <-chan *structs.Event
}
