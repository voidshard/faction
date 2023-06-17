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
	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

const (
	diedInChildbirth = "died in childbirth"
)

// demographicsRand holds random dice for a demographics struct.
// We need quite a few with various average / deviation values so
// this helps keep things tidy.
type demographicsRand struct {
	cfg *config.Demographics
	rng *rand.Rand

	familySize       *stats.Rand
	childbearingAge  *stats.Rand
	childbearingTerm *stats.Rand
	ethosAltruism    *stats.Rand
	ethosAmbition    *stats.Rand
	ethosTradition   *stats.Rand
	ethosPacifism    *stats.Rand
	ethosPiety       *stats.Rand
	ethosCaution     *stats.Rand
	professionLevel  map[string]*stats.Rand
	professionOccur  stats.Normalised
	professionCount  stats.Normalised
	faithLevel       map[string]*stats.Rand
	faithOccur       stats.Normalised
	faithCount       stats.Normalised
	deathCauseReason []string
	deathCauseProb   stats.Normalised
	relationTrust    *stats.Rand
	friendshipsProb  stats.Normalised
}

type metaPeople struct {
	adults    []*structs.Person
	children  []*structs.Person
	skills    []*structs.Tuple
	faith     []*structs.Tuple
	trust     []*structs.Tuple
	relations []*structs.Tuple
	families  []*structs.Family
}

func newMetaPeople() *metaPeople {
	return &metaPeople{
		adults:    []*structs.Person{},
		children:  []*structs.Person{},
		skills:    []*structs.Tuple{},
		faith:     []*structs.Tuple{},
		trust:     []*structs.Tuple{},
		relations: []*structs.Tuple{},
		families:  []*structs.Family{},
	}
}

func (s *Base) randPerson(dice *demographicsRand, areaID string) *structs.Person {
	p := &structs.Person{
		ID:     structs.NewID(),
		Race:   dice.cfg.Race,
		AreaID: areaID,
		Ethos: structs.Ethos{
			Altruism:  dice.ethosAltruism.Int(),
			Ambition:  dice.ethosAmbition.Int(),
			Tradition: dice.ethosTradition.Int(),
			Pacifism:  dice.ethosPacifism.Int(),
			Piety:     dice.ethosPiety.Int(),
			Caution:   dice.ethosCaution.Int(),
		},
		BirthTick: -1 * (dice.childbearingAge.Int() + dice.childbearingTerm.Int()),
		IsMale:    dice.rng.Float64() <= 0.5,
	}
	if dice.rng.Float64() <= dice.cfg.EthosBlackSheepProbability {
		blacksheep(dice.rng, p)
	}
	return p
}

func (s *Base) randFaith(dice *demographicsRand, subject string) []*structs.Tuple {
	data := []*structs.Tuple{}

	count := dice.faithCount.Int()
	if count <= 0 {
		return data
	}

	for i := 0; i < count*2; i++ {
		faith := dice.cfg.Faiths[dice.faithOccur.Int()]
		if len(data) > 0 && faith.IsMonotheistic { // can't add monotheistic faiths if we already have a faith
			continue
		}

		faithDice := dice.faithLevel[faith.ReligionID]

		data = append(data, &structs.Tuple{Subject: subject, Object: faith.ReligionID, Value: faithDice.Int()})
		if faith.IsMonotheistic { // we can't add any more faiths
			break
		}
	}

	return data
}

