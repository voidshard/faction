package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

// Guild represents configuration for randomly creating a guild which involves or revolves
// around some commodities & associated professions
//
// (Commodity based faction).
type Guild struct {
	// Professions that this guild prefers
	Professions []string

	// Probability of this profession being selected (at all)
	Probability float64

	// LandMinCommodityYield notes that the guild needs areas of land that produce the given
	// commodities in order to form.
	//
	// Ie. a Mining / Iron working guild might need
	// 	IRON_ORE: 10
	// 	TIMBER: 100
	// Or whatever.
	//
	// When the guild is formed it will be given enough land to meet this across some
	// number of area(s).
	//
	// Applies to Plots with Commodity set.
	LandMinCommodityYield map[string]int

	// Desired import(s) (Commoidity names)
	Imports []string

	// Desired export(s) (Commoidity names)
	Exports []string
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

	// Explicitly set research topic(s) for this faction.
	// If there are more Research actions than topics, then topics will be randomly chosen to match.
	ResearchTopics []string
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

	// Probability of a given topic being researched by a faction.
	//
	// Where factions have no favoured profession, research will be chosen
	// at random from the given probabilities.
	//
	// Where factions have favoured profession(s) that have matching topics
	// (ie. Smith -> METALLURGY) then a topic will be chosen from the
	// favoured professions topics first.
	// If there are favoured professions but no matching topics, then
	// a topic will be chosen at random from all topics.
	ResearchProbability map[string]float64

	// Guilds are profession based factions; eg. a guild of blacksmiths,
	// a guild of assassins etc.
	Guilds []Guild
	// Probability that a faction is an amalgamation of more than one guild
	// (count by index).
	// ie. [0.2, 0.3, 0.4, ...] means a faction has a 30% chance of including one
	// professions, a 40% chance of including two professions (etc)
	GuildProbability []float64

	// Number of Plots a faction wil own (count by index + 1).
	// Minimum of 1.
	//
	// Note that the Guild 'MinYield' superceeds this.
	// ie. a guild will be given property to meet MinYield first, then will
	// acquire extra land if this isn't met.
	PropertyProbability []float64

	// Size of (non-commodity) plots that should be allocated.
	// Warning; factions will need to find this total land area (in Plot.Size) in
	// order to form.
	PlotSize Distribution

	// AllowEmptyPlotCreation allows the process spawning factions to generate empty plots
	// if not enough Plots are found.
	// If this isn't permitted, we'll error if we can't find enough Plots when spawning a faction.
	AllowEmptyPlotCreation bool

	// AllowCommodityPlotCreation allows the process spawning factions to generate plots
	// that yield Commodities if not enough Plots are found.
	// If this isn't permitted, we'll error if we can't find enough Plots when spawning a faction.
	AllowCommodityPlotCreation bool
}
