package faction

import ()

type Commodity interface {
	// Value at the given location
	Value(area string, ticks int64) int
}

// Harvest is a commodity that is harvested somehow.
// Eg. it could be mined, farmed, found .. whatever
type Harvest interface {
	Commodity

	// + Skill / Skill level?

	// Yield at the given location at the given offset in ticks.
	Yield(area string, ticks int64) int

	//
	Difficulty() int
}

// Craft is a commodity that is produced from some inputs
type Craft interface {
	Commodity

	// + Skill / Skill level?

	// Materials used to manufacture this commodity
	Materials() map[Commodity]int

	//
	Yield() int

	//
	Difficulty() int
}