// randProfession returns a slice of tuples representing a person's skills at various professions.
// The first tuple is the person's preferred profession, which isn't necessarily what they're
// *best* at, but what they probably want to do.
//
// TODO: we should weight this based on personal ethos, currently we assume what they're best
// at weighted heavily towards (*2) non side professions (ie. more dedicated / high skill trades).
// That is, we assume the 'side professions' are what someone worked at when younger, does in order
// to cover costs or whatever while training for their desired profession.
func (s *Base) randProfession(dice *demographicsRand, subject string) []*structs.Tuple {
	data := []*structs.Tuple{}

	count := dice.professionCount.Int()
	if count <= 0 {
		return data
	}

	// preferred profession / trade
	score := -1
	preferrence := -1

	hasPrimaryProfession := false
	for i := 0; i < count*2; i++ {
		prof := dice.cfg.Professions[dice.professionOccur.Int()]
		if hasPrimaryProfession && prof.ValidSideProfession {
			continue
		}
		hasPrimaryProfession = hasPrimaryProfession || prof.ValidSideProfession

		profDice := dice.professionLevel[prof.Name]

		last := &structs.Tuple{Subject: subject, Object: prof.Name, Value: profDice.Int()}

		newScore := last.Value
		if !prof.ValidSideProfession { // implies dedicated trade
			newScore *= 2
		}
		if newScore > score {
			score = newScore
			preferrence = len(data) - 1 // preferred profession index
		}

		data = append(data, last)
		if len(data) >= count {
			break
		}
	}

	if preferrence > 0 {
		// move preffered role to the front
		data[0], data[preferrence] = data[preferrence], data[0]
	}

	return data
}

func (s *Base) spawnFamily(dice *demographicsRand, areaID string) *metaPeople {
	mp := newMetaPeople()

	mum := s.randPerson(dice, areaID)
	dad := s.randPerson(dice, areaID)
	mum.IsMale = false
	dad.IsMale = true
	mp.adults = append(mp.adults, mum, dad)

	eldest := dad.BirthTick
	if mum.BirthTick < eldest { // ie. "is born before"
		eldest = mum.BirthTick
	}

	family := &structs.Family{
		Ethos:               *structs.EthosAverage(&mum.Ethos, &dad.Ethos),
		ID:                  structs.NewID(mum.ID, dad.ID),
		Race:                mum.Race,
		AreaID:              areaID,
		IsChildBearing:      true,
		MaleID:              dad.ID,
		FemaleID:            mum.ID,
		MaxChildBearingTick: eldest + int(dice.cfg.ChildbearingAge.Max),
	}
	mp.families = append(mp.families, family)

	skills := s.randProfession(dice, mum.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		mum.PreferredProfession = skills[0].Object
	}
	skills = s.randProfession(dice, dad.ID)
	if len(skills) > 0 {
		mp.skills = append(mp.skills, skills...)
		dad.PreferredProfession = skills[0].Object
	}

	mumFaiths := s.randFaith(dice, mum.ID)
	dadFaiths := s.randFaith(dice, dad.ID)
	mp.faith = append(mp.faith, mumFaiths...)
	mp.faith = append(mp.faith, dadFaiths...)

	havingAffair := dice.rng.Float64() <= dice.cfg.MarriageAffairProbability
	if havingAffair {
		affair := s.randPerson(dice, areaID)
		mp.adults = append(mp.adults, affair)
		skills = s.randProfession(dice, affair.ID)
		if len(skills) > 0 {
			mp.skills = append(mp.skills, skills...)
			affair.PreferredProfession = skills[0].Object
		}
		mp.faith = append(mp.faith, s.randFaith(dice, affair.ID)...)
		if affair.IsMale {
			mp.relations = append(
				mp.relations,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: int(structs.PersonalRelationLover)},
			)
			mp.trust = append(
				mp.trust,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: dice.relationTrust.Int()},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: dice.relationTrust.Int()},
			)
		} else {
			mp.relations = append(
				mp.relations,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: int(structs.PersonalRelationLover)},
			)
			mp.trust = append(
				mp.trust,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: dice.relationTrust.Int()},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: dice.relationTrust.Int()},
			)
		}
	}

	for i := 0; i < dice.familySize.Int(); i++ {
		if mum.DeathTick > 0 {
			break
		}

		child := s.randPerson(dice, areaID)
		child.BirthTick = -1

		if dice.rng.Float64() <= dice.cfg.DeathInfantMortalityProbability { // check if child dies
			child.DeathMetaReason = diedInChildbirth
			child.DeathTick = 1
		}

		addChildToFamily(dice, child, family)
		addParentChildRelations(dice, mp, child, family)

		if i == 0 && dice.rng.Float64() <= dice.cfg.ChildbearingDeathProbability { // check if mother dies
			mum.DeathTick = 1
			mum.DeathMetaReason = diedInChildbirth
			mum.DeathMetaKey = structs.MetaKeyPerson
			mum.DeathMetaVal = child.ID
			family.IsChildBearing = false
		}

		childFaiths := []*structs.Tuple{}
		if dice.rng.Float64() <= 0.5 && mum.DeathTick <= 0 { // child takes faith from either mum or dad
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

	if mum.DeathTick > 0 || dice.rng.Float64() <= dice.cfg.MarriageDivorceProbability {
		mp.relations = append(
			mp.relations,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationExHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationExWife)},
		)
		mp.trust = append(
			mp.trust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / 2},
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / 2},
		)
		family.IsChildBearing = false
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
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / div},
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / div},
		)
	}

	if len(mp.children) > 1 {
		// add inter-sibling relations
		for i, child := range mp.children {
			for j, sibling := range mp.children {
				if i == j {
					continue
				}
				siblingRelationship(dice.relationTrust, mp, child, sibling)
			}
		}
	}

	if len(mp.children) >= int(dice.cfg.FamilySize.Max) {
		family.IsChildBearing = false
	}

	return mp
}

