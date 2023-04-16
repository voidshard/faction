package sim

import (
	"fmt"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	maxEthosDistance = structs.EthosDistance(
		(&structs.Ethos{}).Sub(structs.MaxEthos),
		(&structs.Ethos{}).Add(structs.MaxEthos),
	)
)

type factionContext struct {
	Summary     *structs.FactionSummary
	Areas       map[string]bool                // map areaID -> bool (in which the faction has influence)
	Governments map[string]*structs.Government // map areaID -> government (only the above areas)
}

// getFactionContext returns (pretty much) everything about a given faction and
// the regions / governments in which it has influence.
func (s *simulationImpl) getFactionContext(factionID string) (*factionContext, error) {
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

// InspireFactionAffiliation adds affiliaton to the given factions in regions they have influence.
func (s *simulationImpl) InspireFactionAffiliation(cfg *config.Affiliation, factionID string) error {
	ctx, err := s.getFactionContext(factionID)
	if err != nil {
		return err
	}

	pf := make([]*db.PersonFilter, len(ctx.Areas))
	i := 0
	for areaID := range ctx.Areas {
		pf[i] = &db.PersonFilter{AreaID: areaID}
		i++
	}

	tf := []*db.TupleFilter{}
	iter := s.dbconn.IterPeople(pf...)
	for {
		people, err := iter.Next()
		if err != nil {
			return err
		}

		for _, p := range people {
			dist := cfg.EthosDistance
			illegal, _ := ctx.Governments[p.AreaID].Outlawed.Factions[factionID]
			if illegal {
				dist *= cfg.OutlawedDistanceMod
			}
			if dist <= 0 {
				continue
			}

			ethDist := structs.EthosDistance(&p.Ethos, &ctx.Summary.Ethos) / maxEthosDistance
			if ethDist > dist {
				continue
			}

			tf = append(tf, &db.TupleFilter{Subject: p.ID})
		}

		if len(tf) < 500 && iter.HasNext() {
			continue
		}

	}

	//	max := cfg.EthosDistance

	return nil
}
