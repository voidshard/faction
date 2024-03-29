package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

type ActionCategory string

const (
	ActionCategoryHostile    ActionCategory = "hostile"
	ActionCategoryUnfriendly ActionCategory = "unfriendly"
	ActionCategoryNeutral    ActionCategory = "" // left blank as it's the default if not specified
	ActionCategoryFriendly   ActionCategory = "friendly"
)

// Action is configuration for a single ActionType (see pkg/structs/action.go)
type Action struct {
	// Ethos that matches this action; factions with matching ethos'
	// are more likely to perform it
	Ethos structs.Ethos

	// If the engine should pick a more specific target (above Faction & Area).
	// Eg. this needs to target a specific plot for Raid.
	Target structs.MetaKey

	// Category of action gives us some hint as to how the action should be is perceived.
	Category ActionCategory

	// Min / Max number of people.
	MinPeople int // min: 1
	MaxPeople int // numbers <= 0 are ignored (considered "no max")

	// JobPriority is used when it comes to assigning jobs to people.
	// Higher priority Jobs are assigned first.
	// Nb. there are other weights at play
	// - people prefer to work for their faction (rather than an overlord)
	// - people prefer to work in their home area
	JobPriority int

	// Restricts the kinds of faction(s) that are allowed to perform this action.
	//
	// If this is set, this action cannot be performed unless the faction meets all
	// conditions in at least one of the slices.
	//
	// Ie. [][][]Condition{
	//    {ConditionA, ConditionB},
	//    {ConditionC},
	// }
	// Means a faction must meet (ConditionA && ConditionB) || (ConditionC)
	Restricted [][]Condition

	// Under what circumstances the user of this action can force people to take part (if any).
	// Probably should be used sparingly but it makes sense for some actions.
	//
	// A Job with Conscription will force people to take part if they're within faction controlled
	// area(s).
	Conscription [][]Condition

	// Weight given to performing this Action
	Probability float64

	// Target denotes that this action requires a target (faction) and the targets min/max trust
	// Ie. we only trade with factions we trust, and only make war on factions we don't.
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

	// MercenaryActions controls if this job is a mercenary job.
	// For this action to succeed there must be a Faction B in the area who is willing to perform the action.
	// Ie. this job is performed by Faction A, who finds a Faction B to that can perform one of these actions
	// on a target of their choosing.
	// Eg. If MercenaryActions = []string{"Assassinate"} then Faction A can hire Faction B to assassinate
	// someone of their choosing from Faction C.
	MercenaryActions []string
}
