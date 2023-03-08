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
			Altruism:  20,
			Ambition:  100,
			Tradition: 40,
			Pacifism:  90,
			Piety:     90,
			Caution:   15,
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
					Min:       5,
					Max:       100,
					Mean:      20,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:                MINER,
				Occurs:              0.2,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Min:       5,
					Max:       100,
					Mean:      40,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:                FISHERMAN,
				Occurs:              0.25,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      20,
					Deviation: 10,
				},
			},
			config.Profession{
				Name:                FORESTER,
				Occurs:              0.1,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      20,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:                WEAVER,
				Occurs:              0.05,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      50,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:   CLOTHIER,
				Occurs: 0.03,
				Distribution: config.Distribution{
					Mean:      75,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:                TANNER,
				Occurs:              0.10,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      10,
					Deviation: 2,
				},
			},
			config.Profession{
				Name:   LEATHERWORKER,
				Occurs: 0.05,
				Distribution: config.Distribution{
					Mean:      25,
					Deviation: 10,
				},
			},
			config.Profession{
				Name:   CARPENTER,
				Occurs: 0.20,
				Distribution: config.Distribution{
					Mean:      50,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:                SMELTER,
				Occurs:              0.05,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      25,
					Deviation: 4,
				},
			},
			config.Profession{
				Name:   SMITH,
				Occurs: 0.005,
				Distribution: config.Distribution{
					Mean:      85,
					Deviation: 10,
				},
			},
			config.Profession{
				Name:                SAILOR,
				Occurs:              0.20,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      30,
					Deviation: 10,
				},
			},
			config.Profession{
				Name:                MERCHANT,
				Occurs:              0.05,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      30,
					Deviation: 8,
				},
			},
			config.Profession{
				Name:                CLERK,
				Occurs:              0.1,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      20,
					Deviation: 10,
				},
			},
			config.Profession{
				Name:   MAGE,
				Occurs: 0.0005,
				Distribution: config.Distribution{
					Mean:      95,
					Deviation: 1,
				},
			},
			config.Profession{
				Name:                HUNTER,
				Occurs:              0.08,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      80,
					Deviation: 5,
				},
			},
			config.Profession{
				Name:   PRIEST,
				Occurs: 0.005,
				Distribution: config.Distribution{
					Mean:      80,
					Deviation: 3,
				},
			},
			config.Profession{
				Name:                SOLDIER,
				Occurs:              0.01,
				ValidSideProfession: true,
				Distribution: config.Distribution{
					Mean:      45,
					Deviation: 10,
				},
			},
			config.Profession{
				Name:   SCRIBE,
				Occurs: 0.002,
				Distribution: config.Distribution{
					Mean:      80,
					Deviation: 4,
				},
			},
		},
		ProfessionProbability: []float64{0.01, 0.2, 0.8, 0.04},
		Faiths: []config.Faith{
			config.Faith{
				ReligionID: religion1,
				Occurs:     0.8,
				Distribution: config.Distribution{
					Mean:      25,
					Deviation: 5,
				},
			},
			config.Faith{
				ReligionID: religion2,
				Occurs:     0.2,
				Distribution: config.Distribution{
					Mean:      75,
					Deviation: 2,
				},
			},
			config.Faith{
				ReligionID:     religion3,
				Occurs:         0.01,
				IsMonotheistic: true,
				Distribution: config.Distribution{
					Mean:      90,
					Deviation: 1,
				},
			},
		},
		FaithProbability: []float64{0.05, 0.4, 0.3, 0.15, 0.05},
	}
}
