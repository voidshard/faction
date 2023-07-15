package base

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Demographic here represents the demographic information for a given
// (race, culture) tuple.
//
// We hold a number of random dice with configured probabilities / spreads
// for any relevant random information about a given member of this demographic.
type Demographic struct {
	rng *rand.Rand

	// the tuple we're talking about
	raceID  string
	cultID  string
	race    *config.Race
	culture *config.Culture

	// random dice covering various aspects of this demographic
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
	lifespan         *stats.Rand
	relationTrust    *stats.Rand
	friendshipsProb  stats.Normalised
}

func mult(in []float64) float64 {
	if in == nil || len(in) == 0 {
		return 1
	}
	m := 1.0
	for _, v := range in {
		m *= v
	}
	return m
}

// MaxFamilySize returns the max children families tend to have.
func (d *Demographic) MaxFamilySize() int {
	return int(d.culture.FamilySize.Max)
}

// AdultMortalityProbability returns the probability that an adult will die of "natural" causes.
func (d *Demographic) AdultMortalityProbability(in ...float64) float64 {
	return d.culture.DeathAdultMortalityProbability * mult(in)
}

// RandomChildbearingTerm returns a random childbearing term (ticks) for a person.
func (d *Demographic) RandomChildbearingTerm() int {
	return d.childbearingTerm.Int()
}

// RandomBecomesPregrant rolls to see if someone becomes pregnant.
func (d *Demographic) RandomBecomesPregnant(in ...float64) bool {
	return d.rng.Float64() < d.culture.ChildbearingProbability*mult(in)
}

// RandomDeathAdultMortality rolls to see if someone dies of "natural" causes as an adult.
// Returns the cause of death and if the person died.
func (d *Demographic) RandomDeathAdultMortality(in ...float64) (string, bool) {
	cause := d.deathCauseReason[d.deathCauseProb.Int()]
	return cause, d.rng.Float64() < d.race.DeathAdultMortalityProbability*mult(in)
}

// RandomAdultDeathInChildbirth rolls to see if an adult dies in childbirth.
func (d *Demographic) RandomAdultDeathInChildbirth(in ...float64) bool {
	return d.rng.Float64() < d.race.ChildbearingDeathProbability*mult(in)
}

// RandomDeathInfantMortality rolls to see if someone dies of "natural" causes as an adult.
func (d *Demographic) RandomDeathInfantMortality(in ...float64) bool {
	return d.rng.Float64() < d.race.DeathInfantMortalityProbability*mult(in)
}

// RandomFamilySize returns a random family size (number of children) for a couple.
func (d *Demographic) RandomFamilySize() int {
	return d.familySize.Int()
}

// RandomTrust returns a random trust value for a person.
func (d *Demographic) RandomTrust() int {
	return d.relationTrust.Int()
}

// RandomIsHavingAffair returns if a person is having an affair on this tick
func (d *Demographic) RandomIsHavingAffair(in ...float64) bool {
	return d.rng.Float64() < d.culture.MarriageAffairProbability*mult(in)
}

// RandomIsMarried returns if a person is married
func (d *Demographic) RandomIsMarried(in ...float64) bool {
	return d.rng.Float64() < d.culture.MarriageProbability*mult(in)
}

// RandomIsDivorced returns if a person is divorced on this tick
func (d *Demographic) RandomIsDivorced(in ...float64) bool {
	return d.rng.Float64() < d.culture.MarriageDivorceProbability*mult(in)
}

// MinParentingAge returns the minimum age at which a person can have children.
func (d *Demographic) MinParentingAge() int {
	return int(d.race.ChildbearingAgeMin + d.race.ChildbearingTerm.Min)
}

// MaxParentingAge returns the maximum age at which a person can have children.
func (d *Demographic) MaxParentingAge() int {
	return int(d.race.ChildbearingAgeMax)
}

// RandomParentingAge returns a random age within expected child-bearing ticks.
func (d *Demographic) RandomParentingAge() int {
	return d.childbearingAge.Int() + d.childbearingTerm.Int()
}

// RandomLifespan returns how long a person will live for.
func (d *Demographic) RandomLifespan() int {
	return int(d.lifespan.Int())
}

// RandomIsMale returns if a person is male or not based on the Race demographic.
func (d *Demographic) RandomIsMale() bool {
	return d.race.IsMaleProbability <= d.rng.Float64()
}

