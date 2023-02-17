package faction

import ()

// Economy is a primitive way for our simulation to understand enough about the world
// economy to make hopefully not irrational decisions.
type Economy interface {
	// Commodity looks up a commodity by name & returns it
	Commodity(commodity string) *Commodity

	// CommodityValue returns the value (or forecasted value) of a commodity in
	// some area at some time offset in ticks (ie. '0' is 'now')
	CommodityValue(commodity, area string, ticks int64) float64

	// CommodityYield returns the yield (or forecasted yield) of a commodity in
	// some area at some time offset in ticks (ie. '0' is 'now').
	// Skill here is the skill named on the Commodity itself.
	CommodityYield(commodity, area string, skill int) float64

	// LandValue returns the value of 1 unit squared of land in the area
	LandValue(area string, ticks int64) float64
}

// Commodity is something of value in the economy.
//
// Raw materials are defined as those with no inputs.
// Everything else is a worked good (a "craft") that has one or more input(s) and
// returns one or more output(s).
//
// Every commodity has a name, a profession and an applicable skill.
//
// Professions are generally the focus of a Faction `ProfessionFocus`
// ie. a guild of farmers might have a focus on professions Farming, Weaving, Tanning
// and produce commodities wheat, textiles, leather.
//
// The simulator doesn't actually need to understand what these professions mean
// it's simply a way of linking commodities -> professions -> factions.
type Commodity struct {
	Name            string // eg. wheat, iron, iron ingots, silk
	Profession      string // eg. farmer, blacksmith
	Skill           string // eg. farming, smithing, weaving, painting
	LandRequirement int    // some amount of land required to perform this task (units squared)

	Requires  map[string]float64 // required input Commodity (names) + their amounts (if any)
	BaseYield map[string]float64 // output Commoditiy (names) + their amounts
}
