/*
random_populate.go - Randomly populate the world with people and families.
*/
package base

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	dg "github.com/voidshard/faction/internal/demographics"
	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	diedInChildbirth = "died in childbirth"
	diedOfOldAge     = "died of old age"

	// "random" is a value 0-randomValueMax that is used to pull rows at random
	// from the db with some probability.
	randomValueMax = 1000000.0
)

func (s *Base) spawnFamily(tick int, areaID, race, culture string) *metaPeople {
	demo := s.dice.MustDemographic(race, culture)

	mp := newMetaPeople()

	// nb. we're creating mum & dad in the past so that our children (if any) are born "now"
	mum := demo.RandomPerson(areaID)
	mum.SetBirthTick(tick - demo.RandomParentingAge())

	dad := demo.RandomPerson(areaID)
	dad.SetBirthTick(tick - demo.RandomParentingAge())

	mum.IsMale = false
	dad.IsMale = true
	mp.adults = append(mp.adults, mum, dad)

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
	}
	mp.families = append(mp.families, family)

	skills := demo.RandomProfession(mum.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		mum.PreferredProfession = skills[0].Object
	}
	skills = demo.RandomProfession(dad.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		dad.PreferredProfession = skills[0].Object
	}

	mumFaiths := demo.RandomFaith(mum.ID)
	dadFaiths := demo.RandomFaith(dad.ID)
	mp.faith = append(mp.faith, mumFaiths...)
	mp.faith = append(mp.faith, dadFaiths...)

	havingAffair := demo.RandomIsHavingAffair(parentingTicks)
	if havingAffair {
		affair := demo.RandomPerson(areaID)
		mp.adults = append(mp.adults, affair)
		skills = demo.RandomProfession(affair.ID)
		if len(skills) > 0 {
			mp.skills = append(mp.skills, skills...)
			affair.PreferredProfession = skills[0].Object
		}
		mp.faith = append(mp.faith, demo.RandomFaith(affair.ID)...)
		if affair.IsMale {
			mp.relations = append(
				mp.relations,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: int(structs.PersonalRelationLover)},
			)
			mp.trust = append(
				mp.trust,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: demo.RandomTrust()},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: demo.RandomTrust()},
			)
		} else {
			mp.relations = append(
				mp.relations,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: int(structs.PersonalRelationLover)},
			)
			mp.trust = append(
				mp.trust,
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

		if demo.RandomDeathInfantMortality() {
			child.DeathMetaReason = diedInChildbirth
			child.DeathTick = tick
			child.DeathMetaKey = structs.MetaKeyPerson
			child.DeathMetaVal = mum.ID
		}

		addChildToFamily(demo, child, family)
		addParentChildRelations(demo, mp, child, family)

		if i == 0 && demo.RandomAdultDeathInChildbirth() {
			mum.DeathTick = tick
			mum.DeathMetaReason = diedInChildbirth
			mum.DeathMetaKey = structs.MetaKeyPerson
			mum.DeathMetaVal = child.ID
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
		mp.faith = append(mp.faith, childFaiths...)
		mp.children = append(mp.children, child)
	}

	if mum.DeathTick > 0 || demo.RandomIsDivorced(parentingTicks) {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationExHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationExWife)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: demo.RandomTrust() / 2},
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: demo.RandomTrust() / 2},
		)
		family.IsChildBearing = false
		if mum.DeathTick > 0 {
			family.WidowedTick = tick
		} else {
			family.DivorceTick = tick
		}
	} else {
		div := 1
		if havingAffair {
			div = 4
		}
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationWife)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: demo.RandomTrust() / div},
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: demo.RandomTrust() / div},
		)
	}

	if len(mp.children) > 1 {
		// add inter-sibling relations
		for i, child := range mp.children {
			for j, sibling := range mp.children {
				if i == j {
					continue
				}
				siblingRelationship(demo, mp, child, sibling)
			}
		}
	}

	if len(mp.children) >= demo.MaxFamilySize() {
		family.IsChildBearing = false
	}

	return mp
}

func siblingRelationship(dice *dg.Demographic, mp *metaPeople, a, b *structs.Person) {
	brel := structs.PersonalRelationSister
	if b.IsMale {
		brel = structs.PersonalRelationBrother
	}

	arel := structs.PersonalRelationSister
	if a.IsMale {
		arel = structs.PersonalRelationBrother
	}

	mp.relations = append(
		mp.relations,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: int(brel)},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: int(arel)},
	)
	mp.trust = append(
		mp.trust,
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: dice.RandomTrust()},
	)
}

