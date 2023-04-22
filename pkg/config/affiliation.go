package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

type Affiliation struct {
	// Possible starting affiliation
	Affiliation Distribution

	// People who are in the faction are randomly given trust towards (or against) one another.
	// Negative values are permitted, and imply faction in-fighting / rivalries.
	Trust Distribution

	// Faith granted to people who are in the faction if the faction has a religion.
	// Ie. 'religion_id' is set
	Faith Distribution

	// The max ethos distance from a faction that a person can have to be considered for affiliation.
	//
	// This is expressed as a distance where `1` is the maximum possible (ie, `1` implies they're
	// diametrically opposed in all ethics) and `0` implies they're identical.
	EthosDistance float64

	// OutlawedWeight narrows the EthosDistance for areas in which a faction is outlawed.
	// Ie. 0.5 implies EthosDistance is halved if the person being considered is in an area
	// whose government has outlawed the faction.
	// This implies the person must be more fanatical in order to consider joining.
	OutlawedWeight float64

	// ReligionWeight is used to grant extra faith if the faction in question is an
	// organised religion (ie. 'IsReligion' is set).
	ReligionWeight float64

	// The number of members a faction should have. (Best effort).
	Members Distribution

	// Allows affiliation calculations to poach people who currently prefer / have ranks in
	// another faction(s). By default this is 0 (the lowest rank) so factions do not poach.
	// If set, people will be considered if their rank in their currently preferred faction
	// is below the rank given here (otherwise, they're skipped).
	PoachBelowRank structs.FactionRank
}
