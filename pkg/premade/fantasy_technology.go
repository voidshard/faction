package premade

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/structs"
)

type FantasyTechnology struct {
	rng       *rand.Rand
	allTopics map[string]*structs.ResearchTopic
}

func (f *FantasyTechnology) Topics() []*structs.ResearchTopic {
	ret := []*structs.ResearchTopic{}
	for _, t := range f.allTopics {
		ret = append(ret, t)
	}
	return ret
}

func (f *FantasyTechnology) Topic(name string) *structs.ResearchTopic {
	t, _ := f.allTopics[name]
	return t
}

func (f *FantasyTechnology) Research(jobID string, currentLevel int) int {
	// Nb. in an actual implementation we probably want to
	// - look at who is performing the research
	// - decide on how much research is done / achieved
	// - add side effects based on research & the level
	// Ie. research into dark magic might have a chance of causing some
	// calamity, or research into a new technology might have a chance
	// of adding new commodities to our economy.
	//
	// For our premade fantasy technology, we'll just rolls some dice ..

	v := f.rng.Float64()
	if v <= 0.01 {
		// some sort of break through
		return currentLevel + f.rng.Intn(100) + f.rng.Intn(100)
	} else if v <= 0.02 {
		// some sort of disaster
		return currentLevel - f.rng.Intn(20)
	}
	return currentLevel + f.rng.Intn(10) + f.rng.Intn(10)
}

func NewFantasyTechnology() *FantasyTechnology {
	return &FantasyTechnology{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		allTopics: map[string]*structs.ResearchTopic{
			// Nb. these aren't all topics in constants.go just the ones
			// that seem to fit the setting
			AGRICULTURE: &structs.ResearchTopic{
				Name:        AGRICULTURE,
				Probability: 0.01,
				Profession:  SCHOLAR,
				Ethos:       structs.Ethos{},
			},
			ASTRONOMY: &structs.ResearchTopic{
				Name:        ASTRONOMY,
				Probability: 0.20,
				Profession:  SCHOLAR,
				Ethos:       structs.Ethos{Piety: structs.MaxEthos / 2},
			},
			WARFARE: &structs.ResearchTopic{
				Name:        WARFARE,
				Probability: 0.05,
				Profession:  SOLDIER,
				Ethos:       structs.Ethos{Pacifism: structs.MinEthos / 2},
			},
			METALLURGY: &structs.ResearchTopic{
				Name:        METALLURGY,
				Probability: 0.04,
				Profession:  SMITH,
				Ethos:       structs.Ethos{},
			},
			PHILOSOPHY: &structs.ResearchTopic{
				Name:        PHILOSOPHY,
				Probability: 0.10,
				Profession:  SCHOLAR,
				Ethos:       structs.Ethos{Tradition: structs.MaxEthos / 2},
			},
			MEDICINE: &structs.ResearchTopic{
				Name:        MEDICINE,
				Probability: 0.10,
				Profession:  PRIEST, // early temples doubled as healing springs (eg Asclepius)
				Ethos:       structs.Ethos{Altruism: structs.MaxEthos / 2},
			},
			MATHEMATICS: &structs.ResearchTopic{
				Name:        MATHEMATICS,
				Probability: 0.08,
				Profession:  SCHOLAR,
				Ethos:       structs.Ethos{},
			},
			LITERATURE: &structs.ResearchTopic{
				Name:        LITERATURE,
				Probability: 0.02,
				Profession:  SCHOLAR,
				Ethos:       structs.Ethos{Tradition: structs.MaxEthos / 2},
			},
			LAW: &structs.ResearchTopic{
				Name:        LAW,
				Probability: 0.10,
				Profession:  SCHOLAR, // training in law / policy / speaking wasn't uncommon
				Ethos:       structs.Ethos{Tradition: structs.MaxEthos / 4, Ambition: structs.MaxEthos / 4},
			},
			ARCHITECTURE: &structs.ResearchTopic{
				Name:        ARCHITECTURE,
				Probability: 0.04,
				Profession:  SCHOLAR, // we know of ancient schools of architecture so .. why not?
				Ethos:       structs.Ethos{Tradition: structs.MaxEthos / 4, Piety: structs.MaxEthos / 4},
			},
			THEOLOGY: &structs.ResearchTopic{
				Name:        THEOLOGY,
				Probability: 0.06,
				Profession:  PRIEST,
				Ethos:       structs.Ethos{Piety: structs.MaxEthos / 2},
			},
			MAGIC_ARCANA: &structs.ResearchTopic{
				Name:        MAGIC_ARCANA,
				Probability: 0.15,
				Profession:  MAGE,
				Ethos:       structs.Ethos{Tradition: structs.MaxEthos / 4, Ambition: structs.MaxEthos / 2},
			},
			MAGIC_OCCULT: &structs.ResearchTopic{
				Name:        MAGIC_OCCULT,
				Probability: 0.04,
				Profession:  MAGE, // mage / priest
				Ethos:       structs.Ethos{Piety: structs.MinEthos / 4, Ambition: structs.MaxEthos / 2},
			},
			ALCHEMY: &structs.ResearchTopic{
				Name:        ALCHEMY,
				Probability: 0.05,
				Profession:  ALCHEMIST, // early practiioners were into medicine & creating gold
				Ethos:       structs.Ethos{Tradition: structs.MaxEthos / 4, Altruism: structs.MaxEthos / 4},
			},
		},
	}
}
