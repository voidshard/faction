package config

type Race struct {
	// Ages min/max (in ticks) at which someone can have children.
	ChildbearingAgeMin           float64
	ChildbearingAgeMax           float64
	ChildbearingTerm             float64
	ChildbearingDeathProbability float64 // probability of (mother's) death during childbirth

	// Probability of a person dying of some natural cause
	DeathInfantMortalityProbability float64            // natural death in childhood
	DeathAdultMortalityProbability  float64            // natural death in adulthood
	DeathCauseNaturalProbability    map[string]float64 // what actually kills someone
}