// RandomPerson returns a random person for this demographic.
func (d *Demographic) RandomPerson(areaID string) *structs.Person {
	birth := -1 * int(d.MinParentingAge())
	return &structs.Person{
		ID:               structs.NewID(),
		Race:             d.raceID,
		Culture:          d.cultID,
		AreaID:           areaID,
		Ethos:            d.RandomEthos(),
		BirthTick:        birth,
		AdulthoodTick:    int(d.race.ChildbearingAgeMin) + birth,
		NaturalDeathTick: d.RandomLifespan() + birth,
		IsMale:           d.RandomIsMale(),
		Random:           int(d.rng.Float64() * structs.PersonRandomMax),
	}
}

// RandomFaith returns a random faith for a person.
func (d *Demographic) RandomFaith(subject string) []*structs.Tuple {
	data := []*structs.Tuple{}
	count := d.faithCount.Int()
	if count <= 0 {
		return data
	}
	for i := 0; i < count*2; i++ {
		faith := d.culture.Faiths[d.faithOccur.Int()]
		if len(data) > 0 && faith.IsMonotheistic { // can't add monotheistic faiths if we already have a faith
			continue
		}
		faithDice := d.faithLevel[faith.ReligionID]
		data = append(data, &structs.Tuple{Subject: subject, Object: faith.ReligionID, Value: faithDice.Int()})
		if faith.IsMonotheistic { // we can't add any more faiths
			break
		}
	}
	return data
}

// RandomEthos returns a random ethos for a person.
func (d *Demographic) RandomEthos() structs.Ethos {
	e := structs.Ethos{
		// values within cultural norms
		Altruism:  d.ethosAltruism.Int(),
		Ambition:  d.ethosAmbition.Int(),
		Tradition: d.ethosTradition.Int(),
		Pacifism:  d.ethosPacifism.Int(),
		Piety:     d.ethosPiety.Int(),
		Caution:   d.ethosCaution.Int(),
	}
	if d.rng.Float64() <= d.culture.EthosBlackSheepProbability {
		tenth := structs.MaxTuple / 10

		// randomly flip some values to extremes
		v := structs.MinTuple + d.rng.Intn(tenth) // very low
		if d.rng.Float64() < 0.5 {
			v = structs.MaxTuple - d.rng.Intn(tenth) // very high
		}
		switch d.rng.Intn(6) {
		case 0:
			e.Altruism = v
		case 1:
			e.Ambition = v
		case 2:
			e.Tradition = v
		case 3:
			e.Pacifism = v
		case 4:
			e.Piety = v
		case 5:
			e.Caution = v
		}
	}
	return e
}

// RandomProfession returns a slice of tuples representing a person's skills at various professions.
// The first tuple is the person's preferred profession, which isn't necessarily what they're
// *best* at, but what they probably want to do.
//
// TODO: we should weight this based on personal ethos, currently we assume what they're best
// at weighted heavily towards (*2) non side professions (ie. more dedicated / high skill trades).
// That is, we assume the 'side professions' are what someone worked at when younger, does in order
// to cover costs or whatever while training for their desired profession.
func (d *Demographic) RandomProfession(subject string) []*structs.Tuple {
	data := []*structs.Tuple{}

	count := d.professionCount.Int()
	if count <= 0 {
		return data
	}

	chosen := map[string]config.Profession{}

	hasPrimaryProfession := false
	for i := 0; i < count*2; i++ {
		prof := d.culture.Professions[d.professionOccur.Int()]
		if hasPrimaryProfession && prof.ValidSideProfession {
			continue
		}
		hasPrimaryProfession = hasPrimaryProfession || prof.ValidSideProfession

		profDice := d.professionLevel[prof.Name]
		last := &structs.Tuple{Subject: subject, Object: prof.Name, Value: profDice.Int()}
		chosen[prof.Name] = prof

		data = append(data, last)
		if len(data) >= count {
			break
		}
	}

	sort.Slice(data, func(i, j int) bool {
		iProf, iok := chosen[data[i].Subject]
		jProf, jok := chosen[data[j].Subject]

		ival := data[i].Value
		jval := data[j].Value

		if iok {
			if !iProf.ValidSideProfession {
				ival *= 2
			}
		}
		if jok {
			if !jProf.ValidSideProfession {
				jval *= 2
			}
		}

		return ival > jval // we want highest first
	})

	return data
}

