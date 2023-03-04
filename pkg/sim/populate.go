package sim

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/structs"
)

// demographicsRand holds random dice for a demographics struct.
// We need quite a few with various average / deviation values so
// this helps keep things tidy.
type demographicsRand struct {
	familySize      *stats.Rand
	childbearingAge *stats.Rand
	ethosAltruism   *stats.Rand
	ethosAmbition   *stats.Rand
	ethosTradition  *stats.Rand
	ethosPacifism   *stats.Rand
	ethosPiety      *stats.Rand
	ethosCaution    *stats.Rand
	professionLevel map[string]*stats.Rand
	professionOccur stats.Normalised
	professionCount stats.Normalised
	faithLevel      map[string]*stats.Rand
	faithOccur      stats.Normalised
	faithCount      stats.Normalised

	relationTrust *stats.Rand
}

type peopleData struct {
	People    []*structs.Person
	Families  []*structs.Family
	Relations []*structs.Tuple
	Trust     []*structs.Tuple
	Skills    []*structs.Tuple
	Faith     []*structs.Tuple
}

func newPeopleData() *peopleData {
	return &peopleData{
		People:    []*structs.Person{},
		Families:  []*structs.Family{},
		Relations: []*structs.Tuple{},
		Trust:     []*structs.Tuple{},
		Skills:    []*structs.Tuple{},
		Faith:     []*structs.Tuple{},
	}
}

func (s *simulationImpl) randPerson(rng *rand.Rand, demo *structs.Demographics, dice *demographicsRand, areaID string) *structs.Person {
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

func (s *simulationImpl) randProfession(demo *structs.Demographics, dice *demographicsRand, subject string) []*structs.Tuple {
	data := []*structs.Tuple{}

	count := dice.professionCount.Random()
	if count <= 0 {
		return data
	}

	hasPrimaryProfession := false
	for i := 0; i < count*5; i++ {
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

func (s *simulationImpl) spawnFamily(rng *rand.Rand, demo *structs.Demographics, dice *demographicsRand, areaID string) *peopleData {
	data := newPeopleData()

	mum := s.randPerson(rng, demo, dice, areaID)
	dad := s.randPerson(rng, demo, dice, areaID)
	family := &structs.Family{
		ID:             structs.NewID(mum.ID, dad.ID),
		AreaID:         areaID,
		IsChildBearing: true,
		MaleID:         dad.ID,
		FemaleID:       mum.ID,
	}

	data.People = append(data.People, mum)
	data.People = append(data.People, dad)
	data.Families = append(data.Families, family)
	data.Skills = append(data.Skills, s.randProfession(demo, dice, mum.ID)...)
	data.Skills = append(data.Skills, s.randProfession(demo, dice, dad.ID)...)

	if rng.Float64() <= demo.MarriageAffairProbability {
		affair := s.randPerson(rng, demo, dice, areaID)
		data.People = append(data.People, affair)
		data.Skills = append(data.Skills, s.randProfession(demo, dice, affair.ID)...)
		if affair.IsMale {
			data.Relations = append(data.Relations, &structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)})
			data.Relations = append(data.Relations, &structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: int(structs.PersonalRelationLover)})
			data.Trust = append(data.Trust, &structs.Tuple{Subject: mum.ID, Object: affair.ID, Value: dice.relationTrust.Int()})
			data.Trust = append(data.Trust, &structs.Tuple{Subject: affair.ID, Object: mum.ID, Value: dice.relationTrust.Int()})
		} else {
			data.Relations = append(data.Relations, &structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: int(structs.PersonalRelationLover)})
			data.Relations = append(data.Relations, &structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: int(structs.PersonalRelationLover)})
			data.Trust = append(data.Trust, &structs.Tuple{Subject: dad.ID, Object: affair.ID, Value: dice.relationTrust.Int()})
			data.Trust = append(data.Trust, &structs.Tuple{Subject: affair.ID, Object: dad.ID, Value: dice.relationTrust.Int()})
		}
	}

	for i := 0; i < dice.familySize.Int(); i++ {
		child := s.randPerson(rng, demo, dice, areaID)

		if i == 0 && rng.Float64() <= demo.ChildbearingDeathProbability {
			mum.DeathTick = 1
			mum.DeathMetaReason = "died in childbirth"
			break
		}
	}

	if mum.DeathTick == 1 || rng.Float64() <= demo.MarriageDivorceProbability {
		data.Relations = append(data.Relations, &structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationExHusband)})
		data.Relations = append(data.Relations, &structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationExWife)})
		data.Trust = append(data.Trust, &structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / 2})
		data.Trust = append(data.Trust, &structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int() / 2})
	} else {
		data.Relations = append(data.Relations, &structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: int(structs.PersonalRelationHusband)})
		data.Relations = append(data.Relations, &structs.Tuple{Subject: dad.ID, Object: mum.ID, Value: int(structs.PersonalRelationWife)})
		data.Trust = append(data.Trust, &structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int()})
		data.Trust = append(data.Trust, &structs.Tuple{Subject: mum.ID, Object: dad.ID, Value: dice.relationTrust.Int()})
	}

	return data
}

