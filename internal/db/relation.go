package db

import (
	"fmt"
)

// Relation here is a name we accept for Tuple(s) and Modifier(s).
type Relation string

const (
	// ## Relation Table Partitioning
	//
	// We could in theory have all Tuples and Modifiers in two giant tables
	// (possible even the same table where Tuples don't use all of the columns)
	// and partition by the Relation (in the case of Postgres at least?) but
	// in practice we disallow querying the table directly in our interface and
	// it's sort of easier conceptually / book keeping wise to have distinct tables
	// which hold distinct data.
	//
	// Ie. an inter-personal relationship (wife / son / ally) between Person X and
	// Person Y has no bearing on the trust between Faction A and Faction B .. and
	// we don't to be able to query them in the same breath .. although they're both
	// tuples and in theory we *could*
	//
	// Thus currently each Relation has a "tuples_" and a "modifiers_" table created
	// for it. Modifiers apply temporary alterations to Tuple values until their expire
	// time is reached.
	//
	//
	// ## Compute Tuples
	//
	// Given Tuple from tuples_NAME
	//    Subject (A), Object (B), Value (V)
	// Modifier(s) from modifiers_NAME
	//    Subject (A), Object (B), Value (X), Expires 10
	//    Subject (A), Object (B), Value (Y), Expires 12
	//    Subject (A), Object (B), Value (Z), Expires 15
	// We compute the final value based on tick T as
	//    [Time]       =     [Final Value]
	//    T < 10       =     V + X + Y + Z
	//    10 <= T < 12 =     V + Y + Z
	//    12 <= T < 15 =     V + Z
	//    15 <= T      =     V
	//
	// Thus ComputeTuples() accepts;
	// - the relation in question
	// - iter token
	// - current time tick (any modifier with this time or less is not considered)
	// - tuple filters (which tuples you're interested in)

	// RelationPFAffiliation holds how closely people affiliate with factions
	RelationPFAffiliation Relation = "affiliation_person_to_faction"

	// RelationFFTrust holds how much factions trust each other
	RelationFFTrust Relation = "trust_faction_to_faction"

	// RelationPProfessionSkill holds how skilled people are in professions
	RelationPProfessionSkill Relation = "skill_person_to_profession"
)

func (r Relation) tupleTable() string {
	return fmt.Sprintf("tuples_%s", r)
}

func (r Relation) modTable() string {
	return fmt.Sprintf("modifiers_%s", r)
}