func siblingRelationship(trust *stats.Rand, mp *metaPeople, a, b *structs.Person) {
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
		&structs.Tuple{Subject: a.ID, Object: b.ID, Value: trust.Int()},
		&structs.Tuple{Subject: b.ID, Object: a.ID, Value: trust.Int()},
	)
}

func randomRelationship(personA, personB string, dice *demographicsRand) ([]*structs.Tuple, []*structs.Tuple) {
	trust := 1
	rel := structs.PersonalRelationCloseFriend

	switch dice.friendshipsProb.Int() {
	case 0: // default set above
		break
	case 1:
		rel = structs.PersonalRelationFriend
	case 2:
		trust = -1
		rel = structs.PersonalRelationEnemy
	case 3:
		trust = -1
		rel = structs.PersonalRelationHatedEnemy
	}

	return []*structs.Tuple{
			&structs.Tuple{Subject: personA, Object: personB, Value: int(rel)},
			&structs.Tuple{Subject: personB, Object: personA, Value: int(rel)},
		}, []*structs.Tuple{
			&structs.Tuple{Subject: personA, Object: personB, Value: trust * dice.relationTrust.Int()},
			&structs.Tuple{Subject: personB, Object: personA, Value: trust * dice.relationTrust.Int()},
		}
}

func (s *Base) demographicDice(name string) (*demographicsRand, error) {
	demo, ok := s.cfg.Demographics[name]
	if !ok {
		return nil, fmt.Errorf("unknown demographics %q", name)
	}
	return newDemographicsRand(demo), nil
}

