package sim

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

func (s *simulationImpl) FaithDemographics(areas ...string) (map[string]*structs.DemographicStats, error) {
	return s.singleDemographic(db.RelationPersonReligionFaith, areas)
}

func (s *simulationImpl) ProfessionDemographics(areas ...string) (map[string]*structs.DemographicStats, error) {
	return s.singleDemographic(db.RelationPersonProfessionSkill, areas)
}

func (s *simulationImpl) AffiliationDemographics(areas ...string) (map[string]*structs.DemographicStats, error) {
	return s.singleDemographic(db.RelationPersonFactionAffiliation, areas)
}

// singleDemographic is a helper function that returns a single demographic, which is essentially a
// weaker version of our internal func.
//
// Internally we can collect more and different stats for internal use, but externally we keep things
// a bit neater, less dangerous and hide internal details.
func (s *simulationImpl) singleDemographic(r db.Relation, areas []string) (map[string]*structs.DemographicStats, error) {
	stats, err := s.dbconn.Demographics([]db.Relation{r}, &db.DemographicQuery{Areas: areas})
	if err != nil {
		return nil, err
	}
	res, _ := stats[r]
	return res, nil
}
