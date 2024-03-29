package structs

// JobState is the state of a job
// - pending on first creation, awaiting people to signup
// - active when enough people have signed up
// - done when the action has been attempted
// - failed when the action could not be attempted / was cancelled
type JobState string

const (
	JobStatePending JobState = "pending" // job is waiting to start (collecting people)
	JobStateReady   JobState = "ready"   // job is ready to start (enough people)
	JobStateActive  JobState = "active"  // job is in progress
	JobStateDone    JobState = "done"    // job is complete
	JobStateFailed  JobState = "failed"  // job failed to start (not enough people / cancelled)
)

var (
	AllJobStates = []JobState{
		JobStatePending,
		JobStateReady,
		JobStateActive,
		JobStateDone,
		JobStateFailed,
	}
)

// Job is what a faction creates when it wishes to perform an Action.
//
// People sympathetic to the faction who don't already have work sign on to 'work'
// jobs. If enough people signon by the time the job is registered to start then
// it goes ahead (ie. the action is attempted).
type Job struct {
	ID          string `db:"id"`
	ParentJobID string `db:"parent_job_id"` // if this job is a sub-job, this is the parent job ID

	SourceFactionID string `db:"source_faction_id"` // ID of the faction posting the job
	SourceAreaID    string `db:"source_area_id"`    // where people will be recruited from

	Action   string `db:"action"`   // action that is due to take place
	Priority int    `db:"priority"` // priority of this job

	Conscription bool `db:"conscription"` // if job is allowed to force people take part

	TargetFactionID string `db:"target_faction_id"` // ID of the faction the action is aimed at
	TargetAreaID    string `db:"target_area_id"`    // where the action will take place

	// key/val pair to hold adv. target metadata (ie. key:PERSON val:PERSON_ID)
	// We only set this if TargetFactionID and TargetAreaID (always set) are not enough
	// ie. we target a specfic person, or a specific building etc within some area & faction
	TargetMetaKey MetaKey `db:"target_meta_key"`
	TargetMetaVal string  `db:"target_meta_val"`

	PeopleMin int `db:"people_min"` // required min number of people (else job fails to kick off)
	PeopleMax int `db:"people_max"` // max number of people that can work this (if any)
	PeopleNow int `db:"people_now"` // people signed up to work this

	TickCreated int `db:"tick_created"` // when the job was created
	TickStarts  int `db:"tick_starts"`  // when the job is due to start
	TickEnds    int `db:"tick_ends"`    // when the job will end

	Secrecy   int  `db:"secrecy"`    // the result of an espionage defence roll (if covert)
	IsIllegal bool `db:"is_illegal"` // action has been outlawed

	State JobState `db:"state"` // current state of the job
}

func (j *Job) ObjectID() string {
	return j.ID
}