func (s *simulationImpl) Populate(total int, demo *structs.Demographics, areas ...string) error {
	// TODO: support passing a 'Namer' to generate names
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	chunk := 500                      // people to create at a time
	dice := newDemographicsRand(demo) // initialise our dice with probabilities
	totalAlive := 0
	prevPeople := map[string]*structs.Person{} // some saved people to inter-link chunks
	for {
		areaID := areas[rng.Intn(len(areas))]

	}
}

func blacksheep(rng *rand.Rand, p *structs.Person) {
	v := -100 + rng.Intn(5)
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

func newDemographicsRand(demo *structs.Demographics) *demographicsRand {
	// professions
	skills := map[string]*stats.Rand{}
	profOccurProb := []float64{}
	for _, profession := range demo.Professions {
		skills[profession.Name] = stats.NewRand(0, 100, float64(profession.Average), float64(profession.Deviation))
		profOccurProb = append(profOccurProb, profession.Occurs)
	}

	// faiths
	faiths := map[string]*stats.Rand{}
	faithOccurProb := []float64{}
	for _, faith := range demo.Faiths {
		faiths[faith.ReligionID] = stats.NewRand(0, 100, float64(faith.Average), float64(faith.Deviation))
		faithOccurProb = append(faithOccurProb, faith.Occurs)
	}

	return &demographicsRand{
		familySize: stats.NewRand(
			0, float64(demo.FamilySizeMax),
			float64(demo.FamilySizeAverage), float64(demo.FamilySizeDeviation),
		),
		childbearingAge: stats.NewRand(
			float64(demo.ChildbearingAgeMin), float64(demo.ChildbearingAgeMax),
			float64(demo.ChildbearingAgeAverage), float64(demo.ChildbearingAgeDeviation),
		),
		ethosAltruism:   stats.NewRand(-100, 100, float64(demo.EthosAverage.Altruism), float64(demo.EthosDeviation.Altruism)),
		ethosAmbition:   stats.NewRand(-100, 100, float64(demo.EthosAverage.Ambition), float64(demo.EthosDeviation.Ambition)),
		ethosTradition:  stats.NewRand(-100, 100, float64(demo.EthosAverage.Tradition), float64(demo.EthosDeviation.Tradition)),
		ethosPacifism:   stats.NewRand(-100, 100, float64(demo.EthosAverage.Pacifism), float64(demo.EthosDeviation.Pacifism)),
		ethosPiety:      stats.NewRand(-100, 100, float64(demo.EthosAverage.Piety), float64(demo.EthosDeviation.Piety)),
		ethosCaution:    stats.NewRand(-100, 100, float64(demo.EthosAverage.Caution), float64(demo.EthosDeviation.Caution)),
		professionLevel: skills,
		professionOccur: stats.NewNormalised(profOccurProb),
		professionCount: stats.NewNormalised(demo.ProfessionProbability),
		faithLevel:      faiths,
		faithOccur:      stats.NewNormalised(faithOccurProb),
		faithCount:      stats.NewNormalised(demo.FaithProbability),
		relationTrust:   stats.NewRand(-100, 100, 75, 15),
	}
}
