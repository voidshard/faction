package faction

// Tuple is a way of storing a value for a generic a->b relationship.
// Ie. we use this to store
// - affiliation of Faction A to Faction B
// - favor of faction A to Faction B
// - trust of faction A to Faction B
// - .. all the same of Person A to Faction B
// - .. similar for Person A to Person B
// So uh, a lot of things really.
//
// Table here represents the kind of relation we're talking about;
// Person->Person
// Faction->Faction
// etc.
type Tuple struct {
	Subject string
	Object  string
	Value   int
}
