package base

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/stats"
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

// InspireFactionAffiliation adds affiliaton to the given factions in regions they have influence.
func (s *Base) InspireFactionAffiliation(cfg *config.Affiliation, factionID string) error {
	ctx, err := s.getFactionContext(factionID)
	if err != nil {
		return err
	}

	wgiterpeople := &sync.WaitGroup{}
	wgiterpeople.Add(1)
	wgcheckpeople := &sync.WaitGroup{}
	wgcheckpeople.Add(1)
	wgsetpeople := &sync.WaitGroup{}
	wgsetpeople.Add(1)

	// control channels
	consider := make(chan []*structs.Person)
	check := make(chan []*structs.Person)
	errors := make(chan error)
	finalerr := make(chan error)
	stopIterPeople := false

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	affdice := stats.NewRand(cfg.Affiliation.Min, cfg.Affiliation.Max, cfg.Affiliation.Mean, cfg.Affiliation.Deviation)
	trustdice := stats.NewRand(cfg.Trust.Min, cfg.Trust.Max, cfg.Trust.Mean, cfg.Trust.Deviation)
	faithdice := stats.NewRand(cfg.Faith.Min, cfg.Faith.Max, cfg.Faith.Mean, cfg.Faith.Deviation)

	go func() {
		// rolls up error(s)
		var ferr error
		for err := range errors {
			if err == nil {
				continue
			}
			stopIterPeople = true // we have an error -> tell iter people to stop
			if ferr == nil {
				ferr = err
			} else {
				ferr = fmt.Errorf("%w %v", ferr, err)
			}
		}
		finalerr <- ferr
	}()

	go func() {
		// iterates 'consider' chunks and determines affiliation & rank(s)
		defer wgsetpeople.Done()

		members := stats.NewRand(cfg.Members.Min, cfg.Members.Max, cfg.Members.Mean, cfg.Members.Deviation).Int()
		done := 0

		for candidates := range consider {
			trust := []*structs.Tuple{}
			rels := []*structs.Tuple{}
			affils := []*structs.Tuple{}
			ranks := []*structs.Tuple{}
			faith := []*structs.Tuple{}

			people := []*structs.Person{}
			for _, p := range candidates {
				affil := affdice.Int()
				desiredRank := rankFromAffiliation(affil)
				rank := ctx.closestOpenRank(desiredRank)

				people = append(people, p)

				if done < members {
					p.PreferredFactionID = factionID
					ranks = append(ranks, &structs.Tuple{Subject: p.ID, Object: factionID, Value: int(rank)})
					done += 1
				} else {
					// we've added full members, so we'll just set some affiliation on the rest
					// ie. this person isn't a full member, but they're sympathetic to the faction
					affil /= rng.Intn(4) + 2
				}
				affils = append(affils, &structs.Tuple{Subject: p.ID, Object: factionID, Value: affil})

				if ctx.Summary.ReligionID == "" {
					continue // we don't need to consider faith
				}

				fth := faithdice.Int()
				if ctx.Summary.IsReligion {
					fth += int(float64(fth) * cfg.ReligionWeight)
				}

				faith = append(faith, &structs.Tuple{Subject: p.ID, Object: ctx.Summary.ReligionID, Value: fth})
			}

			if len(people) > 1 {
				// insert random inter-person relationships
				for i := 0; i < len(people)/2; i++ {
					a := people[rng.Intn(len(people))]
					b := people[rng.Intn(len(people))]
					if a.ID == b.ID {
						continue
					}

					trustlevel := trustdice.Int()
					rel := professionalRelationByTrust(trustlevel)

					rels = append(
						rels,
						&structs.Tuple{Subject: a.ID, Object: b.ID, Value: int(rel)},
						&structs.Tuple{Subject: b.ID, Object: a.ID, Value: int(rel)},
					)
					trust = append(
						trust,
						&structs.Tuple{Subject: a.ID, Object: b.ID, Value: trustlevel + rng.Intn(50) - 25},
						&structs.Tuple{Subject: b.ID, Object: a.ID, Value: trustlevel + rng.Intn(50) - 25},
					)
				}
			}

			errors <- s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
				err = tx.SetTuples(db.RelationPersonReligionFaith, faith...)
				if err != nil {
					return err
				}
				err = tx.SetTuples(db.RelationPersonPersonTrust, trust...)
				if err != nil {
					return err
				}
				err = tx.SetTuples(db.RelationPersonPersonRelationship, rels...)
				if err != nil {
					return err
				}
				err = tx.SetTuples(db.RelationPersonFactionAffiliation, affils...)
				if err != nil {
					return err
				}
				return tx.SetTuples(db.RelationPersonFactionRank, ranks...)
			})
		}
	}()

	go func() {
		// iterates 'check' chunks, verifies faction rank, sends to 'consider' channel
		defer wgcheckpeople.Done()

		for people := range check {
			tf := db.Q()
			for _, p := range people {
				tf.Or(
					db.F(db.Subject, db.Equal, p.ID),
					db.F(db.Object, db.Equal, p.PreferredFactionID),
				)
			}

			token := dbutils.NewToken()
			token.Limit = len(people)

			ranks, _, err := s.dbconn.Tuples(db.RelationPersonFactionRank, token.String(), tf)
			if err != nil {
				errors <- err
				continue
			}

			skip := map[string]bool{} // person id -> bool
			for _, r := range ranks {
				if r.Value >= int(cfg.PoachBelowRank) {
					skip[r.Subject] = true
				}
			}

			todo := []*structs.Person{} // run through people and throw out those who're too high
			for _, p := range people {
				invalid, _ := skip[p.ID]
				if invalid {
					continue
				}
				todo = append(todo, p)
			}

			if len(todo) > 0 {
				consider <- todo
			}
		}
	}()

	go func() {
		// iterates people and chunks into either "consider" or "check" channels
		defer wgiterpeople.Done()

		chunk := 250

		pf := db.Q()
		noProfessions := ctx.Summary.Professions == nil || len(ctx.Summary.Professions) == 0
		if noProfessions {
			pf.Or(
				db.F(db.AreaID, db.In, ctx.Areas),
				db.F(db.IsChild, db.Equal, false),
			)
		} else {
			professions := []string{}
			for p := range ctx.Summary.Professions {
				professions = append(professions, p)
			}
			pf.Or(
				db.F(db.AreaID, db.In, ctx.Areas),
				db.F(db.PreferredProfession, db.In, professions),
				db.F(db.IsChild, db.Equal, false),
			)
		}

		toconsider := []*structs.Person{}
		tocheck := []*structs.Person{}

		var (
			people []*structs.Person
			token  string
			err    error
		)
		for {
			people, token, err = s.dbconn.People(token, pf)
			if err != nil {
				errors <- err
				return
			}
			for _, p := range people {
				if stopIterPeople {
					return
				}

				ethDist := structs.EthosDistance(&p.Ethos, &ctx.Summary.Ethos) / maxEthosDistance
				govt, _ := ctx.Governments[p.AreaID]
				if govt != nil {
					illegal, _ := govt.Outlawed.Factions[factionID]
					if illegal {
						ethDist += ethDist * cfg.OutlawedWeight
					}
				}

				if ethDist > cfg.EthosDistance {
					continue
				}

				if p.IsChild || p.DeathTick > 0 {
					continue // invalid
				} else if p.PreferredFactionID == "" {
					toconsider = append(toconsider, p) // easy case
				} else if cfg.PoachBelowRank > 0 {
					tocheck = append(tocheck, p) // have to check this person's rank
				} else {
					continue // ie. PreferredFaction && PoachBelowRank <= 0 -> no poaching
				}

				if len(toconsider) >= chunk {
					consider <- toconsider
					toconsider = []*structs.Person{}
				} else if len(tocheck) >= chunk {
					check <- tocheck
					tocheck = []*structs.Person{}
				}
			}

			if token == "" {
				break
			}
		}

		if len(toconsider) > 0 {
			consider <- toconsider
		}
		if len(tocheck) > 0 {
			check <- tocheck
		}
	}()

	wgiterpeople.Wait() // stops when we get an error or finish all people
	close(check)        // stop checker routine and wait for it
	wgcheckpeople.Wait()
	close(consider) // stop set routine and wait for it
	wgsetpeople.Wait()
	close(errors) // stop error routine and wait for it
	return <-finalerr
}