func (s *Base) spawnCouple(tick int, areaID, race, culture string, mp *metaPeople) int {
	demo := s.dice.MustDemographic(race, culture)

	alive := 0

	person := demo.RandomPerson(areaID)
	mp.adults = append(mp.adults, person)

	skills := demo.RandomProfession(person.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		person.PreferredProfession = skills[0].Object
	}
	mp.faith = append(mp.faith, demo.RandomFaith(person.ID)...)

	cause, died := demo.RandomDeathAdultMortality()
	if died {
		person.DeathTick = tick
		person.DeathMetaReason = cause
	} else {
		alive++
	}

	if person.DeathTick > 0 {
		return alive
	}

	lover := demo.RandomPerson(areaID)
	lover.IsMale = !person.IsMale
	alive++

	mp.adults = append(mp.adults, lover)
	skills = demo.RandomProfession(lover.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		lover.PreferredProfession = skills[0].Object
	}
	mp.faith = append(mp.faith, demo.RandomFaith(lover.ID)...)

	rel := structs.PersonalRelationLover
	if s.dice.Float64() <= 0.1 {
		rel = structs.PersonalRelationFiance
	}

	mp.relations = append(
		mp.relations,
		&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: int(rel)},
		&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: int(rel)},
	)
	mp.trust = append(
		mp.trust,
		&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: demo.RandomTrust()},
		&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: demo.RandomTrust()},
	)

	return alive
}

func (s *Base) SpawnPopulace(desiredTotal int, race, culture string, areas []string) error {
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

	var finalerr error
	errs := make(chan error)

	go func() {
		for err := range errs {
			if err == nil {
				continue
			}
			if finalerr == nil {
				finalerr = err
			} else {
				finalerr = fmt.Errorf("%w %v", finalerr, err)
			}
		}
	}()

	desiredArea := desiredTotal / len(areas)
	wg := &sync.WaitGroup{}

	// we place people one area at a time because it feels more natural as this results in more
	// links between people within a given area than between areas.
	for _, a := range areas {
		wg.Add(1)
		go func(areaID string) {
			defer wg.Done()

			dice := s.dice.MustDemographic(race, culture) // initialise our dice with probabilities

			prevAdults := []*structs.Person{} // some saved people to inter-link chunks
			prevChildren := []*structs.Person{}
			aliveArea := 0
			for {
				if finalerr != nil || aliveArea >= desiredArea {
					break
				}

				mp := newMetaPeople()

				if dice.RandomIsMarried(float64(dice.MinParentingAge())) {
					// spawn explicit familiy
					mp = s.spawnFamily(tick, areaID, race, culture)
					for _, p := range append(mp.adults, mp.children...) {
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

				if len(prevAdults) > 0 && len(mp.adults) > 0 {
					for _, a := range mp.adults {
						for _, p := range stats.ChooseIndexes(len(prevAdults), s.dice.Intn(3)) {
							relationships, trust := dice.RandomRelationship(a.ID, prevAdults[p].ID)
							mp.relations = append(mp.relations, relationships...)
							mp.trust = append(mp.trust, trust...)
						}
					}
				}
				if len(prevChildren) > 0 && len(mp.children) > 0 {
					for _, a := range mp.children {
						for _, p := range stats.ChooseIndexes(len(prevChildren), s.dice.Intn(2)) {
							relationships, trust := dice.RandomRelationship(a.ID, prevChildren[p].ID)
							mp.relations = append(mp.relations, relationships...)
							mp.trust = append(mp.trust, trust...)
						}
					}
				}

				prevAdults = mp.adults
				prevChildren = mp.children

				errs <- writeMetaPeople(s.dbconn, mp)
			}
		}(a)
	}

	wg.Wait()
	close(errs)

	return finalerr
}

// AdjustPopulation adjusts the population of an area simulating natural life events
// like births, deaths etc.
func (s *Base) AdjustPopulation(area string) error {
	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	// check non-pregnant families for conceptions
	err = s.conceiveChildren(tick, area)
	if err != nil {
		return err
	}

	// check pregnant families for births
	err = s.birthChildren(tick, area)
	if err != nil {
		return err
	}

	// check for death(s)
	err = s.deathCheck(tick, area)
	if err != nil {
		return err
	}

	// check for deaths / marriages / divorces
	return s.lifeEvents(tick, area)
}

func (s *Base) deathCheck(tick int, area string) error {
	// We want to loop once over a given area for all races.
	//
	// So compile the most extreme death values - we'll select slightly more people
	// than needed from races whose death rates are lower, but we'll save doing
	// a per-race iteration .. which is almost certainly cheaper .. probably ..
	w := s.dice.MaxDeathAdultMortalityProbability()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	deathRange := int(w * randomValueMax)
	deathStart := rng.Intn(int(randomValueMax) - deathRange)

	token := (&dbutils.IterToken{Limit: 1000, Offset: 0}).String()
	ff := db.Q(
		// people will die of old age
		db.F(db.AreaID, db.Equal, area),
		db.F(db.DeathTick, db.Equal, 0),
		db.F(db.NaturalDeathTick, db.Less, tick+1),
	)
	if deathRange > 0 {
		ff.Or(
			// people will die of misc reasons
			db.F(db.AreaID, db.Equal, area),
			db.F(db.DeathTick, db.Equal, 0),
			db.F(db.Random, db.Greater, deathStart-1),
			db.F(db.Random, db.Less, deathStart+deathRange+1),
		)
	}
	fmt.Println("deathCheck", ff.String())

	var (
		people []*structs.Person
		err    error
	)
	for {
		people, token, err = s.dbconn.People(token, ff)
		if err != nil {
			return err
		}
		died := []*structs.Person{}

		for _, p := range people {
			dice := s.dice.MustDemographic(p.Race, p.Culture)

			// death range for this demographic
			demoDeathRange := int(dice.AdultMortalityProbability() * randomValueMax)
			if p.NaturalDeathTick > tick && p.Random > demoDeathRange+deathStart {
				continue
			}

			died = append(died, p)
			p.DeathTick = tick

			if p.NaturalDeathTick < tick+1 {
				p.DeathMetaReason = diedOfOldAge
			} else {
				cause, _ := dice.RandomDeathAdultMortality()
				p.DeathMetaReason = cause
			}
		}

		err = s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
			return tx.SetPeople(died...)
		})

		err = s.applyDeath(tick, died...)
		if err != nil {
			return err
		}

		if token == "" {
			break
		}
	}

	return nil
}

