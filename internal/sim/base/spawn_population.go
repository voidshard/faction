/*
random_populate.go - Randomly populate the world with people and families.
*/
package base

import (
	"fmt"

	"github.com/voidshard/faction/internal/log"
	"github.com/voidshard/faction/internal/random/rng"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	diedInChildbirth = "died in childbirth"
	diedOfOldAge     = "died of old age"

	// "random" is a value 0-randomValueMax that is used to pull rows at random
	// from the db with some probability.
	randomValueMax = 1000000.0
)

func (s *Base) spawnFamily(tick int, areaID, race, culture string) *simutil.MetaPeople {
	demo := s.dice.MustDemographic(race, culture)

	mp := simutil.NewMetaPeople()

	// nb. we're creating mum & dad in the past so that our children (if any) are born "now"
	mum := demo.RandomPerson(areaID)
	mum.SetBirthTick(tick - demo.RandomParentingAge())
	_, mumFaiths := simutil.AddSkillAndFaith(s.dice, mp, mum)
	mum.IsMale = false
	mp.Events = append(mp.Events, simutil.NewBirthEvent(mum, ""))

	dad := demo.RandomPerson(areaID)
	dad.SetBirthTick(tick - demo.RandomParentingAge())
	_, dadFaiths := simutil.AddSkillAndFaith(s.dice, mp, dad)
	dad.IsMale = true
	mp.Events = append(mp.Events, simutil.NewBirthEvent(dad, ""))

	mp.Adults = append(mp.Adults, mum, dad)

	eldest := dad.BirthTick
	youngest := mum.BirthTick
	if mum.BirthTick < eldest { // ie. "is born before"
		eldest = mum.BirthTick
		youngest = dad.BirthTick
	}

	// number of ticks the couple could have been having children in
	// (We multiply by this as when generating a family we assume they've been a family
	// for a while now .. since we're spawning them from nothing).
	parentingTicks := float64(tick - youngest - demo.MinParentingAge())
	if parentingTicks < 1.0 {
		parentingTicks = 1.0
	}

	family := &structs.Family{
		Ethos:               *structs.EthosAverage(&mum.Ethos, &dad.Ethos),
		ID:                  structs.NewID(mum.ID, dad.ID),
		Race:                mum.Race,
		Culture:             mum.Culture,
		AreaID:              areaID,
		IsChildBearing:      true,
		MaleID:              dad.ID,
		FemaleID:            mum.ID,
		MaxChildBearingTick: eldest + demo.MaxParentingAge(), // the last tick the couple can have children
		MarriageTick:        tick,
		Random:              int(s.dice.Float64() * structs.FamilyRandomMax),
	}
	mp.Families = append(mp.Families, family)
	mp.Events = append(mp.Events, simutil.NewMarriageEvent(family, tick-demo.RandomChildbearingTerm()))

	havingAffair := demo.RandomIsHavingAffair(parentingTicks)
	if havingAffair {
		affair := demo.RandomPerson(areaID)
		affair.SetBirthTick(tick - demo.RandomParentingAge())
		simutil.AddSkillAndFaith(s.dice, mp, affair)
		mp.Adults = append(mp.Adults, affair)
		mp.Events = append(mp.Events, simutil.NewBirthEvent(affair, ""))

		if affair.IsMale {
			mp.Relations = append(
				mp.Relations,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: int(structs.PersonalRelationLover)},
			)
			mp.Trust = append(
				mp.Trust,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: demo.RandomTrust()},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: demo.RandomTrust()},
			)
		} else {
			mp.Relations = append(
				mp.Relations,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: int(structs.PersonalRelationLover)},
			)
			mp.Trust = append(
				mp.Trust,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: demo.RandomTrust()},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: demo.RandomTrust()},
			)
		}
	}

	for i := 0; i < demo.RandomFamilySize(); i++ {
		if mum.DeathTick > 0 {
			break
		}

		child := demo.RandomPerson(areaID)
		child.SetBirthTick(tick)
		mp.Events = append(mp.Events, simutil.NewBirthEvent(child, family.ID))

		if demo.RandomDeathInfantMortality() {
			child.DeathMetaReason = diedInChildbirth
			child.DeathTick = tick
			child.DeathMetaKey = structs.MetaKeyPerson
			child.DeathMetaVal = mum.ID
			mp.Events = append(mp.Events, simutil.NewDeathEvent(child))
		}

		simutil.AddChildToFamily(demo, child, family)
		simutil.AddParentChildRelations(demo, mp, child, family)

		if i == 0 && demo.RandomAdultDeathInChildbirth() {
			mum.DeathTick = tick
			mum.DeathMetaReason = diedInChildbirth
			mum.DeathMetaKey = structs.MetaKeyPerson
			mum.DeathMetaVal = child.ID
			mp.Events = append(mp.Events, simutil.NewDeathEvent(mum))
			family.IsChildBearing = false
			family.WidowedTick = tick
		}

		childFaiths := []*structs.Tuple{}
		if s.dice.Float64() <= 0.5 && mum.DeathTick <= 0 { // child takes faith from either mum or dad
			for _, f := range mumFaiths {
				childFaiths = append(childFaiths, &structs.Tuple{Subject: child.ID, Object: f.Object, Value: f.Value / 2})
			}
		} else {
			for _, f := range dadFaiths {
				childFaiths = append(childFaiths, &structs.Tuple{Subject: child.ID, Object: f.Object, Value: f.Value / 2})
			}
		}
		mp.Faith = append(mp.Faith, childFaiths...)
		mp.Children = append(mp.Children, child)
	}

	if mum.DeathTick > 0 || demo.RandomIsDivorced(parentingTicks) {
		mp.Relations = append(
			mp.Relations,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationExHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationExWife)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: demo.RandomTrust() / 2},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: demo.RandomTrust() / 2},
		)
		family.IsChildBearing = false
		if mum.DeathTick > 0 {
			family.WidowedTick = tick
		} else {
			family.DivorceTick = tick
			mp.Events = append(mp.Events, simutil.NewDivorceEvent(family, tick))
		}
	} else {
		div := 1
		if havingAffair {
			div = 4
		}
		mp.Relations = append(
			mp.Relations,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationWife)},
		)
		mp.Trust = append(
			mp.Trust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: demo.RandomTrust() / div},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: demo.RandomTrust() / div},
		)
	}

	if len(mp.Children) > 1 {
		// add inter-sibling relations
		for i, child := range mp.Children {
			for j := i + 1; j < len(mp.Children); j++ {
				sibling := mp.Children[j]
				if i == j {
					continue
				}
				simutil.SiblingRelationship(demo, mp, child, sibling)
			}
		}
	}

	if len(mp.Children) >= demo.MaxFamilySize() {
		family.IsChildBearing = false
	}

	return mp
}

