package base

import (
	"math/rand"
	"time"

	stats "github.com/voidshard/faction/internal/random/rng"
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
	ethosAltruism  *stats.Rand
	ethosAmbition  *stats.Rand
	ethosTradition *stats.Rand
	ethosPacifism  *stats.Rand
	ethosPiety     *stats.Rand
	ethosCaution   *stats.Rand
	leaderOccur    stats.Normalised
	leaderList     []structs.LeaderType
	structOccur    stats.Normalised
	structList     []structs.LeaderStructure
	wealth         *stats.Rand
	cohesion       *stats.Rand
	corruption     *stats.Rand
	espOffense     *stats.Rand
	espDefense     *stats.Rand
	milOffense     *stats.Rand
	milDefense     *stats.Rand
	focusOccur     stats.Normalised
	focusCount     stats.Normalised
	focusWeights   []*stats.Rand
	guildOccur     stats.Normalised
	guildCount     stats.Normalised
	propertyCount  stats.Normalised
	areas          []string
	plotSize       *stats.Rand
}

func (fr *factionRand) randResearch(f *metaFaction, count int) []*structs.Tuple {
	if count <= 0 {
		return []*structs.Tuple{}
	}

	areas := []string{} // all areas faction has presence in
	for area := range f.areas {
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
		probs = append(probs, prob)

		ptopics, ok := byProfession[topic.Profession]
		if !ok {
			ptopics = []string{}
		}
		byProfession[topic.Profession] = append(ptopics, topic.Name)
	}

	favouredTopics := map[string]bool{}
	for _, weight := range f.profWeights {
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
	norm := stats.NewNormalised(probs)
	seen := map[int]bool{}
	weight := stats.NewRand(
		// TODO: could expose this in config
		structs.MaxEthos/8,
		structs.MaxEthos,
		structs.MaxEthos/2,
		structs.MaxEthos,
	)

	weights := []*structs.Tuple{}
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
			Subject: f.faction.ID,
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
	fr.guildOccur = stats.NewNormalised(guildProb)
}

// newFactionRand creates a new dice roller for creationg factions based on faction config
// and the available land rights in some area(s).
func newFactionRand(f *config.Faction, tech technology.Technology, yields *yieldRand, areas []string) *factionRand {
	focusOccurProb := []float64{}
	focusWeights := []*stats.Rand{}
	for _, focus := range f.Focuses {
		focusOccurProb = append(focusOccurProb, focus.Probability)
		focusWeights = append(focusWeights, stats.NewRand(focus.Weight.Min, focus.Weight.Max, focus.Weight.Mean, focus.Weight.Deviation))
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
		ethosAltruism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Altruism), float64(f.EthosDeviation.Altruism)),
		ethosAmbition:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Ambition), float64(f.EthosDeviation.Ambition)),
		ethosTradition: stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Tradition), float64(f.EthosDeviation.Tradition)),
		ethosPacifism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Pacifism), float64(f.EthosDeviation.Pacifism)),
		ethosPiety:     stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Piety), float64(f.EthosDeviation.Piety)),
		ethosCaution:   stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Caution), float64(f.EthosDeviation.Caution)),
		leaderOccur:    stats.NewNormalised(leaderProb),
		leaderList:     llist,
		structOccur:    stats.NewNormalised(structureProb),
		structList:     slist,
		wealth:         stats.NewRand(f.Wealth.Min, f.Wealth.Max, f.Wealth.Mean, f.Wealth.Deviation),
		cohesion:       stats.NewRand(f.Cohesion.Min, f.Cohesion.Max, f.Cohesion.Mean, f.Cohesion.Deviation),
		corruption:     stats.NewRand(f.Corruption.Min, f.Corruption.Max, f.Corruption.Mean, f.Corruption.Deviation),
		espOffense:     stats.NewRand(f.EspionageOffense.Min, f.EspionageOffense.Max, f.EspionageOffense.Mean, f.EspionageOffense.Deviation),
		espDefense:     stats.NewRand(f.EspionageDefense.Min, f.EspionageDefense.Max, f.EspionageDefense.Mean, f.EspionageDefense.Deviation),
		milOffense:     stats.NewRand(f.MilitaryOffense.Min, f.MilitaryOffense.Max, f.MilitaryOffense.Mean, f.MilitaryOffense.Deviation),
		milDefense:     stats.NewRand(f.MilitaryDefense.Min, f.MilitaryDefense.Max, f.MilitaryDefense.Mean, f.MilitaryDefense.Deviation),
		focusOccur:     stats.NewNormalised(focusOccurProb),
		focusCount:     stats.NewNormalised(f.FocusProbability),
		focusWeights:   focusWeights,
		guildCount:     stats.NewNormalised(f.GuildProbability),
		propertyCount:  stats.NewNormalised(f.PropertyProbability),
		areas:          areas,
		plotSize:       stats.NewRand(f.PlotSize.Min, f.PlotSize.Max, f.PlotSize.Mean, f.PlotSize.Deviation),
	}

	fr.recalcGuildProb()
	return fr
}
