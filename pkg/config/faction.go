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

	// MinYield is the minimum yield that a guild will accept, that assuming
	// the resource is attached to a LandRight.
	//
	// That is, if the land we're looking at when randomly building factions
	// yields only 1 iron ore, then that may not be enough for a guild of
	// miners to form -- since it's not a large enough part of the economy
	// to form a guild around.
	//
	// A 0 here can make sense if a land yield isn't relevant for guild formation;
	// ie. a silversmiths guild can form even if silver is imported from elsewhere.
	// But a guild of silver-ore miners require silver mines to operate.
	MinYield int
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

	// Probabilities to associate with each type of structure
	LeadershipStructureProbability map[structs.LeaderStructure]float64

	// starting values for wealth / cohesion / corruption.
	Wealth     Distribution
	Cohesion   Distribution
	Corruption Distribution

	// starting values for espionage offense / defense.
	EspionageOffense Distribution
	EspionageDefense Distribution

	// starting values for military offense / defense.
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

	// Number of Plots/LandRights a faction wil own (count by index + 1).
	// Minimum of 1.
	//
	// This can include
	// - LandRight(s)
	// - Plot(s)
	//
	// Note that the Guild 'MinYield' superceeds this.
	// ie. a guild will be given property to meet MinYield first, then will
	// acquire extra land (plots) if this isn't met.
	PropertyProbability []float64

	// Size of plots that can be allocated
	PlotSize Distribution
}
