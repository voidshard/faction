/*
random_government.go - random government generation
*/

package base

import (
	"math"
	"math/rand"
	"time"

	stats "github.com/voidshard/faction/internal/random/rng"
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// factionRand is a helper struct to generate random factions
// using dice / distributions provided in configs

func (s *Base) SpawnGovernment(g *config.Government) (*structs.Government, error) {
	if g.TaxRate.Min < 0 {
		g.TaxRate.Min = 0
	}
	if g.TaxRate.Max > 100 {
		g.TaxRate.Max = 100
	}
	if g.TaxFrequency.Min < 1 {
		g.TaxFrequency.Min = 1
	}

	rate := stats.NewRand(g.TaxRate.Min, g.TaxRate.Max, g.TaxRate.Mean, g.TaxRate.Deviation).Float64()

	govt := &structs.Government{
		ID:           structs.NewID(),
		TaxRate:      math.Round(rate) / 100,
		TaxFrequency: stats.NewRand(g.TaxFrequency.Min, g.TaxFrequency.Max, g.TaxFrequency.Mean, g.TaxFrequency.Deviation).Int(),
		Outlawed:     structs.NewLaws(),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for act, prob := range g.ProbabilityOutlawAction {
		if rng.Float64() <= prob {
			govt.Outlawed.Actions[act] = true
		}
	}
	for item, prob := range g.ProbabilityOutlawCommodity {
		if rng.Float64() <= prob {
			govt.Outlawed.Commodities[item] = true
		}
	}
	for topic, prob := range g.ProbabilityOutlawResearch {
		if rng.Float64() <= prob {
			govt.Outlawed.Research[topic] = true
		}
	}
	for religion, prob := range g.ProbabilityOutlawReligion {
		if rng.Float64() <= prob {
			govt.Outlawed.Religions[religion] = true
		}
	}

	return govt, s.dbconn.SetGovernments(govt)
}
