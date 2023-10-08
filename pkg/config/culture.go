package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Culture
type Culture struct {
	// EthosMean represents the average outlook of members of the population
	EthosMean structs.Ethos
	// EthosDeviation is the standard deviation of the populace from the average (above).
	EthosDeviation structs.Ethos
	// Probability of a person getting a wildly different Ethos ("black sheep")
	EthosBlackSheepProbability float64

	// FamilySizeAverage is the average number of children in a family.
	// A family always has parents ..
	FamilySize Distribution

	// Chances of various em, extra relations
	FriendshipProbability      float64
	FriendshipCloseProbability float64
	EnemyProbability           float64
	EnemyHatedProbability      float64

	// The chance someone will marry (probably very high)
	MarriageProbability float64
	// The chance two people will divorce (hopefully very low)
	MarriageDivorceProbability float64
	// The chance one partner will have an affair (hopefully very low)
	MarriageAffairProbability float64

	// Professions to allocate to people.
	Professions []Profession
	// ProfessionProbability is a list of probabilities such that a given person has the number
	// of professions indicated by the index (index 0 means "" or "no profession").
	//
	// - a person cannot have the same profession multiple times
	// - a person may have only one profession with `ValidSideProfession` False
	// - a person can have any number of (unique) professions with `ValidSideProfession` True
	//
	// Ie. given
	//   {Name: farmer, ValidSideProfession: True}
	//   {Name: miner, ValidSideProfession: True}
	//   {Name: scribe, ValidSideProfession: False}
	//   {Name: priest, ValidSideProfession: False}
	// A person may not have "scribe" and "priest" at the same time, but could have any other
	// combination (including no profession).
	//
	// Eg. given SideProfessions = []float64{0.2, 0.4, 0.3, 0.1}
	// 20% chance of "no profession"
	// 40% chance of 1 profession
	// 30% chance of 2 professions
	// 10% chance of 3 professions
	ProfessionProbability []float64

	// Faiths to allocate to people.
	Faiths           []Faith
	FaithProbability []float64 // likelihood of 0 or more faiths

	// Ages min/max (in ticks) at which someone can have children (rolled once on birth)
	ChildbearingDeathProbability float64 // probability of (mother's) death during childbirth
	ChildbearingProbability      float64 // probability of becoming an expectant parent (per tick)

	// the min / max values are set by 'race'
	ChildbearingAgeMean      float64 // average age at which someone has children
	ChildbearingAgeDeviation float64 // standard deviation of age at which someone has children

	// Probability of a person dying of some natural cause in any given tick
	DeathInfantMortalityProbability float64 // natural death in childhood
	DeathAdultMortalityProbability  float64 // natural death in adulthood

	// what actually kills someone ("naturally" in adulthood)
	DeathCauseNaturalProbability map[string]float64

	// social structure denotes how "classes" are decided (if used)
	SocialStructure SocialStructure
}
