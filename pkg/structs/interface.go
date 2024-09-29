package structs

// Object is an interface that all objects must implement.
//
// Note that the "Get" and "String" funcs are written by gRPC / protobuf so
// we only need to implement the "Set" funcs.
type Object interface {
	GetId() string
	SetId(v string)
	GetEtag() string
	SetEtag(v string)
	SetWorld(v string)
	String() string
}

type Interactive interface {
	SetAllies(v []*Relationship)
	SetEnemies(v []*Relationship)
	GetAllies() []*Relationship
	GetEnemies() []*Relationship
}
