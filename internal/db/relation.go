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

	// RelationFactionCommodityImport holds how much a faction wishes to import a commodity
	// So: <faction_id> <commodity> <import_level>
	RelationFactionCommodityImport Relation = "import_faction_to_commodity"

	// RelationFactionCommodityExport holds how much a faction wishes to export a commodity
	// So: <faction_id> <commodity> <export_level>
	RelationFactionCommodityExport Relation = "export_faction_to_commodity"

	// RelationFactionTopicResearch holds how much research a faction has done on a topic
	// (we use the ActionType as the object)
	// So: <faction_id> <topic> <research_level>
	RelationFactionTopicResearch Relation = "research_faction_to_topic"

	// RelationFactionTopicResearchWeight holds how much weight a faction gives to a topic.
	// Where RelationFactionTopicResearch is the actual research done, this is how likely a
	// faction is to research this topic over other(s) (if any) when performing research.
	// So: <faction_id> <topic> <weight_level>
	RelationFactionTopicResearchWeight Relation = "weight_faction_to_topic"

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

	// RelationFactionActionTypeWeight holds how much weight a faction gives to an action type.
	// This is used to weight the proability of a faction performing an action.
	// Tuples are limited to MAX_TUPLE to MIN_TUPLE, currently -10k to 10k so each point here
	// is a 1/100th (0.01) of a percent increase / decrease in likelihood.
	// So: <faction_id> <action_type> <weight_level>
	RelationFactionActionTypeWeight Relation = "weight_faction_to_action_type"

	// RelationFactionProfessionWeight holds how much a faction desires a profession.
	// This is based on the factions owned land and favored actions.
	RelationFactionProfessionWeight Relation = "weight_faction_to_profession"

	// RelationPersonFactionRank holds the rank of a person in a faction.
	// If not given people are assumed to be at 0 (ie. "not a member")
	// So: <person_id> <faction_id> <rank>
	RelationPersonFactionRank Relation = "rank_person_to_faction"

	// RelationFactionFactionIntelligence holds how much intelligence a faction has on another faction.
	// So: <faction_id> <faction_id> <intelligence_level>
	RelationFactionFactionIntelligence Relation = "intelligence_faction_to_faction"
)

var (
	allRelations = []Relation{
		RelationFactionFactionTrust,
		RelationFactionCommodityImport,
		RelationFactionCommodityExport,
		RelationFactionTopicResearch,
		RelationFactionTopicResearchWeight,
		RelationPersonFactionAffiliation,
		RelationPersonProfessionSkill,
		RelationPersonPersonRelationship,
		RelationPersonReligionFaith,
		RelationPersonPersonTrust,
		RelationFactionActionTypeWeight,
		RelationFactionProfessionWeight,
		RelationPersonFactionRank,
		RelationFactionFactionIntelligence,
	}
)

func (r Relation) tupleTable() string {
	return fmt.Sprintf("tuples_%s", r)
}

func (r Relation) modTable() string {
	return fmt.Sprintf("modifiers_%s", r)
}

// SupportsModifiers returns if the relation supports modifiers,
// ie. the Tuple has a matching Modifer table.
func (r Relation) SupportsModifiers() bool {
	// modifiers complicate queries & add calculations but are a nice way of adding
	// slow burn buffs / debuffs.
	// In general we only support these on super important tuples (where we have
	// a use case) since it requires us to track & query a lot more variables.
	switch r {
	case RelationPersonFactionAffiliation:
		return true
	case RelationFactionActionTypeWeight:
		return true
	case RelationFactionFactionTrust:
		return true
	case RelationPersonPersonTrust:
		return true
	}
	return false
}
