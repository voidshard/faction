package v1

type IntentType string

const (
	Allied     IntentType = "allied"
	Friendly   IntentType = "friendly"
	Neutral    IntentType = "neutral"
	Suspicious IntentType = "suspicious"
	Unfriendly IntentType = "unfriendly"
	Hostile    IntentType = "hostile"
)

// Action is a user defined action that can be taken by a faction, family or actor.
//
// An action is a loose definition of what is involved in performing some large scale faction
// activity. If a faction elects to perform an action a "job" is spawned using the actions definition.
type Action struct {
	Meta `yaml:",inline" json:",inline"`

	// Base sets defaults for each task under steps. Values set in steps overwrite these.
	Base Task `yaml:"Base" json:"Base" validate:"dive"`

	// Steps are the steps that make up this action. Order is important.
	Steps []Step `yaml:"Steps" json:"Steps" validate:"min=1,max=100,dive"`

	// TimeLimit is the maximum number of ticks this action can take to complete.
	// If a job is still running after this many ticks it will be failed.
	TimeLimit int `yaml:"TimeLimit" json:"TimeLimit" validate:"gte=0,lte=1000000000"`
}

type Step struct {
	// Name, mostly for human users.
	Name string `yaml:"Name" json:"Name"`

	// Ticks is the number of ticks this step takes to complete.
	Ticks int `yaml:"Ticks" json:"Ticks" validate:"gte=0,lte=1000000000"`

	// Tasks that will happen as part of this action
	Tasks []Task `yaml:"Tasks" json:"Tasks" validate:"min=1,max=1000,dive"`
}

type Progress struct {
	// How much 'progress' this task takes to complete.
	CompletedAt Distribution `yaml:"CompletedAt" json:"CompletedAt" validate:"dive"`

	// Each tick an actor contributes this much 'progress' towards completing the task.
	// This forms a minimum usefulness per actor helping on the task.
	YieldPerActor float64 `yaml:"YieldPerActor" json:"YieldPerActor" validate:"gte=0,lte=1000000"`

	// Each tick an actor contributes (in addition) this much 'progress' towards completing the task
	// where the value added is
	// progress.YieldByProfession[profession] * actor.Professions[profession].Level
	//
	// Eg.
	// progress.YieldByProfession["miner"] = 0.1
	// actor.Professions["miner"].Level = 50
	// actor will contribute 0.1 * 50 = 5 progress towards the task each tick.
	YieldByProfession map[string]float64 `yaml:"YieldByProfession" json:"YieldByProfession" validate:"min=1,max=1000000,dive,keys,alphanum,endkeys"`

	// Secrecy is a roll to determine how clandestine the task execution is.
	// In 'Support' it indicates if the target even notices.
	// In 'Oppose' it indicates if the target can identify the source.
	Secrecy Distribution `yaml:"Secrecy" json:"Secrecy" validate:"dive"`
}

type Task struct {
	// Support details how the source faction supports this task. How much progress
	// actors can help each tick, what professions matter etc.
	Support Progress `yaml:"Progress" json:"Progress" validate:"required,dive"`

	// Oppose details how a target faction can oppose this task (the inverse of Support).
	// Assuming it is not a friendly action.
	//
	// If the opposition completes the task before the support does the task will fail.
	// Nb. even for unfriendly actions, not all tasks may have opposition defined.
	Oppose *Progress `yaml:"Opposition" json:"Opposition" validate:"dive"`

	// Cost is how much it costs to perform this task. This is deducted from the faction / family / actor
	// that is performing the task. An action will fail if the cost cannot be met.
	Cost Distribution `yaml:"Cost" json:"Cost" validate:"required,dive"`

	// TargetArea of this task.
	// If not set the default target is the HQ area the faction performing the task is in.
	//
	// Query Fields and Values in this query may use the following template variables:
	//
	// $ACTOR - the name of the actor doing the task
	// $FACTION - the faction the actor belongs to performing this task
	// $TARGET_FACTION - the faction that this action is targeting (this may be the same as $FACTION)
	TargetArea *Query `yaml:"TargetArea" json:"TargetArea" validate:"required,dive"`

	// Target of this task.
	// If not set the default target is the faction performing the task.
	//
	// Query Fields and Values in this query may use the following template variables:
	//
	// $ACTOR - the name of the actor doing the task
	// $AREA - the area the task is taking place in (chosen above)
	// $FACTION - the faction the actor belongs to performing this task
	// $TARGET_FACTION - the faction that this action is targeting
	Target *Query `yaml:"Target" json:"Target" validate:"dive"`

	// Who should do the task. Query Fields and Values in this query may use the following template variables:
	//
	// $ACTOR - the name of the actor doing the task
	// $AREA - the area the task is taking place in (chosen above)
	// $FACTION - the faction the actor belongs to performing this task
	// $TARGET - the target of the task (chosen above)
	// $TARGET_FACTION - the faction that this action is targeting
	Actors Query `yaml:"Actors" json:"Actors" validate:"required,dive"`

	// If set, this the 'Actors' query (above) will have to return at least this many
	// actors for the task to be valid. Otherwise our action will fail.
	ActorsRequired *Distribution `yaml:"ActorsRequired" json:"ActorsRequired" validate:"dive"`

	// Descriptions of the task. May use the following template variables:
	//
	// $ACTOR - the name of the actor doing the task
	// $AREA - the area the task is taking place in (chosen above)
	// $FACTION - the faction the actor belongs to performing this task
	// $TARGET - the target of the task (chosen above)
	// $TARGET_FACTION - the faction the target belongs to
	PastTenseSuccess []string `yaml:"PastTenseSuccess" json:"PastTenseSuccess" validate:"min=1,max=100,dive"`
	PastTenseFail    []string `yaml:"PastTenseFail" json:"PastTenseFail" validate:"min=1,max=100,dive"`
	PresentTense     []string `yaml:"PresentTense" json:"PresentTense" validate:"min=1,max=100,dive"`
	FutureTense      []string `yaml:"FutureTense" json:"FutureTense" validate:"min=1,max=100,dive"`

	// Type indicates how the action should be perceived by other factions.
	Type string `yaml:"Type" json:"Type" validate:"required,oneof=allied friendly neutral suspicious unfriendly hostile"`
}
