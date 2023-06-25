/*filters.go contains generic query / filter implementations*/
package db

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

	// person
	PreferredProfession Field = "preferred_profession"
	PreferredFactionID  Field = "preferred_faction_id"
	IsChild             Field = "is_child"
	BirthFamilyID       Field = "birth_family_id"

	// modifiers
	TickExpires Field = "tick_expires"

	// family
	IsChildBearing Field = "is_child_bearing"
	PregnancyEnd   Field = "pregnancy_end"

	// job
	SourceFactionID Field = "source_faction_id"
	TargetMetaKey   Field = "target_meta_key"
	TargetMetaVal   Field = "target_meta_val"
	Secrecy         Field = "secrecy"
	State           Field = "state"

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
