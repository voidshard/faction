package structs

// Job is what a faction creates when it wishes to perform an Action.
//
// People sympathetic to the faction who don't already have work sign on to 'work'
// jobs. If enough people signon by the time the job is registered to start then
// it goes ahead (ie. the action is attempted).
type Job struct {
	ID string

	Source      string // ID of the faction posting the job
	SourceEthos *Ethos // ethos of source faction
	SourceArea  string // where people will be recruited from

	Action  ActionType // action that is due to take place
	Secrecy int        // the result of an espionage defence roll (if covert)

	Target      string // ID of the target faction (if any)
	TargetEthos *Ethos // ethos of target faction (if any)
	TargetArea  string // where the action will take place

	TargetMetaKey string // key/val pair to hold adv. target metadata
	TargetMetaVal string

	PeopleMin int // required min number of people (else job fails to kick off)
	PeopleMax int // max number of people that can work this (if any)

	TickCreated  int // when the job was created
	TickStarts   int // when the job is due to start
	TickDuration int // how long the job should last

	IsIllegal bool // action has been outlawed
	IsReady   bool // action has all people it needs to start
	IsDone    bool // job is completed
}