func (s *Base) SpawnPopulace(desiredTotal int, name string, areas ...string) error {
	if len(areas) < 1 {
		return nil
	}

	demo, ok := s.cfg.Demographics[name]
	if !ok {
		return fmt.Errorf("unknown demographics %q", name)
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

			dice := newDemographicsRand(demo) // initialise our dice with probabilities
			prevAdults := []*structs.Person{} // some saved people to inter-link chunks
			prevChildren := []*structs.Person{}
			aliveArea := 0
			for {
				if finalerr != nil || aliveArea >= desiredArea {
					break
				}

				mp := newMetaPeople()

				family := dice.rng.Float64() <= demo.MarriageProbability
				if family {
					mp = s.spawnFamily(dice, areaID)
					for _, p := range append(mp.adults, mp.children...) {
						if p.DeathTick <= 0 {
							aliveArea++
						}
					}
				} else {
					for i := 0; i < dice.rng.Intn(5)+1; i++ {
						person := s.randPerson(dice, areaID)
						mp.adults = append(mp.adults, person)

						skills := s.randProfession(dice, person.ID)
						if len(skills) > 0 {
							mp.skills = append(mp.skills, skills...)
							person.PreferredProfession = skills[0].Object
						}
						mp.faith = append(mp.faith, s.randFaith(dice, person.ID)...)

						if dice.rng.Float64() < demo.DeathAdultMortalityProbability {
							person.DeathTick = 1
							person.DeathMetaReason = dice.deathCauseReason[dice.deathCauseProb.Int()]
						} else {
							aliveArea++
						}

						if dice.rng.Float64() > demo.MarriageProbability || person.DeathTick > 0 {
							continue
						}

						lover := s.randPerson(dice, areaID)
						lover.IsMale = !person.IsMale

						mp.adults = append(mp.adults, lover)
						skills = s.randProfession(dice, lover.ID)
						if len(skills) > 0 {
							mp.skills = append(mp.skills, skills...)
							lover.PreferredProfession = skills[0].Object
						}
						mp.faith = append(mp.faith, s.randFaith(dice, lover.ID)...)

						rel := structs.PersonalRelationLover
						if dice.rng.Float64() <= 0.1 {
							rel = structs.PersonalRelationFiance
						}

						mp.relations = append(
							mp.relations,
							&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: int(rel)},
							&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: int(rel)},
						)
						mp.trust = append(
							mp.trust,
							&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: dice.relationTrust.Int()},
							&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: dice.relationTrust.Int()},
						)
					}
				}

				if len(prevAdults) > 0 && len(mp.adults) > 0 {
					for _, a := range mp.adults {
						for _, p := range stats.ChooseIndexes(len(prevAdults), dice.rng.Intn(3)) {
							relationships, trust := randomRelationship(a.ID, prevAdults[p].ID, dice)
							mp.relations = append(mp.relations, relationships...)
							mp.trust = append(mp.trust, trust...)
						}
					}
				}
				if len(prevChildren) > 0 && len(mp.children) > 0 {
					for _, a := range mp.children {
						for _, p := range stats.ChooseIndexes(len(prevChildren), dice.rng.Intn(2)) {
							relationships, trust := randomRelationship(a.ID, prevChildren[p].ID, dice)
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

func writeMetaPeople(conn *db.FactionDB, p *metaPeople) error {
	return conn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetPeople(append(p.adults, p.children...)...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonProfessionSkill, p.skills...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonReligionFaith, p.faith...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonPersonTrust, p.trust...)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationPersonPersonRelationship, p.relations...)
		if err != nil {
			return err
		}
		return tx.SetFamilies(p.families...)
	})
}

// blacksheep randomly changes a single attribute of a person to some extreme value.
func blacksheep(rng *rand.Rand, p *structs.Person) {
	v := structs.MinTuple + rng.Intn(5)
	if rng.Float64() < 0.5 {
		v = 96 + rng.Intn(5)
	}

	switch rng.Intn(6) {
	case 0:
		p.Altruism = v
	case 1:
		p.Ambition = v
	case 2:
		p.Tradition = v
	case 3:
		p.Pacifism = v
	case 4:
		p.Piety = v
	case 5:
		p.Caution = v
	}
}

// AdjustPopulation adjusts the population of an area simulating natural life events
// like births, deaths etc.
func (s *Base) AdjustPopulation(area string) error {
	tick, err := s.dbconn.Tick()
	if err != nil {
		return err
	}

	// 1. families who can have children may do so
	err = s.growFamilies(tick, area)
	if err != nil {
		return err
	}

	// 2.a people who're too old may die
	//  .b people can randomly die of misc reasons (disease etc)

	return err
}

func (s *Base) lifeEvents(tick int, area string) error {
	return nil
}

