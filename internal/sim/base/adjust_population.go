package base

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/internal/sim/simutil"
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
		mp := simutil.NewMetaPeople()

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

			simutil.AddChildToFamily(dice, child, f)
			simutil.AddParentChildRelations(dice, mp, child, f)

			mp.Children = append(mp.Children, child)

			mp.Events = append(mp.Events, simutil.NewBirthEvent(child, f.ID))

			if dice.RandomDeathInfantMortality() {
				child.DeathMetaReason = diedInChildbirth
				child.DeathTick = tick
				child.DeathMetaKey = structs.MetaKeyPerson
				child.DeathMetaVal = f.FemaleID
				mp.Events = append(mp.Events, simutil.NewDeathEvent(child))
			}

			if dice.RandomAdultDeathInChildbirth() {
				mp.Events = append(mp.Events, simutil.NewDeathEventWithCause(
					tick,
					f.FemaleID,
					structs.MetaKeyPerson,
					child.ID,
				))
				f.IsChildBearing = false
				f.WidowedTick = tick
				mp.Relations = append(mp.Relations,
					&structs.Tuple{Subject: f.MaleID, Object: f.FemaleID, Value: int(structs.PersonalRelationExWife)},
					&structs.Tuple{Subject: f.FemaleID, Object: f.MaleID, Value: int(structs.PersonalRelationExHusband)},
				)
			}

			mp.Families = append(mp.Families, f)
		}

		// write new children & changes to families
		err = simutil.WriteMetaPeople(s.dbconn, mp)
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
		mp := simutil.NewMetaPeople()

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
				mp.Families = append(mp.Families, f)
			}
		}

		// write new children & changes to families
		err = simutil.WriteMetaPeople(s.dbconn, mp)
		if err != nil {
			return err
		}

		if token == "" {
			break
		}
	}

	return nil
}
