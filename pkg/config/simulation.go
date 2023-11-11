package config

// Simulation Configuration to instantiate a Simulation object.
//
// Our main concern here are things like;
// - base professions / skills we can assign people
// - base probabilities of various professions within a demographic
// - base Action probabilities
// - base probabilities of various forms of government
type Simulation struct {
	// Database configuration information
	Database *Database

	// Queue configuration information
	Queue *Queue

	// Actions that do not appear in this are not permitted.
	Actions map[string]*Action

	// Race name -> Race
	Races map[string]*Race

	// Culture name -> Culture
	Cultures map[string]*Culture

	// Global settings that control how the simulation runs
	Settings *SimulationSettings
}

// SimulationSettings are settings that control how the simulation runs & global values / weights
// for internal calculations.
type SimulationSettings struct {
	// When planning jobs, this is the multiplier used when a faction believes it's survival
	// is at stake. (Probably, a large one).
	// Ie. if we're at risk of running out of money, this is applied to all "wealth" generating
	// actions when we're planning our next move.
	SurvivalGoalWeight float64

	// If our faction has less than this, we'll consider growth a priority.
	SurvivalMinPeople int

	// When planning jobs, this is the multiplier used to apply to faction goals when survival is
	// not at stake.
	GoalWeight float64
}

// DefaultSimulationSettings returns a new SimulationSettings object with default values.
func DefaultSimulationSettings() *SimulationSettings {
	return &SimulationSettings{
		SurvivalGoalWeight: 10.0,
		GoalWeight:         3.0,
		SurvivalMinPeople:  50,
	}
}
