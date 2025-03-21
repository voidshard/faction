package v1

type Actor struct {
	Meta `json:",inline" yaml:",inline"`

	Firstname string `json:"Firstname" yaml:"Firstname" validate:"alphanum-or-empty"`
	Lastname  string `json:"Lastname" yaml:"Lastname" validate:"alphanum-or-empty"`

	Race    string `json:"Race" yaml:"Race" validate:"alphanum"`
	Culture string `json:"Culture" yaml:"Culture" validate:"alphanum"`

	Area string `json:"Area" yaml:"Area" validate:"uuid4-or-empty"`

	Ethos       map[string]float64              `json:"Ethos" yaml:"Ethos" validate:"max=100,dive,keys,alphanum,endkeys,gte=0,lte=100"`
	Professions map[string]ActorValueProfession `json:"Professions" yaml:"Professions" validate:"max=100,dive,keys,alphanum,endkeys"`
	Ranks       map[string]ActorValueRank       `json:"Ranks" yaml:"Ranks" validate:"max=100,dive,keys,alphanum,endkeys"`
}

type ActorValueProfession struct {
	Name  string `json:"Name" yaml:"Name" validate:"alphanum"`
	Level int    `json:"Level" yaml:"Level" validate:"gte=0,lte=1000000"`
}

type ActorValueRank struct {
	Name  string `json:"Name" yaml:"Name" validate:"alphanum"`
	Level int    `json:"Level" yaml:"Level" validate:"gte=0,lte=100"`
}

func (x *Actor) New(in interface{}) (Object, error) {
	i := &Actor{}
	err := unmarshalObject(in, i)
	i.Kind = "actor"
	return i, err
}

/*
func (x *Actor) SetAllies(v []*Relationship) {
	x.Allies = v
}

func (x *Actor) SetEnemies(v []*Relationship) {
	x.Enemies = v
}
*/
