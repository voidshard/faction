package v1

// Query is a query to search for results.
type Query struct {
	// Filter results based on object fields
	Filter `json:",inline" yaml:",inline"`

	// Score is used to rank results returned from Filter(s).
	Score []Score `yaml:"Score" json:"Score" validate:"min=0,max=100,dive"`

	// RandomWeight adds randomness in the scoring of results.
	RandomWeight float64 `yaml:"RandomWeight" json:"RandomWeight" validate:"gte=0,lte=100"`

	// Limit is the number of results to return.
	Limit int64 `yaml:"Limit" json:"Limit" validate:"gte=1,lte=5000"`

	// Kind is the kind of object to search for.
	Kind string `yaml:"Kind" json:"Kind" validate:"alphanum"`
}

func NewQuery() *Query {
	return &Query{
		Filter: Filter{
			All: []Match{},
			Any: []Match{},
			Not: []Match{},
		},
		Score: []Score{},
	}
}

// Match is a determines what a Query will return as a result.
type Filter struct {
	// All is a list of filters that must all be true for the query to match a result.
	All []Match `yaml:"All" json:"All" validate:"required,min=0,max=100,dive"`

	// Any is a list of filters where any one of them can be true for the query to match a result.
	Any []Match `yaml:"Any" json:"Any" validate:"min=0,max=100,dive"`

	// Not is a list of filters that must all be false for the query to match a result.
	Not []Match `yaml:"Not" json:"Not" validate:"min=0,max=100,dive"`
}

// Match is a comparison to make on a field.
type Match struct {
	// Field is the field to compare. This is a lowercased dot-separated flattened path to the field.
	Field string `yaml:"Field" json:"Field" validate:"alphanumsymbol"`

	// Op is the operation to perform on the field. Where "" is 'eq' (equals) as a default.
	Op string `yaml:"Op" json:"Op" validate:"oneof=eq '' lt gt"`

	// Value is the value to compare the field to.
	Value interface{} `yaml:"Value" json:"Value" validate:"required"`
}

// Score is a weight to apply to a match.
type Score struct {
	Match  `yaml:",inline" json:",inline"`
	Weight float64 `yaml:"Weight" json:"Weight" validate:"gte=0,lte=100"`
}
