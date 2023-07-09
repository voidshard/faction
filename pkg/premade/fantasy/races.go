package premade

import (
	"github.com/voidshard/faction/pkg/config"
)

// DemographicsHuman returns a set of demographics for a fantasy human
// population.
//
// In general users are expected to create their own demographics, but this
// should help lazy people get started / test things out.
//
// Seriously please load these from a config or something .. this is just for
// playing around.
func RaceHuman() *config.Race {
	return &config.Race{
		// Probabilities of inter-personal relations per tick
		Biomes: map[string]float64{ // biome -> preference weight
			WOODLAND:         0.8,
			FOREST:           0.4,
			JUNGLE:           0.2,
			MOUNTAINS:        0.05,
			HILLS:            0.5,
			SAVANNAH:         0.8,
			GRASSLAND:        0.7,
			DESERT:           0.05,
			LAKE:             0.8,
			SWAMP:            0.1,
			MARSH:            0.1,
			COAST:            0.8,
			SUBARCTIC:        0.01,
			VOLCANIC_FERTILE: 0.6,
		},
		ChildbearingAgeMin: 13 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeMax: 45 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingTerm: config.Distribution{ // 9 months-ish
			Min:       8 * 30 * DEFAULT_TICKS_PER_DAY,
			Max:       10 * 30 * DEFAULT_TICKS_PER_DAY,
			Mean:      9 * 30 * DEFAULT_TICKS_PER_DAY,
			Deviation: 30 * DEFAULT_TICKS_PER_DAY,
		},
		ChildbearingDeathProbability:    0.10, // actual medieval rate probably higher (on Birth)
		DeathInfantMortalityProbability: 0.2,  // Probability of death in childbirth
		DeathAdultMortalityProbability:  0.02 / float64(DEFAULT_TICKS_PER_YEAR),
		Lifespan: config.Distribution{
			// Nb. this is fantasy, where "magic" to grant healing is not unheard of.
			// Medieval life expectancy is considerably more dire
			Min:       2 * DEFAULT_TICKS_PER_YEAR,
			Max:       90 * DEFAULT_TICKS_PER_YEAR,
			Mean:      55 * DEFAULT_TICKS_PER_YEAR,
			Deviation: 15 * DEFAULT_TICKS_PER_YEAR,
		},
	}
}
