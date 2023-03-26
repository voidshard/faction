package sim

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// demographicsRand holds random dice for a demographics struct.
// We need quite a few with various average / deviation values so
// this helps keep things tidy.
type demographicsRand struct {
	familySize       *stats.Rand
	childbearingAge  *stats.Rand
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

func (s *simulationImpl) randPerson(rng *rand.Rand, demo *config.Demographics, dice *demographicsRand, areaID string) *structs.Person {
	p := &structs.Person{
		ID:     structs.NewID(),
		Race:   demo.Race,
		AreaID: areaID,
		Ethos: structs.Ethos{
			Altruism:  dice.ethosAltruism.Int(),
			Ambition:  dice.ethosAmbition.Int(),
			Tradition: dice.ethosTradition.Int(),
			Pacifism:  dice.ethosPacifism.Int(),
			Piety:     dice.ethosPiety.Int(),
			Caution:   dice.ethosCaution.Int(),
		},
		BirthTick: -1 * dice.childbearingAge.Int(),
		IsMale:    rng.Float64() <= 0.5,
	}
	if rng.Float64() <= demo.EthosBlackSheepProbability {
		blacksheep(rng, p)
	}
	return p
}

func (s *simulationImpl) randFaith(demo *config.Demographics, dice *demographicsRand, subject string) []*structs.Tuple {
	data := []*structs.Tuple{}

	count := dice.faithCount.Random()
	if count <= 0 {
		return data
	}

	for i := 0; i < count*2; i++ {
		faith := demo.Faiths[dice.faithOccur.Random()]
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

func (s *simulationImpl) randProfession(demo *config.Demographics, dice *demographicsRand, subject string) []*structs.Tuple {
	data := []*structs.Tuple{}

	count := dice.professionCount.Random()
	if count <= 0 {
		return data
	}

	hasPrimaryProfession := false
	for i := 0; i < count*2; i++ {
		prof := demo.Professions[dice.professionOccur.Random()]
		if hasPrimaryProfession && prof.ValidSideProfession {
			continue
		}
		hasPrimaryProfession = hasPrimaryProfession || prof.ValidSideProfession

		profDice := dice.professionLevel[prof.Name]

		data = append(data, &structs.Tuple{Subject: subject, Object: prof.Name, Value: profDice.Int()})
		if len(data) >= count {
			break
		}
	}

	return data
}

func (s *simulationImpl) spawnFamily(rng *rand.Rand, demo *config.Demographics, dice *demographicsRand, areaID string, pump *db.Pump) ([]*structs.Person, []*structs.Person) {
	mum := s.randPerson(rng, demo, dice, areaID)
	dad := s.randPerson(rng, demo, dice, areaID)

	mum.IsMale = false
	dad.IsMale = true
	family := &structs.Family{
		ID:             structs.NewID(mum.ID, dad.ID),
		AreaID:         areaID,
		IsChildBearing: true,
		MaleID:         dad.ID,
		FemaleID:       mum.ID,
	}
	adults := []*structs.Person{mum, dad}

	mumFaiths := s.randFaith(demo, dice, mum.ID)
	dadFaiths := s.randFaith(demo, dice, dad.ID)

	pump.SetTuples(db.RelationPersonProfessionSkill, s.randProfession(demo, dice, mum.ID)...)
	pump.SetTuples(db.RelationPersonProfessionSkill, s.randProfession(demo, dice, dad.ID)...)
	pump.SetTuples(db.RelationPersonReligionFaith, mumFaiths...)
	pump.SetTuples(db.RelationPersonReligionFaith, dadFaiths...)

	havingAffair := rng.Float64() <= demo.MarriageAffairProbability
	if havingAffair {
		affair := s.randPerson(rng, demo, dice, areaID)
		adults = append(adults, affair)

		pump.SetTuples(db.RelationPersonProfessionSkill, s.randProfession(demo, dice, affair.ID)...)
		pump.SetTuples(db.RelationPersonReligionFaith, s.randFaith(demo, dice, affair.ID)...)
		if affair.IsMale {
			pump.SetTuples(
				db.RelationPersonPersonTrust,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: int(structs.PersonalRelationLover)},
			)
			pump.SetTuples(
				db.RelationPersonPersonRelationship,
				&structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: dice.relationTrust.Int()},
				&structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: dice.relationTrust.Int()},
			)
		} else {
			pump.SetTuples(db.RelationPersonPersonTrust,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: int(structs.PersonalRelationLover)},
			)
			pump.SetTuples(
				db.RelationPersonPersonRelationship,
				&structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: dice.relationTrust.Int()},
				&structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: dice.relationTrust.Int()},
			)
		}
	}

	children := []*structs.Person{}
	for i := 0; i < dice.familySize.Int(); i++ {
		if mum.DeathTick > 0 {
			break
		}

		child := s.randPerson(rng, demo, dice, areaID)
		child.Ethos = *structs.EthosAverage(&child.Ethos, &mum.Ethos, &dad.Ethos) // average out parents

		child.BirthTick = -1 * rng.Intn(5) // TODO: add config around this
		if child.BirthTick < mum.BirthTick {
			child.BirthTick = 0
		}

		rel := structs.PersonalRelationDaughter
		if child.IsMale {
			rel = structs.PersonalRelationSon
		}

		pump.SetTuples(
			db.RelationPersonPersonRelationship,
			&structs.Tuple{Subject: mum.ID, Object: child.ID, Value: int(rel)},
			&structs.Tuple{Subject: dad.ID, Object: child.ID, Value: int(rel)},
			&structs.Tuple{Subject: child.ID, Object: mum.ID, Value: int(structs.PersonalRelationMother)},
			&structs.Tuple{Subject: child.ID, Object: dad.ID, Value: int(structs.PersonalRelationFather)},
		)
		pump.SetTuples(
			db.RelationPersonPersonTrust,
			&structs.Tuple{Subject: mum.ID, Object: child.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: dad.ID, Object: child.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: child.ID, Object: mum.ID, Value: dice.relationTrust.Int()},
			&structs.Tuple{Subject: child.ID, Object: dad.ID, Value: dice.relationTrust.Int()},
		)

		if i == 0 && rng.Float64() <= demo.ChildbearingDeathProbability { // check if mother dies
			mum.DeathTick = 1
			mum.DeathMetaReason = "died in childbirth"
			mum.DeathMetaKey = structs.MetaKeyPerson
			mum.DeathMetaVal = child.ID
		}

		if rng.Float64() <= demo.DeathInfantMortalityProbability { // check if child dies
			child.DeathMetaReason = dice.deathCauseReason[dice.deathCauseProb.Random()]
			child.DeathTick = 1
		}

		childFaiths := []*structs.Tuple{}
		if rng.Float64() <= 0.5 && mum.DeathTick <= 0 { // child takes faith from either mum or dad
			for _, f := range mumFaiths {
				childFaiths = append(childFaiths, &structs.Tuple{Subject: child.ID, Object: f.Object, Value: f.Value / 2})
			}
		} else {
			for _, f := range dadFaiths {
				childFaiths = append(childFaiths, &structs.Tuple{Subject: child.ID, Object: f.Object, Value: f.Value / 2})
			}
		}
		pump.SetTuples(db.RelationPersonReligionFaith, childFaiths...)
		children = append(children, child)
	}

	if mum.DeathTick > 0 || rng.Float64() <= demo.MarriageDivorceProbability {
		pump.SetTuples(
			db.RelationPersonPersonRelationship,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationExHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationExWife)},
		)
		pump.SetTuples(
			db.RelationPersonPersonTrust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / 2},
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / 2},
		)
	} else {
		div := 1
		if havingAffair {
			div = 4
		}
		pump.SetTuples(db.RelationPersonPersonRelationship,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationHusband)},
			&structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationWife)},
		)
		pump.SetTuples(db.RelationPersonPersonTrust,
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / div},
			&structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / div},
		)
	}

	pump.SetPeople(adults...)
	pump.SetPeople(children...)
	pump.SetFamilies(family)

	return adults, children
}

