package v1

// ActionSelection provides a way to weight the likelihood of various actions being taken.
type ActionSelection struct {
	// Tags that this user has applied either to self (if target-id is this) or others.
	// A Tag that would overwrite an existing tag before it's TTL expires is ignored.
	//
	// Eg. If some target is tagged with ["hostile"] then this user might be more likely to select
	// actions that are tagged with ["hostile"].
	//
	// map[target-id] -> map[tag] -> Tag
	//
	// That is, Tags helps us keep track of what the subject is thinking, goals it has in mind etc.
	Tags map[string]map[string]Tag `yaml:"Tags" json:"Tags" validate:"min=0,max=1000,dive,keys,valid_id,endkeys,dive,keys,alphanum,endkeys,dive"`

	// Actions map action id -> weight
	//
	// This decides how likely various actions are to be chosen and what tags (if any) can
	// influence their selection (and if so by how much).
	Actions map[string]ActionWeight `yaml:"Actions" json:"Actions" validate:"min=0,max=100,dive,keys,alphanum,endkeys,dive"`

	// Hook is some callback that takes the subject (eg. faction, actor) and returns
	// a set of Tags that should be applied to the subject.
	Hook Hook `yaml:"Hook" json:"Hook"`
}

type Hook struct {
	// Protocol is the protocol to use for the hook where '' indicates a local builtin
	Protocol string `yaml:"Protocol" json:"Protocol" validate:"oneof=http ''"`

	// Address is the address to send the hook to, or the name of a builtin.
	// If not set a default builtin will be used.
	Address string `yaml:"Address" json:"Address"`
}

type Tag struct {
	// TTL is the time to live for this tag in ticks
	TTL int `yaml:"TTL" json:"TTL" validate:"gte=0,lte=1000000000"`

	// Metadata specific to this tag
	Metadata map[string]interface{} `yaml:"Metadata" json:"Metadata" validate:"min=0,max=100,dive,keys,alphanumsymbol,endkeys`
}

type ActionWeight struct {
	// Probability (base) of this action being selected.
	Probability float64 `yaml:"Probability" json:"Probability" validate:"gte=0,lte=1"`

	// Tags are user set strings that can be used to influnce the probability of this action.
	// Eg. {"growth": 1.5} makes this action 1.5 times more likely if there is a 'growth' tag on the user.
	Tags map[string]float64 `yaml:"Tags" json:"Tags" validate:"min=1,max=100,dive,keys,alphanum,endkeys,gte=0,lte=100"`
}
