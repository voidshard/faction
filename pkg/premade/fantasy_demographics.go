package premade

import (
	"github.com/voidshard/faction/pkg/structs"
)

// DemographicsFantasyHuman returns a set of demographics for a fantasy human
// population.
// In general users are expected to create their own demographics, but this
// should help lazy people get started.
func DemographicsFantasyHuman() *structs.Demographics {
	return &structs.Demographics{
		FamilySizeAverage:            8, // children
		FamilySizeDeviation:          4,
		Race:                         "human",
		ChildbearingAgeMin:           13 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeMax:           45 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeAverage:       22 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeDeviation:     5 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingDeathProbability: 0.02,
		EthosAverage:                 &structs.Ethos{},
		EthosDeviation: &structs.Ethos{
			Altruism:  20,
			Ambition:  100,
			Tradition: 40,
			Pacifism:  90,
			Piety:     90,
			Caution:   15,
		},
		EthosBlackSheepProbability: 0.02, // Probability given at least one radical ethos change
		DeathCauseNaturalProbability: map[string]float64{
			"malaria":       0.09, // natural diseases (50%)
			"pox":           0.08,
			"polio":         0.07,
			"dysentery":     0.07,
			"plague":        0.06,
			"measles":       0.05,
			"typoid":        0.03,
			"scarlet fever": 0.03,
			"flu":           0.01,
			"ergotism":      0.01, // end diseases
			"accidental":    0.15,
			"war":           0.15,
			"starvation":    0.1,
			"suicide":       0.02,
			"executed":      0.02,
			"assassination": 0.01,
			"poisoning":     0.01,
			"eaten":         0.04,
		},
	}
}
