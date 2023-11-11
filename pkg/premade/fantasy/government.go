package fantasy

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// Government returns a fantasy government configuration.
//
// Well, for a human type setup.
// We could have a much more extreme setup for Drow or something.
func Government() *config.Government {
	return &config.Government{
		ProbabilityOutlawAction: map[string]float64{
			Bribe:           0.6,
			Propoganda:      0.1,
			Research:        0.01,
			Excommunicate:   0.02,
			ConcealSecrets:  0.05,
			GatherSecrets:   0.2,
			SpreadRumors:    0.1,
			Assassinate:     0.98,
			HireMercenaries: 0.97,
			HireSpies:       0.95,
			Frame:           0.96,
			Raid:            0.94,
			Steal:           0.98,
			Pillage:         0.98,
			Blackmail:       0.9,
			Kidnap:          0.95,
			ShadowWar:       0.98,
			Crusade:         0.2,
			War:             0.99,
		},
		ProbabilityOutlawCommodity: map[string]float64{
			OPIUM: 0.95,
		},
		ProbabilityOutlawResearch: map[string]float64{
			MAGIC_OCCULT: 0.98,
			MAGIC_ARCANA: 0.15,
			ALCHEMY:      0.10,
		},
		TaxFrequency:          config.Distribution{Min: DEFAULT_TICKS_PER_DAY * 90, Max: DEFAULT_TICKS_PER_DAY * 90, Mean: DEFAULT_TICKS_PER_DAY * 90, Deviation: 0},
		TaxRate:               config.Distribution{Min: 1, Max: 15, Mean: 10, Deviation: 5},
		ActionWeight:          config.Distribution{Min: structs.MaxEthos / 5, Max: structs.MaxEthos, Mean: structs.MaxEthos / 3, Deviation: structs.MaxEthos / 2},
		MilitaryOffenseBonus:  0.25,
		MilitaryDefenseBonus:  0.25,
		EspionageOffenseBonus: 0.10,
		EspionageDefenseBonus: 0.25,
	}
}
