package base

import (
	"math/rand"
	"time"

	"github.com/voidshard/faction/internal/stats"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// demographicsRand holds random dice for a demographics struct.
//
// We need quite a few with various average / deviation values so
// this helps keep things tidy.
//
// Ok well, tidy-er.
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
	lifespan         *stats.Rand
	relationTrust    *stats.Rand
	friendshipsProb  stats.Normalised
}

// spawnTickes returns some number of ticks we can assume have passed.
//
// When we create people / families from thin air we assume a certain amount of ticks have
// passed before this point (you know so they can like, exist and stuff).
func (d *demographicsRand) spawnTicks() float64 {
	return d.cfg.ChildbearingAge.Min + d.cfg.ChildbearingTerm.Min
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
		lifespan: stats.NewRand(
			demo.Lifespan.Min, demo.Lifespan.Max,
			demo.Lifespan.Mean, demo.Lifespan.Deviation,
		),
		friendshipsProb: stats.NewNormalised([]float64{
			demo.FriendshipCloseProbability,
			demo.FriendshipProbability,
			demo.EnemyProbability,
			demo.EnemyHatedProbability,
		}),
	}
}
