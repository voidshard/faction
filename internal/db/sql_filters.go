/* sql_filters.go turns filters into SQL queries. */
package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/voidshard/faction/internal/dbutils"
)

func sqlWhereFromFilters(in []*Filter, offset int) (string, []interface{}, error) {
	var (
		ands  []string
		args  []interface{}
		where string
	)

	for i, f := range in {
		segment, fargs, err := f.sqlQuery(i + offset)
		if err != nil {
			return "", nil, err
		}

		args = append(args, fargs...)
		ands = append(ands, segment)
	}

	if len(ands) != 0 { // at least one subject, object must be passed in
		where = fmt.Sprintf("WHERE %s", strings.Join(ands, " AND "))
	}

	return where, args, nil
}

func sqlSummationTuplesFromModifiers(r Relation, tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf(`SELECT
		    subject, object, sum(value) as value
		FROM %s
		%s 
		GROUP BY subject, object
		ORDER BY subject
		LIMIT $%d OFFSET $%d;`,
		r.modTable(), where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromModifierFilters(r Relation, tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT
		    subject, object, value, tick_expires, meta_key, meta_val, meta_reason
		FROM %s %s LIMIT $%d OFFSET $%d;`,
		r.modTable(), where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromTupleFilters(r Relation, tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(
		"SELECT subject, object, value FROM %s %s LIMIT $%d OFFSET $%d;",
		r.tupleTable(), where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromGovernmentFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(
		"SELECT id, tax_rate, tax_frequency FROM %s %s LIMIT $%d OFFSET $%d;",
		tableGovernments, where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlLawsFromGovernmentIDs(in []string) (string, []interface{}) {
	inGovtIDs, args := sqlIn(in)
	return fmt.Sprintf(`SELECT government_id, meta_key, meta_val, illegal
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

func sqlFromRouteFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT 
	    source_area_id, target_area_id, travel_time
	    FROM %s %s LIMIT $%d OFFSET $%d;`,
		tableRoutes, where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromJobFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT id, parent_job_id,
		source_faction_id, source_area_id,
		action,
		target_area_id, target_meta_key, target_meta_val,
		people_min,
		people_max,
		tick_created, tick_starts, tick_ends,
		secrecy,
		is_illegal,
		state
	    FROM %s
	    %s
	    ORDER BY id LIMIT $%d OFFSET $%d`, tableJobs, where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromFamilyFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT id,
		ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
		area_id,
		faction_id,
		is_child_bearing,
		max_child_bearing_tick,
		pregnancy_end,
		male_id,
		female_id,
		ma_grandma_id, ma_grandpa_id, pa_grandma_id, pa_grandpa_id,
		number_of_children
	    FROM %s %s ORDER BY id LIMIT $%d OFFSET $%d;`,
		tableFamilies,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromFactionFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT id, name,
		home_area_id,
		ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
		action_frequency_ticks,
		leadership,
		structure,
		wealth,
		cohesion,
		is_covert,
		government_id,
		is_government,
		religion_id,
		is_religion,
		is_member_by_birth,
		parent_faction_id,
		parent_faction_relation
	    FROM %s %s ORDER BY id LIMIT $%d OFFSET $%d;`,
		tableFactions,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromPlotFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT
		id, area_id, faction_id, size, commodity, yield
	    FROM %s
	    %s
	    ORDER BY id LIMIT $%d OFFSET $%d;`,
		tablePlots,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromPersonFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT
		    id, birth_family_id,
		    first_name, last_name,
		    ethos_altruism, ethos_ambition, ethos_tradition, ethos_pacifism, ethos_piety, ethos_caution,
		    area_id, job_id,
		    preferred_profession, preferred_faction_id,
		    birth_tick, death_tick,
		    is_male, is_child,
		    death_meta_reason, death_meta_key, death_meta_val
		FROM %s
		%s
		ORDER BY id ASC LIMIT $%d OFFSET $%d;`,
		tablePeople,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}

func sqlFromAreaFilters(tk *dbutils.IterToken, in []*Filter) (string, []interface{}, error) {
	where, args, err := sqlWhereFromFilters(in, 0)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf(`SELECT id, government_id
		FROM %s
		%s 
		ORDER BY id ASC LIMIT $%d OFFSET $%d;`,
		tableAreas,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset), nil
}
