package base

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	stats "github.com/voidshard/faction/internal/random/rng"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	maxEthosDistance = structs.EthosDistance(
		(&structs.Ethos{}).Add(structs.MinEthos),
		(&structs.Ethos{}).Add(structs.MaxEthos),
	)
)

// InspireFactionAffiliation adds affiliaton to the given factions in regions they have influence.
func (s *Base) InspireFactionAffiliation(cfg *config.Affiliation, factionID string) error {
	// This is going to be .. complex.
	//
	// We need to:
	// 1. figure out the faction context; the areas in which it operates and their governments
	// 2. iterate over people in those areas who
	//   a. have no profession (if the faction doesn't have preferred professions)
	//   -or-
	//   b. have a preferred profession that the faction is interested in
	// 3. for each person (if their ethos isn't too far from the faction), we bucket them into
	//    either "consider" or "check"
	//  a. "consider" are people who currently have no major affiliation (eg. no preferred faction) so we
	//     can forward them straight along
	//  b. "check" are people who have a preferred faction, but we need to check their rank in case
	//     we can poach them. People from "check" will be moved to "consider" if they're below the rank
	//     threshold or dropped
	// 4. for each person sent to "consider" we assign them a faction affiliation & rank based on faction
	//    ethos, religion, random dice etc, until the faction has enough people or we run out of people
	//
	// This entire process is made more fun as we do everything in batches of people
	// (as always) to make our DB calls efficient.

	ctx, err := simutil.NewFactionContext(s.dbconn, factionID)
	if err != nil {
		return err
	}

	tick, err := s.dbconn.Tick()
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

		for people := range consider {
			trust := []*structs.Tuple{}
			rels := []*structs.Tuple{}
			affils := []*structs.Tuple{}
			ranks := []*structs.Tuple{}
			faith := []*structs.Tuple{}
			events := []*structs.Event{}

			for _, p := range people {
				affil := affdice.Int()
				desiredRank := simutil.RankFromAffiliation(affil)
				rank := ctx.ClosestOpenRank(desiredRank)

				if done < members {
					prev := p.PreferredFactionID
					p.PreferredFactionID = factionID
					ranks = append(ranks, &structs.Tuple{Subject: p.ID, Object: factionID, Value: int(rank)})
					events = append(events,
						simutil.NewFactionChangeEvent(p, tick, prev),
						simutil.NewFactionPromotionEvent(p, tick, factionID),
					)
					affil += (structs.MaxEthos / 50 * int(rank)) // higher rank -> higher affiliation
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
					rel := simutil.ProfessionalRelationByTrust(trustlevel)

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
				err = tx.SetTuples(db.RelationPersonFactionRank, ranks...)
				if err != nil {
					return err
				}
				return tx.SetEvents(events...)
			})
		}
	}()

	go func() {
		// iterates 'check' chunks, verifies faction rank, sends to 'consider' channel
		defer wgcheckpeople.Done()

		for people := range check {
			tf := db.Q().DisableSort()
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
				db.F(db.AdulthoodTick, db.Less, tick),
				db.F(db.DeathTick, db.Equal, 0),
			)
		} else {
			professions := []string{}
			for p := range ctx.Summary.Professions {
				professions = append(professions, p)
			}
			pf.Or(
				db.F(db.AreaID, db.In, ctx.Areas),
				db.F(db.PreferredProfession, db.In, professions),
				db.F(db.AdulthoodTick, db.Less, tick),
				db.F(db.DeathTick, db.Equal, 0),
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

				if p.AdulthoodTick < tick || p.DeathTick > 0 {
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
