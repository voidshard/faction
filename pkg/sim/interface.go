package sim

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Simulation is our raison d'être and provides a single interface for any action a
// user might want to perform.
type Simulation interface {
	// SetGovernments upserts government(s).
	//
	// Nb. remember that a Government for us isn't a faction, it's the set of laws / rules.
	// A faction marked "IsGovernment" with the matching GovernmentID is the actual
	// government Faction.
	SetGovernments(g ...*structs.Government) error

	// SetFactions upserts faction(s).
	SetFactions(f ...*structs.Faction) error

	// Factions returns factions by their ID(s).
	Factions(ids ...string) ([]*structs.Faction, error)

	// SetAreas upserts area(s).
	SetAreas(a ...*structs.Area) error

	// SetPlots upserts plot(s).
	//
	// Plots might be a single building, land + buildings, an open tract of land
	// or a combination of the above.
	SetPlots(p ...*structs.Plot) error

	// SetRoutes upserts links between two areas
	SetRoutes(r ...*structs.Route) error

	// SpawnGovernment creates a new government.
	//
	// This does not create a faction nor grant the government any area(s).
	//
	// [Setup Function] Used to seed initial governments; governments are crafted from thin air,
	// Probably this should be followed by SetAreaGovernment to assign this government area(s) of
	// influence.
	SpawnGovernment(g *config.Government) (*structs.Government, error)

	// SetAreaGovernment sets the government ID for the given area(s)
	//
	// That is, these lands are marked as under the sway of the given government.
	//
	// This probably follows SpawnGovernment, but could also represent a change in leadership
	// if an area is conquered or similar.
	SetAreaGovernment(governmentID string, areas ...string) error

	// SpawnFactions creates `number` of new faction(s) within the given area(s).
	//
	// Note that factions settings will depend on landrights in these areas and government.
	// That is, if there are no mines in the area, then it makes no sense to spawn mining based
	// factions.
	//
	// Similarly the government laws in these area(s) will dictate that some factions are marked
	// `covert` (as their favoured activities are illegal) or not.
	//
	// Note that none of these factions will be marked as the Government itself, so you'll
	// want to explicitly pick one to set that, or spawn a new faction for that reason.
	SpawnFactions(number int, f *config.Faction, areas ...string) ([]*structs.Faction, error)

	// FactionSummaries returns a summary of the given faction(s).
	//
	// These can contain a lot of data, so it's probably best to limit the number of factions
	// you pass in here.
	FactionSummaries(factionIDs ...string) ([]*structs.FactionSummary, error)

	// SpawnPopulace adds people to area(s) based on a general 'Demographics' outline.
	//
	// [Setup Function] Used to seed initial populations; people are crafted from thin air,
	//
	// - You can Populate the same area multiple times, this operation is strictly
	// additive.
	// - People are spread evenly over the given areas.
	// - At least one area ID is required.
	// - If you're aiming for unevenly distributed population centres (ie. large cities)
	// then probably you want to call this again for those (or dedicate a call for each).
	// - The people count here is the number of people we want *alive* at the end
	// we can create more people than this; as we sometimes create people dead
	// in order to simulate some past family tragedy or something.
	SpawnPopulace(people int, demo *config.Demographics, areas ...string) error

	// Demographics for the given area(s).
	Demographics(areas ...string) (*structs.Demographics, error)

	// InspireFactionAffiliation will assign faction affliations to people with some probability.
	// This ties together a lot of information & is expensive to compute.
	//
	// If the faction is a religion or has a religion, then faith is also added.
	//
	// This doesn't create people, though it can modify them, so it should be called after
	// SpawnPopulace.
	InspireFactionAffiliation(cfg *config.Affiliation, factionID string) error

	// Tick advances the simulation by one 'tick' and returns the current tick.
	// This kicks off a full simulation loop asyncrhonously.
	// Results are viewed via the Events() channel as required.
	//
	// The simulation uses "ticks" to track time, what that means depends on your settings.
	// The simulation starts at tick 1.
	Tick() (int, error)

	//
	//	Events() <-chan *structs.Event
}
