package base

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/random/rng"
	"github.com/voidshard/faction/internal/sim/simutil"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/technology"
)

// factionRand is a helper struct to generate random factions
// using dice / distributions provided in configs
//
// Nb. probably this should be moved to internal/random/faction
// .. most likely after splitting out the dice / distributions from
// the bits that need db data (yields, governments).
//
// Probably it's too annoying to do nicely to be worth it?
type factionRand struct {
	yieldRand      *yieldRand
	cfg            *config.Faction
	tech           technology.Technology
	rng            *rand.Rand
	ethosAltruism  *rng.Rand
	ethosAmbition  *rng.Rand
	ethosTradition *rng.Rand
	ethosPacifism  *rng.Rand
	ethosPiety     *rng.Rand
	ethosCaution   *rng.Rand
	leaderOccur    rng.Normalised
	leaderList     []structs.LeaderType
	structOccur    rng.Normalised
	structList     []structs.LeaderStructure
	wealth         *rng.Rand
	cohesion       *rng.Rand
	corruption     *rng.Rand
	espOffense     *rng.Rand
	espDefense     *rng.Rand
	milOffense     *rng.Rand
	milDefense     *rng.Rand
	focusOccur     rng.Normalised
	focusCount     rng.Normalised
	focusWeights   []*rng.Rand
	guildOccur     rng.Normalised
	guildCount     rng.Normalised
	propertyCount  rng.Normalised
	areas          []string
	plotSize       *rng.Rand
}

// randResearch returns a list of research topic weights for a faction.
func (fr *factionRand) randResearch(f *simutil.MetaFaction, count int, requiredTopics []string) []*structs.Tuple {
	if count <= 0 {
		return []*structs.Tuple{}
	}

	weight := rng.NewRand(
		// TODO: could expose this in config
		structs.MaxEthos/8,
		structs.MaxEthos,
		structs.MaxEthos/2,
		structs.MaxEthos,
	)
	weights := []*structs.Tuple{}

	excludeTopics := map[string]bool{}
	for _, topic := range requiredTopics {
		_, ok := excludeTopics[topic]
		if ok {
			continue
		}
		excludeTopics[topic] = true
		weights = append(weights, &structs.Tuple{
			Subject: f.Faction.ID,
			Object:  topic,
			Value:   weight.Int(),
		})
	}
	if len(weights) >= count {
		// we don't need any more topics
		return weights
	}

	areas := []string{} // all areas faction has presence in
	for area := range f.Areas {
		areas = append(areas, area)
	}

	topics := fr.tech.Topics(areas...) // all possible topics

	probs := []float64{}
	byProfession := map[string][]string{}
	for _, topic := range topics {
		prob, ok := fr.cfg.ResearchProbability[topic.Name]
		if !ok {
			prob = 0.0
		}

		_, ok = excludeTopics[topic.Name] // since this topic is done already
		if ok {
			prob = 0.0
		}

		probs = append(probs, prob)

		ptopics, ok := byProfession[topic.Profession]
		if !ok {
			ptopics = []string{}
		}
		byProfession[topic.Profession] = append(ptopics, topic.Name)
	}

	favouredTopics := map[string]bool{}
	for _, weight := range f.ProfWeights {
		//  consider favoured professions for the faction
		ptopics, ok := byProfession[weight.Object]
		if !ok { // profession has no research topics
			continue
		}

		// professions that have matching research topic(s)
		for _, topic := range ptopics {
			favouredTopics[topic] = true
		}
	}
	if len(favouredTopics) > 0 {
		// set non-favoured topics to 0 probability
		for i, topic := range topics {
			safe, _ := favouredTopics[topic.Name]
			if !safe {
				probs[i] = 0.0
			}
		}
	}

	// finally, we can choose actual research topics
	norm := rng.NewNormalised(probs)
	seen := map[int]bool{}

	for i := 0; i < count; i++ {
		choice := norm.Int()
		if choice < 0 {
			break
		}

		_, ok := seen[choice]
		if ok {
			continue
		}
		seen[choice] = true

		topic := topics[choice]
		weights = append(weights, &structs.Tuple{
			Subject: f.Faction.ID,
			Object:  topic.Name,
			Value:   weight.Int(),
		})
	}

	return weights
}

