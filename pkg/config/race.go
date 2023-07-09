package config

// Race imparts information about the species / race / breed of a being.
// These are physical limits & traits that are not simply cultural.
//
// Ie. a human can't live underwater, but a mermaid can, a human baby takes
// 9 months to be born, but a kitten takes 2.
//
// The companion struct 'Culture' is meant to impart information about the
// culture someone lives in.
// Ie. a human female might be capable of having a child at age 13, but culturally
// it might be considered odd to have a child before 18 (or whatever).
type Race struct {
	// Biomes that this race can survive in
	Biomes map[string]float64

	// Ages min/max (in ticks) at which someone can have children.
	ChildbearingAgeMin float64
	ChildbearingAgeMax float64

	// How long it takes to have a child
	ChildbearingTerm Distribution

	// probability of (mother's) death during childbirth (naturally)
	ChildbearingDeathProbability float64

	// Probability of a person dying of some natural cause
	DeathInfantMortalityProbability float64 // natural death in childhood
	DeathAdultMortalityProbability  float64 // natural death in adulthood

	// what actually kills someone ("naturally" in adulthood)
	DeathCauseAdultMortalityProbability map[string]float64 // what actually kills someone

	// how long someone lives (in ticks) before death of "old age"
	Lifespan Distribution

	// probability of being male
	IsMaleProbability float64
}