// RandomRelationship returns a random relationship between two people.
//
// Returns:
//   - a slice of tuples representing the relationship between the two people
//   - a slice of tuples representing the trust between the two people
func (d *Demographic) RandomRelationship(personA, personB string) ([]*structs.Tuple, []*structs.Tuple) {
	trust := 1
	rel := structs.PersonalRelationCloseFriend
	switch d.friendshipsProb.Int() {
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
			&structs.Tuple{Subject: personA, Object: personB, Value: trust * d.relationTrust.Int()},
			&structs.Tuple{Subject: personB, Object: personA, Value: trust * d.relationTrust.Int()},
		}
}

func newDemographic(rid, cid string, race *config.Race, culture *config.Culture) *Demographic {
	// professions
	skills := map[string]*stats.Rand{}
	profOccurProb := []float64{}
	for _, profession := range culture.Professions {
		skills[profession.Name] = stats.NewRand(10, structs.MaxTuple, profession.Mean, profession.Deviation)
		profOccurProb = append(profOccurProb, profession.Occurs)
	}

	// faiths
	faiths := map[string]*stats.Rand{}
	faithOccurProb := []float64{}
	for _, faith := range culture.Faiths {
		faiths[faith.ReligionID] = stats.NewRand(10, structs.MaxTuple, faith.Mean, faith.Deviation)
		faithOccurProb = append(faithOccurProb, faith.Occurs)
	}

	// deaths
	deathProb := []float64{}
	deathReason := []string{}
	for cause, prob := range culture.DeathCauseNaturalProbability {
		deathProb = append(deathProb, prob)
		deathReason = append(deathReason, cause)
	}

	// all the dice we might want to roll for this combination, plus the configs themselves
	return &Demographic{
		raceID:  rid,
		cultID:  cid,
		race:    race,
		culture: culture,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		familySize: stats.NewRand(
			0, culture.FamilySize.Max,
			culture.FamilySize.Mean, culture.FamilySize.Deviation,
		),
		childbearingAge: stats.NewRand(
			race.ChildbearingAgeMin, race.ChildbearingAgeMax,
			culture.ChildbearingAgeMean, culture.ChildbearingAgeDeviation,
		),
		childbearingTerm: stats.NewRand(
			race.ChildbearingTerm.Min, race.ChildbearingTerm.Max,
			math.Max(math.Min(race.ChildbearingTerm.Mean, race.ChildbearingTerm.Max), race.ChildbearingTerm.Min),
			race.ChildbearingTerm.Deviation,
		),
		ethosAltruism:    stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(culture.EthosMean.Altruism), float64(culture.EthosDeviation.Altruism)),
		ethosAmbition:    stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(culture.EthosMean.Ambition), float64(culture.EthosDeviation.Ambition)),
		ethosTradition:   stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(culture.EthosMean.Tradition), float64(culture.EthosDeviation.Tradition)),
		ethosPacifism:    stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(culture.EthosMean.Pacifism), float64(culture.EthosDeviation.Pacifism)),
		ethosPiety:       stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(culture.EthosMean.Piety), float64(culture.EthosDeviation.Piety)),
		ethosCaution:     stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(culture.EthosMean.Caution), float64(culture.EthosDeviation.Caution)),
		professionLevel:  skills,
		professionOccur:  stats.NewNormalised(profOccurProb),
		professionCount:  stats.NewNormalised(culture.ProfessionProbability),
		faithLevel:       faiths,
		faithOccur:       stats.NewNormalised(faithOccurProb),
		faithCount:       stats.NewNormalised(culture.FaithProbability),
		relationTrust:    stats.NewRand(20, structs.MaxTuple, structs.MaxTuple/2, structs.MaxTuple/4),
		deathCauseReason: deathReason,
		deathCauseProb:   stats.NewNormalised(deathProb),
		lifespan: stats.NewRand(
			race.Lifespan.Min, race.Lifespan.Max,
			race.Lifespan.Mean, race.Lifespan.Deviation,
		),
		friendshipsProb: stats.NewNormalised([]float64{
			culture.FriendshipCloseProbability,
			culture.FriendshipProbability,
			culture.EnemyProbability,
			culture.EnemyHatedProbability,
		}),
	}
}