func (fr *factionRand) randLandForGuild(g *config.Guild) []*structs.Plot {
	options, ok := fr.yieldRand.professionLand[g.Profession]
	if !ok {
		return nil
	}

	index := 0
	total := 0
	for _, l := range options {
		index++
		total += l.Yield
		if total >= g.MinYield {
			break
		}
	}

	if total < g.MinYield {
		return nil
	}

	defer fr.recalcGuildProb()

	if index >= len(options) {
		// we've assigned all available land
		delete(fr.yieldRand.professionLand, g.Profession)
		delete(fr.yieldRand.professionYield, g.Profession)
		return options
	}

	// we've assigned some subset of available land
	fr.yieldRand.professionLand[g.Profession] = options[index+1:]
	ftotal := fr.yieldRand.professionYield[g.Profession]
	fr.yieldRand.professionYield[g.Profession] = ftotal - total

	return options[:index+1]
}

func (fr *factionRand) recalcGuildProb() {
	// as we hand out land, it becomes increasingly unlikely that some commodities
	// form the backbone of a(nother) guild.
	guildProb := []float64{}
	if fr.yieldRand != nil {
		for _, guild := range fr.cfg.Guilds {
			total, _ := fr.yieldRand.professionYield[guild.Profession]
			prob := 0.0
			if total >= guild.MinYield {
				prob = guild.Probability
			}
			guildProb = append(guildProb, prob)
		}
	}
	fr.guildOccur = rng.NewNormalised(guildProb)
}

// newFactionRand creates a new dice roller for creationg factions based on faction config
// and the available land rights in some area(s).
func newFactionRand(f *config.Faction, tech technology.Technology, yields *yieldRand, areas []string) *factionRand {
	focusOccurProb := []float64{}
	focusWeights := []*rng.Rand{}
	for _, focus := range f.Focuses {
		focusOccurProb = append(focusOccurProb, focus.Probability)
		focusWeights = append(focusWeights, rng.NewRand(focus.Weight.Min, focus.Weight.Max, focus.Weight.Mean, focus.Weight.Deviation))
	}

	leaderProb := []float64{}
	llist := []structs.LeaderType{}
	for leader, prob := range f.LeadershipProbability {
		leaderProb = append(leaderProb, prob)
		llist = append(llist, leader)
	}

	structureProb := []float64{}
	slist := []structs.LeaderStructure{}
	for structure, prob := range f.LeadershipStructureProbability {
		structureProb = append(structureProb, prob)
		slist = append(slist, structure)
	}

	fr := &factionRand{
		yieldRand:      yields,
		cfg:            f,
		tech:           tech,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		ethosAltruism:  rng.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Altruism), float64(f.EthosDeviation.Altruism)),
		ethosAmbition:  rng.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Ambition), float64(f.EthosDeviation.Ambition)),
		ethosTradition: rng.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Tradition), float64(f.EthosDeviation.Tradition)),
		ethosPacifism:  rng.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Pacifism), float64(f.EthosDeviation.Pacifism)),
		ethosPiety:     rng.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Piety), float64(f.EthosDeviation.Piety)),
		ethosCaution:   rng.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Caution), float64(f.EthosDeviation.Caution)),
		leaderOccur:    rng.NewNormalised(leaderProb),
		leaderList:     llist,
		structOccur:    rng.NewNormalised(structureProb),
		structList:     slist,
		wealth:         rng.NewRand(f.Wealth.Min, f.Wealth.Max, f.Wealth.Mean, f.Wealth.Deviation),
		cohesion:       rng.NewRand(f.Cohesion.Min, f.Cohesion.Max, f.Cohesion.Mean, f.Cohesion.Deviation),
		corruption:     rng.NewRand(f.Corruption.Min, f.Corruption.Max, f.Corruption.Mean, f.Corruption.Deviation),
		espOffense:     rng.NewRand(f.EspionageOffense.Min, f.EspionageOffense.Max, f.EspionageOffense.Mean, f.EspionageOffense.Deviation),
		espDefense:     rng.NewRand(f.EspionageDefense.Min, f.EspionageDefense.Max, f.EspionageDefense.Mean, f.EspionageDefense.Deviation),
		milOffense:     rng.NewRand(f.MilitaryOffense.Min, f.MilitaryOffense.Max, f.MilitaryOffense.Mean, f.MilitaryOffense.Deviation),
		milDefense:     rng.NewRand(f.MilitaryDefense.Min, f.MilitaryDefense.Max, f.MilitaryDefense.Mean, f.MilitaryDefense.Deviation),
		focusOccur:     rng.NewNormalised(focusOccurProb),
		focusCount:     rng.NewNormalised(f.FocusProbability),
		focusWeights:   focusWeights,
		guildCount:     rng.NewNormalised(f.GuildProbability),
		propertyCount:  rng.NewNormalised(f.PropertyProbability),
		areas:          areas,
		plotSize:       rng.NewRand(f.PlotSize.Min, f.PlotSize.Max, f.PlotSize.Mean, f.PlotSize.Deviation),
	}

	fr.recalcGuildProb()
	return fr
}
