package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func FactionFantasy() *config.Faction {
	dev := &structs.Ethos{}
	return &config.Faction{
		EthosMean:      structs.Ethos{},
		EthosDeviation: *dev.Add(structs.MaxEthos),
		LeadershipProbability: map[structs.LeaderType]float64{
			structs.LeaderTypeSingle:  0.4,
			structs.LeaderTypeCouncil: 0.3,
			structs.LeaderTypeDual:    0.15,
			structs.LeaderTypeTriad:   0.1,
			structs.LeaderTypeAll:     0.05,
		},
		Wealth: config.Distribution{
			// min here is around the base land value
			Min:       120000000,  // 12,000 gp
			Max:       1000000000, // 100,000 gp
			Mean:      150000000,  // 15,000 gp
			Deviation: 900000000,  // 90,000 gp
		},
		Cohesion: config.Distribution{
			Min:       8000,
			Max:       10000,
			Mean:      9000,
			Deviation: 1000,
		},
		Corruption: config.Distribution{
			Min:       0,
			Max:       3000,
			Mean:      500,
			Deviation: 2000,
		},
		EspionageOffense: config.Distribution{
			Min:       0,
			Max:       2000,
			Mean:      1000,
			Deviation: 1000,
		},
		EspionageDefense: config.Distribution{
			Min:       0,
			Max:       4000,
			Mean:      2000,
			Deviation: 2000,
		},
		AreaProbability:     []float64{0.65, 0.2, 0.1, 0.05},
		PropertyProbability: []float64{0.40, 0.3, 0.25, 0.10, 0.05},
		Focuses: []config.Focus{
			config.Focus{
				Actions: []structs.ActionType{ // hedonists
					structs.ActionTypeFestival,
				},
				Probability: 0.10,
				Weight: config.Distribution{
					Min:       1000,
					Max:       5000,
					Mean:      2500,
					Deviation: 2500,
				},
				EspionageOffenseBonus: 0.0,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{
				Actions: []structs.ActionType{ // harvesters
					structs.ActionTypeHarvest,
				},
				Probability: 0.20,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      5000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.1,
				EspionageDefenseBonus: 0.1,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{
				Actions: []structs.ActionType{ // crafters
					structs.ActionTypeCraft,
				},
				Probability: 0.50,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      5000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.1,
				EspionageDefenseBonus: 0.1,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{
				Actions: []structs.ActionType{ // warfare
					structs.ActionTypePillage,
					structs.ActionTypeWar,
					structs.ActionTypeBribe,
					structs.ActionTypeRaid,
				},
				Probability: 0.10,
				Weight: config.Distribution{
					Min:       1000,
					Max:       5000,
					Mean:      3000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.0,
				EspionageDefenseBonus: 0.25,
				MilitaryOffenseBonus:  0.5,
				MilitaryDefenseBonus:  0.5,
			},
			config.Focus{
				Actions: []structs.ActionType{ // pirates
					structs.ActionTypeRaid,
					structs.ActionTypeBribe,
					structs.ActionTypeKidnap,
				},
				Probability: 0.01,
				Weight: config.Distribution{
					Min:       1000,
					Max:       5000,
					Mean:      3000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: -0.25,
				EspionageDefenseBonus: -0.25,
				MilitaryOffenseBonus:  0.25,
				MilitaryDefenseBonus:  0.25,
			},
			config.Focus{
				Actions: []structs.ActionType{ // removing people
					structs.ActionTypeAssassinate,
					structs.ActionTypeFrame,
				},
				Probability: 0.02,
				Weight: config.Distribution{
					Min:       2000,
					Max:       6000,
					Mean:      4000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.50,
				EspionageDefenseBonus: 0.50,
				MilitaryOffenseBonus:  -0.50,
				MilitaryDefenseBonus:  -0.50,
			},
			config.Focus{
				Actions: []structs.ActionType{ // secrets
					structs.ActionTypeGatherSecrets,
					structs.ActionTypeConcealSecrets,
				},
				Probability: 0.05,
				Weight: config.Distribution{
					Min:       2000,
					Max:       6000,
					Mean:      4000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.50,
				EspionageDefenseBonus: 0.50,
				MilitaryOffenseBonus:  -0.25,
				MilitaryDefenseBonus:  -0.25,
			},
			config.Focus{
				Actions: []structs.ActionType{ // research institution
					structs.ActionTypeResearch,
				},
				Probability: 0.15,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      5000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.25,
				EspionageDefenseBonus: 0.50,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{
				Actions: []structs.ActionType{ // cult
					structs.ActionTypeKidnap,
					structs.ActionTypeResearch,
					structs.ActionTypeConcealSecrets,
					structs.ActionTypeGatherSecrets,
				},
				Probability: 0.02,
				Weight: config.Distribution{
					Min:       4000,
					Max:       8000,
					Mean:      6000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.75,
				EspionageDefenseBonus: 0.75,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{
				Actions: []structs.ActionType{ // theft
					structs.ActionTypeBlackmail,
					structs.ActionTypeSteal,
					structs.ActionTypeKidnap,
					structs.ActionTypeConcealSecrets,
				},
				Probability: 0.2,
				Weight: config.Distribution{
					Min:       1000,
					Max:       5000,
					Mean:      3000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.25,
				EspionageDefenseBonus: 0.25,
				MilitaryOffenseBonus:  -0.25,
				MilitaryDefenseBonus:  -0.2,
			},
			config.Focus{ // corporate corruption
				Actions: []structs.ActionType{
					structs.ActionTypeBribe,
					structs.ActionTypeBlackmail,
				},
				Probability: 0.2,
				Weight: config.Distribution{
					Min:       1500,
					Max:       5500,
					Mean:      3000,
					Deviation: 1500,
				},
				EspionageOffenseBonus: 0.2,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{ // information corruption
				Actions: []structs.ActionType{
					structs.ActionTypeSpreadRumors,
					structs.ActionTypePropoganda,
				},
				Probability: 0.2,
				Weight: config.Distribution{
					Min:       2500,
					Max:       6500,
					Mean:      4000,
					Deviation: 2500,
				},
				EspionageOffenseBonus: 0.2,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  0.0,
				MilitaryDefenseBonus:  0.0,
			},
			config.Focus{ // expansion
				Actions: []structs.ActionType{
					structs.ActionTypeRecruit,
					structs.ActionTypeExpand,
				},
				Probability: 0.3,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      5000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.0,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  0.1,
				MilitaryDefenseBonus:  0.1,
			},
			config.Focus{ // trade
				Actions: []structs.ActionType{
					structs.ActionTypeTrade,
				},
				Probability: 0.30,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      5000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: -0.2,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  -0.2,
				MilitaryDefenseBonus:  -0.2,
			},
		},
		FocusProbability: []float64{0.0, 0.05, 0.65, 0.25, 0.05},
	}
}
