package structs

// Tuple is a way of storing a value for a generic a->b relationship.
// Ie. we use this to store
// - affiliation of Faction A to Faction B
// - favor of faction A to Faction B
// - trust of faction A to Faction B
// - .. all the same of Person A to Faction B
// - .. similar for Person A to Person B
// - family relations
// So uh, a lot of things really.
//
// Table here represents the kind of relation we're talking about;
// Person->Person
// Faction->Faction
// etc.
type Tuple struct {
	Subject string `db:"subject"`
	Object  string `db:"object"`
	Value   int    `db:"value"`
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

	Expires    int    `db:"expires"`     // modifiers expire at some Tick
	MetaType   string `db:"meta_type"`   // information about what MetaID refers to
	MetaID     string `db:"meta_id"`     // ID of what caused this modifier
	MetaReason string `db:"meta_reason"` // human readable reason string (eg. "bribe")
}
