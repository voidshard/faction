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

	// Probability that a government marks this as illegal
	// (ie. when a government is formed)
	IllegalProbability float64

	// used as Key on Event object when the job with this Action has
	// a state change.
	// If not set, the string(ActionType) is used.
	EventKey string

	// % of initial cost refunded if cancelled
	// (Ie. the action fails to gather enough people to kick off)
	RefundOnCancel float64

	// % of initial cost refunded on failure
	// (Ie. the action kicks off, but fails)
	RefundOnFailure float64
}
