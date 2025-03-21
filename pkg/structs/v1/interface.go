package v1

// Object is the base interface for all objects
type Object interface {
	New(interface{}) (Object, error)

	GetKind() string
	GetId() string
	GetEtag() string
	GetWorld() string
	GetController() string
	GetLabels() map[string]string
	GetAttributes() map[string]float64

	SetId(v string)
	SetEtag(v string)
	SetWorld(v string)
}

/*
type Relationship struct{} // TEMP -> possibly remove later

type Interactive interface {
	SetAllies(v []*Relationship)
	SetEnemies(v []*Relationship)
	GetAllies() []*Relationship
	GetEnemies() []*Relationship
}
*/
