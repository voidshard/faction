package fantasy

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/pkg/structs"
)

type Technology struct {
	rng       *rand.Rand
	allTopics map[string]*structs.ResearchTopic
	topicList []*structs.ResearchTopic
}

func (f *Technology) Topics(areas ...string) []*structs.ResearchTopic {
	if f.topicList != nil {
		return f.topicList
	}

	f.topicList = []*structs.ResearchTopic{}
	for _, t := range f.allTopics {
		f.topicList = append(f.topicList, t)
	}

	return f.topicList
}

func (f *Technology) Topic(name string) *structs.ResearchTopic {
	t, _ := f.allTopics[name]
	return t
}

func NewTechnology() *Technology {
	return &Technology{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		allTopics: map[string]*structs.ResearchTopic{
			// Nb. these aren't all topics in constants.go just the ones
			// that seem to fit the setting
			AGRICULTURE: &structs.ResearchTopic{
				Name:       AGRICULTURE,
				Profession: SCHOLAR,
				Ethos:      structs.Ethos{},
			},
			ASTRONOMY: &structs.ResearchTopic{
				Name:       ASTRONOMY,
				Profession: SCHOLAR,
				Ethos:      structs.Ethos{Piety: structs.MaxEthos / 2},
			},
			WARFARE: &structs.ResearchTopic{
				Name:       WARFARE,
				Profession: SOLDIER,
				Ethos:      structs.Ethos{Pacifism: structs.MinEthos / 2},
			},
			METALLURGY: &structs.ResearchTopic{
				Name:       METALLURGY,
				Profession: SMITH,
				Ethos:      structs.Ethos{},
			},
			PHILOSOPHY: &structs.ResearchTopic{
				Name:       PHILOSOPHY,
				Profession: SCHOLAR,
				Ethos:      structs.Ethos{Tradition: structs.MaxEthos / 2},
			},
			MEDICINE: &structs.ResearchTopic{
				Name:       MEDICINE,
				Profession: PRIEST, // early temples doubled as healing springs (eg Asclepius)
				Ethos:      structs.Ethos{Altruism: structs.MaxEthos / 2},
			},
			MATHEMATICS: &structs.ResearchTopic{
				Name:       MATHEMATICS,
				Profession: SCHOLAR,
				Ethos:      structs.Ethos{},
			},
			LITERATURE: &structs.ResearchTopic{
				Name:       LITERATURE,
				Profession: SCHOLAR,
				Ethos:      structs.Ethos{Tradition: structs.MaxEthos / 2},
			},
			LAW: &structs.ResearchTopic{
				Name:       LAW,
				Profession: SCHOLAR, // training in law / policy / speaking wasn't uncommon
				Ethos:      structs.Ethos{Tradition: structs.MaxEthos / 4, Ambition: structs.MaxEthos / 4},
			},
			ARCHITECTURE: &structs.ResearchTopic{
				Name:       ARCHITECTURE,
				Profession: SCHOLAR, // we know of ancient schools of architecture so .. why not?
				Ethos:      structs.Ethos{Tradition: structs.MaxEthos / 4, Piety: structs.MaxEthos / 4},
			},
			THEOLOGY: &structs.ResearchTopic{
				Name:       THEOLOGY,
				Profession: PRIEST,
				Ethos:      structs.Ethos{Piety: structs.MaxEthos / 2},
			},
			MAGIC_ARCANA: &structs.ResearchTopic{
				Name:       MAGIC_ARCANA,
				Profession: MAGE,
				Ethos:      structs.Ethos{Tradition: structs.MaxEthos / 4, Ambition: structs.MaxEthos / 2},
			},
			MAGIC_OCCULT: &structs.ResearchTopic{
				Name:       MAGIC_OCCULT,
				Profession: MAGE, // mage / priest
				Ethos:      structs.Ethos{Piety: structs.MinEthos / 4, Ambition: structs.MaxEthos / 2},
			},
			ALCHEMY: &structs.ResearchTopic{
				Name:       ALCHEMY,
				Profession: ALCHEMIST, // early practiioners were into medicine & creating gold
				Ethos:      structs.Ethos{Tradition: structs.MaxEthos / 4, Altruism: structs.MaxEthos / 4},
			},
		},
	}
}
