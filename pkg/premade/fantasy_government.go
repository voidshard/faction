package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

// GovernmentFantasy returns a fantasy government configuration.
//
// Well, for a human type setup.
// We could have a much more extreme setup for Drow or something.
func GovernmentFantasy() *config.Government {
	return &config.Government{
		ProbabilityOutlawAction: map[structs.ActionType]float64{
			structs.ActionTypeBribe:          0.6,
			structs.ActionTypePropoganda:     0.1,
			structs.ActionTypeResearch:       0.01,
			structs.ActionTypeExcommunicate:  0.02,
			structs.ActionTypeConcealSecrets: 0.05,
			structs.ActionTypeGatherSecrets:  0.2,
			structs.ActionTypeSpreadRumors:   0.1,
			structs.ActionTypeAssassinate:    0.98,
			structs.ActionTypeFrame:          0.96,
			structs.ActionTypeRaid:           0.94,
			structs.ActionTypeSteal:          0.98,
			structs.ActionTypePillage:        0.98,
			structs.ActionTypeBlackmail:      0.9,
			structs.ActionTypeKidnap:         0.95,
			structs.ActionTypeShadowWar:      0.98,
			structs.ActionTypeCrusade:        0.2,
			structs.ActionTypeWar:            0.99,
		},
		ProbabilityOutlawCommodity: map[string]float64{
			OPIUM: 0.95,
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
