package sim

import (
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

	// SetLandRights upserts land rights.
	//
	// A land right is the legally awarded right to work some particular resource
	// in a given area.
	SetLandRights(l ...*structs.LandRight) error

	// SetPlots upserts plot(s).
	//
	// Plots here refers to plots of land / buildings. Some centre of operations
	// for the faction in a given area that isn't associated with a LandRight.
	// (ie. a factory for smelting iron ore into ingots, a shop for selling produce).
	//
	// The other distinction here; a landright represents a limited resource
	// which is ultimately controlled by a government (who can grant / revoke rights)
	// where as a plot is something a faction can own privately, which a government
	// doesn't have any particular rights over.
	//
	// A faction can go out & buy a plot assuming it has enough money, but one
	// must lobby / bribe / otherwise deal with the govt to get a land right.
	SetPlots(p ...*structs.Plot) error

	// SetRoutes upserts links between two areas
	SetRoutes(r ...*structs.Route) error

	// SetGoverningFaction sets the governing faction for the given area(s)
	// including any resources / landrights in those areas.
	SetGoverningFaction(factionID string, areas ...string) error

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
	// we can create more people than this; as we sometimes create people dead (as part
	// of a family) in order to simulate some past family tragedy or something.
	//
	// TODO: consider a func to determine current demographics given an area id(s)
	SpawnPopulace(people int, demo *structs.Demographics, areas ...string) error

	// InspireFactionAffiliation will assign faction affliations to people with some probability.
	//
	// [Setup Function] To grant people initial affiliation to factions.
	//
	// The affiliation level is dictated by (min, max, mean, deviation).
	// If someone gains affiliation is controlled by (probability, minEthosDistance, maxEthosDistance).
	//
	// Ethos distance is defined in pkg/structs/ethos.go - in short a distance of 0-100 is very close;
	// it means across all ethics the person is very similar to the factions, differing by at most half
	// (-100 -> 100) on one ethic. For very loose / non extreme factions probably higher max ethos
	// distances are ok. For super tight nit cult type factions, the max ethos distance should probably
	// be less.
	//
	// Factions only inspire affiliation in areas they have influence (they must control
	// a Plot [building] or LandRight [land + work rights] in the area).
	//
	// (To recap: people are more like to work for factions with whom they have high
	// affiliation. People gain affiliation working for a given faction. The people
	// with the highest Affiliation are considered the faction leader(s)).
	//
	// If the faction is a religion or has a religion, then faith is also added for the person.
	InspireFactionAffiliation(factions []*structs.Faction, min, max, mean, deviation, probability, minEthosDistance, maxEthosDistance float64) error

	// Tick advances the simulation by one 'tick' and returns the current tick.
	// This kicks off a full simulation loop asyncrhonously.
	Tick() (int, error)

	//
	//	Events() <-chan *structs.Event
}
