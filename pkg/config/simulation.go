package config

// Simulation Configuration to instantiate a Simulation object.
//
// Our main concern here are things like;
// - base professions / skills we can assign people
// - base probabilities of various professions within a demographic
// - base Action probabilities
// - base probabilities of various forms of government
type Simulation struct {
	// Database configuration information
	Database *Database

	// Database is an interface to some .. database (whoa!).
	// The sim module only outlines what it needs this to do. Implementation(s) are
	// under internal/db

	// DefaultCraftCommodities are commodities (see economy.go) that can be produced
	// using ActionTypeCraft (see action.go) by a faction(s).
	// These are available to any faction that is into crafting in general.
	//
	// That is, anyone in any area can "Craft" to output these commodities assuming
	// they have a requisite resources and the will to do so.
	DefaultCraftCommodities []string

	// TODO: queue (fan out for faction job calculations)
	// TODO: graph? (future w/ path calculations for trade)
	// TODO: event sink (allow caller to collect decisions)
}