// applyDeath enacts from hooks after someone has died - marking families as non child bearing etc.
//
// This isn't needed for deaths via childbirth, since we set families in that function (since we
// already have a the families to hand).
func (s *Base) applyDeath(tick int, in ...*structs.Person) error {
	if len(in) == 0 {
		return nil
	}

	female := []string{}
	male := []string{}
	femaleDead := map[string]bool{}

	for _, d := range in {
		if d.IsMale {
			male = append(male, d.ID)
		} else {
			female = append(female, d.ID)
			femaleDead[d.ID] = true
		}
	}

	token := (&dbutils.IterToken{Limit: 1000, Offset: 0}).String()
	var (
		families []*structs.Family
		err      error
	)
	q := db.Q(
		db.F(db.MaleID, db.In, male),
	).Or(
		db.F(db.FemaleID, db.In, female),
	)
	for {
		families, token, err = s.dbconn.Families(token, q)
		if err != nil {
			return err
		}

		update := []*structs.Family{}
		rels := []*structs.Tuple{}

		for _, f := range families {
			if f.WidowedTick > 0 {
				// implies the other partner is already dead
				continue
			}

			f.IsChildBearing = false
			f.WidowedTick = tick

			// if a mother has died, the child has too :(
			_, ok := femaleDead[f.FemaleID]
			if ok {
				f.PregnancyEnd = 0
			}

			rels = append(
				rels,
				&structs.Tuple{Subject: f.MaleID, Object: f.FemaleID, Value: int(structs.PersonalRelationExWife)},
				&structs.Tuple{Subject: f.FemaleID, Object: f.MaleID, Value: int(structs.PersonalRelationExHusband)},
			)
			update = append(update, f)
		}

		s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
			err = tx.SetFamilies(update...)
			if err != nil {
				return err
			}
			return tx.SetTuples(db.RelationPersonPersonRelationship, rels...)
		})

		if token == "" {
			break
		}
	}

	return nil
}

