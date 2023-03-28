package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Guild represents configuration for randomly creating a guild.
// (Profession based faction).
type Guild struct {
	// Profession that this guild prefers
	Profession string

	// Probability of this profession being selected (at all)
	Probability float64

	// Min number of skill (of members) for us to consider a guild.
	//
	// Ie. low skill with high members might imply a relatively
	// easy to join guild (eg. a guild of farmers).
	//
	// A high skill, low members might imply some faction of
	// elite warriors, mages etc.
	MinSkill int

	// Min number of members for us to even consider forming a guild
	// See. MinSkill above
	MinMembers int
}

// Focus represents a set of actions that a faction prefers to perform.
// A faction can have multiple focuses.
//
// A Focus is a set of (at least one) actions;
// Eg.
// a Temple might prefer Festival, Excommunicate, Charity
// a Guild of Thieves might prefer Steal, Kidnap
//
// This way we can form unusual factions; eg. a Temple of Thieves
type Focus struct {
	// Actions included in this focus
	//
	// :Warning:
	// Keep in mind GovernmentOnly / ReligionOnly actions noted in
	// pkg/structs/action.go
	// Such actions can in theory appear here, but will do nothing if the faction cannot
	// actually perform them (ie. a non-religion faction cannot Excommunicate).
	// Probably it's best to avoid adding such actions here .. since it wont achieve anything.
	Actions []structs.ActionType

	// Probability that this focus is chosen
	Probability float64

	// Action weight to apply to these actions.
	// Ie. how much to weight the faction in favour of these.
	Weight Distribution

	// Bonuses to apply to given values if this Focus is used
	EspionageOffenseBonus float64
	EspionageDefenseBonus float64
	MilitaryOffenseBonus  float64
	MilitaryDefenseBonus  float64
}

// Faction represents configuration for randomly creating a faction.
type Faction struct {
	// Ethos settings
	EthosMean      structs.Ethos
	EthosDeviation structs.Ethos

	// Probabilities to associate with each type of leadership
	LeadershipProbability map[structs.LeaderType]float64

	// starting values for wealth / cohesion / corruption.
	Wealth     Distribution
	Cohesion   Distribution
	Corruption Distribution

	// starting values for espionage offense / defense.
	// Factions marked "IsCovert" (illegal factions) will be given another 25% buff to both
	EspionageOffense Distribution
	EspionageDefense Distribution

	// starting values for military offense / defense.
	//
	// Government factions will be given another 25% buff to both.
	// Trade guilds will recieve a -25% debuff to both.
	MilitaryOffense Distribution
	MilitaryDefense Distribution

	// Actions that can be selected as a "focus"
	Focuses []Focus
	// Probability of a faction having some number of focuses (count by index).
	//
	// Eg. given [0.5, 0.3, 0.2], a faction has a 50% chance of having 0 focuses,
	// a 30% chance of having 1 focus, and a 20% chance of having 2 focuses.
	//
	// A 'focus' is a action(s) that the given faction prefers to perform
	// ("business as usual"). A faction with an illegal focus is covert by definition.
	//
	// Nb. you almost certainly want a faction to have at least one focus
	FocusProbability []float64

	// Guilds are profession based factions; eg. a guild of blacksmiths,
	// a guild of assassins etc.
	Guilds []Guild
	// Probability that a faction is an amalgamation of more than one guild
	// (count by index).
	// ie. [0.2, 0.3, 0.4, ...] means a faction has a 30% chance of including one
	// professions, a 40% chance of including two professions (etc)
	GuildProbability []float64

	// Probability that a faction will be in some number of areas (count by index + 1).
	// Minimum of 1.
	//
	// Ie.
	// - [0.5, 0.3, 0.2] means a faction has a 50% chance of being in 1 area (index 0 + 1)
	// and a 30% chance of being in 2 areas (index 1 + 1) etc.
	//
	// Factions are marked IsCovert if one of their Focus actions are illegal in at least
	// one area they are in.
	AreaProbability []float64

	// Number of Plots/LandRights a faction will have in each area (count by index + 1).
	// Minimum of 1.
	//
	// This can include
	// - LandRight(s)
	// - Plot(s)
	PropertyProbability []float64
}
