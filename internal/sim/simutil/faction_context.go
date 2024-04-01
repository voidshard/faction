package simutil

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/random/rng"
	"github.com/voidshard/faction/pkg/structs"
)

var rnggen = rand.New(rand.NewSource(time.Now().UnixNano()))

// FactionContext stores full faction information plus data about all areas
// and governments that rule those areas
type FactionContext struct {
	Summary         *structs.FactionSummary
	Areas           map[string]*structs.Area       // map areaID -> nil (in which the faction has influence)
	Governments     map[string]*structs.Government // map areaID -> government (only the above areas)
	LocalGovernment *structs.Government            // the government of the area the faction HQ is in
	Land            *structs.LandSummary

	cachedTargetAreas map[string][]string
	areaIDs           []string
	openRanks         *structs.DemographicRankSpread
	rels              *TrustRelations
	dbconn            *db.FactionDB

	researchTopics []string
	research       rng.Normalised
}

// RelationWeight returns a number representing how much this faction prefers peace / caution over ambition / war.
// Applied to Trust to generally make more aggressive / unfriendly factions more likely to throw the first punch,
// and more peaceful factions more likely to try to avoid conflict.
func (f *FactionContext) RelationWeight() int64 {
	return (f.Summary.Faction.Ethos.Caution / 50) + (f.Summary.Faction.Ethos.Altruism / 20) - (f.Summary.Faction.Ethos.Ambition / 20) + (f.Summary.Faction.Ethos.Pacifism / 50)
}

func (f *FactionContext) RandomArea(purpose string) string {
	if len(f.areaIDs) == 0 {
		return f.Summary.Faction.HomeAreaID
	}
	return f.areaIDs[rnggen.Intn(len(f.areaIDs))]
}

func (f *FactionContext) AreaIDs() []string {
	return f.areaIDs
}

func (f *FactionContext) RandomResearch() string {
	if len(f.researchTopics) == 0 {
		return ""
	} else if len(f.researchTopics) == 1 {
		return f.researchTopics[0]
	}
	return f.researchTopics[f.research.Int()]
}

func (f *FactionContext) Relations() *TrustRelations {
	if f.rels == nil {
		fr := NewTrustRelations()
		w := f.RelationWeight()
		for factionID, trust := range f.Summary.Trust {
			fr.Add(factionID, trust+w)
		}
		f.rels = fr
	}
	return f.rels
}

func (f *FactionContext) AllGovernments() []*structs.Government {
	seen := map[string]bool{}

	govs := []*structs.Government{}
	for _, gov := range f.Governments {
		_, ok := seen[gov.ID]
		if ok {
			continue
		}
		govs = append(govs, gov)
		seen[gov.ID] = true
	}

	return govs
}

func (f *FactionContext) ClosestOpenRank(desired structs.FactionRank) structs.FactionRank {
	if f.openRanks == nil {
		f.openRanks = AvailablePositions(f.Summary.Ranks, f.Summary.Faction.Leadership, f.Summary.Faction.Structure)
	}

	nearest, ok := ClosestRank(f.openRanks, desired)
	if !ok {
		f.openRanks = AvailablePositions(f.Summary.Ranks, f.Summary.Faction.Leadership, f.Summary.Faction.Structure)
		nearest, _ = ClosestRank(f.openRanks, desired)
	}

	f.Summary.Ranks.Add(nearest, 1)
	f.openRanks.Add(nearest, -1)

	return nearest
}

func NewFactionContext(dbconn *db.FactionDB, factionID string) (*FactionContext, error) {
	if !dbutils.IsValidID(factionID) {
		return nil, fmt.Errorf("invalid faction id %s", factionID)
	}

	// look up the faction summary, we only need a few fields so we'll limit it to those
	summaries, err := dbconn.FactionSummary([]db.Relation{
		db.RelationFactionProfessionWeight,
		db.RelationPersonFactionRank,
		db.RelationFactionTopicResearch,
		db.RelationFactionTopicResearchWeight,
		db.RelationFactionFactionTrust,
	}, factionID)
	if err != nil {
		return nil, err
	} else if len(summaries) == 0 {
		return nil, fmt.Errorf("faction %s not found", factionID)
	}

	// lookup where a faction has influence
	fareas, err := dbconn.FactionAreas(false, factionID)
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
	areaGovs, err := dbconn.AreaGovernments(areaIDs...)
	if err != nil {
		return nil, err
	}

	gov, _ := areaGovs[summaries[0].Faction.HomeAreaID] // can be nil

	// research
	researchTopics := []string{}
	prob := []float64{}
	for topic, weight := range summaries[0].Research {
		researchTopics = append(researchTopics, topic)
		prob = append(prob, float64(weight))
	}

	// land
	land, err := dbconn.LandSummary(nil, []string{factionID})
	if err != nil {
		return nil, err
	}

	return &FactionContext{
		Summary:           summaries[0],
		Areas:             areas,
		Land:              land,
		Governments:       areaGovs,
		LocalGovernment:   gov,
		cachedTargetAreas: map[string][]string{},
		areaIDs:           areaIDs,
		dbconn:            dbconn,
		researchTopics:    researchTopics,
		research:          rng.NewNormalised(prob),
	}, nil
}