func (s *Base) lifeEvents(tick int, area string) error {
	// a. people who're too old may die
	// b. people can randomly die of misc reasons (disease etc)
	// c. people can get married / divorced
	//token := (&dbutils.IterToken{Limit: 1000, Offset: 0}).String()
	//ff := db.Q(
	//		db.F(db.AreaID, db.Equal, area),
	//	)

	return nil
}

func (s *Base) conceiveChildren(tick int, area string) error {
	token := (&dbutils.IterToken{Limit: 1000, Offset: 0}).String()
	ff := db.Q( // all childbearing families in the area not expecting a birth
		db.F(db.AreaID, db.Equal, area),
		db.F(db.IsChildBearing, db.Equal, true),
		db.F(db.PregnancyEnd, db.Equal, 0),
	)
	var (
		families []*structs.Family
		err      error
	)
	for {
		mp := newMetaPeople()

		families, token, err = s.dbconn.Families(token, ff)
		if err != nil {
			return err
		}
		for _, f := range families {
			dice := s.dice.MustDemographic(f.Race, f.Culture)

			modified := false
			// check if the family is too old to have children
			if f.MaxChildBearingTick <= tick || f.NumberOfChildren >= dice.MaxFamilySize() {
				modified = true
				f.IsChildBearing = false
			}

			// check if the family is expecting a baby
			if f.PregnancyEnd <= 0 && f.IsChildBearing && dice.RandomBecomesPregnant() {
				modified = true
				f.PregnancyEnd = tick + dice.RandomChildbearingTerm()
			}

			if modified {
				mp.families = append(mp.families, f)
			}
		}

		// write new children & changes to families
		err = writeMetaPeople(s.dbconn, mp)
		if err != nil {
			return err
		}

		if token == "" {
			break
		}
	}

	return nil
}

// applyChildbirthEffects updates People in one of two cases;
// - an expectant mother has died in childbirth (mother needs to be marked as dead)
// - a new child is born (set up relationships with siblings)
func (s *Base) applyChildbirthEffects(tick int, mothersDiedInChildbirth map[string]string, familiesNewSibling map[string]*structs.Person) error {
	if len(mothersDiedInChildbirth)+len(familiesNewSibling) == 0 {
		return nil
	}

	// ok, now we need to find people we have to update
	families := []string{}
	mothers := []string{}
	for personID := range mothersDiedInChildbirth {
		mothers = append(mothers, personID)
	}
	for familyID := range familiesNewSibling {
		families = append(families, familyID)
	}

	pf := db.Q()
	if len(families) > 0 {
		pf.Or(db.F(db.BirthFamilyID, db.In, families))
	}
	if len(mothers) > 0 {
		pf.Or(db.F(db.ID, db.In, mothers)) // mothers who died in childbirth
	}

	// modify the people we need to (either siblings or mothers)
	var (
		people []*structs.Person
		ptoken string
		err    error
	)
	for {
		mp := newMetaPeople()

		people, ptoken, err = s.dbconn.People(ptoken, pf)
		if err != nil {
			return err
		}
		for _, p := range people {
			// there are currently two reasons we have a person here
			// 1. it's a mother we want to register as died in childbirth
			// 2. it's a family member (child) that has a new sibling
			childID, ok := mothersDiedInChildbirth[p.ID]
			if ok {
				p.DeathTick = tick
				p.DeathMetaReason = diedInChildbirth
				p.DeathMetaKey = structs.MetaKeyPerson
				p.DeathMetaVal = childID
				mp.adults = append(mp.adults, p)
				continue
			}

			child, ok := familiesNewSibling[p.BirthFamilyID]
			if !ok {
				continue
			}

			if p.ID == child.ID {
				continue
			}

			if !s.dice.IsValidDemographic(p.Race, p.Culture) {
				return fmt.Errorf("invalid demographic not found: [race] %s, [culture] %s", p.Race, p.Culture)
			}

			demo := s.dice.MustDemographic(p.Race, p.Culture)
			siblingRelationship(demo, mp, p, child)
			mp.children = append(mp.children, p)
		}

		err = writeMetaPeople(s.dbconn, mp)
		if err != nil {
			return err
		}

		if ptoken == "" {
			break
		}
	}

	return nil
}

