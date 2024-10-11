package structs

func (x *World) SetId(v string) {
	x.Id = v
}

func (x *World) SetEtag(v string) {
	x.Etag = v
}

func (x *World) SetWorld(v string) {}

func (x *World) GetWorld() string {
	return x.Id // duh!
}
