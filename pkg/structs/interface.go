package structs

// Marshalable defines how we can convert objects to and from YAML and JSON.
type Marshalable interface {
	MarshalYAML() ([]byte, error)
	UnmarshalYAML([]byte) error

	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// Object is an interface that all objects must implement.
//
// Note that the "Get" and "String" funcs are written by gRPC / protobuf so
// we only need to implement the "Set" funcs.
type Object interface {
	Marshalable

	Kind() string
	GetId() string
	SetId(v string)
	GetEtag() string
	SetEtag(v string)
	SetWorld(v string)
	GetWorld() string
	String() string
}

type Interactive interface {
	Object

	SetAllies(v []*Relationship)
	SetEnemies(v []*Relationship)
	GetAllies() []*Relationship
	GetEnemies() []*Relationship
}
