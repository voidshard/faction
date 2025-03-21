package v1

// Race imparts information about the species / race / breed of a being.
// These are physical limits & traits that are not simply cultural.
//
// Ie. a human can't live underwater, but a mermaid can, a human baby takes
// 9 months to be born, but a kitten takes 2.
//
// The companion struct 'Culture' is meant to impart information about the
// culture someone lives in.
// Ie. a human female might be capable of having a child at age 13, but culturally
// it might be considered odd to have a child before 18 (or whatever).
type Race struct {
	Meta `yaml:",inline" json:",inline"`

	// Base set default values for all Castes, values set in Castes overwrite these
	// where conflicts arise.
	Base RaceCaste `yaml:"Default" json:"Default" validate:"required,dive"`

	// 'Castes' are physically distinct types of this Race.
	// It might cover a race of insects where there are 'royal' 'soldier' 'worker' castes
	// or be more straight forward genders (or a single Asexual caste).
	Castes map[string]RaceCaste `yaml:"Castes" json:"Castes" validate:"min=1,max=100,dive,keys,alphanum,endkeys"`
}

func (x *Race) New(in interface{}) (Object, error) {
	i := &Race{}
	err := unmarshalObject(in, i)
	i.Kind = "race"
	return i, err
}

type RaceCaste struct {
	// Meta sets values on all Actors using this Caste
	Set SetMeta `yaml:"Set" json:"Set" validate:"required,dive"`

	// Probability of someone being born into this caste
	Probability float64 `yaml:"Probability" json:"Probability" validate:"gte=0,lte=1"`

	// Lifespan of someone in this caste
	Lifespan Distribution `yaml:"Lifespan" json:"Lifespan" validate:"dive"`
}
