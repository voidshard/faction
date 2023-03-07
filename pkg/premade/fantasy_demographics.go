package premade

import (
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
func DemographicsFantasyHuman() *structs.Demographics {
	return &structs.Demographics{
		FamilySizeAverage:            8, // children
		FamilySizeDeviation:          4,
		FamilySizeMax:                18,
		FriendshipProbability:        0.3,
		FriendshipCloseProbability:   0.1,
		EnemyProbability:             0.15,
		EnemyHatedProbability:        0.5,
		MarriageProbability:          0.8,
		MarriageDivorceProbability:   0.10,
		MarriageAffairProbability:    0.02,
		Race:                         "human",
		ChildbearingAgeMin:           13 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeMax:           45 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeAverage:       22 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingAgeDeviation:     5 * DEFAULT_TICKS_PER_YEAR,
		ChildbearingDeathProbability: 0.10,
		EthosAverage:                 &structs.Ethos{},
		EthosDeviation: &structs.Ethos{
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
		Professions: []*structs.Profession{
			&structs.Profession{
				Name:                FARMER,
				Occurs:              5.0,
				Average:             20,
				Deviation:           5,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                MINER,
				Occurs:              0.2,
				Average:             40,
				Deviation:           5,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                FISHERMAN,
				Occurs:              0.25,
				Average:             20,
				Deviation:           10,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                FORESTER,
				Occurs:              0.1,
				Average:             20,
				Deviation:           5,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                WEAVER,
				Occurs:              0.05,
				Average:             50,
				Deviation:           5,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                CLOTHIER,
				Occurs:              0.03,
				Average:             75,
				Deviation:           5,
				ValidSideProfession: false,
			},
			&structs.Profession{
				Name:                TANNER,
				Occurs:              0.10,
				Average:             10,
				Deviation:           2,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                LEATHERWORKER,
				Occurs:              0.05,
				Average:             25,
				Deviation:           10,
				ValidSideProfession: false,
			},
			&structs.Profession{
				Name:                CARPENTER,
				Occurs:              0.20,
				Average:             50,
				Deviation:           5,
				ValidSideProfession: false,
			},
			&structs.Profession{
				Name:                SMELTER,
				Occurs:              0.05,
				Average:             25,
				Deviation:           4,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                SMITH,
				Occurs:              0.005,
				Average:             85,
				Deviation:           10,
				ValidSideProfession: false,
			},
			&structs.Profession{
				Name:                SAILOR,
				Occurs:              0.20,
				Average:             30,
				Deviation:           10,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                MERCHANT,
				Occurs:              0.05,
				Average:             30,
				Deviation:           8,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                CLERK,
				Occurs:              0.1,
				Average:             20,
				Deviation:           10,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                MAGE,
				Occurs:              0.0005,
				Average:             95,
				Deviation:           1,
				ValidSideProfession: false,
			},
			&structs.Profession{
				Name:                HUNTER,
				Occurs:              0.08,
				Average:             80,
				Deviation:           5,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                PRIEST,
				Occurs:              0.005,
				Average:             80,
				Deviation:           3,
				ValidSideProfession: false,
			},
			&structs.Profession{
				Name:                SOLDIER,
				Occurs:              0.01,
				Average:             45,
				Deviation:           10,
				ValidSideProfession: true,
			},
			&structs.Profession{
				Name:                SCRIBE,
				Occurs:              0.002,
				Average:             80,
				Deviation:           4,
				ValidSideProfession: false,
			},
		},
		ProfessionProbability: []float64{0.01, 0.2, 0.8, 0.04},
		Faiths: []*structs.Faith{
			&structs.Faith{
				ReligionID:     religion1,
				Occurs:         0.8,
				Average:        25,
				Deviation:      5,
				IsMonotheistic: false,
			},
			&structs.Faith{
				ReligionID:     religion2,
				Occurs:         0.2,
				Average:        75,
				Deviation:      2,
				IsMonotheistic: false,
			},
			&structs.Faith{
				ReligionID:     religion3,
				Occurs:         0.01,
				Average:        90,
				Deviation:      1,
				IsMonotheistic: true,
			},
		},
		FaithProbability: []float64{0.05, 0.4, 0.3, 0.15, 0.05},
	}
}
