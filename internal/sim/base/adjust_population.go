package base

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

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

	// check for marriages / divorces
	return s.lifeEvents(tick, area)
}

func (s *Base) lifeEvents(tick int, area string) error {
	return nil
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

		if token == "" {
			break
		}
	}

	return nil
}

// applyDeath enacts from hooks after someone has died - marking families as non child bearing etc.
//
// This isn't needed for deaths via childbirth, since we set families in that function (since we

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
