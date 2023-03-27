/*
random_faction.go - random faction / government generation
*/
package sim

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/db"
	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// factionRand is a helper struct to generate random factions
// using dice / distributions provided in configs
type factionRand struct {
	cfg            *config.Faction
	ethosAltruism  *stats.Rand
	ethosAmbition  *stats.Rand
	ethosTradition *stats.Rand
	ethosPacifism  *stats.Rand
	ethosPiety     *stats.Rand
	ethosCaution   *stats.Rand
	leaderOccur    stats.Normalised
	leaderList     []structs.LeaderType
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
	areaCount      stats.Normalised
	propertyCount  stats.Normalised
}

// faction + associated metadata used during creation
type metaFaction struct {
	faction       *structs.Faction
	actionWeights []*structs.Tuple

	professions []string
	focuses     []*config.Action
}

func newFactionRand(f *config.Faction) *factionRand {
	focusOccurProb := []float64{}
	focusWeights := []*stats.Rand{}
	for _, focus := range f.Focuses {
		focusOccurProb = append(focusOccurProb, focus.Probability)
		focusWeights = append(focusWeights, stats.NewRand(focus.Weight.Min, focus.Weight.Max, focus.Weight.Mean, focus.Weight.Deviation))
	}

	guildProb := []float64{}
	for _, guild := range f.Guilds {
		guildProb = append(guildProb, guild.Probability)
	}

	leaderProb := []float64{}
	llist := []structs.LeaderType{}
	for leader, prob := range f.LeadershipProbability {
		leaderProb = append(leaderProb, prob)
		llist = append(llist, leader)
	}

	return &factionRand{
		cfg:            f,
		ethosAltruism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Altruism), float64(f.EthosDeviation.Altruism)),
		ethosAmbition:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Ambition), float64(f.EthosDeviation.Ambition)),
		ethosTradition: stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Tradition), float64(f.EthosDeviation.Tradition)),
		ethosPacifism:  stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Pacifism), float64(f.EthosDeviation.Pacifism)),
		ethosPiety:     stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Piety), float64(f.EthosDeviation.Piety)),
		ethosCaution:   stats.NewRand(structs.MinEthos, structs.MaxEthos, float64(f.EthosMean.Caution), float64(f.EthosDeviation.Caution)),
		leaderOccur:    stats.NewNormalised(leaderProb),
		leaderList:     llist,
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
		guildOccur:     stats.NewNormalised(guildProb),
		guildCount:     stats.NewNormalised(f.GuildProbability),
		areaCount:      stats.NewNormalised(f.AreaProbability),
		propertyCount:  stats.NewNormalised(f.PropertyProbability),
	}
}

// randFaction spits out a random faction in isolation of the rest of the world.
// For determining guilds (rather than simple action focuses) we'll need a more complex
// inter-faction + world consideration (in a larger func), this one just gets us started.
func (s *simulationImpl) randFaction(fr *factionRand) *metaFaction {
	// start with a lot of randomly inserted fields
	mf := &metaFaction{
		faction: &structs.Faction{
			Ethos: structs.Ethos{ // random ethos to start with
				Altruism:  fr.ethosAltruism.Int(),
				Ambition:  fr.ethosAmbition.Int(),
				Tradition: fr.ethosTradition.Int(),
				Pacifism:  fr.ethosPacifism.Int(),
				Piety:     fr.ethosPiety.Int(),
				Caution:   fr.ethosCaution.Int(),
			},
			ID:               structs.NewID(),
			Leadership:       fr.leaderList[fr.leaderOccur.Int()],
			Wealth:           fr.wealth.Int(),
			Cohesion:         fr.cohesion.Int(),
			Corruption:       fr.corruption.Int(),
			EspionageOffense: fr.espOffense.Int(),
			EspionageDefense: fr.espDefense.Int(),
			MilitaryOffense:  fr.milOffense.Int(),
			MilitaryDefense:  fr.milDefense.Int(),
		},
		professions:   []string{},
		focuses:       []*config.Action{},
		actionWeights: []*structs.Tuple{},
	}

	// consider action focuses
	seen := map[int]bool{}
	actionsEthos := []*structs.Ethos{}
	for i := 0; i < fr.focusCount.Int(); i++ {
		choice := fr.focusOccur.Int()
		_, ok := seen[choice]
		if ok {
			continue
		}
		seen[choice] = true

		focus := fr.cfg.Focuses[choice]
		weight := fr.focusWeights[choice]

		for _, act := range focus.Actions {
			mf.actionWeights = append(mf.actionWeights, &structs.Tuple{
				Subject: mf.faction.ID,
				Object:  string(act),
				Value:   weight.Int(),
			})

			actionCfg, ok := s.cfg.Actions[act]
			if !ok {
				continue
			}

			mf.focuses = append(mf.focuses, actionCfg)
			actionsEthos = append(actionsEthos, &actionCfg.Ethos)
			for prof := range actionCfg.ProfessionWeights {
				mf.professions = append(mf.professions, prof)
			}
		}
	}

	// favoured actions influence faction's starting ethos
	mf.faction.Ethos = *structs.EthosAverageNonZero(actionsEthos...)

	return mf
}

