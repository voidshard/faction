/* sql_filters.go turns filters into SQL queries. */
package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/voidshard/faction/internal/dbutils"
	"github.com/voidshard/faction/pkg/structs"
)

var (
	fieldTypes = map[Op][]isValid{
		Equal:    {isInt, isString, isBool},
		NotEqual: {isInt, isString},
		In:       {isListID},
		Greater:  {isInt},
		Less:     {isInt},
	}
	colChecks = map[Field][]isValid{
		ID:                 {isID, isListID},
		JobID:              {isID, isListID},
		AreaID:             {isID, isListID},
		GovernmentID:       {isID, isListID},
		FactionID:          {isID, isListID},
		SourceAreaID:       {isID, isListID},
		TargetAreaID:       {isID, isListID},
		PreferredFactionID: {isID, isListID},
		BirthFamilyID:      {isID, isListID},
		SourceFactionID:    {isID, isListID},
		TargetMetaKey:      {isMetaKey},
	}
	metaKeys = map[string]bool{}
)

func init() {
	for _, key := range []structs.MetaKey{
		structs.MetaKeyPerson,
		structs.MetaKeyPlot,
		structs.MetaKeyResearch,
		structs.MetaKeyFaction,
		structs.MetaKeyReligion,
		structs.MetaKeyGovernment,
		structs.MetaKeyFamily,
		structs.MetaKeyCommodity,
		structs.MetaKeyAction,
		structs.MetaKeyArea,
		structs.MetaKeyJob,
		structs.MetaKeyRoute,
	} {
		metaKeys[string(key)] = true
	}
}

func (q *Query) sqlQuery(offset int) (string, []interface{}, error) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, fset := range q.filters {
		if fset == nil || len(fset) == 0 {
			continue
		}

		ands := []string{}
		for _, f := range fset {
			segment, fargs, err := f.sqlQuery(offset + len(args))
			if err != nil {
				return "", nil, err
			}

			args = append(args, fargs...)
			ands = append(ands, segment)
		}
		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 { // at least one subject, object must be passed in
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return where, args, nil
}

func (f *Filter) validate() error {
	if f.Field == "" {
		return fmt.Errorf("invalid filter field: %s", f.Field)
	}
	if f.Op == "" {
		return fmt.Errorf("invalid filter op: %s", f.Op)
	}
	if f.Value == nil {
		return fmt.Errorf("invalid filter value: %v", f.Value)
	}

	// check the value vs. the operation we want to do
	checks, ok := fieldTypes[f.Op]
	if !ok {
		return fmt.Errorf("invalid filter operation '%s'", f.Op)
	}
	valid := false
	for _, check := range checks {
		if check(f.Value) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid filter value %v (failed check for '%s' operation)", f.Value, f.Op)
	}

	// check the value vs. the column we're talking about
	checks, ok = colChecks[f.Field]
	if !ok {
		// we don't have checks for all columns, that's ok
		return nil
	}

	valid = false
	for _, check := range checks {
		if check(f.Value) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid filter value %v for field %s", f.Value, f.Field)
	}
	return nil
}

func (f *Filter) sqlQuery(offset int) (string, []interface{}, error) {
	err := f.validate()
	if err != nil {
		return "", nil, err
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

	col := f.sqlColumn()
	if col == "" {
		return "", nil, fmt.Errorf("invalid filter column: %s", f.Field)
	}

	return fmt.Sprintf("%s %s %s", col, f.Op, placeholder), args, nil
}

func (f *Filter) sqlColumn() string {
	switch f.Field {
	case ID:
		return "id"
	case BirthFamilyID:
		return "birth_family_id"
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
		return "is_child_bearing"
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

func isBool(v interface{}) bool {
	_, ok := v.(bool)
	return ok
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

func isMetaKey(v interface{}) bool {
	i, ok := v.(string)
	if !ok {
		return false
	}
	_, ok = metaKeys[i]
	return ok
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

func sqlSummationTuplesFromModifiers(r Relation, tk *dbutils.IterToken, q *Query) (string, []interface{}, error) {
	where, args, err := q.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if q.sort {
		order = "ORDER BY subject"
	}

	return fmt.Sprintf(`SELECT
		    subject, object, sum(value) as value
		FROM %s
		%s 
		GROUP BY subject, object
		%s
		LIMIT $%d OFFSET $%d;`,
		r.modTable(), where, order, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromModifierFilters(r Relation, tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY subject"
	}

	return fmt.Sprintf(`SELECT * FROM %s %s %s LIMIT $%d OFFSET $%d;`,
		r.modTable(), where, order, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromTupleFilters(r Relation, tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY subject"
	}

	return fmt.Sprintf(
		"SELECT * FROM %s %s %s LIMIT $%d OFFSET $%d;",
		r.tupleTable(), where, order, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromGovernmentFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(
		"SELECT * FROM %s %s %s LIMIT $%d OFFSET $%d;",
		tableGovernments, where, order, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlLawsFromGovernmentIDs(in []string) (string, []interface{}) {
	inGovtIDs, args := sqlIn(in)
	return fmt.Sprintf(`SELECT *
	    FROM %s WHERE government_id IN (%s);
	`, tableLaws, inGovtIDs), args
}

func sqlIn(in []string) (string, []interface{}) {
	bindvars := make([]string, len(in))
	args := make([]interface{}, len(in))
	for i, id := range in {
		args[i] = id
		bindvars[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(bindvars, ", "), args
}

func sqlFromRouteFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY source_area_id"
	}

	return fmt.Sprintf(`SELECT * 
	    FROM %s %s %s LIMIT $%d OFFSET $%d;`,
		tableRoutes, where, order, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromJobFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(`SELECT *
	    FROM %s
	    %s
	    %s
	    LIMIT $%d OFFSET $%d`, tableJobs, where, order, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromFamilyFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(`SELECT * FROM %s %s %s LIMIT $%d OFFSET $%d;`,
		tableFamilies,
		where,
		order,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromFactionFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(`SELECT * 
	    FROM %s %s %s LIMIT $%d OFFSET $%d;`,
		tableFactions,
		where,
		order,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromPlotFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(`SELECT *
	    FROM %s
	    %s
	    %s LIMIT $%d OFFSET $%d;`,
		tablePlots,
		where,
		order,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromPersonFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(`SELECT *
		FROM %s
		%s
		%s LIMIT $%d OFFSET $%d;`,
		tablePeople,
		where,
		order,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromAreaFilters(tk *dbutils.IterToken, in *Query) (string, []interface{}, error) {
	where, args, err := in.sqlQuery(0)
	if err != nil {
		return "", nil, err
	}

	order := ""
	if in.sort {
		order = "ORDER BY id"
	}

	return fmt.Sprintf(`SELECT *
		FROM %s
		%s 
		%s LIMIT $%d OFFSET $%d;`,
		tableAreas,
		where,
		order,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}
