package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

func Faction() *config.Faction {
	dev := &structs.Ethos{}
	return &config.Faction{
		EthosMean:      structs.Ethos{},
		EthosDeviation: *dev.Add(structs.MaxEthos / 4), // ethos is further weighted by favoured actions
		LeadershipProbability: map[structs.LeaderType]float64{
			structs.LeaderTypeSingle:  0.4,
			structs.LeaderTypeCouncil: 0.3,
			structs.LeaderTypeDual:    0.15,
			structs.LeaderTypeTriad:   0.1,
			structs.LeaderTypeAll:     0.05,
		},
		LeadershipStructureProbability: map[structs.LeaderStructure]float64{
			structs.LeaderStructurePyramid: 0.60,
			structs.LeaderStructureLoose:   0.35,
			structs.LeaderStructureCell:    0.05,
		},
		Wealth: config.Distribution{
			// min here is around the base land value
			Min:       20000000,  // 2,000 gp
			Max:       200000000, // 20,000 gp
			Mean:      12000000,  // 12,000 gp
			Deviation: 100000000, // 10,000 gp
		},
		Cohesion: config.Distribution{
			Min:       8000, // higher min to begin with so factions don't immediately split
			Max:       10000,
			Mean:      9000,
			Deviation: 1000,
		},
		Corruption: config.Distribution{
			Min:       10,
			Max:       3000,
			Mean:      500,
			Deviation: 2000,
		},
		EspionageOffense: config.Distribution{
			Min:       100,
			Max:       2000,
			Mean:      1000,
			Deviation: 1000,
		},
		EspionageDefense: config.Distribution{
			Min:       100,
			Max:       4000,
			Mean:      2000,
			Deviation: 2000,
		},
		MilitaryOffense: config.Distribution{
			Min:       100,
			Max:       500,
			Mean:      100,
			Deviation: 400,
		},
		MilitaryDefense: config.Distribution{
			Min:       100,
			Max:       2000,
			Mean:      200,
			Deviation: 1800,
		},
		PropertyProbability: []float64{0.40, 0.3, 0.25, 0.10, 0.05},
		Focuses: []config.Focus{
			config.Focus{
				Actions: []structs.ActionType{ // hedonists
					structs.ActionTypeFestival,
				},
				Probability: 0.05,
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
				Actions: []structs.ActionType{ // warfare
					structs.ActionTypePillage,
					structs.ActionTypeWar,
					structs.ActionTypeBribe,
					structs.ActionTypeRaid,
				},
				Probability: 0.02,
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
				Probability: 0.06,
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
					structs.ActionTypeHireMercenaries,
				},
				Probability: 0.1,
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
				Probability: 0.15,
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
					structs.ActionTypeHireMercenaries,
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
					structs.ActionTypeHireMercenaries,
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
		Guilds: []config.Guild{
			config.Guild{
				Profession:  MERCHANT,
				Probability: 0.5,
			},
			config.Guild{
				Profession:  SMITH,
				Probability: 0.2,
			},
			config.Guild{
				Profession:  ALCHEMIST,
				Probability: 0.1,
			},
			config.Guild{
				Profession:  MAGE,
				Probability: 0.02,
			},
			config.Guild{
				Profession:  CLOTHIER,
				Probability: 0.05,
			},
			config.Guild{
				Profession:  CARPENTER,
				Probability: 0.05,
			},
			config.Guild{
				Profession:  SCRIBE,
				Probability: 0.05,
			},
			config.Guild{
				Profession:  MINER,
				Probability: 0.9,
				MinYield:    100,
			},
			config.Guild{
				Profession:  FARMER,
				Probability: 0.9,
				MinYield:    200,
			},
		},
		GuildProbability: []float64{0.30, 0.55, 0.1, 0.05},
		PlotSize: config.Distribution{
			Min:       100,
			Max:       10000,
			Mean:      500,
			Deviation: 9500,
		},
		ResearchProbability: map[string]float64{
			AGRICULTURE:  0.01,
			ASTRONOMY:    0.2,
			WARFARE:      0.05,
			METALLURGY:   0.04,
			PHILOSOPHY:   0.1,
			MEDICINE:     0.1,
			MATHEMATICS:  0.03,
			LITERATURE:   0.02,
			LAW:          0.1,
			ARCHITECTURE: 0.04,
			THEOLOGY:     0.06,
			MAGIC_ARCANA: 0.15,
			MAGIC_OCCULT: 0.03,
			ALCHEMY:      0.03,
		},
	}
}
