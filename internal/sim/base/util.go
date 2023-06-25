package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

type yieldRand struct {
	professionYield map[string]int             // profession -> yield
	professionLand  map[string][]*structs.Plot // profession -> land
}

func newYieldRand() *yieldRand {
	return &yieldRand{
		professionYield: map[string]int{},
		professionLand:  map[string][]*structs.Plot{},
	}
}

// areaYields returns the total yield of all land in the given areas, used to
// determine guild type factions.
func (s *Base) areaYields(in *db.Query, includeOwned bool) (*yieldRand, error) {
	var (
		land  []*structs.Plot
		token string
		err   error
	)

	yield := newYieldRand()
	for {
		land, token, err = s.dbconn.Plots(token, in)
		if err != nil {
			return nil, err
		}

		for _, lr := range land {
			if !includeOwned && lr.FactionID != "" {
				// someone already runs this
				continue
			}

			commodity := s.eco.Commodity(lr.Commodity)
			if commodity == nil {
				continue // ?
			}

			prof := commodity.Profession
			if prof == "" {
				continue // ?
			}

			total, _ := yield.professionYield[prof]
			yield.professionYield[prof] = lr.Yield + total

			lands, ok := yield.professionLand[prof]
			if !ok {
				lands = []*structs.Plot{}
			}
			yield.professionLand[prof] = append(lands, lr)
		}

		if token == "" {
			break
		}
	}

	return yield, nil
}

// areaGovernments returns
// 1. a map of area id to government id
// 2. a map of government id to government
func (s *Base) areaGovernments(in *db.Query) (map[string]string, map[string]*structs.Government, error) {
	areaToGovt := map[string]string{}
	govtToGovt := map[string]*structs.Government{}

	var (
		areas []*structs.Area
		token string
		err   error
	)

	for {
		areas, token, err = s.dbconn.Areas(token, in)
		if err != nil {
			return nil, nil, err
		}

		for _, area := range areas {
			if dbutils.IsValidID(area.GovernmentID) {
				areaToGovt[area.ID] = area.GovernmentID
				govtToGovt[area.GovernmentID] = nil
			} else {
				areaToGovt[area.ID] = ""
			}
		}

		if token == "" {
			break
		}
	}

	if len(govtToGovt) == 0 {
		return areaToGovt, govtToGovt, nil
	}

	gids := []string{}
	for gid := range govtToGovt {
		gids = append(gids, gid)
	}
	gf := db.Q(db.F(db.ID, db.In, gids))

	var governments []*structs.Government

	for {
		governments, token, err = s.dbconn.Governments(token, gf)
		if err != nil {
			return nil, nil, err
		}

		for _, govt := range governments {
			govtToGovt[govt.ID] = govt
		}

		if token == "" {
			break
		}
	}

	return areaToGovt, govtToGovt, nil
}