func (s *simulationImpl) SpawnGovernment(g *config.Government, f *config.Faction, areas ...string) (*structs.Faction, *structs.Government, error) {
	// roll some dice
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
	if g.TaxRate.Min < 0 {
		g.TaxRate.Min = 0
	}
	if g.TaxRate.Max > 100 {
		g.TaxRate.Max = 100
	}

	tax := stats.NewRand(g.TaxFrequency.Min, g.TaxFrequency.Max, g.TaxFrequency.Mean, g.TaxFrequency.Deviation)
	rate := stats.NewRand(g.TaxRate.Min, g.TaxRate.Max, g.TaxRate.Mean, g.TaxRate.Deviation)
	fr := newFactionRand(f)
	fact := s.randFaction(fr)

	govt := &structs.Government{
		ID:           structs.NewID(),
		TaxRate:      rate.Float64() / 100,
		TaxFrequency: tax.Int(),
		Outlawed:     laws,
	}

	// apply buffs
	fact.faction.MilitaryDefense += int(float64(fact.faction.MilitaryDefense) * g.MilitaryDefenseBonus)
	fact.faction.MilitaryOffense += int(float64(fact.faction.MilitaryOffense) * g.MilitaryOffenseBonus)
	fact.faction.EspionageDefense += int(float64(fact.faction.EspionageDefense) * g.EspionageDefenseBonus)
	fact.faction.EspionageOffense += int(float64(fact.faction.EspionageOffense) * g.EspionageOffenseBonus)

	// set some government specific fields
	fact.faction.IsGovernment = true
	fact.faction.GovernmentID = govt.ID

	// apply government action weights
	weight := stats.NewRand(g.ActionWeight.Min, g.ActionWeight.Max, g.ActionWeight.Mean, g.ActionWeight.Deviation)
	fact.actionWeights = append(
		fact.actionWeights,
		&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeRevokeLand), Value: weight.Int()},
		&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeGrantLand), Value: weight.Int()},
	)

	// add some uh, flavor action weights based on government ethos.
	if fact.faction.Ambition <= structs.MinEthos/4 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypePropoganda), Value: weight.Int()},
		)
	}
	if fact.faction.Altruism >= structs.MaxEthos/2 || fact.faction.Piety >= structs.MaxEthos/4 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeFestival), Value: weight.Int()},
		)
	} else if fact.faction.Altruism >= structs.MaxEthos/4 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeCharity), Value: weight.Int()},
		)
	}
	if fact.faction.Pacifism <= structs.MinEthos/2 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeShadowWar), Value: weight.Int()},
		)
	} else if fact.faction.Pacifism <= structs.MinEthos/4 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeRaid), Value: weight.Int()},
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeAssassinate), Value: weight.Int()},
		)
	}
	if fact.faction.Caution >= structs.MaxEthos/4 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeGatherSecrets), Value: weight.Int()},
		)
	}
	if fact.faction.Tradition <= structs.MinEthos/2 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeFrame), Value: weight.Int()},
		)
	} else if fact.faction.Tradition <= structs.MinEthos/4 {
		fact.actionWeights = append(
			fact.actionWeights,
			&structs.Tuple{Subject: fact.faction.ID, Object: string(structs.ActionTypeSpreadRumors), Value: weight.Int()},
		)
	}

	// whack everything in
	err := s.dbconn.InTransaction(func(tx db.ReaderWriter) error {
		err := tx.SetGovernments(govt)
		if err != nil {
			return err
		}
		err = tx.SetTuples(db.RelationFactionActionTypeWeight, fact.actionWeights...)
		if err != nil {
			return err
		}
		return tx.SetFactions(fact.faction)
	})

	return fact.faction, govt, err
}
