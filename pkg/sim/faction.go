package sim

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

type factionRand struct {
	ethosAltruism  *stats.Rand
	ethosAmbition  *stats.Rand
	ethosTradition *stats.Rand
	ethosPacifism  *stats.Rand
	ethosPiety     *stats.Rand
	ethosCaution   *stats.Rand
	focusOccur     stats.Normalised
	focusCount     stats.Normalised
	guildOccur     stats.Normalised
	guildCount     stats.Normalised
	areaCount      stats.Normalised
	propertyCount  stats.Normalised
}

func newFactionRand(f *config.Faction) *factionRand {
	focusOccurProb := []float64{}
	for _, focus := range f.Focuses {
		focusOccurProb = append(focusOccurProb, focus.Probability)
	}

	guildProb := []float64{}
	for _, guild := range f.Guilds {
		guildProb = append(guildProb, guild.Probability)
	}

	return &factionRand{
		ethosAltruism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Altruism), float64(f.EthosDeviation.Altruism)),
		ethosAmbition:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Ambition), float64(f.EthosDeviation.Ambition)),
		ethosTradition: stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Tradition), float64(f.EthosDeviation.Tradition)),
		ethosPacifism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Pacifism), float64(f.EthosDeviation.Pacifism)),
		ethosPiety:     stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Piety), float64(f.EthosDeviation.Piety)),
		ethosCaution:   stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Caution), float64(f.EthosDeviation.Caution)),
		focusOccur:     stats.NewNormalised(focusOccurProb),
		focusCount:     stats.NewNormalised(f.FocusProbability),
		guildOccur:     stats.NewNormalised(guildProb),
		guildCount:     stats.NewNormalised(f.GuildProbability),
		areaCount:      stats.NewNormalised(f.AreaProbability),
		propertyCount:  stats.NewNormalised(f.PropertyProbability),
	}
}

func (s *simulationImpl) randFaction(fr *factionRand) *structs.Faction {
	eth := structs.Ethos{
		Altruism:  fr.ethosAltruism.Int(),
		Ambition:  fr.ethosAmbition.Int(),
		Tradition: fr.ethosTradition.Int(),
		Pacifism:  fr.ethosPacifism.Int(),
		Piety:     fr.ethosPiety.Int(),
		Caution:   fr.ethosCaution.Int(),
	}

	return &structs.Faction{
		Ethos: structs.Ethos{},
	}
}

func (s *simulationImpl) SpawnGovernment(g *config.Government, f *config.Faction) (*structs.Faction, *structs.Government, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	laws := structs.NewLaws()
	for action, prob := range g.ProbabilityOutlawAction {
		if rng.Float64() < prob {
			laws.Actions[action] = true
		}
	}
	for commodity, prob := range g.ProbabilityOutlawCommodity {
		if rng.Float64() < prob {
			laws.Commodities[commodity] = true
		}
	}

	if g.TaxFrequency.Min < 1 {
		g.TaxFrequency.Min = 1
	}
	if g.TaxRate.Min < 1 {
		g.TaxRate.Min = 1
	}
	if g.TaxRate.Max > 100 {
		g.TaxRate.Max = 100
	}

	tax := stats.NewRand(g.TaxFrequency.Min, g.TaxFrequency.Max, g.TaxFrequency.Mean, g.TaxFrequency.Deviation)
	rate := stats.NewRand(g.TaxRate.Min, g.TaxRate.Max, g.TaxRate.Mean, g.TaxRate.Deviation)

	govt := &structs.Government{
		ID:           structs.NewID(),
		TaxRate:      rate.Float64() / 100,
		TaxFrequency: tax.Int(),
		Outlawed:     laws,
	}

	fr := newFactionRand(f)
}
