package v1

type World struct {
	Meta `json:",inline" yaml:",inline"`

	Tick uint64 `json:"Tick" yaml:"Tick" validate:"gte=0"`
}

func (x *World) New(in interface{}) (Object, error) {
	i := &World{}
	err := unmarshalObject(in, i)
	i.Kind = "world"
	return i, err
}

func (x *World) GetWorld() string {
	return x.GetId()
}

func (x *World) SetWorld(v string) {
	x.SetId(v)
}
