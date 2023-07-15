package base

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/structs"
)

// yieldRand is a struct for storing the various amount(s) of commodities across some area(s).
// Useful for metadata about what land might be useful for.
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
