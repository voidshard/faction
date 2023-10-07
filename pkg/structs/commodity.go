package structs

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
	Name       string // eg. wheat, iron, iron ingots, silk
	Profession string // eg. farmer, blacksmith (cannot be "" - meaning "no profession")

	// Required input Commodity (names) + their amounts (if any) to `Craft` 1 unit of
	// this commodity.
	//
	// Nb. this only makes sense for stuff that isn't harvested or mined directly.
	// Ie. harvesting wheat doesn't have a requirement, but producing flour requires
	// some units of wheat to hand.
	//
	// Commodities can be created from different recipies, ie.
	// we can collect meat from wild_game or domestic livestock.
	Requires []map[string]float64
}
