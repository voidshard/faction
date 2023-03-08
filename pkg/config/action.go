package config

// Outcomes covers the possible outcomes and how to determine them
// based on RNG & various weights and measures.
type Outcomes struct {
	// probabilities (0-1) for each outcome
	//
	// Decided as follows:
	// - a random roll 0-1 is made and multiplied by `BaseRollWeight`
	// ie. a `BaseRollWeight` of .75 means the highest the base roll can be is .75.
	// Effectively this means the action outcome is less RNG and more skill / person based.
	// - for each person aiding in the action we add `PersonWeight`
	// - for each person with a profession that appears in `ProfessionWeights`
	// we add `ProfessionWeights[profession] * skill / 100`
	// ie. if a person has 50% skill in a profession, and the profession weight is 0.4,
	// then that person adds 0.2 to the roll.
	// - we add the parity weight(s)
	CriticalSuccess float64
	Success         float64
	Draw            float64
	Failure         float64
	CriticalFailure float64

	// BaseRollWeight is the maximum value of the base roll (0-1).
	// This effectively caps the RNG aspect of the outcome decision.
	BaseRollWeight float64

	// Amount to add per person who aids in the action, regardless of skills.
	// Ie. man power is needed rather than high skills.
	PersonWeight float64

	// Amount to add per person skill in given profession, weighted by their
	// actual skill level /100 (ie. someone with 100 "Smithing" adds the
	// entire weight).
	//
	// If someone has multiple skills that appear in this map, each of
	// them will apply.
	ProfessionWeights map[string]float64

	// Weight added based on the ratio of espionage offense / defense.
	//
	// Ie. We're comparing the attack & defense of the relevant factions.
	// If attack > defense then we add to the roll:
	// 	parity weight * attackers offense / defenders defense
	// If defense > attack then we subtract from the roll:
	// 	parity weight * defenders defense / attackers offense
	//
	// Eg.
	//   attack: 50, defense 100, parity weight 0.1
	//   Since defense > attack, we'll subtract from the roll:
	//   0.1 * 100 / 50 = 0.2
	//
	// These weights have a massive effect where the difference between
	// the attacker & defender is large, making some gaps unassialable.
	EspionageParityWeight float64

	// Weight added based on the ratio of military offense / defense.
	// See `EspionageParityWeight` for more details.
	MilitaryParityWeight float64
}

// Action is configuration for a single ActionType (see pkg/structs/action.go)
type Action struct {
	// cost in whatever units the economy is using
	Cost Distribution

	// time to spend reading action / gathering people, in ticks
	TimeToPrepare Distribution

	// how long an action will take to do after it's ready, in ticks
	TimeToEnact Distribution

	// probability that a government marks this as illegal
	// (ie. when a government is formed)
	IllegalProbability float64

	// information about how various outcomes are decided
	Outcomes Outcomes

	// if set, used as Key on Event object when the job with this Action has
	// a state change
	EventKey string
}
