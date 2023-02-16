package faction

import ()

// Commodity is anything that can be bought or sold
type Commodity interface {
	// Value at the given location
	Value(area string, ticks int64) int
}

// profession is some job that creates some Commodity.
// Used for following Harvest & Craft.
type profession interface {
	// Skill that applies to this profession
	Skill() string

	// Difficulty rating (applied to given Harvest, Craft).
	// In *general* we anticipate that Crafts require higher and
	// rarer Skill(s) than Harvest(ables).
	Difficulty() int

	// Land required in order to perform this craft.
	// Metalworking is probably fairly compact.
	// Tanning hides will take a fair amount of room ..
	// Return should be in land units ^2 (ie. an area).
	LandRequired() int
}

// Harvest is a commodity that is harvested somehow.
// Eg. it could be mined, farmed, found .. whatever
type Harvest interface {
	Commodity
	profession

	// Yield at the given location at the given offset in ticks.
	// Nb. this should be a yield expected of one person in a `tick`
	Yield(area string, ticks int64) int
}

// Craft represents a process that
// - requires one or more input Commodities
// - yields one or more output Commodities
//
// Note that this (along with Harvest) effectively forms a chain, where Harvest
// goods are those that are raw or starting points in a production chain.
//
// Ie.
//
//	You harvest Eggs
//	You harvest Grain & 'craft' to yield Flour
//	You harvest Milk & 'craft' to make Butter
//	You use Egg, Milk, Butter, Flour & 'craft' to yield Custard.
//
// Nb. this is only an example, please don't make custard without sugar, we're not advocating
// barbarism here.
type Craft interface {
	profession

	// Material input(s)
	Materials() map[Commodity]int

	// Material output(s)
	// Nb. this should be a yield expected of one person in a `tick`
	Yield() map[Commodity]int
}
