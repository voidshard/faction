package sim

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/economy"
	"github.com/voidshard/faction/pkg/queue"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/technology"
)

// Simulation is our raison d'être and provides a single interface for any action a
// user might want to perform.
type Simulation interface {
	// -- Config functions --

	// Options to set extention interfaces.
	SetTechnology(tech technology.Technology) error
	SetEconomy(eco economy.Economy) error
	SetQueue(q queue.Queue) error

	// -- Crud functions --

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

	// SetAreaGovernment sets the government ID for the given area(s)
	//
	// That is, these lands are marked as under the sway of the given government.
	//
	// This probably follows SpawnGovernment, but could also represent a change in leadership
	// if an area is conquered or similar.
	SetAreaGovernment(governmentID string, areas ...string) error

	// -- Info functions --

	// FactionSummaries returns a summary of the given faction(s).
	//
	// These can contain a lot of data, so it's probably best to limit the number of factions
	// you pass in here.
	FactionSummaries(factionIDs ...string) ([]*structs.FactionSummary, error)

	// Demographics for the given area(s).
	Demographics(areas ...string) (*structs.Demographics, error)

	// -- Simulation functions --

	// SpawnGovernment creates a new government.
	//
	// This does not create a faction nor grant the government any area(s).
	//
	// [Setup Function] Used to seed initial governments; governments are crafted from thin air,
	// Probably this should be followed by SetAreaGovernment to assign this government area(s) of
	// influence.
	SpawnGovernment(g *config.Government) (*structs.Government, error)

	// SpawnFaction creates a faction within the given area(s).
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
	SpawnFaction(f *config.Faction, areas ...string) (*structs.Faction, error)

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
	// - The number created is approximate
	SpawnPopulace(people int, race, culture string, areas ...string) error

	// InspireFactionAffiliation will assign faction affliation & rank(s) to people with
	// some probability.
	//
	// If the faction is a religion or has a religion, then faith is also added.
	//
	// This doesn't create people, though it can modify them, so it should be called after
	// SpawnPopulace.
	InspireFactionAffiliation(cfg *config.Affiliation, factionID string) error

	// Tick updates internal time by one tick, returns the current tick.
	Tick() (int, error)

	// PlanFactionJobs tries to figure out given the current context / climate around the
	// faction what Jobs they wish to enact.
	// More rarely factions may also cancel Jobs they no longer wish to enact if for example they
	// feel a Job will not go ahead or priorities change.
	//
	// We return any extra Jobs created and/or cancelled as a result of this (if any).
	//
	// Note that a faction may decide not to queue up extra work for many reasons, eg.
	// - they don't have enough money to support their desired Job(s)
	// - they currently have enough / too much work on
	// - their desired Job(s) are already in progress, and they don't believe they have
	// enough followers / resources to support more
	//
	// This function doesn't actually enact Jobs.
	PlanFactionJobs(factionID string) ([]*structs.Job, error)

	// AdjustPopulation accounts for natural deaths, births etc in an area
	AdjustPopulation(areaID string) error

	// AssignJobs tries to find people to fill any unfilled Jobs belonging to the given faction.
	AssignJobs(factionID string) error

	// -- Event / Async handling --

	// FireEvents will fire any events that have been queued up for the current tick.
	//
	// Rather than immediately fire events as they occur in a given tick, we store them in
	// the DB for firing in batches. We don't really need real time events, and this is way
	// more efficient to process.
	//
	// This should be called as the last function in a given tick and gives us the chance
	// to perform post-processing based on what has happened over the last tick / apply various
	// side effects and what not.
	FireEvents() error

	// StartProcessingEvents will start a background process to process events as they are fired
	// in batches.
	StartProcessingEvents() error

	// StopProcessingEvents will stop the background process.
	StopProcessingEvents() error
}
