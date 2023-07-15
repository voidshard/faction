package base

import (
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

// factionContext stores full faction information plus data about all areas
// and governments that rule those areas
type factionContext struct {
	Summary     *structs.FactionSummary
	Areas       map[string]bool                // map areaID -> bool (in which the faction has influence)
	Governments map[string]*structs.Government // map areaID -> government (only the above areas)

	openRanks *structs.DemographicRankSpread
}

func (f *factionContext) closestOpenRank(desired structs.FactionRank) structs.FactionRank {
	if f.openRanks == nil {
		f.openRanks = availablePositions(f.Summary.Ranks, f.Summary.Leadership, f.Summary.Structure)
	}

	nearest, ok := closestRank(f.openRanks, desired)
	if !ok {
		f.openRanks = availablePositions(f.Summary.Ranks, f.Summary.Leadership, f.Summary.Structure)
		nearest, _ = closestRank(f.openRanks, desired)
	}

	f.Summary.Ranks.Add(nearest, 1)
	f.openRanks.Add(nearest, -1)

	return nearest
}

// getFactionContext returns (pretty much) everything about a given faction and
// the regions / governments in which it has influence.
func (s *Base) getFactionContext(factionID string) (*factionContext, error) {
	if !dbutils.IsValidID(factionID) {
		return nil, fmt.Errorf("invalid faction id %s", factionID)
	}

	// look up the faction summary, we only need a few fields so we'll limit it to those
	summaries, err := s.dbconn.FactionSummary([]db.Relation{
		db.RelationFactionProfessionWeight,
		db.RelationPersonFactionRank,
	}, factionID)
	if err != nil {
		return nil, err
	} else if len(summaries) == 0 {
		return nil, fmt.Errorf("faction %s not found", factionID)
	}

	// lookup where a faction has influence
	fareas, err := s.dbconn.FactionAreas(factionID)
	if err != nil {
		return nil, err
	} else if len(fareas) == 0 {
		return nil, fmt.Errorf("faction %s has no areas of influence", factionID)
	}
	areas, ok := fareas[factionID] // map areaID -> bool
	if !ok {
		return nil, fmt.Errorf("faction %s has no areas of influence", factionID)
	}

	// lookup the government(s) of areas in which faction has influence
	areaIDs := make([]string, len(areas))
	count := 0
	for areaID := range areas {
		areaIDs[count] = areaID
		count++
	}
	areaGovs, err := s.dbconn.AreaGovernments(areaIDs...)
	if err != nil {
		return nil, err
	}

	return &factionContext{
		Summary:     summaries[0],
		Areas:       areas,
		Governments: areaGovs,
	}, nil
}
