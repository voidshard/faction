/* sql_filters.go turns filters into SQL queries. */
package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/voidshard/faction/internal/dbutils"
)

func sqlSummationTuplesFromModifiers(r Relation, tk *dbutils.IterToken, in []*ModifierFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if f.Subject != "" {
			args = append(args, f.Subject)
			ands = append(ands, fmt.Sprintf("subject = $%d", len(args)))
		}
		if f.Object != "" {
			args = append(args, f.Object)
			ands = append(ands, fmt.Sprintf("object = $%d", len(args)))
		}
		if f.MinTickExpires > 0 {
			args = append(args, f.MinTickExpires)
			ands = append(ands, fmt.Sprintf("tick_expires >= $%d", len(args)))
		}
		if f.MaxTickExpires > 0 {
			args = append(args, f.MaxTickExpires)
			ands = append(ands, fmt.Sprintf("tick_expires <= $%d", len(args)))
		}
		if f.MetaKey != "" {
			args = append(args, f.MetaKey)
			ands = append(ands, fmt.Sprintf("meta_key = $%d", len(args)))
		}
		if f.MetaVal != "" {
			args = append(args, f.MetaVal)
			ands = append(ands, fmt.Sprintf("meta_val = $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) != 0 { // at least one subject, object must be passed in
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT
		    subject, object, sum(value) as value
		FROM %s
		%s 
		GROUP BY subject, object
		ORDER BY subject
		LIMIT $%d OFFSET $%d;`,
		r.modTable(), where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset)
}

func sqlWhereFromModifierFilters(in []*ModifierFilter, offset int) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if f.Subject != "" {
			args = append(args, f.Subject)
			ands = append(ands, fmt.Sprintf("subject = $%d", len(args)+offset))
		}
		if f.Object != "" {
			args = append(args, f.Object)
			ands = append(ands, fmt.Sprintf("object = $%d", len(args)+offset))
		}
		if f.MinTickExpires > 0 {
			args = append(args, f.MinTickExpires)
			ands = append(ands, fmt.Sprintf("tick_expires >= $%d", len(args)+offset))
		}
		if f.MaxTickExpires > 0 {
			args = append(args, f.MaxTickExpires)
			ands = append(ands, fmt.Sprintf("tick_expires <= $%d", len(args)+offset))
		}
		if f.MetaKey != "" {
			args = append(args, f.MetaKey)
			ands = append(ands, fmt.Sprintf("meta_key = $%d", len(args)+offset))
		}
		if f.MetaVal != "" {
			args = append(args, f.MetaVal)
			ands = append(ands, fmt.Sprintf("meta_val = $%d", len(args)+offset))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return where, args
}

func sqlFromModifierFilters(r Relation, tk *dbutils.IterToken, in []*ModifierFilter) (string, []interface{}) {
	where, args := sqlWhereFromModifierFilters(in, 0)
	return fmt.Sprintf(`SELECT
		    subject, object, value, tick_expires, meta_key, meta_val, meta_reason
		FROM %s %s LIMIT $%d OFFSET $%d;`,
		r.modTable(), where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset)
}

func sqlWhereFromTupleFilters(in []*TupleFilter, offset int) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if f.Subject != "" {
			args = append(args, f.Subject)
			ands = append(ands, fmt.Sprintf("subject = $%d", len(args)+offset))
		}
		if f.Object != "" {
			args = append(args, f.Object)
			ands = append(ands, fmt.Sprintf("object = $%d", len(args)+offset))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return where, args
}

func sqlFromTupleFilters(r Relation, tk *dbutils.IterToken, in []*TupleFilter) (string, []interface{}) {
	where, args := sqlWhereFromTupleFilters(in, 0)
	return fmt.Sprintf(
		"SELECT subject, object, value FROM %s %s LIMIT $%d OFFSET $%d;",
		r.tupleTable(), where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromGovernmentFilters(tk *dbutils.IterToken, in []*GovernmentFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(
		"SELECT id, tax_rate, tax_frequency FROM %s %s LIMIT $%d OFFSET $%d;",
		tableGovernments, where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset)
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

func sqlFromRouteFilters(tk *dbutils.IterToken, in []*RouteFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.SourceAreaID) {
			args = append(args, f.SourceAreaID)
			ands = append(ands, fmt.Sprintf("source_area_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.TargetAreaID) {
			args = append(args, f.TargetAreaID)
			ands = append(ands, fmt.Sprintf("target_area_id = $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT 
	    source_area_id, target_area_id, travel_time
	    FROM %s %s LIMIT $%d OFFSET $%d;`,
		tableRoutes, where, len(args)+1, len(args)+2,
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromJobFilters(tk *dbutils.IterToken, in []*JobFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.SourceFactionID) {
			args = append(args, f.SourceFactionID)
			ands = append(ands, fmt.Sprintf("source_faction_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.SourceAreaID) {
			args = append(args, f.SourceAreaID)
			ands = append(ands, fmt.Sprintf("source_area_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.TargetAreaID) {
			args = append(args, f.TargetAreaID)
			ands = append(ands, fmt.Sprintf("target_area_id = $%d", len(args)))
		}
		if f.TargetMetaKey != "" {
			args = append(args, f.TargetMetaKey)
			ands = append(ands, fmt.Sprintf("target_meta_key = $%d", len(args)))
		}
		if f.TargetMetaVal != "" {
			args = append(args, f.TargetMetaVal)
			ands = append(ands, fmt.Sprintf("target_meta_val = $%d", len(args)))
		}
		if f.MinSecrecy >= 0 {
			args = append(args, f.MinSecrecy)
			ands = append(ands, fmt.Sprintf("secrecy >= $%d", len(args)))
		}
		if f.MaxSecrecy > 0 {
			args = append(args, f.MaxSecrecy)
			ands = append(ands, fmt.Sprintf("secrecy <= $%d", len(args)))
		}
		if f.State != "" {
			args = append(args, f.State)
			ands = append(ands, fmt.Sprintf("state = $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
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
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromFamilyFilters(tk *dbutils.IterToken, in []*FamilyFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.FactionID) {
			args = append(args, f.FactionID)
			ands = append(ands, fmt.Sprintf("faction_id = $%d", len(args)))
		}
		if f.OnlyChildBearing {
			ands = append(ands, fmt.Sprintf("is_child_bearing = 1"))
		}
		if dbutils.IsValidID(f.MaleID) {
			args = append(args, f.MaleID)
			ands = append(ands, fmt.Sprintf("male_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.FemaleID) {
			args = append(args, f.FemaleID)
			ands = append(ands, fmt.Sprintf("female_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.AreaID) {
			args = append(args, f.AreaID)
			ands = append(ands, fmt.Sprintf("area_id = $%d", len(args)))
		}
		if f.PregnancyEnd > 0 {
			args = append(args, f.PregnancyEnd)
			ands = append(ands, fmt.Sprintf("pregnancy_end = $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT id
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
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromFactionFilters(tk *dbutils.IterToken, in []*FactionFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.ReligionID) {
			args = append(args, f.ReligionID)
			ands = append(ands, fmt.Sprintf("religion_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.GovernmentID) {
			args = append(args, f.GovernmentID)
			ands = append(ands, fmt.Sprintf("government_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.ParentFactionID) {
			args = append(args, f.ParentFactionID)
			ands = append(ands, fmt.Sprintf("parent_faction_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.HomeAreaID) {
			args = append(args, f.HomeAreaID)
			ands = append(ands, fmt.Sprintf("home_area_id = $%d", len(args)))
		}
		if f.MinEthos != nil {
			args = append(args, f.MinEthos.Altruism)
			ands = append(ands, fmt.Sprintf("ethos_altruism >= $%d", len(args)))
			args = append(args, f.MinEthos.Ambition)
			ands = append(ands, fmt.Sprintf("ethos_ambition >= $%d", len(args)))
			args = append(args, f.MinEthos.Tradition)
			ands = append(ands, fmt.Sprintf("ethos_tradition >= $%d", len(args)))
			args = append(args, f.MinEthos.Pacifism)
			ands = append(ands, fmt.Sprintf("ethos_pacifism >= $%d", len(args)))
			args = append(args, f.MinEthos.Piety)
			ands = append(ands, fmt.Sprintf("ethos_piety >= $%d", len(args)))
			args = append(args, f.MinEthos.Caution)
			ands = append(ands, fmt.Sprintf("ethos_caution >= $%d", len(args)))
		}
		if f.MaxEthos != nil {
			args = append(args, f.MaxEthos.Altruism)
			ands = append(ands, fmt.Sprintf("ethos_altruism <= $%d", len(args)))
			args = append(args, f.MaxEthos.Ambition)
			ands = append(ands, fmt.Sprintf("ethos_ambition <= $%d", len(args)))
			args = append(args, f.MaxEthos.Tradition)
			ands = append(ands, fmt.Sprintf("ethos_tradition <= $%d", len(args)))
			args = append(args, f.MaxEthos.Pacifism)
			ands = append(ands, fmt.Sprintf("ethos_pacifism <= $%d", len(args)))
			args = append(args, f.MaxEthos.Piety)
			ands = append(ands, fmt.Sprintf("ethos_piety <= $%d", len(args)))
			args = append(args, f.MaxEthos.Caution)
			ands = append(ands, fmt.Sprintf("ethos_caution <= $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
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
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromPlotFilters(tk *dbutils.IterToken, in []*PlotFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.FactionID) {
			args = append(args, f.FactionID)
			ands = append(ands, fmt.Sprintf("faction_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.AreaID) {
			args = append(args, f.AreaID)
			ands = append(ands, fmt.Sprintf("area_id = $%d", len(args)))
		}
		if f.HasCommodity {
			ands = append(ands, "(commodity <> '')")
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
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
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromPersonFilters(tk *dbutils.IterToken, in []*PersonFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.JobID) {
			args = append(args, f.JobID)
			ands = append(ands, fmt.Sprintf("job_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.AreaID) {
			args = append(args, f.AreaID)
			ands = append(ands, fmt.Sprintf("area_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.BirthFamilyID) {
			args = append(args, f.BirthFamilyID)
			ands = append(ands, fmt.Sprintf("birth_family_id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.PreferredFactionID) {
			args = append(args, f.PreferredFactionID)
			ands = append(ands, fmt.Sprintf("preferred_faction_id = $%d", len(args)))
		}
		if f.PreferredProfession != "" {
			args = append(args, f.PreferredProfession)
			ands = append(ands, fmt.Sprintf("preferred_profession = $%d", len(args)))
		}
		if f.Race != "" {
			args = append(args, f.Race)
			ands = append(ands, fmt.Sprintf("race = $%d", len(args)))
		}
		if !f.IncludeDead { // add this by default
			ands = append(ands, fmt.Sprintf("death_tick <= 0"))
		}
		if !f.IncludeChildren { // add this by default
			ands = append(ands, fmt.Sprintf("is_child = 0"))
		}
		if f.MinEthos != nil {
			args = append(args, f.MinEthos.Altruism)
			ands = append(ands, fmt.Sprintf("ethos_altruism >= $%d", len(args)))
			args = append(args, f.MinEthos.Ambition)
			ands = append(ands, fmt.Sprintf("ethos_ambition >= $%d", len(args)))
			args = append(args, f.MinEthos.Tradition)
			ands = append(ands, fmt.Sprintf("ethos_tradition >= $%d", len(args)))
			args = append(args, f.MinEthos.Pacifism)
			ands = append(ands, fmt.Sprintf("ethos_pacifism >= $%d", len(args)))
			args = append(args, f.MinEthos.Piety)
			ands = append(ands, fmt.Sprintf("ethos_piety >= $%d", len(args)))
			args = append(args, f.MinEthos.Caution)
			ands = append(ands, fmt.Sprintf("ethos_caution >= $%d", len(args)))
		}
		if f.MaxEthos != nil {
			args = append(args, f.MaxEthos.Altruism)
			ands = append(ands, fmt.Sprintf("ethos_altruism <= $%d", len(args)))
			args = append(args, f.MaxEthos.Ambition)
			ands = append(ands, fmt.Sprintf("ethos_ambition <= $%d", len(args)))
			args = append(args, f.MaxEthos.Tradition)
			ands = append(ands, fmt.Sprintf("ethos_tradition <= $%d", len(args)))
			args = append(args, f.MaxEthos.Pacifism)
			ands = append(ands, fmt.Sprintf("ethos_pacifism <= $%d", len(args)))
			args = append(args, f.MaxEthos.Piety)
			ands = append(ands, fmt.Sprintf("ethos_piety <= $%d", len(args)))
			args = append(args, f.MaxEthos.Caution)
			ands = append(ands, fmt.Sprintf("ethos_caution <= $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
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
	), append(args, tk.Limit, tk.Offset)
}

func sqlFromAreaFilters(tk *dbutils.IterToken, in []*AreaFilter) (string, []interface{}) {
	var (
		ors   []string
		args  []interface{}
		where string
	)

	for _, f := range in {
		ands := []string{}

		if dbutils.IsValidID(f.ID) {
			args = append(args, f.ID)
			ands = append(ands, fmt.Sprintf("id = $%d", len(args)))
		}
		if dbutils.IsValidID(f.GovernmentID) {
			args = append(args, f.GovernmentID)
			ands = append(ands, fmt.Sprintf("government_id = $%d", len(args)))
		}

		if len(ands) > 0 {
			ors = append(ors, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
		}
	}

	if len(ors) > 0 {
		where = fmt.Sprintf("WHERE %s", strings.Join(ors, " OR "))
	}

	return fmt.Sprintf(`SELECT id, government_id
		FROM %s
		%s 
		ORDER BY id ASC LIMIT $%d OFFSET $%d;`,
		tableAreas,
		where,
		len(args)+1,
		len(args)+2,
	), append(args, tk.Limit, tk.Offset)
}
