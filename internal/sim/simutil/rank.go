package simutil

import (
	"github.com/voidshard/faction/pkg/structs"
)

// ClosestRank returns the closest rank to the desired rank that is available.
// We iterate over all slots (eventually) defaulting upwards.
// If nothing is available, we return Associate.
func ClosestRank(d *structs.DemographicRankSpread, desired structs.FactionRank) (structs.FactionRank, bool) {
	if desired == structs.FactionRank_Associate { // there's always a slot for Associate
		return desired, true
	}
	if d.Count(desired) > 0 { // there is a slot at the given rank -> yay!
		return desired, true
	}

	min := int(structs.FactionRank_Associate)
	max := int(structs.FactionRank_Ruler)

	for i := 1; i <= max; i++ {
		j := int(desired) + i
		k := int(desired) - i

		if j >= min && j <= max && d.Count(structs.FactionRank(j)) > 0 {
			return structs.FactionRank(j), true
		}
		if k >= min && k <= max && d.Count(structs.FactionRank(k)) > 0 {
			return structs.FactionRank(k), true
		}
	}

	return structs.FactionRank_Associate, false
}

func RankFromAffiliation(a int) structs.FactionRank {
	space := structs.MaxTuple / int(structs.FactionRank_Ruler)
	for i := 1; i < int(structs.FactionRank_Ruler); i++ {
		if a < i*space {
			return structs.FactionRank(i - 1)
		}
	}
	return structs.FactionRank_Ruler
}

// AvailablePositions returns open positions given the current positions taken, leadership type & faction structure.
func AvailablePositions(d *structs.DemographicRankSpread, ltype structs.FactionLeadership, structure structs.FactionStructure) *structs.DemographicRankSpread {
	minZero := func(i int64) int64 {
		if i < 0 {
			return 0
		}
		return i
	}

	rulers := ltype.Rulers()
	rulerSlots := minZero(rulers - d.Ruler)

	var ds *structs.DemographicRankSpread

	switch structure {
	case structs.FactionStructure_Pyramid:
		// ie. the number of positions available for a rank is always
		// {people in lower rank / 3} - {people in this rank}
		//
		// That is, if there are currently 15 Journeyman, then there are 5 Adept positions.
		// Thus the number of *open* positions is 5 minus the current number of Adepts.
		//
		// So if 15 journeyman and 2 adepts, the number of open adept positions is 15 / 3 - 2 = 3
		ds = &structs.DemographicRankSpread{
			Ruler:       rulerSlots,
			Elder:       minZero(d.GrandMaster/3 - d.Elder),
			GrandMaster: minZero(d.Expert/3 - d.GrandMaster),
			Expert:      minZero(d.Adept/3 - d.Expert),
			Adept:       minZero(d.Journeyman/3 - d.Adept),
			Journeyman:  minZero(d.Novice/3 - d.Journeyman),
			Novice:      minZero(d.Apprentice/3 - d.Novice),
			Apprentice:  minZero(d.Associate/3 - d.Apprentice),
			Associate:   1, // there is always a spot for a new recruit
		}
	case structs.FactionStructure_Cell:
		ds = &structs.DemographicRankSpread{
			Ruler:       rulerSlots,
			Elder:       minZero(rulers - d.Elder),           // each ruler slot has a single second in command
			GrandMaster: minZero(rulers*2 - d.GrandMaster),   // each elder has two squad commanders
			Expert:      minZero(rulers*2*2 - d.Expert),      // each squad has 2 experts
			Adept:       minZero(rulers*2*10 - d.Adept),      // each squad has 10 adepts
			Journeyman:  minZero(rulers*2*25 - d.Journeyman), // each squad has 25 journeyman
			Novice:      1,                                   // open slots for junior ranks
			Apprentice:  1,
			Associate:   1,
		}
	case structs.FactionStructure_Loose:
		// no one cares, there's always an open slot if you've got the skills
		ds = &structs.DemographicRankSpread{
			Ruler:       rulerSlots,
			Elder:       1,
			GrandMaster: 1,
			Master:      1,
			Expert:      1,
			Adept:       1,
			Journeyman:  1,
			Novice:      1,
			Apprentice:  1,
			Associate:   1,
		}
	}

	return ds
}