func (s *Base) growFamilies(tick int, area string) error {
	// ultimately we iterate families that can have children in chunks
	// - for each chunk
	//  - for each family check if it
	//    - has a birth
	//    - has reached it's max children
	//  - write changes
	//  - for each person in {siblings of children or mothers who died in this chunk
	//    - if it's a mother, mark mother as dead
	//    - if it's a child mark as sibling
	//  - write changes
	//
	// The batched approach looks ugly but should make our code more efficient
	// that a nicer looking more query intensive approach

	token := (&dbutils.IterToken{Limit: 500, Offset: 0}).String()
	ff := []*db.FamilyFilter{&db.FamilyFilter{OnlyChildBearing: true, AreaID: area}}
	var (
		families []*structs.Family
		err      error
	)

	for {
		mp := newMetaPeople()

		mothersDiedInChildbirth := map[string]string{}     // mother id -> child id
		familiesNewSibling := map[string]*structs.Person{} // family id -> sibling

		pf := []*db.PersonFilter{}

		families, token, err = s.dbconn.Families(token, ff...)
		if err != nil {
			return nil
		}
		for _, f := range families {
			dice, err := s.demographicDice(f.Race)
			if err != nil {
				return nil
			}

			modified := false

			// check if a baby is born
			if f.PregnancyEnd == tick {
				modified = true
				f.PregnancyEnd = 0 // reset

				child := s.randPerson(dice, area)
				child.BirthTick = tick
				mp.children = append(mp.children, child)

				if dice.rng.Float64() <= dice.cfg.DeathInfantMortalityProbability { // check if child dies
					child.DeathMetaReason = diedInChildbirth
					child.DeathTick = tick
				}

				addChildToFamily(dice, child, f)
				addParentChildRelations(dice, mp, child, f)

				if f.NumberOfChildren > 1 {
					familiesNewSibling[f.ID] = child
					pf = append(pf, &db.PersonFilter{BirthFamilyID: f.ID, IncludeChildren: true})
				}

				if dice.rng.Float64() <= dice.cfg.ChildbearingDeathProbability { // check if mother dies
					mothersDiedInChildbirth[f.FemaleID] = child.ID
					pf = append(pf, &db.PersonFilter{ID: f.FemaleID})
					f.IsChildBearing = false
				}
			}

			// check if the family is too old to have children
			if f.MaxChildBearingTick >= tick || f.NumberOfChildren >= int(dice.cfg.FamilySize.Max) {
				modified = true
				f.IsChildBearing = false
			}

			// check if the family is expecting a baby
			if f.PregnancyEnd <= 0 && f.IsChildBearing && dice.rng.Float64() <= dice.cfg.ChildbearingProbability {
				modified = true
				f.PregnancyEnd = tick + dice.childbearingTerm.Int()
			}

			// if we changed it, add to batch to be written
			if modified {
				mp.families = append(mp.families, f)
			}
		}

		// write new children & changes to families
		err = writeMetaPeople(s.dbconn, mp)
		if err != nil {
			return err
		}

		// we probably should allow configuring this
		trust := stats.NewRand(20, structs.MaxTuple, structs.MaxTuple/2, structs.MaxTuple/4)

		// modify the people we need to (either siblings or mothers)
		var (
			people []*structs.Person
			ptoken string
		)
		for {
			mp = newMetaPeople()

			people, ptoken, err = s.dbconn.People(ptoken, pf...)
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

				siblingRelationship(trust, mp, p, child)
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

		if token == "" {
			break
		}
	}

	return nil
}

func addChildToFamily(dice *demographicsRand, child *structs.Person, f *structs.Family) {
	child.Ethos = *structs.EthosAverage(&child.Ethos, &f.Ethos)
	child.BirthFamilyID = f.ID
	child.IsChild = true
	f.NumberOfChildren++
	if f.NumberOfChildren >= int(dice.cfg.FamilySize.Max) {
		f.IsChildBearing = false
	}
}

// addParentChildRelations adds parent & grandparent relationships
func addParentChildRelations(dice *demographicsRand, mp *metaPeople, child *structs.Person, f *structs.Family) {
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
		&structs.Tuple{Subject: f.FemaleID, Object: child.ID, Value: dice.relationTrust.Int()},
		&structs.Tuple{Subject: f.MaleID, Object: child.ID, Value: dice.relationTrust.Int()},
		&structs.Tuple{Subject: child.ID, Object: f.FemaleID, Value: dice.relationTrust.Int()},
		&structs.Tuple{Subject: child.ID, Object: f.MaleID, Value: dice.relationTrust.Int()},
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
			&structs.Tuple{Subject: f.MaGrandmaID, Object: child.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandmaID, Value: dice.relationTrust.Int()},
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
			&structs.Tuple{Subject: f.MaGrandpaID, Object: child.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: child.ID, Object: f.MaGrandpaID, Value: dice.relationTrust.Int()},
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
			&structs.Tuple{Subject: f.PaGrandmaID, Object: child.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandmaID, Value: dice.relationTrust.Int()},
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
			&structs.Tuple{Subject: f.PaGrandpaID, Object: child.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: child.ID, Object: f.PaGrandpaID, Value: dice.relationTrust.Int()},
		)
	}
}

