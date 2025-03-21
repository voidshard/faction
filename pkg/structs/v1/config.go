package v1

type Config struct {
	Meta `json:",inline" yaml:",inline"`

	Data map[string]string `json:"Data" yaml:"Data"`
}

func (x *Config) New(in interface{}) (Object, error) {
	i := &Config{}
	err := unmarshalObject(in, i)
	i.Kind = "config"
	return i, err
}
