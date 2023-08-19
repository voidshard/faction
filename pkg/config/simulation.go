package config

import (
	"github.com/voidshard/faction/pkg/structs"
)

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
	Actions map[structs.ActionType]*Action

	// Race name -> Race
	Races map[string]*Race

	// Culture name -> Culture
	Cultures map[string]*Culture

	// TODO: graph? (future w/ path calculations for trade)
	// TODO: event sink (allow caller to collect decisions)
}
