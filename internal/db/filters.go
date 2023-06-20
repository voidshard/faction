package db

import (
	"fmt"
	"strings"

	"github.com/voidshard/faction/internal/dbutils"
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

	// area, faction
	GovernmentID Field = "government_id"

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

	// modifiers
	TickExpires Field = "tick_expires"

	// family
	IsChildBearing Field = "child_bearing"
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

var (
	fieldTypes = map[Op][]isValid{
		Equal:    {isInt, isString},
		NotEqual: {isInt, isString},
		In:       {isListID},
		Greater:  {isInt},
		Less:     {isInt},
	}
)

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

func (f *Filter) sqlQuery(offset int) (string, []interface{}, error) {
	if f.Field == "" {
		return "", nil, fmt.Errorf("invalid filter field: %s", f.Field)
	}
	if f.Op == "" {
		return "", nil, fmt.Errorf("invalid filter op: %s", f.Op)
	}
	if f.Value == nil {
		return "", nil, fmt.Errorf("invalid filter value: %s", f.Value)
	}

	checks, ok := fieldTypes[f.Op]
	if !ok {
		return "", nil, fmt.Errorf("invalid filter op: %s", f.Op)
	}
	valid := false
	for _, check := range checks {
		valid = check(f.Value)
		if valid {
			break
		}
	}
	if !valid {
		return "", nil, fmt.Errorf("invalid filter value: %v", f.Value)
	}

	args := []interface{}{}
	placeholder := fmt.Sprintf("$%d", offset+1)

	if f.Op == In {
		values := f.Value.([]string) // we've already checked the type
		placeholders := []string{}
		for i, v := range values {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+offset+1))
			args = append(args, v)
		}
		placeholder = fmt.Sprintf("(%s)", strings.Join(placeholders, ","))
	} else {
		args = append(args, f.Value)
	}

	return fmt.Sprintf("%s %s %s", f.sqlColumn(), f.Op, placeholder), args, nil
}

func (f *Filter) sqlColumn() string {
	switch f.Field {
	case ID:
		return "id"
	case JobID:
		return "job_id"
	case AreaID:
		return "area_id"
	case GovernmentID:
		return "government_id"
	case FactionID:
		return "faction_id"
	case SourceAreaID:
		return "source_area_id"
	case TargetAreaID:
		return "target_area_id"
	case EthosAltruism:
		return "ethos_altruism"
	case EthosAmbition:
		return "ethos_ambition"
	case EthosTradition:
		return "ethos_tradition"
	case EthosPacificism:
		return "ethos_pacificism"
	case EthosPiety:
		return "ethos_piety"
	case EthosCaution:
		return "ethos_caution"
	case TickExpires:
		return "tick_expires"
	case IsChildBearing:
		return "child_bearing"
	case PregnancyEnd:
		return "pregnancy_end"
	case SourceFactionID:
		return "source_faction_id"
	case TargetMetaKey:
		return "target_meta_key"
	case TargetMetaVal:
		return "target_meta_val"
	case Secrecy:
		return "secrecy"
	case State:
		return "state"
	case Subject:
		return "subject"
	case Object:
		return "object"
	}
	return ""
}

func isInt(v interface{}) bool {
	_, ok := v.(int)
	return ok
}

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

func isID(v interface{}) bool {
	i, ok := v.(string)
	if !ok {
		return false
	}
	return dbutils.IsValidID(i)
}

func isListID(v interface{}) bool {
	ls, ok := v.([]string)
	if !ok {
		return false
	}
	for _, i := range ls {
		valid := dbutils.IsValidID(i)
		if !valid {
			return false
		}
	}
	return true
}

func isListString(v interface{}) bool {
	_, ok := v.([]string)
	return ok
}