func (s *Base) spawnCouple(tick int, areaID, race, culture string, mp *simutil.MetaPeople) int {
	demo := s.dice.MustDemographic(race, culture)

	alive := 0

	person := demo.RandomPerson(areaID)
	person.SetBirthTick(tick - demo.RandomParentingAge())
	mp.Events = append(mp.Events, simutil.NewBirthEvent(person, ""))

	mp.Adults = append(mp.Adults, person)
	simutil.AddSkillAndFaith(s.dice, mp, person)

	cause, died := demo.RandomDeathAdultMortality()
	if died {
		person.DeathTick = tick
		person.DeathMetaReason = cause
		mp.Events = append(mp.Events, simutil.NewDeathEvent(person))
	} else {
		alive++
	}

	if person.DeathTick > 0 {
		return alive
	}

	lover := demo.RandomPerson(areaID)
	lover.SetBirthTick(tick - demo.RandomParentingAge())
	mp.Events = append(mp.Events, simutil.NewBirthEvent(lover, ""))
	lover.IsMale = !person.IsMale
	alive++

	mp.Adults = append(mp.Adults, lover)
	simutil.AddSkillAndFaith(s.dice, mp, lover)

	rel := structs.PersonalRelationLover
	if s.dice.Float64() <= 0.1 {
		rel = structs.PersonalRelationFiance
	}

	mp.Relations = append(
		mp.Relations,
		&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: int(rel)},
		&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: int(rel)},
	)
	mp.Trust = append(
		mp.Trust,
		&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: demo.RandomTrust()},
		&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: demo.RandomTrust()},
	)

	return alive
}

func (s *Base) SpawnPopulace(desiredTotal int, race, culture string, areas ...string) error {
	if len(areas) < 1 {
		return nil
	}

	if !s.dice.IsValidDemographic(race, culture) {
		return fmt.Errorf("invalid demographic not found: [race] %s, [culture] %s", race, culture)
	}

	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	log.Info().Str("race", race).Str("culture", culture).Int("desired", desiredTotal).Msg()("spawning populace")

	desiredArea := desiredTotal / len(areas)

	// we place people one area at a time because it feels more natural as this results in more
	// links between people within a given area than between areas.
	for _, areaID := range areas {
		dice := s.dice.MustDemographic(race, culture) // initialise our dice with probabilities

		prevAdults := []*structs.Person{} // some saved people to inter-link chunks
		prevChildren := []*structs.Person{}
		aliveArea := 0
		for {
			if aliveArea >= desiredArea {
				break
			}

			mp := simutil.NewMetaPeople()

			// we'll have 1 in 20 unmarried (just to force some variety)
			if aliveArea%20 > 0 && dice.RandomIsMarried(float64(dice.MinParentingAge())) {
				// spawn explicit familiy
				mp = s.spawnFamily(tick, areaID, race, culture)
				for _, p := range append(mp.Adults, mp.Children...) {
					if p.DeathTick <= 0 {
						aliveArea++
					}
				}
			} else {
				// spawn random couples
				for i := 0; i < s.dice.Intn(5)+1; i++ {
					aliveArea += s.spawnCouple(tick, areaID, race, culture, mp)
				}
			}

			// randomly add a few relationships
			if len(prevAdults) > 0 && len(mp.Adults) > 0 {
				seen := map[string]bool{}
				for _, a := range mp.Adults {
					for _, p := range rng.ChooseIndexes(len(prevAdults), s.dice.Intn(3)) {
						k := fmt.Sprintf("%s-%s", a.ID, prevAdults[p].ID)
						_, ok := seen[k]
						if ok {
							continue
						}
						seen[k] = true
						relationships, trust := dice.RandomRelationship(a.ID, prevAdults[p].ID)
						mp.Relations = append(mp.Relations, relationships...)
						mp.Trust = append(mp.Trust, trust...)
					}
				}
			}
			if len(prevChildren) > 0 && len(mp.Children) > 0 {
				seen := map[string]bool{}
				for _, a := range mp.Children {
					for _, p := range rng.ChooseIndexes(len(prevChildren), s.dice.Intn(2)) {
						k := fmt.Sprintf("%s-%s", a.ID, prevChildren[p].ID)
						_, ok := seen[k]
						if ok {
							continue
						}
						seen[k] = true
						relationships, trust := dice.RandomRelationship(a.ID, prevChildren[p].ID)
						mp.Relations = append(mp.Relations, relationships...)
						mp.Trust = append(mp.Trust, trust...)
					}
				}
			}

			prevAdults = mp.Adults
			prevChildren = mp.Children

			err := simutil.WriteMetaPeople(s.dbconn, mp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
