package v1

// Distribution represents a random distribution of values.
type Distribution struct {
	Min       float64 `json:"Min" yaml:"Min" validate:"ltfield=Mean"`
	Max       float64 `json:"Max" yaml:"Max" validate:"gtfield=Mean"`
	Mean      float64 `json:"Mean" yaml:"Mean" validate:"gtfield=Min,ltfield=Max"`
	Deviation float64 `json:"Deviation" yaml:"Deviation" validate:"gte=0"`
}