func (s *Base) birthChildren(tick int, area string) error {
	token := (&dbutils.IterToken{Limit: 1000, Offset: 0}).String()
	ff := db.Q( // all childbearing families in the area expecting a birth
		db.F(db.AreaID, db.Equal, area),
		db.F(db.PregnancyEnd, db.Less, tick+1),
		db.F(db.PregnancyEnd, db.Greater, 0),
	)
	var (
		families []*structs.Family
		err      error
	)

	for {
		mp := newMetaPeople()

		mothersDiedInChildbirth := map[string]string{}     // mother id -> child id
		familiesNewSibling := map[string]*structs.Person{} // family id -> sibling

		families, token, err = s.dbconn.Families(token, ff)
		if err != nil {
			return err
		}
		for _, f := range families {
			if !s.dice.IsValidDemographic(f.Race, f.Culture) {
				return fmt.Errorf("invalid demographic not found: [race] %s, [culture] %s", f.Race, f.Culture)
			}
			dice := s.dice.MustDemographic(f.Race, f.Culture)

			f.PregnancyEnd = 0 // reset

			child := dice.RandomPerson(area)
			child.SetBirthTick(tick)

			addChildToFamily(dice, child, f)
			addParentChildRelations(dice, mp, child, f)

			mp.children = append(mp.children, child)

			if dice.RandomDeathInfantMortality() {
				child.DeathMetaReason = diedInChildbirth
				child.DeathTick = tick
				child.DeathMetaKey = structs.MetaKeyPerson
				child.DeathMetaVal = f.FemaleID
			}

			if f.NumberOfChildren > 1 {
				familiesNewSibling[f.ID] = child
			}

			if dice.RandomAdultDeathInChildbirth() {
				mothersDiedInChildbirth[f.FemaleID] = child.ID
				f.IsChildBearing = false
				f.WidowedTick = tick
				mp.relations = append(mp.relations,
					&structs.Tuple{Subject: f.MaleID, Object: f.FemaleID, Value: int(structs.PersonalRelationExWife)},
					&structs.Tuple{Subject: f.FemaleID, Object: f.MaleID, Value: int(structs.PersonalRelationExHusband)},
				)
			}

			mp.families = append(mp.families, f)
		}

		// write new children & changes to families
		err = writeMetaPeople(s.dbconn, mp)
		if err != nil {
			return err
		}

		// update people affected by the above changes
		err = s.applyChildbirthEffects(tick, mothersDiedInChildbirth, familiesNewSibling)
		if err != nil {
			return err
		}

		if token == "" {
			break
		}
	}

	return nil
}

func addChildToFamily(dice *dg.Demographic, child *structs.Person, f *structs.Family) {
	child.Ethos = *structs.EthosAverage(&child.Ethos, &f.Ethos)
	child.BirthFamilyID = f.ID
	child.IsChild = true
	child.Race = f.Race
	f.NumberOfChildren++
	if f.NumberOfChildren >= dice.MaxFamilySize() {
		f.IsChildBearing = false
	}
}

// addParentChildRelations adds parent & grandparent relationships
func addParentChildRelations(dice *dg.Demographic, mp *metaPeople, child *structs.Person, f *structs.Family) {
	rel := structs.PersonalRelationDaughter
	grel := structs.PersonalRelationGranddaughter
	if child.IsMale {
		rel = structs.PersonalRelationSon
		grel = structs.PersonalRelationGrandson
	}

	// parents
	mp.relations = append(
		mp.relations,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: int(rel)},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: int(rel)},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: int(structs.PersonalRelationMother)},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: int(structs.PersonalRelationFather)},
	)
	mp.trust = append(
		mp.trust,
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: dice.RandomTrust()},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: dice.RandomTrust()},
	)

	// grandparents (first generation families wont have grandparents)
	if f.MaGrandmaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: int(structs.PersonalRelationGrandmother)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: dice.RandomTrust()},
		)
	}
	if f.MaGrandpaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: int(structs.PersonalRelationGrandfather)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: dice.RandomTrust()},
		)
	}
	if f.PaGrandmaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: int(structs.PersonalRelationGrandmother)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: dice.RandomTrust()},
		)
	}
	if f.PaGrandpaID != "" {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: int(grel)},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: int(structs.PersonalRelationGrandfather)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: dice.RandomTrust()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: dice.RandomTrust()},
		)
	}
}
