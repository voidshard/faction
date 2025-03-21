package v1

// Meta is the base struct for all objects supplying common metadata fields
type Meta struct {
	Id         string `json:"_id" yaml:"_id" validate:"valid_id"`
	Etag       string `json:"_etag" yaml:"_etag" validate:"uuid4-or-empty"`
	Kind       string `json:"_kind" yaml:"_kind" validate:"alphanum,required"`
	Controller string `json:"_controller" yaml:"_controller" validate:"alphanum-or-empty"`

	World string `json:"World" yaml:"World" validate:"alphanum-if-non-global"`

	Labels     map[string]string  `json:"Labels" yaml:"Labels" validate:"max=500,dive,keys,alphanumsymbol,endkeys"`
	Attributes map[string]float64 `json:"Attributes" yaml:"Attributes" validate:"max=500,dive,keys,alphanumsymbol,endkeys"`
}

func (x *Meta) GetKind() string {
	return x.Kind
}

func (x *Meta) GetId() string {
	return x.Id
}

func (x *Meta) SetId(v string) {
	x.Id = v
}

func (x *Meta) GetEtag() string {
	return x.Etag
}

func (x *Meta) SetEtag(v string) {
	x.Etag = v
}

func (x *Meta) GetWorld() string {
	return x.World
}

func (x *Meta) SetWorld(v string) {
	x.World = v
}

func (x *Meta) GetLabels() map[string]string {
	return x.Labels
}

func (x *Meta) GetAttributes() map[string]float64 {
	return x.Attributes
}

func (x *Meta) GetController() string {
	if x.Controller == "" {
		return "default"
	}
	return x.Controller
}

// SetMeta is used to set metadata on given objects.
//
// Ie. On a CultureCaste we set these on all Actors using this Caste.
type SetMeta struct {
	Labels     map[string]string  `json:"Labels" yaml:"Labels" validate:"max=50"`
	Attributes map[string]float64 `json:"Attributes" yaml:"Attributes" validate:"max=50"`
	Controller string             `json:"Controller" yaml:"Controller" validate:"alphanum-or-empty"`
}
