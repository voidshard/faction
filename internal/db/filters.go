/*filters.go contains generic query / filter implementations*/
package db

import (
	"fmt"
	"strings"
)

type Field string
type Op string

type isValid func(interface{}) bool

const (
	Equal    Op = "="
	NotEqual Op = "<>"
	In       Op = "in"
	Greater  Op = ">"
	Less     Op = "<"
)

const (
	// pretty much everything
	ID Field = "id"

	JobID  Field = "job_id"
	AreaID Field = "area_id"

	// event
	Type Field = "type"
	Tick Field = "tick"

	// area
	Biome Field = "biome"

	// person & family
	Race    Field = "race"
	Culture Field = "culture"

	// area, faction
	GovernmentID Field = "government_id"

	// area
	Commodity Field = "commodity"

	// plot
	FactionID Field = "faction_id"

	// route, job
	SourceAreaID Field = "source_area_id"
	TargetAreaID Field = "target_area_id"

	// faction, person, family
	EthosAltruism   Field = "ethos_altruism"
	EthosAmbition   Field = "ethos_ambition"
	EthosTradition  Field = "ethos_tradition"
	EthosPacificism Field = "ethos_pacificism"
	EthosPiety      Field = "ethos_piety"
	EthosCaution    Field = "ethos_caution"

	// faction
	ParentFactionID       Field = "parent_faction_id"
	ParentFactionRelation Field = "parent_faction_relation"

	// person
	PreferredProfession Field = "preferred_profession"
	PreferredFactionID  Field = "preferred_faction_id"
	AdulthoodTick       Field = "adulthood_tick"
	BirthFamilyID       Field = "birth_family_id"
	BirthTick           Field = "birth_tick"
	DeathTick           Field = "death_tick"
	NaturalDeathTick    Field = "natural_death_tick"
	Random              Field = "random"

	// modifiers
	TickExpires Field = "tick_expires"

	// family
	IsChildBearing Field = "is_child_bearing"
	PregnancyEnd   Field = "pregnancy_end"
	MaleID         Field = "male_id"
	FemaleID       Field = "female_id"

	// job
	SourceFactionID Field = "source_faction_id"
	TargetFactionID Field = "target_faction_id"
	TargetMetaKey   Field = "target_meta_key"
	TargetMetaVal   Field = "target_meta_val"
	Secrecy         Field = "secrecy"
	JobState        Field = "state"
	TickEnds        Field = "tick_ends"
	ActionType      Field = "action_type"

	// tuple, modifiers
	Subject Field = "subject"
	Object  Field = "object"
)

type Query struct {
	sort    bool // by ID
	filters [][]*Filter
}

func Q(filters ...*Filter) *Query {
	return &Query{filters: [][]*Filter{filters}, sort: true}
}

func (q *Query) String() string {
	ors := make([]string, len(q.filters))
	for i, fset := range q.filters {
		ands := make([]string, len(fset))
		for j, f := range fset {
			ands[j] = fmt.Sprintf("\t\t%s", f.String())
		}
		ors[i] = strings.Join(ands, "\n\t\tAND\n ")
	}
	return fmt.Sprintf("Query:\n%s", strings.Join(ors, "\n\tOR\n "))
}

func (q *Query) DisableSort() *Query {
	q.sort = false
	return q
}

func (q *Query) Or(filters ...*Filter) *Query {
	q.filters = append(q.filters, filters)
	return q
}

type Filter struct {
	Field Field
	Op    Op
	Value interface{}
}

func F(field Field, op Op, value interface{}) *Filter {
	return &Filter{
		Field: field,
		Op:    op,
		Value: value,
	}
}

func (f *Filter) String() string {
	return fmt.Sprintf("%s %s %v", f.Field, f.Op, f.Value)
}
