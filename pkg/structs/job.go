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
	ID string

	SourceFactionID string // ID of the faction posting the job
	SourceAreaID    string // where people will be recruited from

	Action ActionType // action that is due to take place

	TargetFactionID string // ID of the target faction (if any)
	TargetAreaID    string // where the action will take place

	TargetMetaKey JobTargetMetaKey // key/val pair to hold adv. target metadata (ie. key:PERSON val:PERSON_ID)
	TargetMetaVal string

	PeopleMin int // required min number of people (else job fails to kick off)
	PeopleMax int // max number of people that can work this (if any)

	TickCreated  int // when the job was created
	TickStarts   int // when the job is due to start
	TickDuration int // how long the job should last

	Secrecy   int  // the result of an espionage defence roll (if covert)
	IsIllegal bool // action has been outlawed

	State JobState // current state of the job
}
