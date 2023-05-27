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

	// Actions that do not appear in this are not permitted.
	Actions map[structs.ActionType]*Action

	// Demographics match name -> Demographics.
	Demographics map[string]*Demographics

	// TODO: queue (fan out for faction job calculations)
	// TODO: graph? (future w/ path calculations for trade)
	// TODO: event sink (allow caller to collect decisions)
}
