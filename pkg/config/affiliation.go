package config

type Affiliation struct {
	// Possible starting affiliation
	Distribution Distribution

	// The max ethos distance from a faction that a person can have to be considered for affiliation.
	// This is expressed as a distance where `1` is the maximum possible (ie, `1` implies they're
	// diametrically opposed in all ethics).
	EthosDistance float64

	// OutlawedDistanceMod is a multiplier to EthosDistance if the faction is outlawed.
	// This can be used to restrict people from affiliating with outlawed factions unless they're
	// much more closely aligned with them. (Or you know, make it easier instead if you wanted ..).
	//
	// Used as a straight multiplier to EthosDistance.
	OutlawedDistanceMod float64

	// ProfessionWeight gives people with higher skill in faction profession(s) (if any) higher affiliation
	// & rank in a given faction.
	// In general this makes sense for guild or highly skilled factions, but does not well suit other looser
	// or more welcoming factions (ie. a church / temple).
	ProfessionWeight float64

	// The minimum number of members a faction must have to be considered for affiliation.
	// (Best effort).
	//
	// If needed, people will be selected & forcibly modified to fit the faction if no
	// suitable people can be found.
	//
	// We do not take people from factions in which they already have a rank
	// above Associate (the lowest), nor are people created from scratch.
	//
	// Useful when kicking off a faction.
	MinMembers int
}