// closestRank returns the closest rank to the desired rank that is available.
// We iterate over all slots (eventually) defaulting upwards.
// If nothing is available, we return Associate.
func closestRank(d *structs.DemographicRankSpread, desired structs.FactionRank) (structs.FactionRank, bool) {
	if desired == structs.FactionRankAssociate { // there's always a slot for Associate
		return desired, true
	}
	if d.Count(desired) > 0 { // there is a slot at the given rank -> yay!
		return desired, true
	}

	min := int(structs.FactionRankAssociate)
	max := int(structs.FactionRankRuler)

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

	return structs.FactionRankAssociate, false
}

func rankFromAffiliation(a int) structs.FactionRank {
	space := structs.MaxTuple / int(structs.FactionRankRuler)
	for i := 1; i < int(structs.FactionRankRuler); i++ {
		if a < i*space {
			return structs.FactionRank(i - 1)
		}
	}
	return structs.FactionRankRuler
}

func professionalRelationByTrust(a int) structs.PersonalRelation {
	if a < structs.MaxTuple/10 && a > structs.MinTuple/10 {
		// the neutral zone
		return structs.PersonalRelationColleague
	} else if a < 0 {
		if a < structs.MinTuple/2 {
			return structs.PersonalRelationHatedEnemy
		}
		return structs.PersonalRelationEnemy
	} else { // a > 0
		if a > structs.MaxTuple/2 {
			return structs.PersonalRelationCloseFriend
		}
		return structs.PersonalRelationFriend
	}
}

// availablePositions returns open positions given the current positions taken, leadership type & faction structure.
func availablePositions(d *structs.DemographicRankSpread, ltype structs.LeaderType, structure structs.LeaderStructure) *structs.DemographicRankSpread {
	minZero := func(i int) int {
		if i < 0 {
			return 0
		}
		return i
	}

	rulers := ltype.Rulers()
	rulerSlots := minZero(rulers - d.Ruler)

	var ds *structs.DemographicRankSpread

	switch structure {
	case structs.LeaderStructurePyramid:
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
	case structs.LeaderStructureCell:
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
	case structs.LeaderStructureLoose:
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
