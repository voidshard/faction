package fantasy

import (
	"github.com/voidshard/faction/pkg/config"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	religion1 = structs.NewID("1")
	religion2 = structs.NewID("2")
	religion3 = structs.NewID("3")
)

// CultureHuman returns a set of demographics for a fantasy human
// population.
//
// In general users are expected to create their own demographics, but this
// should help lazy people get started / test things out.
//
// Seriously please load these from a config or something .. this is just for
// playing around.
func CultureHuman() *config.Culture {
	return &config.Culture{
		FamilySize: config.Distribution{
			Min:       1,
			Max:       18,
			Mean:      8,
			Deviation: 4,
		},
		// Probabilities of inter-personal relations per tick
		FriendshipProbability:        0.3 / float64(DEFAULT_TICKS_PER_YEAR),
		FriendshipCloseProbability:   0.1 / float64(DEFAULT_TICKS_PER_YEAR),
		EnemyProbability:             0.20 / float64(DEFAULT_TICKS_PER_YEAR),
		EnemyHatedProbability:        0.35 / float64(DEFAULT_TICKS_PER_YEAR),
		MarriageProbability:          0.6 / float64(DEFAULT_TICKS_PER_YEAR*10),
		MarriageDivorceProbability:   0.02 / float64(DEFAULT_TICKS_PER_YEAR*25),
		MarriageAffairProbability:    0.04 / float64(DEFAULT_TICKS_PER_YEAR*15),
		ChildbearingDeathProbability: 0.10, // actual medieval rate probably higher (on Birth)
		ChildbearingProbability:      .2,   //      1 / (float64(DEFAULT_TICKS_PER_YEAR) * 2), // very roughly every 2 years
		EthosMean:                    structs.Ethos{},
		EthosDeviation: structs.Ethos{
			Altruism:  1000,
			Ambition:  7000,
			Tradition: 4000,
			Pacifism:  3000,
			Piety:     7000,
			Caution:   1500,
		},
		EthosBlackSheepProbability:      0.02, // Probability given at least one radical ethos change (on Birth)
		DeathInfantMortalityProbability: 0.2,  // Probability of death in childbirth
		DeathAdultMortalityProbability:  0.02 / float64(DEFAULT_TICKS_PER_YEAR),
		DeathCauseNaturalProbability: map[string]float64{
			"malaria":       0.09, // natural diseases (~50%)
			"pox":           0.08,
			"polio":         0.07,
			"dysentery":     0.07,
			"plague":        0.06,
			"measles":       0.05,
			"typoid":        0.03,
			"scarlet fever": 0.03,
			"flu":           0.02, // end diseases
			"accidental":    0.20,
			"starvation":    0.1,
			"suicide":       0.02,
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
		SocialStructure: config.SocialStructure{
			Base: 0,
			Classes: map[string]int{
				"elite":  100,
				"upper":  75,
				"middle": 50,
				"lower":  25,
				"serf":   0,
			},
			Profession: map[string]float64{
				NOBLE:     0.05,  // 2000 "noble" puts one in the elite class, 1500 in the upper class
				MERCHANT:  0.009, // high skill as a merchant puts one in the upper class
				MAGE:      0.009,
				SCHOLAR:   0.008, // high skill as a scholar puts one in the upper class
				ALCHEMIST: 0.008,
				ASSASSIN:  0.006,
				PRIEST:    0.006,
				SMITH:     0.005, // high skill as a smith puts one in the high-middle class
				SMELTER:   0.004,
				TANNER:    0.004, // high skill as a tanner puts one in the lower-middle class
				SPY:       0.004,
				CARPENTER: 0.004,
				SOLDIER:   0.003, // high skill as a soldier puts one in the lower class
				CLERK:     0.003,
				WEAVER:    0.002,
				THIEF:     -0.002, // high skill as a thief puts one in the lower class
				FARMER:    -0.001,
			},
			Faith: map[string]float64{
				religion1: 0.002,
				religion2: 0.001,
				religion3: -0.001,
			},
		},
	}
}