func newDemographicsRand(demo *config.Demographics) *demographicsRand {
	// professions
	skills := map[string]*stats.Rand{}
	profOccurProb := []float64{}
	for _, profession := range demo.Professions {
		skills[profession.Name] = stats.NewRand(10, structs.MaxTuple, profession.Mean, profession.Deviation)
		profOccurProb = append(profOccurProb, profession.Occurs)
	}

	// faiths
	faiths := map[string]*stats.Rand{}
	faithOccurProb := []float64{}
	for _, faith := range demo.Faiths {
		faiths[faith.ReligionID] = stats.NewRand(10, structs.MaxTuple, faith.Mean, faith.Deviation)
		faithOccurProb = append(faithOccurProb, faith.Occurs)
	}

	// deaths
	deathProb := []float64{}
	deathReason := []string{}
	for cause, prob := range demo.DeathCauseNaturalProbability {
		deathProb = append(deathProb, prob)
		deathReason = append(deathReason, cause)
	}

	return &demographicsRand{
		cfg: demo,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		familySize: stats.NewRand(
			0, demo.FamilySize.Max,
			demo.FamilySize.Mean, demo.FamilySize.Deviation,
		),
		childbearingAge: stats.NewRand(
			demo.ChildbearingAge.Min, demo.ChildbearingAge.Max,
			demo.ChildbearingAge.Mean, demo.ChildbearingAge.Deviation,
		),
		childbearingTerm: stats.NewRand(
			demo.ChildbearingTerm.Min, demo.ChildbearingTerm.Max,
			demo.ChildbearingTerm.Mean, demo.ChildbearingTerm.Deviation,
		),
		ethosAltruism:    stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(demo.EthosMean.Altruism), float64(demo.EthosDeviation.Altruism)),
		ethosAmbition:    stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(demo.EthosMean.Ambition), float64(demo.EthosDeviation.Ambition)),
		ethosTradition:   stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(demo.EthosMean.Tradition), float64(demo.EthosDeviation.Tradition)),
		ethosPacifism:    stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(demo.EthosMean.Pacifism), float64(demo.EthosDeviation.Pacifism)),
		ethosPiety:       stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(demo.EthosMean.Piety), float64(demo.EthosDeviation.Piety)),
		ethosCaution:     stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(demo.EthosMean.Caution), float64(demo.EthosDeviation.Caution)),
		professionLevel:  skills,
		professionOccur:  stats.NewNormalised(profOccurProb),
		professionCount:  stats.NewNormalised(demo.ProfessionProbability),
		faithLevel:       faiths,
		faithOccur:       stats.NewNormalised(faithOccurProb),
		faithCount:       stats.NewNormalised(demo.FaithProbability),
		relationTrust:    stats.NewRand(20, structs.MaxTuple, structs.MaxTuple/2, structs.MaxTuple/4),
		deathCauseReason: deathReason,
		deathCauseProb:   stats.NewNormalised(deathProb),
		friendshipsProb: stats.NewNormalised([]float64{
			demo.FriendshipCloseProbability,
			demo.FriendshipProbability,
			demo.EnemyProbability,
			demo.EnemyHatedProbability,
		}),
	}
}
