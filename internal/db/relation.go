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

	// RelationPFAffiliation holds how closely people affiliate with factions
	// So: <person_id> <faction_id> <affiliation_level>
	RelationPersonFactionAffiliation Relation = "affiliation_person_to_faction"

	// RelationFFTrust holds how much factions trust each other
	// So: <faction_id> <faction_id> <trust_level>
	RelationFactionFactionTrust Relation = "trust_faction_to_faction"

	// RelationFactionTopicResearch holds how much research a faction has done on a topic
	// (we use the ActionType as the object)
	// So: <faction_id> <topic> <research_level>
	RelationFactionTopicResearch Relation = "research_faction_to_topic"

	// RelationPProfessionSkill holds how skilled people are in professions
	// So: <person_id> <profession> <skill_level>
	RelationPersonProfessionSkill Relation = "skill_person_to_profession"

	// RelationPersonPersonRelationships holds how people relate to each other
	// So: <person_id> <person_id> <PersonalRelation> (see structs/family.go)
	RelationPersonPersonRelationship Relation = "relationship_person_to_person"

	// RelationPersonReligionFaith holds how much faith a person has in a religion
	// So: <person_id> <religion_id> <faith_level>
	RelationPersonReligionFaith Relation = "faith_person_to_religion"

	// RelationPersonPersonTrust holds how much trust a person has in another person
	// So: <person_id> <person_id> <trust_level>
	RelationPersonPersonTrust Relation = "trust_person_to_person"
)

var (
	allRelations = []Relation{
		RelationFactionFactionTrust,
		RelationFactionTopicResearch,
		RelationPersonFactionAffiliation,
		RelationPersonProfessionSkill,
		RelationPersonPersonRelationship,
		RelationPersonReligionFaith,
		RelationPersonPersonTrust,
	}
)

func (r Relation) tupleTable() string {
	return fmt.Sprintf("tuples_%s", r)
}

func (r Relation) modTable() string {
	return fmt.Sprintf("modifiers_%s", r)
}

func (r Relation) supportsModifiers() bool {
	// modifiers complicate queries & add calculations but are a nice way of adding
	// slow burn buffs / debuffs.
	// In general we only support these on super important tuples.
	switch r {
	case RelationPersonFactionAffiliation:
		return true
	case RelationFactionFactionTrust:
		return true
	case RelationPersonPersonTrust:
		return true
	}
	return false
}
