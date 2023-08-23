package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Action is configuration for a single ActionType (see pkg/structs/action.go)
type Action struct {
	// Ethos that matches this action; factions with matching ethos'
	// are more likely to perform it
	Ethos structs.Ethos

	// Min / Max number of people.
	MinPeople int // min: 1
	MaxPeople int // numbers <= 0 are ignored (considered "no max")

	// Weight given to performing this Action
	Probability float64

	// Target denotes that this action requires a target (faction) and the targets min/max trust
	// Ie. we only trade with factions we trust, and only make war on factions we don't.
	Target         bool
	TargetMinTrust int
	TargetMaxTrust int

	// Cost in whatever units the economy is using.
	// Additional explicit cost to performing the action.
	Cost Distribution

	// Time to spend reading action / gathering people, in ticks
	TimeToPrepare Distribution

	// How long an action will take to do after it's ready, in ticks
	TimeToExecute Distribution

	// Multiplied with the users EspionageDefense in order to determine
	// how covert the job is (factions with higher EspionageOffense can detect
	// these actions)
	SecrecyWeight float64

	// Amount to add per person who aids in the action, regardless of skills.
	// Ie. man power is needed rather than high skills.
	PersonWeight float64

	// Amount to add per person skill in given profession, weighted by their
	// actual skill level /100 (ie. someone with 100 "Smithing" adds the
	// entire weight to the roll).
	//
	// If someone has multiple skills that appear in this map, each of
	// them applies.
	ProfessionWeights map[string]float64

	// % of initial cost refunded on failure
	// (Ie. the action kicks off, but fails)
	RefundOnFailure float64

	// Goals are used so our code can determine which action might best suit the factions
	// current goal(s). What we mean here is what immediate goal does this action help with?
	// Ie. a faction low on money might look to "Wealth"
	Goals []structs.Goal
}