func randomRelationship(pump *db.Pump, personA, personB string, rng *rand.Rand, dice *demographicsRand) {
	trust := 1
	rel := structs.PersonalRelationCloseFriend

	switch dice.friendshipsProb.Random() {
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

	pump.SetTuples(
		db.RelationPersonPersonRelationship,
		&structs.Tuple{Subject: personA, Object: personB, Value: int(rel)},
		&structs.Tuple{Subject: personB, Object: personA, Value: int(rel)},
	)
	pump.SetTuples(
		db.RelationPersonPersonTrust,
		&structs.Tuple{Subject: personA, Object: personB, Value: trust * dice.relationTrust.Int()},
		&structs.Tuple{Subject: personB, Object: personA, Value: trust * dice.relationTrust.Int()},
	)
}

func (s *simulationImpl) SpawnPopulace(desiredTotal int, demo *config.Demographics, areas ...string) error {
	// TODO: support passing a 'Namer' to generate names
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pump := s.dbconn.NewPump()
	defer pump.Close()

	var err error
	go func() {
		for e := range pump.Errors() {
			if e == nil {
				continue
			}
			if err == nil {
				err = e
			} else {
				err = fmt.Errorf("%v %w", err, e)
			}
		}
	}()

	wg := &sync.WaitGroup{}

	desiredArea := desiredTotal / len(areas)

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
				if err != nil || aliveArea >= desiredArea {
					break
				}

				var adults, children []*structs.Person

				family := rng.Float64() <= demo.MarriageProbability
				if family {
					adults, children = s.spawnFamily(rng, demo, dice, areaID, pump)

					for _, p := range append(adults, children...) {
						if p.DeathTick <= 0 {
							aliveArea++
						}
					}
				} else {
					for i := 0; i < rng.Intn(5)+1; i++ {
						person := s.randPerson(rng, demo, dice, areaID)
						pump.SetTuples(db.RelationPersonProfessionSkill, s.randProfession(demo, dice, person.ID)...)
						pump.SetTuples(db.RelationPersonReligionFaith, s.randFaith(demo, dice, person.ID)...)

						if rng.Float64() < demo.DeathAdultMortalityProbability {
							person.DeathTick = 1
							person.DeathMetaReason = dice.deathCauseReason[dice.deathCauseProb.Random()]
						} else {
							aliveArea++
						}

						pump.SetPeople(person)
						adults = append(adults, person)

						if rng.Float64() > demo.MarriageProbability || person.DeathTick > 0 {
							continue
						}

						lover := s.randPerson(rng, demo, dice, areaID)
						lover.IsMale = !person.IsMale

						pump.SetTuples(db.RelationPersonProfessionSkill, s.randProfession(demo, dice, lover.ID)...)
						pump.SetTuples(db.RelationPersonReligionFaith, s.randFaith(demo, dice, lover.ID)...)

						rel := structs.PersonalRelationLover
						if rng.Float64() <= 0.1 {
							rel = structs.PersonalRelationFiance
						}

						pump.SetTuples(
							db.RelationPersonPersonRelationship,
							&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: int(rel)},
							&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: int(rel)},
						)
						pump.SetTuples(
							db.RelationPersonPersonTrust,
							&structs.Tuple{Subject: person.ID, Object: lover.ID, Value: dice.relationTrust.Int()},
							&structs.Tuple{Subject: lover.ID, Object: person.ID, Value: dice.relationTrust.Int()},
						)
						pump.SetPeople(lover)
						adults = append(adults, lover)
					}
				}

				if len(prevAdults) > 0 && len(adults) > 0 {
					for _, a := range adults {
						for _, p := range stats.ChooseIndexes(len(prevAdults), rng.Intn(3)) {
							randomRelationship(pump, a.ID, prevAdults[p].ID, rng, dice)
						}
					}
				}
				if len(prevChildren) > 0 && len(children) > 0 {
					for _, a := range children {
						for _, p := range stats.ChooseIndexes(len(prevChildren), rng.Intn(2)) {
							randomRelationship(pump, a.ID, prevChildren[p].ID, rng, dice)
						}
					}
				}

				prevAdults = adults
				prevChildren = children
			}
		}(a)
	}

	wg.Wait()
	return err
}

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
		familySize: stats.NewRand(
			0, demo.FamilySize.Max,
			demo.FamilySize.Mean, demo.FamilySize.Deviation,
		),
		childbearingAge: stats.NewRand(
			demo.ChildbearingAge.Min, demo.ChildbearingAge.Max,
			demo.ChildbearingAge.Mean, demo.ChildbearingAge.Deviation,
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
