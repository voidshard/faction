package fantasy

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
					structs.ActionTypeEnslave,
					structs.ActionTypeHireSpies,
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
				Actions: []structs.ActionType{ // removing people (directly or reputation)
					structs.ActionTypeAssassinate,
					structs.ActionTypeFrame,
					structs.ActionTypeSpreadRumors,
					structs.ActionTypeHireMercenaries,
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
					structs.ActionTypeResearch,
				},
				ResearchTopics: []string{MAGIC_ARCANA}, // ie. MAGIC_ARCANA + another research topic
				Probability:    0.10,
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
				Actions: []structs.ActionType{ // standard fantasy temple / healing institution
					structs.ActionTypeResearch,
					structs.ActionTypeResearch,
					structs.ActionTypeCharity,
				},
				ResearchTopics: []string{MEDICINE, THEOLOGY},
				Probability:    0.08,
				Weight: config.Distribution{
					Min:       1000,
					Max:       3000,
					Mean:      2000,
					Deviation: 1000,
				},
				EspionageOffenseBonus: 0.25,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  -0.25,
				MilitaryDefenseBonus:  -0.25,
			},
			config.Focus{
				Actions: []structs.ActionType{ // knightly order
					structs.ActionTypeResearch,
					structs.ActionTypeResearch,
					structs.ActionTypeRaid,
					structs.ActionTypePillage,
					structs.ActionTypeRecruit,
				},
				ResearchTopics: []string{WARFARE, THEOLOGY},
				Probability:    0.05,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      4000,
					Deviation: 3000,
				},
				EspionageOffenseBonus: 0.0,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  0.5,
				MilitaryDefenseBonus:  0.5,
			},
			config.Focus{
				Actions: []structs.ActionType{ // cult
					structs.ActionTypeKidnap,
					structs.ActionTypeResearch,
					structs.ActionTypeConcealSecrets,
					structs.ActionTypeGatherSecrets,
				},
				ResearchTopics: []string{MAGIC_OCCULT},
				Probability:    0.02,
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
					structs.ActionTypeHireSpies,
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
					structs.ActionTypeFrame,
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
			config.Focus{ // expansion, unfriendly style
				Actions: []structs.ActionType{
					structs.ActionTypeRecruit,
					structs.ActionTypeExpand,
					structs.ActionTypeEnslave,
					structs.ActionTypeHireMercenaries,
					structs.ActionTypeHireSpies,
				},
				Probability: 0.2,
				Weight: config.Distribution{
					Min:       3000,
					Max:       7000,
					Mean:      5000,
					Deviation: 2000,
				},
				EspionageOffenseBonus: 0.0,
				EspionageDefenseBonus: 0.0,
				MilitaryOffenseBonus:  0.2,
				MilitaryDefenseBonus:  0.1,
			},
			config.Focus{ // trade
				Actions: []structs.ActionType{
					structs.ActionTypeTrade,
					structs.ActionTypeHireSpies,
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
		FocusProbability: []float64{0.0, 0.60, 0.25, 0.15},
		Guilds: []config.Guild{
			// Nb. guilds are factions that focus on some good(s) or production chain(s).
			// Where guilds are "vertical integrations" or harvest their own inputs, no Imports are needed.
			config.Guild{ // iron / steel working + charcoal production
				// this sweeping vertical integration is probably exclusively state run / appointed
				Professions: []string{SMITH, MINER, SMELTER, CLERK, FORESTER},
				Probability: 0.01,
				LandMinCommodityYield: map[string]int{
					IRON_ORE: 250,
					TIMBER:   500,
				},
				Exports: []string{IRON_TOOLS, STEEL_ARMOUR, STEEL_WEAPON, STEEL_TOOLS},
			},
			config.Guild{ // a more classic "metal smith" guild
				// this sweeping vertical integration is probably exclusively state run / appointed
				Professions:           []string{SMITH, MERCHANT, CLERK},
				Probability:           0.08,
				LandMinCommodityYield: map[string]int{},
				Imports:               []string{IRON_INGOT, STEEL_INGOT, TIMBER},
				Exports:               []string{IRON_TOOLS, STEEL_ARMOUR, STEEL_WEAPON, STEEL_TOOLS},
			},
			config.Guild{ // textile production (flax -> linen)
				Professions: []string{FARMER, CLOTHIER, WEAVER, MERCHANT, CLERK},
				Probability: 0.40, // textile production & trade was absurdly common
				LandMinCommodityYield: map[string]int{
					FLAX: 400,
				},
				Exports: []string{LINEN}, // pre-made clothing was not so common
			},
			config.Guild{
				Professions: []string{TANNER, LEATHERWORKER},
				Probability: 0.10,
				LandMinCommodityYield: map[string]int{
					WILD_GAME: 50,
					FODDER:    200,
				},
				Exports: []string{LEATHER, MEAT},
			},
			config.Guild{ // wood working
				Professions: []string{CARPENTER, FORESTER},
				Probability: 0.10,
				LandMinCommodityYield: map[string]int{
					TIMBER: 400,
				},
				Exports: []string{WOODEN_FURNITURE, WOODEN_TOOLS},
			},
			config.Guild{ // farming, with a side focus on high value crops
				Professions: []string{FARMER, THIEF},
				Probability: 0.01,
				LandMinCommodityYield: map[string]int{
					OPIUM: 10,
					WHEAT: 500,
				},
				Exports: []string{OPIUM, WHEAT},
			},
			config.Guild{ // farming (rare, selling food stuffs long distance was expensive & awkward)
				Professions: []string{FARMER, SAILOR, MERCHANT},
				Probability: 0.02,
				LandMinCommodityYield: map[string]int{
					WHEAT: 1000,
				},
				Exports: []string{WHEAT},
			},
		},
		GuildProbability: []float64{0.15, 0.65, 0.15},
		PlotSize: config.Distribution{ // metres squared
			Min:       100,
			Max:       500,
			Mean:      250,
			Deviation: 250,
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
		AllowEmptyPlotCreation:     true,  // allow non-commodity yielding plots to be created if not enough can be found
		AllowCommodityPlotCreation: false, // require commodity yielding plots to be found in the DB
	}
}
