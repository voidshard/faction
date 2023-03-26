package premade

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	religion1 = structs.NewID("1")
	religion2 = structs.NewID("2")
	religion3 = structs.NewID("3")
)

// DemographicsFantasyHuman returns a set of demographics for a fantasy human
// population.
//
// In general users are expected to create their own demographics, but this
// should help lazy people get started / test things out.
//
// Seriously please load these from a config or something .. this is just for
// playing around.
func DemographicsFantasyHuman() *config.Demographics {
	return &config.Demographics{
		FamilySize: config.Distribution{
			Min:       1,
			Max:       18,
			Mean:      8,
			Deviation: 4,
		},
		FriendshipProbability:      0.3,
		FriendshipCloseProbability: 0.1,
		EnemyProbability:           0.15,
		EnemyHatedProbability:      0.5,
		MarriageProbability:        0.8,
		MarriageDivorceProbability: 0.10,
		MarriageAffairProbability:  0.02,
		Race:                       "human",
		ChildbearingAge: config.Distribution{
			Min:       13 * DEFAULT_TICKS_PER_YEAR,
			Max:       45 * DEFAULT_TICKS_PER_YEAR,
			Mean:      22 * DEFAULT_TICKS_PER_YEAR,
			Deviation: 5 * DEFAULT_TICKS_PER_YEAR,
		},
		ChildbearingDeathProbability: 0.10,
		EthosMean:                    structs.Ethos{},
		EthosDeviation: structs.Ethos{
			Altruism:  2000,
			Ambition:  10000,
			Tradition: 4000,
			Pacifism:  9000,
			Piety:     9000,
			Caution:   1500,
		},
		EthosBlackSheepProbability:      0.02, // Probability given at least one radical ethos change
		DeathInfantMortalityProbability: 0.4,
		DeathAdultMortalityProbability:  0.1,
		DeathCauseNaturalProbability: map[string]float64{
			"malaria":       0.09, // natural diseases (50%)
			"pox":           0.08,
			"polio":         0.07,
			"dysentery":     0.07,
			"plague":        0.06,
			"measles":       0.05,
			"typoid":        0.03,
			"scarlet fever": 0.03,
			"flu":           0.02, // end diseases
			"accidental":    0.15,
			"war":           0.15,
			"starvation":    0.1,
			"suicide":       0.02,
			"execution":     0.02,
			"assassination": 0.01,
			"poisoning":     0.01,
			"animal attack": 0.04,
		},
		Professions: []config.Profession{
			config.Profession{
				Name:                FARMER,
				Occurs:              5.0,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Min:       500,
					Max:       10000,
					Mean:      2000,
					Deviation: 8000,
				},
			},
			config.Profession{
				Name:                MINER,
				Occurs:              0.2,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Min:       500,
					Max:       10000,
					Mean:      2000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:                FISHERMAN,
				Occurs:              0.25,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:                FORESTER,
				Occurs:              0.1,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:                WEAVER,
				Occurs:              0.05,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:   CLOTHIER,
				Occurs: 0.03,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:                TANNER,
				Occurs:              0.10,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:   LEATHERWORKER,
				Occurs: 0.05,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:   CARPENTER,
				Occurs: 0.20,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:                SMELTER,
				Occurs:              0.05,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:   SMITH,
				Occurs: 0.005,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:                SAILOR,
				Occurs:              0.20,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:                MERCHANT,
				Occurs:              0.05,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      3000,
					Deviation: 5000,
				},
			},
			config.Profession{
				Name:                CLERK,
				Occurs:              0.1,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2000,
					Deviation: 4000,
				},
			},
			config.Profession{
				Name:   MAGE,
				Occurs: 0.0005,
				Distribution: config.Distribution{
					Mean:      4000,
					Deviation: 6000,
				},
			},
			config.Profession{
				Name:                HUNTER,
				Occurs:              0.08,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      3000,
					Deviation: 6000,
				},
			},
			config.Profession{
				Name:   PRIEST,
				Occurs: 0.005,
				Distribution: config.Distribution{
					Mean:      4000,
					Deviation: 6000,
				},
			},
			config.Profession{
				Name:                SOLDIER,
				Occurs:              0.01,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:   SCRIBE,
				Occurs: 0.002,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:   THIEF,
				Occurs: 0.009,
				Distribution: config.Distribution{
					Mean:      1000,
					Deviation: 7500,
				},
			},
			config.Profession{
				Name:   ASSASSIN,
				Occurs: 0.003,
				Distribution: config.Distribution{
					Mean:      4000,
					Deviation: 5500,
				},
			},
			config.Profession{
				Name:   ALCHEMIST,
				Occurs: 0.008,
				Distribution: config.Distribution{
					Mean:      4000,
					Deviation: 5500,
				},
			},
		},
		ProfessionProbability: []float64{0.01, 0.2, 0.8, 0.04},
		Faiths: []config.Faith{
			config.Faith{
				ReligionID: religion1,
				Occurs:     0.8,
				Distribution: config.Distribution{
					Mean:      2500,
					Deviation: 8500,
				},
			},
			config.Faith{
				ReligionID: religion2,
				Occurs:     0.2,
				Distribution: config.Distribution{
					Mean:      7500,
					Deviation: 2000,
				},
			},
			config.Faith{
				ReligionID:     religion3,
				Occurs:         0.01,
				IsMonotheistic: true,
				Distribution: config.Distribution{
					Mean:      9000,
					Deviation: 500,
				},
			},
		},
		FaithProbability: []float64{0.05, 0.4, 0.3, 0.15, 0.05},
	}
}
