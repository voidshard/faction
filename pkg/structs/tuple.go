package structs

import "fmt"

const (
	MaxTuple = 10000
	MinTuple = -10000
)

// Tuple is a way of storing a value for a generic a->b relationship,
// or lists of items.
//
// Ie. we use this to store
// - affiliation of Faction A to Faction B
// - favor of faction A to Faction B
// - trust of faction A to Faction B
// - .. all the same of Person A to Faction B
// - .. similar for Person A to Person B
// - family relations
// So uh, a lot of things really.
//
// A specific table represents the kind of relation we're talking about;
// (hence we use the Relation in the table name).
// Person->Person
// Faction->Faction
// etc.
type Tuple struct {
	Subject string `db:"subject"`
	Object  string `db:"object"`
	Value   int    `db:"value"`
}

func (t *Tuple) ObjectID() string {
	return fmt.Sprintf("%s-%s", t.Subject, t.Object)
}

// Modifier is some temporary modifier to a Tuple(s).
// Eg. A Tuple might be:
//
//	[Table:Affiliation]
//	Subject (Faction A), Object (Faction B), Value 10
//
// With a Modifier:
//
//	[Table:Affiliation]
//	Subject (Faction A), Object (Faction B), Value -5
//	Expires 100 (Tick)
//	MetaID (ActionID), MetaType (Action), MetaReason ("Attempted infiltration of B leadership")
//
// Tl;dr when we look up the value of the Tuple (Affiliation, A, B) we apply non expired modifiers.
type Modifier struct {
	Tuple

	TickExpires int     `db:"tick_expires"` // modifiers expire at some Tick
	MetaKey     MetaKey `db:"meta_key"`     // information about what MetaVal refers to
	MetaVal     string  `db:"meta_val"`     // ID of what caused this modifier
	MetaReason  string  `db:"meta_reason"`  // human readable reason string (eg. "bribe")
}
