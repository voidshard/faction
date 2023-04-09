package sim

import (
	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/premade"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	minEthosDistance = 0.0
	maxEthosDistance = structs.EthosDistance((&structs.Ethos{}).Sub(structs.MaxEthos), (&structs.Ethos{}).Add(structs.MaxEthos))
)

// simulationImpl implements Simulation
type simulationImpl struct {
	cfg  *config.Simulation
	eco  Economy
	tech Technology

	dbconn *db.FactionDB
}

// New Simulation, the main doo-da
func New(cfg *config.Simulation, opts ...simOption) (Simulation, error) {
	// apply default settings
	if cfg == nil {
		cfg = &config.Simulation{}
	}
	if cfg.Database == nil {
		cfg.Database = config.DefaultDatabase()
	}

	// setup sim using settings
	dbconn, err := db.New(cfg.Database)
	if err != nil {
		return nil, err
	}

	me := &simulationImpl{
		cfg:    cfg,
		dbconn: dbconn,
		// default economy / tech
		eco:  premade.NewFantasyEconomy(),
		tech: premade.NewFantasyTechnology(),
	}

	// apply given options
	for _, opt := range opts {
		err = opt(me)
		if err != nil {
			return nil, err
		}
	}

	return me, nil
}

func (s *simulationImpl) Factions(ids ...string) ([]*structs.Faction, error) {
	if len(ids) == 0 {
		return nil, nil // we don't want to return everything
	}
	f := make([]*db.FactionFilter, len(ids))
	for i, id := range ids {
		f[i].ID = id
	}
	out, _, err := s.dbconn.Factions("", f...)
	return out, err
}

func (s *simulationImpl) FactionSummaries(in ...string) ([]*structs.FactionSummary, error) {
	return s.dbconn.FactionSummary(in...)
}

/*
// InspireFactionAffiliation adds affiliaton to the given factions in regions they have influence.
//
// Nb.
//   - people can gain affiliation with all given factions, or none.
//   - people can only gain affiliation if they're in areas where the factions have influence
//     (ie. they control a Plot(s))
//   - if you want finer grain control, consider calling this on a per-faction basis,
//     we're strictly doing inserts so it should be fine to call simultaneously.
func (s *simulationImpl) InspireFactionAffiliation(factions []*structs.Faction, affdist *config.Distribution, probability, minEthosDist, maxEthosDist float64) error {
	minEthosDist = math.Max(minEthosDistance, minEthosDist)
	maxEthosDist = math.Min(maxEthosDistance, maxEthosDist)

	// 1. work out where the factions have influence in the world
	fids := make([]string, len(factions))
	for i, f := range factions {
		fids[i] = f.ID
	}
	zones, err := s.dbconn.FactionAreas(fids...)
	if err != nil {
		return err
	}
	fmt.Println(zones)

	// 2. find people in the given areas, whose ethos is close to the faction ethos
	pfilters := []*db.PersonFilter{}
	for _, faction := range factions {
		areas, ok := zones[faction.ID]
		if !ok {
			// faction controls no plots or landrights
			continue
		}

		// a slightly wider range than actually desired, but don't want to do pythag
		// using the database .. although .. technically we probably *could*
		minEth := faction.Ethos.Sub(int(maxEthosDist / 2))
		maxEth := faction.Ethos.Add(int(maxEthosDist / 2))

		for areaID := range areas {
			pfilters = append(pfilters, &db.PersonFilter{
				EthosFilter: db.EthosFilter{
					MinEthos: minEth,
					MaxEthos: maxEth,
				},
				AreaID: areaID,
			})
			fmt.Println("area", areaID, "min", minEth, "max", maxEth)
		}
	}

	// 3. now run through everyone and decide which faction(s) they join
	var (
		affilDice = stats.NewRand(affdist.Min, affdist.Max, affdist.Mean, affdist.Deviation)
		rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
		people    []*structs.Person
		token     string
		pumpErr   error
	)

	for {
		people, token, err = s.dbconn.People(token, pfilters...)
		if err != nil {
			return err
		}

		for _, person := range people {
			for _, faction := range factions {
				dist := structs.EthosDistance(&person.Ethos, &faction.Ethos)
				if dist < minEthosDist || dist > maxEthosDist {
					continue // ethos is too far away
				}

				if rng.Float64() > probability {
					continue // person doesn't wish to join
				}

				// roll the dice & pump in the new tuple
				pump.SetTuples(
					db.RelationPersonFactionAffiliation,
					&structs.Tuple{
						Subject: person.ID,
						Object:  faction.ID,
						Value:   affilDice.Int(),
					},
				)

				if faction.ReligionID == "" {
					continue
				}

				// if the faction is a religion, add faith
				value := affilDice.Int()
				if faction.IsReligion {
					// the faction doesn't simply *have* a religion, it *is* a religion
					value += 20
				}
				pump.SetTuples(
					db.RelationPersonFactionAffiliation,
					&structs.Tuple{
						Subject: person.ID,
						Object:  faction.ReligionID,
						Value:   value,
					},
				)
			}
		}

		if token == "" || pumpErr != nil {
			break
		}
	}

	return pumpErr
}
*/

func (s *simulationImpl) SetAreas(in ...*structs.Area) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetAreas(in...)
	})
}

func (s *simulationImpl) SetPlots(in ...*structs.Plot) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetPlots(in...)
	})
}

func (s *simulationImpl) SetGovernments(in ...*structs.Government) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetGovernments(in...)
	})
}

func (s *simulationImpl) SetFactions(in ...*structs.Faction) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetFactions(in...)
	})
}

func (s *simulationImpl) SetAreaGovernment(governmentID string, areas ...string) error {
	return s.dbconn.SetAreaGovernment(governmentID, areas)
}

func (s *simulationImpl) SetRoutes(in ...*structs.Route) error {
	return s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		return tx.SetRoutes(in...)
	})
}

func (s *simulationImpl) Tick() (tick int, err error) {
	s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		tick, err = tx.Tick()
		if err != nil {
			return err
		}
		tick += 1
		return tx.SetTick(tick)
	})
	return
}

func (s *simulationImpl) Demographics(areas ...string) (*structs.Demographics, error) {
	return s.dbconn.Demographics(&db.DemographicQuery{Areas: areas})
}
