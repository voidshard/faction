package faction

// Settings to instantiate a Simulation object.
// Our main concern here are things like;
// - base professions / skills we can assign people
// - base probabilities of various professions within a demographic
// - base Action probabilities
// - base probabilities of various forms of government
type Settings struct {

	// DefaultCraftCommodities are commodities (see economy.go) that can be produced
	// using ActionTypeCraft (see action.go) by a faction(s).
	// These are available to any faction that is into craft in general.
	//
	// That is, anyone in any area can "Craft" to output these commodities assuming
	// they have a requisite resources and the will to do so.
	//
	// Phrased another way, these might be considered default technologies that
	// everyone knows how to make -- to save you specifying them per government.
	DefaultCraftCommodities []string

	// database
	// queue
	// graph
	// event sink
}
