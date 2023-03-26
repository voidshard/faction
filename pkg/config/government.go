package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Government is the configuration for randomly creating a government.
type Government struct {
	// Probability that the given action will be outlawed
	ProbabilityOutlawAction map[structs.ActionType]float64

	// Probability that the given commodity will be outlawed
	ProbabilityOutlawCommodity map[string]float64

	// How often (in ticks) the government will collect taxes
	// Min: 1
	TaxFrequency Distribution

	// Rate (converted to a %, so this should be 0-100)
	// Min: 1
	// Max: 100
	TaxRate Distribution
}
