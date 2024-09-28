package structs

type Object interface {
	GetId() string
	SetId(v string)
	GetEtag() string
	SetEtag(v string)
	SetWorld(v string)
}

type Interactive interface {
	SetAllies(v []*Relationship)
	SetEnemies(v []*Relationship)
	GetAllies() []*Relationship
	GetEnemies() []*Relationship
}
