package structs

type JobTargetMetaKey string

const (
	JobTargetPerson    JobTargetMetaKey = "person"     // job targets specific person (ie. assassination)
	JobTargetPlot      JobTargetMetaKey = "plot"       // job taret specific plot (ie. raid)
	JobTargetLandRight JobTargetMetaKey = "land-right" // job targets specific land right (ie. mining)
)

type JobState string

const (
	JobStatePending JobState = "pending" // job is waiting to start (collecting people)
	JobStateActive  JobState = "active"  // job is in progress
	JobStateDone    JobState = "done"    // job is complete
	JobStateFailed  JobState = "failed"  // job failed to start (not enough people / cancelled)
)

// Job is what a faction creates when it wishes to perform an Action.
//
// People sympathetic to the faction who don't already have work sign on to 'work'
// jobs. If enough people signon by the time the job is registered to start then
// it goes ahead (ie. the action is attempted).
type Job struct {
	ID string `db:"id"`

	SourceFactionID string `db:"source_faction_id"` // ID of the faction posting the job
	SourceAreaID    string `db:"source_area_id"`    // where people will be recruited from

	Action ActionType `db:"action"` // action that is due to take place

	TargetFactionID string `db:"target_faction_id"` // ID of the target faction (if any)
	TargetAreaID    string `db:"target_area_id"`    // where the action will take place

	// key/val pair to hold adv. target metadata (ie. key:PERSON val:PERSON_ID)
	TargetMetaKey JobTargetMetaKey `db:"target_meta_key"`
	TargetMetaVal string           `db:"target_meta_value"`

	PeopleMin int `db:"people_min"` // required min number of people (else job fails to kick off)
	PeopleMax int `db:"people_max"` // max number of people that can work this (if any)

	TickCreated  int `db:"tick_created"`  // when the job was created
	TickStarts   int `db:"tick_starts"`   // when the job is due to start
	TickDuration int `db:"tick_duration"` // how long the job should last

	Secrecy   int  `db:"secrecy"`    // the result of an espionage defence roll (if covert)
	IsIllegal bool `db:"is_illegal"` // action has been outlawed

	State JobState `db:"state"` // current state of the job
}
