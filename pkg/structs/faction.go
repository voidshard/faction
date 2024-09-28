package structs

func (x *Faction) SetId(v string) {
	x.Id = v
}

func (x *Faction) SetEtag(v string) {
	x.Etag = v
}

func (x *Faction) SetWorld(v string) {
	x.World = v
}

func (x *Faction) SetAllies(v []*Relationship) {
	x.Allies = v
}

func (x *Faction) SetEnemies(v []*Relationship) {
	x.Enemies = v
}
